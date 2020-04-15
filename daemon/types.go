package daemon

import (
	"sync"

	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
)

type server struct {
}

type Daemon struct {
	logPath  string
	md       *MetaData
	oldmd    *MetaData
	root     *analysis.BGPBST
	oldroot  *analysis.BGPBST
	inUpdate bool
}

type MetaData struct {
	AsPathMap sync.Map
	PrefixMap sync.Map
	TaskList  [][]*analysis.BGPInfo
}

func NewMetaData() *MetaData {
	md := &MetaData{
		AsPathMap: sync.Map{},
		TaskList:  make([][]*analysis.BGPInfo, 16),
	}
	return md
}

func NewDaemon(logPath string, md *MetaData, root *analysis.BGPBST) *Daemon {
	return &Daemon{
		logPath: logPath,
		md:      md,
		root:    root,
	}
}
