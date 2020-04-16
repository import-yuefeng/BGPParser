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
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"unsafe"

	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
)

// func (d *Daemon) ParseBGPData(fileList []string, parserWC int) *analysis.BGPBST {
// 	return parseBGPData(fileList, parserWC)
// }

func (d *Daemon) ParseRIBData(files []string) {
	for _, file := range files {
		parseRIBData(file)
	}
}

func readBGPData(fileList []string, ch chan *string) error {
	defer close(ch)
	var reader *bufio.Reader
	for _, fileName := range fileList {
		file, err := os.Open(fileName)
		if err != nil {
			log.Traceln(err)
			return err
		}
		reader = bufio.NewReader(file)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					log.Infoln("File read ok!")
					file.Close()
					break
				}
				log.Warnln("Read file error!", err)
			}
			tmp := *(*string)(unsafe.Pointer(&line))
			ch <- &tmp
		}
	}
	return nil
}

func parseRIBData(file string) {
	savePath := file + ".txt"
	cmd := exec.Command("bgpdump", "-m", "-O", savePath, file)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("can not obtain stdout pipe for command:%s\n", err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Fatalln("the command is err,", err)
		return
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatalln("readall stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalln("wait:", err.Error())
		return
	}
	log.Infof("stdout:\n\n %s\n", bytes)
	log.Infof("save file into: %s\n", savePath)
}

func (md *MetaData) addPrefix(a *analysis.BGPInfo) {
	if len(a.Prefix) == 0 {
		return
	}
	if tmp, exist := md.PrefixMap.LoadOrStore(a.Prefix[0], a); exist {
		if r, ok := tmp.(*analysis.BGPInfo); ok {
			r.Hashcode += a.Hashcode
		}
	}
}

func (md *MetaData) addAspath(a *analysis.BGPInfo) {
	if len(a.Prefix) == 0 {
		return
	}
	if tmp, exist := md.AsPathMap.LoadOrStore(a.Hashcode, a); exist {
		if r, ok := tmp.(*analysis.BGPInfo); ok {
			r.Prefix = append(r.Prefix, a.Prefix...)
		}
	} else {
		b := a.Hashcode[0]
		idx := int(b) % len(md.TaskList)
		md.TaskList[idx] = append(md.TaskList[idx], a)
	}
}

func (md *MetaData) parseBGPData(fileList []string, parserWC int) *analysis.BGPBST {

	ch := make(chan *string, 0)
	var wg sync.WaitGroup
	wg.Add(parserWC * 300)
	go readBGPData(fileList, ch)
	for i := 0; i < parserWC*300; i++ {
		go func(md *MetaData) {
			for data := range ch {
				if data == nil || *data == "" {
					continue
				}
				binfo := analysis.NewBGPInfo(*data)
				*data = ""
				data = nil
				md.addPrefix(binfo)
			}
			wg.Done()
			return
		}(md)
	}
	wg.Wait()

	md.PrefixMap.Range(func(k, v interface{}) bool {
		if t, ok := v.(*analysis.BGPInfo); ok {
			go func(*analysis.BGPInfo) {
				t.Hashcode = analysis.PackagingHashcode(t.Hashcode)
				md.addAspath(t)
			}(t)
		}
		return true
	})

	runtime.GC()
	root := analysis.NewBGPBST()
	wg.Add(len(md.TaskList))
	for idx, _ := range md.TaskList {
		go func(taskList []*analysis.BGPInfo) {
			for _, task := range taskList {
				root.Insert(task)
			}
			wg.Done()
		}(md.TaskList[idx])
	}
	wg.Wait()
	return root
}
