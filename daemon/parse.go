// MIT License

// Copyright (c) 2019 Yuefeng Zhu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package daemon

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	gobgpdump "github.com/CSUNetSec/gobgpdump"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
)

func (d *Daemon) Read(fileName string, ch chan *string) error {
	if err := readBGPData(fileName, ch); err != nil {
		return err
	}
	return nil
}

func (d *Daemon) Parse(configFile gobgpdump.ConfigFile) {
	parseBGPRAWData(configFile)
}

func readBGPData(fileName string, ch chan *string) error {
	bgpFP, err := os.Open(fileName)
	if err != nil {
		log.Traceln(err)
		return err
	}
	defer bgpFP.Close()
	defer close(ch)
	reader := bufio.NewReader(bgpFP)
	var segment strings.Builder
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
			line = ""
		} else {
			tmp := segment.String()
			segment.Reset()
			ch <- &tmp
		}
	}
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

func (md *MetaData) addAspath(a *analysis.SimpleBGPInfo) {
	if len(a.Prefix) == 0 {
		return
	}
	if tmp, exist := md.AsPathMap.LoadOrStore(a.Hashcode, a); exist {
		if r, ok := tmp.(*analysis.SimpleBGPInfo); ok {
			r.Prefix = append(r.Prefix, a.Prefix[0])
		}
	} else {
		b, c, d, e := a.Hashcode[0], a.Hashcode[1], a.Hashcode[2], a.Hashcode[3]
		idx := (int(b)<<32 | int(c)<<16 | int(d)<<8 | int(e)) % len(md.TaskList)
		md.TaskList[idx] = append(md.TaskList[idx], a)
	}
}

func (md *MetaData) parseBGPData(fileName string, parserWC int) *analysis.BGPBST {

	ch := make(chan *string, 0)
	var wg sync.WaitGroup
	wg.Add(parserWC * 1000)
	go readBGPData(fileName, ch)

	for i := 0; i < parserWC*1000; i++ {
		go func(md *MetaData) {
			for data := range ch {
				if data == nil || len(*data) == 0 {
					continue
				}
				a1 := analysis.NewBGPInfo(*data)
				*data = ""
				data = nil
				sBGPInfo := a1.AnalysisBGPData()
				md.addAspath(sBGPInfo)
			}
			wg.Done()
			return
		}(md)
	}
	wg.Wait()
	runtime.GC()
	root := analysis.NewBGPBST()
	wg.Add(len(md.TaskList))
	for idx, _ := range md.TaskList {
		go func(taskList []*analysis.SimpleBGPInfo) {
			log.Infoln(len(taskList))
			for idx, _ := range taskList {
				root.Insert(taskList[idx])
			}
			wg.Done()
		}(md.TaskList[idx])
	}
	wg.Wait()
	return root
}
