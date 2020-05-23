package utils

import (
	"net"
)

func IsIP(ipaddr string) bool {
	if _, _, err := net.ParseCIDR(ipaddr); err != nil {
		if ip := net.ParseIP(ipaddr); ip == nil {
			return false
		}
	}
	return true
}

func IsIPV4(ipaddr string) bool {
	ip := net.ParseIP(ipaddr)
	if ok := ip.To4(); ok != nil {
		return true
	}
	return false
}

func IsIPV6(ipaddr string) bool {
	return !IsIPV4(ipaddr)
}
