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
	bindingAddress  *net.UDPAddr
	resolverAddress *net.UDPAddr
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

func (server *DnsServer) SetResolver(ipAddress string, port uint16) {
	parsedIp, err := netip.ParseAddr(ipAddress)
	if err != nil {
		panic(err)
	}

	server.resolverAddress = net.UDPAddrFromAddrPort(netip.AddrPortFrom(parsedIp, port))
}

func (server *DnsServer) forwardSlicedPacketAndGetResponse(
	requestPacket *DnsPacket,
) (*DnsPacket, error) {
	conn, err := net.DialUDP("udp4", nil, server.resolverAddress)
	if err != nil {
		return nil, err
	}

	requestRawBytes := requestPacket.Bytes()

	finalPacket := &DnsPacket{}
	resolverBuf := make([]byte, 512)

	for i, question := range requestPacket.Questions {
		header := &DnsHeader{}

		copy(header[:], requestPacket.Header[:])

		derefedQuestion := &DnsQuestion{
			DomainName: question.DomainName.DerefDomainName(requestRawBytes),
			QueryType:  question.QueryType,
			QueryClass: question.QueryClass,
		}

		singleQuestionPacket := DnsPacket{
			Header:    header,
			Questions: []*DnsQuestion{derefedQuestion},
		}
		header.SetQuestionCount(1)

		println(
			singleQuestionPacket.DumpPacket(fmt.Sprintf("REQ-FWD-%d", i)),
		)
		conn.Write(singleQuestionPacket.Bytes())

		responseSize, err := conn.Read(resolverBuf)
		if err != nil {
			return nil, err
		}

		receivedSingleAnswer, err := ParsePacketFromBytes(resolverBuf[:responseSize])
		if err != nil {
			return nil, err
		}
		println(
			receivedSingleAnswer.DumpPacket(fmt.Sprintf("RES-FWD-%d", i)),
		)

		finalPacket.Header = receivedSingleAnswer.Header
		finalPacket.Questions = append(finalPacket.Questions, receivedSingleAnswer.Questions...)
		finalPacket.Answers = append(finalPacket.Answers, receivedSingleAnswer.Answers...)
	}

	finalPacket.Header.SetQuestionCount(uint16(len(finalPacket.Questions)))
	finalPacket.Header.SetAnswerRecordCount(uint16(len(finalPacket.Answers)))

	return finalPacket, nil
}

func (server *DnsServer) StartServer() {
	udpConn, err := net.ListenUDP("udp", server.bindingAddress)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	fmt.Printf(
		"Listening on address %s/UDP\n",
		server.bindingAddress.String(),
	)
	fmt.Printf("Using forwarding server %s\n", server.resolverAddress.String())

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

		responsePacket := &DnsPacket{}
		if server.resolverAddress != nil {
			fmt.Printf("Forwarding request to %s\n", server.resolverAddress.String())

			packet, err := server.forwardSlicedPacketAndGetResponse(receivedPacket)
			if err != nil {
				println(err)
			} else {
				responsePacket = packet
			}
		} else {
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

			packet := &DnsPacket{
				Header: dnsHeader,
			}

			for _, receivedQuestion := range receivedPacket.Questions {
				dnsQuestion := &DnsQuestion{
					DomainName: receivedQuestion.DomainName.DerefDomainName(receivedData),
					QueryType:  receivedQuestion.QueryType,
					QueryClass: receivedQuestion.QueryClass,
				}
				packet.AppendQuestionIncrementCount(dnsQuestion)

				dnsAnswer := &DnsAnswer{
					DomainName:  dnsQuestion.DomainName,
					RecordType:  dnstype.HostAddress,
					RecordClass: dnsclass.Internet,
					TimeToLive:  60,
					Length:      4,
					Data:        []byte{8, 8, 8, 8},
				}
				packet.AppendAnswerIncrementCount(dnsAnswer)
			}

			responsePacket = packet
		}

		println(responsePacket.DumpPacket("RES"))

		_, err = udpConn.WriteToUDP(responsePacket.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
