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
	"net"
	"strconv"
	"strings"

	"github.com/import-yuefeng/BGPParser/utils"
	log "github.com/sirupsen/logrus"
)

// Search search iptree by input ipaddr
func (r *BGPBST) Search(ipaddr string) (ResPrefix []string, isExist error) {
	if !utils.IsIP(ipaddr) {
		return []string{}, errors.New("IP address invalid")
	}
	r.inBackup.RLock()
	defer r.inBackup.RUnlock()

	ResPrefix = make([]string, 0)

	var rangeNumber int
	var curNode *IPAddr

	bs := splitIPAddr(ipaddr)

	if utils.IsIPV4(ipaddr) {
		rangeNumber = MINIPV4CIDR
		curNode = r.v4root
	} else if utils.IsIPV6(ipaddr) {
		rangeNumber = MINIPV6CIDR
		curNode = r.root
	}

	for i := 0; i < rangeNumber; i++ {
		if curNode == nil {
			break
		}
		if curNode.Prefix != "" {
			ResPrefix = append(ResPrefix, curNode.Prefix)
		}
		if bs[i] == 0 {
			curNode = curNode.Left
		} else {
			curNode = curNode.Right
		}
	}
	if len(ResPrefix) == 0 {
		return ResPrefix, errors.New("Not found")
	}
	return ResPrefix, nil
}

// splitV4Address convert string ipv4 ip address to binary array
func splitV4Address(bs *[]byte, ipaddr string) error {
	count := 0
	ip := net.ParseIP(ipaddr)

	for index := 12; index < 16; index++ {
		flag := 1 << 7
		if len(ip) <= index {
			break
		}
		cur := int(ip[index])
		for i := 0; i < 8; i++ {
			if cur&flag != 0 {
				(*bs)[count] = byte(1)
			} else {
				(*bs)[count] = byte(0)
			}
			flag >>= 1
			count++
		}
	}
	return nil
}

// splitV6Address convert string ipv6 ip address to binary array
func splitV6Address(bs *[]byte, ipaddr string) error {
	/*
		IPv6二进位制下为128位长度，以16位为一组，每组以冒号“:”隔开，可以分为8组，每组以4位十六进制方式表示。
		例如：2001:0db8:86a3:08d3:1319:8a2e:0370:7344 是一个合法的IPv6地址。
		类似于IPv4的点分十进制，同样也存在点分十六进制的写法，
		将8组4位十六进制地址的冒号去除后，每位以点号“.”分组，
		例如：2001:0db8:85a3:08d3:1319:8a2e:0370:7344
		则记为2.0.0.1.0.d.b.8.8.5.a.3.0.8.d.3.1.3.1.9.8.a.2.e.0.3.7.0.7.3.4.4，
		其倒序写法用于ip6.arpa子域名记录IPv6地址与域名的映射。

		每项数字前导的0可以省略，省略后前导数字仍是0则继续，例如下组IPv6是等价的。
		2001:0DB8:02de:0000:0000:0000:0000:0e13
		2001:DB8:2de:0000:0000:0000:0000:e13
		2001:DB8:2de:000:000:000:000:e13
		2001:DB8:2de:00:00:00:00:e13
		2001:DB8:2de:0:0:0:0:e13
		可以用双冒号“::”表示一组0或多组连续的0，但只能出现一次：
		如果四组数字都是零，可以被省略。遵照以上省略规则，下面这两组IPv6都是相等的。
		2001:DB8:2de:0:0:0:0:e13
		2001:DB8:2de::e13
		2001:0DB8:0000:0000:0000:0000:1428:57ab
		2001:0DB8:0000:0000:0000::1428:57ab
		2001:0DB8:0:0:0:0:1428:57ab
		2001:0DB8:0::0:1428:57ab
		2001:0DB8::1428:57ab
	**/
	count := 0
	ip := net.ParseIP(ipaddr)

	for index := 0; index < 6; index++ {
		flag := 1 << 7
		if len(ip) <= index {
			break
		}
		cur := int(ip[index])
		for i := 0; i < 8; i++ {
			if cur&flag != 0 {
				(*bs)[count] = byte(1)
			} else {
				(*bs)[count] = byte(0)
			}
			flag >>= 1
			count++
		}
	}
	return nil
}

func splitIPAddr(ipaddr string) []byte {
	var bs []byte
	var v4 bool
	if utils.IsIPV4(ipaddr) {
		bs = make([]byte, 24)
		v4 = true
	} else if utils.IsIPV6(ipaddr) {
		bs = make([]byte, 48)
	}
	if v4 {
		splitV4Address(&bs, ipaddr)
	} else {
		splitV6Address(&bs, ipaddr)
	}
	return bs
}

// Insert new item on iptree by input ipaddr
func (r *BGPBST) Insert(b *BGPInfo) {

	r.inBackup.RLock()
	defer r.inBackup.RUnlock()
	for idx, _ := range b.Prefix {
		ipSegment := b.Prefix[idx]
		if !utils.IsIP(ipSegment) {
			continue
		}
		if _, _, err := net.ParseCIDR(ipSegment); err != nil {
			continue
		}

		tmp := strings.Split(ipSegment, "/")
		if len(tmp) <= 1 {
			log.Warnln("syntaxError: ", tmp)
			return
		}
		ipaddr := tmp[0]
		cidr, err := strconv.Atoi(tmp[1])
		if err != nil {
			log.Warnln("cidr2int error: ", err, tmp)
		}
		if (cidr > 24 && utils.IsIPV4(ipaddr)) || (cidr > 48 && utils.IsIPV6(ipaddr)) || cidr <= 0 {
			continue
		}

		var curNode *IPAddr

		bs := splitIPAddr(ipaddr)

		if utils.IsIPV4(ipaddr) {
			curNode = r.v4root
		} else if utils.IsIPV6(ipaddr) {
			curNode = r.root
		}

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
				} else {
					if curNode.Right == nil {
						curNode.Right = NewIPAddr()
					}
					next := curNode.Right
					curNode.lock.Unlock()
					curNode = next
				}
			} else {
				curNode.Hashcode = b.Hashcode
				curNode.Prefix = ipSegment
				curNode.lock.Unlock()
			}
		}

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

// EncodeIPTree by morris
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

// GetRoot return IPTree root
func (r *BGPBST) GetRoot() *IPAddr {
	return r.root
}

// SetRoot set IPTree root by input
func (r *BGPBST) SetRoot(root *IPAddr) {
	r.root = root
	return
}

// GetID return IPAddr.id
func (i *IPAddr) GetID() string {
	return i.id
}
