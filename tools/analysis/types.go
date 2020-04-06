package analysis

import (
	"sync"
)

type AspathList []string

type SimpleBGPInfo struct {
	Prefix   []string
	Hashcode string
}

type BGPInfo struct {
	Aspath     AspathList
	Prefix     []string
	Aspath2str string
	Hashcode   string
	isSorted   bool
	content    string
}

type Bit uint8

type BGPBST struct {
	root     *IPAddr
	inBackup sync.RWMutex
}

// type IPSegmentHashcode string

type IPAddr struct {
	bit         Bit
	Left, Right *IPAddr
	Hashcode    string
	lock        sync.Mutex
	id          string
	inValid     bool
}

var bgpInfoFree = sync.Pool{
	New: func() interface{} { return new(BGPInfo) },
}

func CleanBuf(buf *BGPInfo) {
	buf.content = ""
	buf.isSorted = false
	buf.Hashcode = ""
	buf.Aspath2str = ""
	if len(buf.Prefix) != 0 {
		buf.Prefix = buf.Prefix[:1]
	}
	if len(buf.Aspath) != 0 {
		buf.Aspath = buf.Aspath[:1]
	}
}

func newBGPInfo(content string) *BGPInfo {
	buf := bgpInfoFree.Get().(*BGPInfo)
	CleanBuf(buf)
	buf.content = content
	return buf
}

func NewBGPInfo(content string) *BGPInfo {
	return newBGPInfo(content)
}
