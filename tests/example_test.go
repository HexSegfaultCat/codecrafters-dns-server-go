package dns_test

import (
	"fmt"
	"net"
	"testing"

	dns "github.com/codecrafters-io/dns-server-starter-go/app/dns"
)

func TestDnsHeader(t *testing.T) {
	var port uint16 = 2053

	dnsServer := dns.NewServer("127.0.0.1", port)
	go dnsServer.StartServer()

	dialer := net.Dialer{
		Timeout: 100_000_000, // 100 ms
	}

	connection, err := dialer.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	_, err = connection.Write([]byte{1, 2, 3})
	if err != nil {
		t.Error(err)
	}

	responseBuffer := make([]byte, 16)
	bytesReadCount, err := connection.Read(responseBuffer)
	if err != nil {
		t.Error(err)
	}

	if bytesReadCount != 12 {
		t.Errorf("Expected to receive %d bytes, but got %d", 12, bytesReadCount)
	}

	packetIdentifier := (uint16(responseBuffer[0]) << 8) | uint16(responseBuffer[1])
	responseIndicator := responseBuffer[2] >> 7
	otherFields := responseBuffer[3:]

	if packetIdentifier != 1234 {
		t.Errorf("Expected PacketIdentifier to be %d, but got %d", 1234, packetIdentifier)
	}
	if responseIndicator != 1 {
		t.Errorf("Expected ResponseIndicator to be %d, but got %d", 1, responseIndicator)
	}

	for i, v := range otherFields {
		if v != 0 {
			t.Errorf("Value at index %d should be 0 but is %d", i, v)
		}
	}
}
