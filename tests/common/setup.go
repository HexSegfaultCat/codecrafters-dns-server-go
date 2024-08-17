package common

import (
	"app/dns"
	"fmt"
	"net"
)

var currentPort uint16 = 2053

func InitializeDnsServer() (net.Conn, error) {
	var port uint16 = currentPort
	currentPort += 1

	dnsServer := dns.NewServer("127.0.0.1", port)
	go dnsServer.StartServer()

	dialer := net.Dialer{
		Timeout: 100_000_000, // 100 ms
	}

	return dialer.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
}
