package main

import (
	"strconv"
	"strings"
)

func tvid(ip string) int {
	m := strings.Split(ip, ":") // remove port if present
	ip = m[0]

	switch ip {
	case "192.168.81.63":
		return 0
	case "192.168.88.30":
		return 1
	case "192.168.123.118":
		return 2
	case "192.168.104.210":
		return 3
	case "192.168.93.19":
		return 4
	case "192.168.85.179":
		return 5
	}

	return -1
}

func tvidstring(ip string) string {
	return strconv.Itoa(tvid(ip) + 1)
}

func fmtip(ip string) string {
	tv := tvidstring(ip)

	if tv == "0" {
		m := strings.Split(ip, ":") // remove port if present
		ip = m[0]
		return ip
	}

	return "TV" + tv
}
