package utils

import "unsafe"

func Strip(s_ string, chars_ string) string {
	s, chars := []byte(s_), []byte(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[byte]bool)
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r {
		return ""
	}
	res := s[start : end+1]
	return *(*string)(unsafe.Pointer(&res))
}
