package daemon

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	gobgpdump "github.com/CSUNetSec/gobgpdump"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
)

type MetaData struct {
	AsPathMap sync.Map
}

func readBGPData(fileName string, ch chan *string, cancel context.CancelFunc) error {
	bgpFP, err := os.Open(fileName)
	if err != nil {
		log.Traceln(err)
		return err
	}

	defer bgpFP.Close()
	defer cancel()
	defer close(ch)

	reader := bufio.NewReader(bgpFP)

	var segment bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Infoln("File read ok!")
				return nil
			}
			log.Warnln("Read file error!", err)
			return err
		}
		if !strings.Contains(line, "MRT") {
			if _, err := segment.WriteString(line); err != nil {
				log.Traceln(err)
				return err
			}
		} else {
			tmp := segment.String()
			segment.Reset()
			ch <- &tmp
		}
	}
}

func (md *MetaData) addAspath(a1 *analysis.BGPInfo) {
	if tmp, ok := md.AsPathMap.LoadOrStore(a1.Hashcode, a1); ok {
		if res, ok := tmp.(*analysis.BGPInfo); ok {
			res.Prefix = append(res.Prefix, a1.Prefix[0])
		}
	}
}

func (md *MetaData) parseBGPData(fileName string, parserWC int) *analysis.BGPBST {

	ch := make(chan *string, 0)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(parserWC * 1000)
	go readBGPData(fileName, ch, cancel)

	for i := 0; i < parserWC*1000; i++ {
		go func(md *MetaData, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					if len(ch) == 0 {
						wg.Done()
						return
					}
				case line, ok := <-ch:
					if ok {
						if len(*line) == 0 {
							continue
						}
						a1 := analysis.NewBGPInfo(*line)
						line = nil
						a1.AnalysisBGPData()
						md.addAspath(a1)
					}
				}
			}
		}(md, ctx)
	}
	wg.Wait()
	root := analysis.NewBGPBST()
	md.AsPathMap.Range(func(k, v interface{}) bool {
		if t, ok := v.(*analysis.BGPInfo); ok {
			root.Insert(t)
			log.Infoln(t.Prefix)
		}
		return true
	})

	return root
}

func parseBGPRAWData(configFile gobgpdump.ConfigFile) {
	dc, err := gobgpdump.GetDumpConfig(configFile)
	if err != nil {
		log.Traceln(err)
		return
	}
	dumpStart := time.Now()
	wg := &sync.WaitGroup{}

	for w := 0; w < dc.GetWorkers(); w++ {
		wg.Add(1)
		go gobgpdump.DumpWorker(dc, wg)
	}

	wg.Wait()
	dc.SummarizeAndClose(dumpStart)
}
