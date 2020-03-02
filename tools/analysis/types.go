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

func (b *BGPInfo) AnalysisBGPData() *SimpleBGPInfo {
	b.FindPrefix()
	b.FindAsPath()
	b.SortASpathBySize()
	b.ConvertHashcode()
	res := &SimpleBGPInfo{
		Prefix:   b.Prefix,
		Hashcode: b.Hashcode,
	}
	bgpInfoFree.Put(b)
	return res
}
