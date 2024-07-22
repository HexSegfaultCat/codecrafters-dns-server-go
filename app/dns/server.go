package dns

import (
	"fmt"
	"net"
	"net/netip"
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
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		dnsHeader := HeaderSection{
			PacketIdentifier:       1234,
			QueryResponseIndicator: 1,
			OperationCode:          0,
			AuthoritativeAnswer:    0,
			Truncation:             0,
			RecursionDesired:       0,
			RecursionAvailable:     0,
			Reserved:               0,
			ResponseCode:           0,
			QuestionCount:          0,
			AnswerRecordCount:      0,
			AuthorityRecordCount:   0,
			AdditionalRecordCount:  0,
		}

		response := dnsHeader.ToBytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
