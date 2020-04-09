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
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
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
