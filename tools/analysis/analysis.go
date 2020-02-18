package analysis

// package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

func (r *BGPBST) Search(ipaddr string) (Hashcode string, isExist error) {
	cur := r.root
	bs := getBitIPAddr(ipaddr)
	for i := 0; i < 24; i++ {
		if cur == nil || cur.Hashcode == "" {
			return "", errors.New("Not find")
		}
		if bs[i] == 0 {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return cur.Hashcode, nil
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
			bs[count] = byte(cur & flag)
			flag >>= 1
			count++
		}
	}
	return bs
}

func (r *BGPBST) Insert(b *BGPInfo) {
	root := r.root

	for _, ipSegment := range b.Prefix {
		tmp := strings.Split(ipSegment, "/")
		if len(tmp) == 0 {
			return
		}
		ipv4Address := tmp[0]
		cur := root
		bs := getBitIPAddr(ipv4Address)
		for i := 0; i < 24; i++ {
			if bs[i] == 0 {
				if cur.Left == nil {
					cur.Left = NewIPAddr(0)
				}
				cur = cur.Left
			} else {
				if cur.Right == nil {
					cur.Right = NewIPAddr(1)
				}
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
