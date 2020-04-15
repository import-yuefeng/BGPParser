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
	"fmt"
	"io"
)

func PackagingHashcode(hashcode string) string {
	tmp := sha1.New()
	res := []byte(hashcode)
	SortHashcode(res)
	hashcode = string(res)
	io.WriteString(tmp, hashcode)
	return fmt.Sprintf("%x", tmp.Sum(nil))
}

func (b *BGPInfo) ConvertHashcode() {
	tmp := sha1.New()
	io.WriteString(tmp, b.AsPath)
	b.Hashcode = fmt.Sprintf("%x", tmp.Sum(nil))
	return
}
