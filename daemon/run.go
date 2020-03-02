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
	"unsafe"

	gobgpdump "github.com/CSUNetSec/gobgpdump"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
)

type MetaData struct {
	AsPathMap sync.Map
}

func readBGPData(fileName string, ch chan *[]byte, cancel context.CancelFunc) error {
	bgpFP, err := os.Open(fileName)
	if err != nil {
		log.Traceln(err)
		return err
	}

	defer bgpFP.Close()
	defer cancel()
	defer close(ch)

	reader := bufio.NewReader(bgpFP)

	// var segment strings.Builder
	var segment bytes.Buffer
	// segment.Grow(2429)
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
			tmp := segment.Next(segment.Len())
			segment.Reset()
			ch <- &tmp
		}
	}
}

func (md *MetaData) addAspath(a *analysis.SimpleBGPInfo) {
	if len(a.Prefix) == 0 {
		return
	}
	if tmp, ok := md.AsPathMap.LoadOrStore(a.Hashcode, a); ok {
		if res, ok := tmp.(*analysis.SimpleBGPInfo); ok {
			res.Prefix = append(res.Prefix, a.Prefix[0])
		}
	}
}

func (md *MetaData) parseBGPData(fileName string, parserWC int) *analysis.BGPBST {

	ch := make(chan *[]byte, 0)
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
				case data, ok := <-ch:
					if ok {
						if len(*data) == 0 {
							continue
						}
						line := *(*string)(unsafe.Pointer(data))
						*data = (*data)[:0]
						a1 := analysis.NewBGPInfo(line)
						sBGPInfo := a1.AnalysisBGPData()
						md.addAspath(sBGPInfo)
					}
				}
			}
		}(md, ctx)
	}
	wg.Wait()
	root := analysis.NewBGPBST()
	md.AsPathMap.Range(func(k, v interface{}) bool {
		if t, ok := v.(*analysis.SimpleBGPInfo); ok {
			go root.Insert(t)
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
