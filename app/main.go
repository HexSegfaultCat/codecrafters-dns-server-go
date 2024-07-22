package main

import (
	dns "github.com/codecrafters-io/dns-server-starter-go/app/dns"
)

func main() {
	dnsServer := dns.NewServer("127.0.0.1", 2053)

	dnsServer.StartServer()
}
