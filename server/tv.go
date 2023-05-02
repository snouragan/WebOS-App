package main

import (
	"errors"
	"strconv"
	"strings"
)

func tvid(ip string) int {
	m := strings.Split(ip, ":") // remove port if present
	ip = m[0]

	switch ip {
	case "192.168.1.67":
		return 0
	case "192.168.1.132":
		return 1
	case "192.168.1.202":
		return 2
	case "192.168.1.203":
		return 3
	case "192.168.1.204":
		return 4
	case "192.168.1.135":
		return 5
	}

	return -1
}

func idtv(tv int) (ip string) {
	switch tv {
	case 0:
		return "192.168.1.67"
	case 1:
		return "192.168.1.132"
	case 2:
		return "192.168.1.202"
	case 3:
		return "192.168.1.203"
	case 4:
		return "192.168.1.204"
	case 5:
		return "192.168.1.135"
	}

	panic(errors.New("not a tv number"))
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
