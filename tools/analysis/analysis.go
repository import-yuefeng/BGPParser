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
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (r *BGPBST) Search(ipaddr string) (ResPrefix []string, isExist error) {
	r.inBackup.RLock()
	defer r.inBackup.RUnlock()
	cur := r.root
	bs := getBitIPAddr(ipaddr)
	log.Infoln("Search: ", bs)
	ResPrefix = make([]string, 0)
	for i := 0; i < 24; i++ {
		if cur == nil {
			break
		}
		if cur.Prefix != "" {
			log.Infoln(i, cur.Prefix)
			ResPrefix = append(ResPrefix, cur.Prefix)
		}
		if bs[i] == 0 {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	if len(ResPrefix) == 0 {
		return ResPrefix, errors.New("Not found")
	}
	return ResPrefix, nil
}

func getBitIPAddr(ipaddr string) []byte {
	bs := make([]byte, 24)
	count := 0
	ip := strings.Split(ipaddr, ".")
	for index := 0; index < 3; index++ {
		flag := 1 << 7
		if len(ip) <= index {
			break
		}
		cur, err := strconv.Atoi(ip[index])
		if err != nil {
			log.Traceln(err)
			return bs
		}
		for i := 0; i < 8; i++ {
			if cur&flag != 0 {
				bs[count] = byte(1)
			} else {
				bs[count] = byte(0)
			}
			flag >>= 1
			count++
		}
	}
	return bs
}

func (r *BGPBST) Insert(b *BGPInfo) {
	r.inBackup.RLock()
	defer r.inBackup.RUnlock()
	root := r.root
	// if b == nil || b.Prefix == nil || len(b.Prefix) == 0 {
	// 	return
	// }
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Warnln(err)
	// 		log.Warnln(b.Prefix)
	// 	}
	// }()
	for idx, _ := range b.Prefix {
		ipSegment := b.Prefix[idx]
		if len(ipSegment) == 0 || ipSegment == "" {
			continue
		}
		tmp := strings.Split(ipSegment, "/")
		if len(tmp) <= 1 {
			log.Warnln("syntaxError: ", tmp)
			return
		}
		ipv4Address := tmp[0]
		cidr, err := strconv.Atoi(tmp[1])
		if cidr > 24 {
			continue
		}
		if cidr <= 0 {
			return
		}
		if err != nil {
			log.Warnln("error: ", err, tmp)
			continue
		}
		cur := root
		bs := getBitIPAddr(ipv4Address)
		for i := 0; i < cidr; i++ {
			cur.lock.Lock()
			if bs[i] == 0 {
				if cur.Left == nil {
					cur.Left = NewIPAddr(0)
				}
				cur.lock.Unlock()
				cur = cur.Left
			} else {
				if cur.Right == nil {
					cur.Right = NewIPAddr(1)
				}
				cur.lock.Unlock()
				cur = cur.Right
			}
		}
		cur.Hashcode = b.Hashcode
		cur.Prefix = ipSegment
	}
	b = nil
	return
}

// NewBGPInfo returns the BGPInfo struct.
func NewBGPInfo(content string) *BGPInfo {
	/*
		content(string) Example:
			TABLE_DUMP2|12/12/19 14:00:00|B|217.192.89.50|
			3303|1.0.0.0/24|3303 13335|IGP

			idx(5) = 223.255.254.0/24
			idx(6) = 1299 7473 3758 55415
		**/
	r := strings.Split(content, "|")
	res := &BGPInfo{
		Prefix: []string{r[5]},
		AsPath: r[6],
	}
	res.ConvertHashcode()
	return res
}

func (r *BGPBST) EncodeIPTree() {
	r.inOrderEncodeByMorris()
}

func (r *BGPBST) inOrderEncodeByMorris() {
	r.inBackup.Lock()
	defer r.inBackup.Unlock()
	if r == nil {
		r.inBackup.Unlock()
		return
	}
	count := 0
	cur := r.root
	for cur != nil {
		if cur.Left == nil {
			cur.id = strconv.Itoa(count)
			count++
			cur = cur.Right
		} else {
			tmp := cur.Left
			for tmp.Right != nil && tmp.Right != cur {
				tmp = tmp.Right
			}
			if tmp.Right == nil {
				tmp.Right = cur
				cur = cur.Left
			} else if tmp.Right == cur {
				tmp.Right = nil
				cur.id = strconv.Itoa(count)
				count++
				cur = cur.Right
			}
		}
	}
	return
}

func (r *BGPBST) GetRoot() *IPAddr {
	return r.root
}

func (r *BGPBST) SetRoot(root *IPAddr) {
	r.root = root
	return
}

func (i *IPAddr) GetID() string {
	return i.id
}

func (i *IPAddr) Getbit() Bit {
	return i.bit
}
