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

package marshal

import (
	"bufio"
	"io"
	"os"
	"strings"

	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	utils "github.com/import-yuefeng/BGPParser/utils"
	log "github.com/sirupsen/logrus"
)

func Marshal(root *analysis.BGPBST, path string) {
	if root == nil {
		return
	}
	if len(path) == 0 {
		path = "iptree"
	}
	preOrderList, inOrderList := marshal(root)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Warnln(err)
		return
	}
	for _, v := range preOrderList {
		f.WriteString(WriteIPAddrNode(v))
	}
	f.WriteString("\n")
	for _, v := range inOrderList {
		f.WriteString(WriteIPAddrNode(v))
	}
}

func WriteIPAddrNode(i *analysis.IPAddr) string {
	if i.Hashcode != "" {
		return i.GetID() + "|" + i.Hashcode + " "
	}
	return i.GetID() + " "
}

func Unmarshal(path string) *analysis.BGPBST {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		log.Warnln("Open file error: ", err)
		return nil
	}
	reader := bufio.NewReader(f)
	point := 0
	var preOrderList, inOrderList []string
	for {
		line, err := reader.ReadString('\n')
		point++
		if point&1 == 1 {
			preOrderList = strings.Split(line, " ")
			tmp := preOrderList[len(preOrderList)-1]
			preOrderList[len(preOrderList)-1] = utils.Strip(tmp, "\n")
			log.Infoln(preOrderList)
		} else if point&1 == 0 {
			inOrderList = strings.Split(line, " ")
			tmp := inOrderList[len(inOrderList)-1]
			inOrderList[len(inOrderList)-1] = utils.Strip(tmp, "\n")
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warnln("Read file error!", err)
			return nil
		}
	}
	r := unmarshal(preOrderList, inOrderList)
	bst := &analysis.BGPBST{}
	bst.SetRoot(r)
	return bst
}

func marshal(root *analysis.BGPBST) (preOrderList, inOrderList []*analysis.IPAddr) {
	if root == nil {
		return []*analysis.IPAddr{}, []*analysis.IPAddr{}
	}
	preOrderList = preOrder(root)
	inOrderList = inOrder(root)
	return preOrderList, inOrderList
}

func unmarshal(preOrder, inOrder []string) *analysis.IPAddr {
	if len(preOrder) != len(inOrder) || len(preOrder) == 0 {
		return nil
	}
	id := preOrder[0]
	root := analysis.NewIPAddr(0)

	if len(preOrder) == 1 {
		return root
	}
	point := 0
	for i, v := range inOrder {
		if v == id {
			point = i
			break
		}
	}
	if point >= len(inOrder) || inOrder[point] != id {
		return nil
	}
	if point > 0 {
		root.Left = unmarshal(preOrder[1:], inOrder[:point])
	}
	if point < len(inOrder)-1 {
		root.Right = unmarshal(preOrder[point+1:], inOrder[point+1:])
	}
	return root
}

func preOrder(r *analysis.BGPBST) []*analysis.IPAddr {
	if r == nil {
		return []*analysis.IPAddr{}
	}
	res := make([]*analysis.IPAddr, 0)
	isExist := make(map[*analysis.IPAddr]struct{})
	stack := utils.NewStack()
	root := r.GetRoot()
	stack.Push(root)
	for !stack.IsEmpty() {
		t, err := stack.Pop()
		if err != nil {
			log.Fatalln(err)
			return []*analysis.IPAddr{}
		}
		if r, ok := t.(*analysis.IPAddr); ok {
			if _, exist := isExist[r]; exist {
				res = append(res, r)
			} else {
				isExist[r] = struct{}{}
				if r.Right != nil {
					stack.Push(r.Right)
				}
				if r.Left != nil {
					stack.Push(r.Left)
				}
				if r != nil {
					stack.Push(r)
				}
			}
		} else {
			log.Fatalln(ok)
			return []*analysis.IPAddr{}
		}
	}
	isExist = nil
	stack.Reset()
	return res
}

func inOrder(r *analysis.BGPBST) []*analysis.IPAddr {
	if r == nil {
		return []*analysis.IPAddr{}
	}
	res := make([]*analysis.IPAddr, 0)
	isExist := make(map[*analysis.IPAddr]struct{})
	stack := utils.NewStack()
	root := r.GetRoot()
	stack.Push(root)
	for !stack.IsEmpty() {
		t, err := stack.Pop()
		if err != nil {
			log.Fatalln(err)
			return []*analysis.IPAddr{}
		}
		if r, ok := t.(*analysis.IPAddr); ok {
			if _, exist := isExist[r]; exist {
				res = append(res, r)
			} else {
				isExist[r] = struct{}{}
				if r.Right != nil {
					stack.Push(r.Right)
				}
				if r != nil {
					stack.Push(r)
				}
				if r.Left != nil {
					stack.Push(r.Left)
				}
			}
		} else {
			log.Fatalln(ok)
			return []*analysis.IPAddr{}
		}
	}
	isExist = nil
	stack.Reset()
	return res
}

func printByLevel(r *analysis.BGPBST) [][]*analysis.IPAddr {
	if r == nil {
		return [][]*analysis.IPAddr{}
	}
	root := r.GetRoot()
	res := make([][]*analysis.IPAddr, 0)
	q := utils.NewQueue()
	q.Push(root)
	for !q.IsEmpty() {
		size := q.Size()
		res = append(res, []*analysis.IPAddr{})
		for size > 0 {
			t, err := q.Pop()
			size--
			if err != nil {
				log.Warnln(err)
				return res
			}
			if v, ok := t.(*analysis.IPAddr); !ok {
				log.Warnln(ok)
				return res
			} else {
				if v == nil {
					continue
				}
				res[len(res)-1] = append(res[len(res)-1], v)
				if v.Left != nil {
					q.Push(v.Left)
				}
				if v.Right != nil {
					q.Push(v.Right)
				}
			}
		}
	}
	return res
}

func PrintBGPBST(root *analysis.BGPBST) {
	treeNode := printByLevel(root)
	log.Infoln("treeNode: ", treeNode)
	for i := 0; i < len(treeNode); i++ {
		log.Infoln()
		for j := 0; j < len(treeNode[i]); j++ {
			if treeNode[i][j] == nil || treeNode[i][j].Hashcode == "" {
				continue
			}
			log.Infof("%s   ", treeNode[i][j].Hashcode[:4])
		}
	}
}
