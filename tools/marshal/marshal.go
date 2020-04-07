package marshal

import (
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	utils "github.com/import-yuefeng/BGPParser/utils"
	log "github.com/sirupsen/logrus"
)

func Marshal(root *analysis.BGPBST) {
	// preOrderList, inOrderList := marshal(root)

}

func Unmarshal(preOrderList, inOrderList []*analysis.IPAddr) *analysis.BGPBST {
	return nil
}

func marshal(root *analysis.BGPBST) (preOrderList, inOrderList []*analysis.IPAddr) {
	if root == nil {
		return []*analysis.IPAddr{}, []*analysis.IPAddr{}
	}
	preOrderList = preOrder(root)
	inOrderList = inOrder(root)
	return preOrderList, inOrderList
}

func unmarshal(preOrder, inOrder []*analysis.IPAddr) *analysis.BGPBST {
	return nil
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
