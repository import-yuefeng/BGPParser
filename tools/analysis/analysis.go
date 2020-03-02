package analysis

// package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Bit uint8

type BGPBST struct {
	root *IPAddr
}

// type IPSegmentHashcode string

type IPAddr struct {
	bit         Bit
	Left, Right *IPAddr
	Hashcode    string
	lock        sync.Mutex
}

func NewIPAddr(bit Bit) *IPAddr {
	return &IPAddr{
		bit: bit,
	}
}

func NewBGPBST() *BGPBST {
	return &BGPBST{
		root: NewIPAddr(0),
	}
}

func (r *BGPBST) Search(ipaddr string) (ResHashcode []string, isExist error) {
	cur := r.root
	bs := getBitIPAddr(ipaddr)
	log.Println(bs)
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
	if (cur == nil || cur.Hashcode == "") && len(ResHashcode) == 0 {
		return ResHashcode, errors.New("Not find")
	}
	return ResHashcode, nil
}

func getBitIPAddr(ipaddr string) []byte {
	bs := make([]byte, 24)
	count := 0
	ip := strings.Split(ipaddr, ".")
	for index := 0; index < 3; index++ {
		flag := 1 << 7
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
	root := r.root

	for _, ipSegment := range b.Prefix {
		tmp := strings.Split(ipSegment, "/")
		if len(tmp) <= 1 {
			return
		}
		ipv4Address := tmp[0]
		cidr, err := strconv.Atoi(tmp[1])
		if cidr > 24 {
			cidr = 24
		}
		if err != nil {
			log.Warnln(err)
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
}

func a(root *IPAddr) {
	if root == nil {
		return
	}
	fmt.Printf("%d ", root.bit)
	if root.Left == nil && root.Right == nil {
		fmt.Printf("%s ", root.Hashcode)
	}
	a(root.Left)
	a(root.Right)
}
