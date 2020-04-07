package analysis

// package main

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

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

func (r *BGPBST) Search(ipaddr string) (ResHashcode []string, isExist error) {
	r.inBackup.RLock()
	defer r.inBackup.RUnlock()
	cur := r.root
	bs := getBitIPAddr(ipaddr)
	log.Infoln("Search: ", bs)
	ResHashcode = make([]string, 0)
	for i := 0; i < 24; i++ {
		if cur == nil {
			break
		}
		if cur.Hashcode != "" {
			ResHashcode = append(ResHashcode, cur.Hashcode)
		}
		if bs[i] == 0 {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	log.Infoln(ResHashcode)
	if len(ResHashcode) == 0 {
		return ResHashcode, errors.New("Not found")
	}
	return ResHashcode, nil
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

func (r *BGPBST) Insert(b *SimpleBGPInfo) {
	r.inBackup.RLock()
	defer r.inBackup.RUnlock()
	root := r.root
	for _, ipSegment := range b.Prefix {
		tmp := strings.Split(ipSegment, "/")
		if len(tmp) <= 1 {
			log.Warnln("syntaxError: ", tmp)
			return
		}
		ipv4Address := tmp[0]
		cidr, err := strconv.Atoi(tmp[1])
		if cidr > 24 {
			// log.Infoln("cidr: ", cidr)
			cidr = 24
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
	}
	return
}

func (r *BGPBST) GetRoot() *IPAddr {
	return r.root
}

func (b *BGPInfo) AnalysisBGPData() *SimpleBGPInfo {
	b.FindPrefix()
	b.FindAsPath()
	b.SortASpathBySize()
	b.ConvertHashcode()
	res := &SimpleBGPInfo{
		Hashcode: b.Hashcode,
		Prefix:   b.Prefix,
	}
	bgpInfoFree.Put(b)
	return res
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
