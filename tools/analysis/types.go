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
	Prefix      string
	lock        sync.Mutex
	id          string
	inValid     bool
}

var bgpInfoFree = sync.Pool{
	New: func() interface{} { return new(BGPInfo) },
}

func NewIPAddr(bit Bit) *IPAddr {
	return &IPAddr{
		bit:  bit,
		lock: sync.Mutex{},
	}
}

func NewBGPBST() *BGPBST {
	root := &BGPBST{
		root:     NewIPAddr(0),
		inBackup: sync.RWMutex{},
	}
	return root
	// t := &SimpleBGPInfo{Hashcode: "000000", Prefix: []string{"0.0.0.0/24"}, }
	// root.Insert()
}

func NewBGPInfo(content string) *BGPInfo {
	return newBGPInfo(content)
}

func newBGPInfo(content string) *BGPInfo {
	buf := bgpInfoFree.Get().(*BGPInfo)
	CleanBuf(buf)
	buf.content = content
	return buf
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
