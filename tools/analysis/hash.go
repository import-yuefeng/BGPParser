package analysis

import (
	"crypto/sha1"
	"fmt"
	"io"
	"errors"
)

func (b *BGPInfo) aspath2sha1() string {
    tmp := sha1.New()
    io.WriteString(tmp, b.Aspath2str)
    return fmt.Sprintf("%x", tmp.Sum(nil))
}

func (b *BGPInfo) aspath2hashcode() string {
	sha1c := b.aspath2sha1()
	b.Hashcode = sha1c
	return b.Hashcode
}

func (b *BGPInfo) ConvertHashcode() error {
	if !b.isSorted {
		errors.New("data is not sorted")
	}
	b.aspath2hashcode()
	return nil
}

func CompareAspath(b1, b2 *BGPInfo) (error, bool) {
	if !b1.isSorted || !b2.isSorted {
		return errors.New("data is not sorted"), false
	}
	return nil, b1.aspath2hashcode() == b2.aspath2hashcode()
}
