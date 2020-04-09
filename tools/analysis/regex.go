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
	"regexp"
	// log "github.com/sirupsen/logrus"
)

const (
	PREFIX_ADDRESS = "PREFIX: (?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\/([0-9]{1,3})"
	IPV4_ADDRESS   = "(?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\/([0-9]{1,3})"
	AS_PATH        = "AS-Path: \\(\\[([\\d\\s]+)\\]\\)"
)

var (
	PREFIX_ADDRESS_REGEXP = regexp.MustCompile(PREFIX_ADDRESS)
	AS_PATH_REGEXP        = regexp.MustCompile(AS_PATH)
	IPV4_ADDRESS_REGEXP   = regexp.MustCompile(IPV4_ADDRESS)
)

func (b *BGPInfo) FindAsPath() []string {
	if len(b.content) == 0 {
		return []string{}
	}
	tmp := AS_PATH_REGEXP.FindAllStringSubmatch(b.content, -1)
	for i := 0; i < len(tmp); i++ {
		for _, v := range tmp[i][1:] {
			// if len(v) > 0 {
			// 	v = v[9:]
			b.Aspath = append(b.Aspath, v)
			// }
		}
	}
	return b.Aspath
}

func (b *BGPInfo) FindPrefix() {
	tmp := PREFIX_ADDRESS_REGEXP.FindAllString(b.content, 1)
	if len(tmp) > 0 {
		b.Prefix = IPV4_ADDRESS_REGEXP.FindAllString(tmp[0], 1)
	} else {
		b.Prefix = []string{""}
	}
}
