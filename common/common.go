package common

import (
	"net"
	"strings"
)

type CollectEnty struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func GetLocalIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
