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
