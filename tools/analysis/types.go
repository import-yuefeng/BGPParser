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

type hashcodeSlice []byte

type BGPInfo struct {
	Prefix   []string
	Hashcode string
	AsPath   string
}

type Bit uint8

const (
	MINIPV4CIDR int = 24
	MINIPV6CIDR int = 48
	IPTreeLEVEL int = 128
)

type BGPBST struct {
	root     *IPAddr
	v4root   *IPAddr
	inBackup sync.RWMutex
}

type IPAddr struct {
	id          string
	Left, Right *IPAddr
	Hashcode    string
	Prefix      string
	lock        sync.Mutex
}

func NewIPAddr() *IPAddr {
	return &IPAddr{
		lock: sync.Mutex{},
	}
}

func NewBGPBST() *BGPBST {
	root := &BGPBST{
		root:     NewIPAddr(),
		v4root:   NewIPAddr(),
		inBackup: sync.RWMutex{},
	}
	root.v4root.Hashcode = "ipv4Root"
	/*
		IPv4转译地址
		::ffff:x.x.x.x/96－用于IPv4映射地址。
	**/
	var initIPTree []byte = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
	}

	bs := make([]byte, 96)
	count := 95
	for index := 0; index < 12; index++ {
		flag := 1
		if len(initIPTree) <= index {
			break
		}
		cur := int(initIPTree[index])
		for i := 0; i < 8; i++ {
			if cur&flag == 1 {
				bs[count] = byte(1)
			} else if cur&flag == 0 {
				bs[count] = byte(0)
			}
			cur >>= 1
			count--
		}
	}
	curNode := root.root
	cidr := 95
	for i := 0; i < cidr; i++ {
		curNode.lock.Lock()
		if i < cidr-1 {
			if bs[i] == 0 {
				if curNode.Left == nil {
					curNode.Left = NewIPAddr()
				}
				next := curNode.Left
				curNode.lock.Unlock()
				curNode = next
			} else if bs[i] == 1 {
				if curNode.Right == nil {
					curNode.Right = NewIPAddr()
				}
				next := curNode.Right
				curNode.lock.Unlock()
				curNode = next
			}
		} else {
			if bs[i] == 0 {
				curNode.Left = root.v4root
			} else if bs[i] == 1 {
				curNode.Right = root.v4root
			}
			curNode.lock.Unlock()
		}
	}
	return root
}
