package analysis
// package main


import (
	"strings"
	"sort"
)

func (a AspathList) Len() int {
	return len(a)
}
func (a AspathList) Less(i, j int) bool {
	_aspath1 := len(strings.Split(a[i], " "))
	_aspath2 := len(strings.Split(a[j], " "))
	if 	_aspath1 != _aspath2 {
		return _aspath1 < _aspath2
	}
	return a[i] < a[j]
}
func (a AspathList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}


func (b *BGPInfo) SortASpathBySize() {
	sort.Sort(b.Aspath)
	b.isSorted = true
	for _, v := range b.Aspath {
		b.Aspath2str += v
	}
	return
}
