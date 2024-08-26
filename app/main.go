package main

import (
	"app/dns"

	"os"
	"strconv"
	"strings"
)

func main() {
	resolver := ""
	for i, arg := range os.Args {
		if arg == "--resolver" && len(os.Args) > (i+1) {
			resolver = os.Args[i+1]
		}
	}

	dnsServer := dns.NewServer("127.0.0.1", 2053)
	if resolver != "" {
		parts := strings.Split(resolver, ":")

		ipAddress := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		dnsServer.SetResolver(ipAddress, uint16(port))
	}

	dnsServer.StartServer()
}
