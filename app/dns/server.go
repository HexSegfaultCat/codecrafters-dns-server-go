package dns

import (
	"fmt"
	"net"
	"net/netip"

	. "app/dns/packet"
	"app/dns/packet/query_class"
	"app/dns/packet/query_type"
	. "app/dns/packet/section"
)

type DnsServer struct {
	bindingAddress *net.UDPAddr
}

func NewServer(ipAddress string, port uint16) *DnsServer {
	parsedIp, err := netip.ParseAddr(ipAddress)
	if err != nil {
		panic(err)
	}

	return &DnsServer{
		bindingAddress: net.UDPAddrFromAddrPort(netip.AddrPortFrom(parsedIp, port)),
	}
}

func (server *DnsServer) StartServer() {
	udpConn, err := net.ListenUDP("udp", server.bindingAddress)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
		}

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s\n", size, source)

		receivedPacket, err := ParsePacketFromBytes(receivedData)
		if err != nil {
			println(err)
		}

		println(receivedPacket.DumpPacket(true))

		dnsHeader := &DnsHeader{}
		dnsHeader.SetPacketIdentifier(1234)
		dnsHeader.SetQueryResponseIndicator(true)

		dnsQuestion := &DnsQuestion{
			QueryType:  qtype.HostAddress,
			QueryClass: qclass.Internet,
		}
		dnsQuestion.SetDomainName("codecrafters.io")

		responsePacket := &DnsPacket{
			Header: dnsHeader,
		}
		responsePacket.AppendQuestionIncrementCount(dnsQuestion)

		println(responsePacket.DumpPacket(false))

		_, err = udpConn.WriteToUDP(responsePacket.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
