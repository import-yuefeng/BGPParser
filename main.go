package main

import (
	"fmt"
	"time"
	"sync"
	gobgpdump "github.com/import-yuefeng/BGPParser/gobgpdump"
)

var (
	configFile gobgpdump.ConfigFile
)


func main() {
	configFile.Do = "BGPData.txt"
	configFile.Wc = 8
	dumpFileList := []string{"rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2", "rib.20191212.1400.bz2"}
	parseBGPFile(dumpFileList)
}

func parseBGPFile(dumpFileList []string) {
	dc, err := gobgpdump.GetDumpConfig(configFile, dumpFileList)
	if err != nil {
		fmt.Println(err)
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

