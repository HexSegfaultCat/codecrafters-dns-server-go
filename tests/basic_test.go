package dnstest

import (
	"testing"

	"tests/common"

	"app/dns/packet"
	"app/dns/packet/section"
)

func TestDnsHeader(t *testing.T) {
	connection, err := common.InitializeDnsServer()
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	// Setup
	packetRequest := packet.DnsPacket{
		Header: &section.DnsHeader{},
	}
	_, err = connection.Write(packetRequest.Bytes())
	if err != nil {
		t.Error(err)
	}

	// Verify
	responseBuffer := make([]byte, 16)
	bytesReadCount, err := connection.Read(responseBuffer)
	if err != nil {
		t.Error(err)
	}

	if bytesReadCount <= 12 {
		t.Errorf("Expected to receive more than %d bytes, but got %d", 12, bytesReadCount)
	}

	packetIdentifier := (uint16(responseBuffer[0]) << 8) | uint16(responseBuffer[1])
	responseIndicator := responseBuffer[2] >> 7
	questionCount := responseBuffer[5] // bytes [4][5] - big-endian
	answerCount := responseBuffer[7]   // bytes [6][7] - big-endian

	if packetIdentifier != 1234 {
		t.Errorf("Expected PacketIdentifier to be %d, but got %d", 1234, packetIdentifier)
	}
	if responseIndicator != 1 {
		t.Errorf("Expected ResponseIndicator to be %d, but got %d", 1, responseIndicator)
	}

	if questionCount != 1 {
		t.Errorf("Expected QuestionCount to be %d, but got %d", 1, questionCount)
	}
	if answerCount != 1 {
		t.Errorf("Expected AnswerCount to be %d, but got %d", 1, answerCount)
	}
}

func TestDnsQuestion(t *testing.T) {
	connection, err := common.InitializeDnsServer()
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	// Setup
	question := &section.DnsQuestion{
		QueryType:  1,
		QueryClass: 1,
	}
	question.DomainName.SetDomainName("www.example.com")

	packetRequest := packet.DnsPacket{
		Header: &section.DnsHeader{},
	}
	packetRequest.AppendQuestionIncrementCount(question)

	_, err = connection.Write(packetRequest.Bytes())
	if err != nil {
		t.Error(err)
	}

	// Verify
	responseBuffer := make([]byte, 512)
	bytesReadCount, err := connection.Read(responseBuffer)
	if err != nil {
		t.Error(err)
	}

	if bytesReadCount <= 12 {
		t.Errorf("Expected to receive more than %d bytes, but got %d", 12, bytesReadCount)
	}

	packetResponse, err := packet.ParsePacketFromBytes(responseBuffer[:bytesReadCount])
	if err != nil {
		t.Error(err)
	}

	headerQuestionCount := packetResponse.Header.QuestionCount()
	questionsLength := len(packetResponse.Questions)
	if headerQuestionCount != 1 || questionsLength != 1 {
		t.Error("Question count mismatch")
	}
}
