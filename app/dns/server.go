package dns

import (
	"fmt"
	"net"
	"net/netip"

	. "app/dns/packet"
	"app/dns/packet/common/dns_class"
	"app/dns/packet/common/dns_type"
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

		println(receivedPacket.DumpPacket("REQ"))

		dnsHeader := &DnsHeader{}
		dnsHeader.SetPacketIdentifier(receivedPacket.Header.PacketIdentifier())
		dnsHeader.SetQueryResponseIndicator(true)
		dnsHeader.SetOperationCode(receivedPacket.Header.OperationCode())
		dnsHeader.SetRecursionDesired(receivedPacket.Header.RecursionDesired())

		if receivedPacket.Header.OperationCode() == 0 {
			dnsHeader.SetResponseCode(0)
		} else {
			dnsHeader.SetResponseCode(4)
		}

		responsePacket := &DnsPacket{
			Header: dnsHeader,
		}

		for _, receivedQuestion := range receivedPacket.Questions {
			dnsQuestion := &DnsQuestion{
				DomainName: receivedQuestion.DomainName.DerefDomainName(receivedData),
				QueryType:  receivedQuestion.QueryType,
				QueryClass: receivedQuestion.QueryClass,
			}
			responsePacket.AppendQuestionIncrementCount(dnsQuestion)

			dnsAnswer := &DnsAnswer{
				DomainName:  dnsQuestion.DomainName,
				RecordType:  dnstype.HostAddress,
				RecordClass: dnsclass.Internet,
				TimeToLive:  60,
				Length:      4,
				Data:        []byte{8, 8, 8, 8},
			}
			responsePacket.AppendAnswerIncrementCount(dnsAnswer)
		}

		println(responsePacket.DumpPacket("RES"))

		_, err = udpConn.WriteToUDP(responsePacket.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
