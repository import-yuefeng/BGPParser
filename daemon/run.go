package daemon

import (
	log "github.com/sirupsen/logrus"
	"os"
	"io"
	"time"
	"sync"
	"bufio"
	"bytes"
	"strings"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	gobgpdump "github.com/CSUNetSec/gobgpdump"
)

type MetaData struct {
	AsPathMap map[string]*analysis.BGPInfo
	rw *sync.RWMutex
}

func readBGPData(fileName string, ch chan *string) {
	bgpFP, err := os.Open(fileName)
	if err != nil {
		log.Traceln(err)
		return
	}

	defer bgpFP.Close()
	reader := bufio.NewReader(bgpFP)

	var segment bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Infoln("File read ok!")
				break
			} else {
				log.Warnln("Read file error!", err)
				return
			}
		}
		if !strings.Contains(line, "MRT") {
			if _, err := segment.WriteString(line); err != nil {
				log.Traceln(err)
				return
			}
		} else {
			tmp := segment.String()
			segment.Reset()
			ch <- &tmp
		}
	}
}

func (a *MetaData) addAspath(a1 *analysis.BGPInfo) {
	a.rw.Lock()
	if res, ok := a.AsPathMap[a1.Hashcode]; ok {
		res.Prefix = append(res.Prefix, a1.Prefix[0])
	} else {
		a.AsPathMap[a1.Hashcode] = a1
	}
	a.rw.Unlock()
}

func (md *MetaData) parseBGPData(fileName string, parserWC int) {

	ch := make(chan *string, parserWC*1000)
	for i:=0; i<parserWC*1000; i++ {
		go func(md *MetaData) {
			for {
				if line, ok := <- ch; ok {
					if len(*line) == 0 {
						continue
					}
					a1 := analysis.NewBGPInfo(*line)
					a1.AnalysisBGPData()
					// BGPInfoChannel<-a1
					md.addAspath(a1)
				}
			}
		}(md)
	}

	readBGPData(fileName, ch)
	root := analysis.NewBGPBST()
	if len(ch) == 0 {
		md.rw.Lock()
		for _, v := range md.AsPathMap {
			// fmt.Println(i, v.Prefix)
			root.Insert(v)
		}
		md.rw.Unlock()		
	}
	log.Infoln(root.Search("1.1.1.1"))
	log.Infoln(root.Search("114.114.114.114.114"))
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

