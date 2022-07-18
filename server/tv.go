package main

import (
	"strconv"
	"strings"
)

func tvid(ip string) int {
	m := strings.Split(ip, ":") // remove port if present
	ip = m[0]

	switch ip {
	case "192.168.1.21":
		return 0
	case "192.168.1.20":
		return 1
	case "192.168.1.19":
		return 2
	case "192.168.1.18":
		return 3
	case "192.168.1.17":
		return 4
	case "192.168.1.16":
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
