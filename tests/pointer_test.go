package dnstest

import (
	"testing"

	"tests/common"

	"app/dns/packet"
	"app/dns/packet/common/domain_name"
	"app/dns/packet/section"
)

func TestBasicPointer(t *testing.T) {
	port := common.InitializeDnsServer()
	connection, err := common.CreateConnection(port)
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	// Setup
	questionBase := &dnssection.DnsQuestion{
		QueryType:  1,
		QueryClass: 1,
	}
	questionBase.DomainName.SetDomainName("example.com")

	questionBaseAbsoluteIndex := byte(12)

	question2 := &dnssection.DnsQuestion{
		DomainName: []dnsname.NamePart{
			[]byte{2, 'q', '2'},
			[]byte{0b11 << 6, questionBaseAbsoluteIndex},
		},
		QueryType:  1,
		QueryClass: 1,
	}
	question3 := &dnssection.DnsQuestion{
		DomainName: []dnsname.NamePart{
			[]byte{2, 'q', '3'},
			[]byte{0b11 << 6, questionBaseAbsoluteIndex},
		},
		QueryType:  1,
		QueryClass: 1,
	}

	packetRequest := packet.DnsPacket{
		Header: &dnssection.DnsHeader{},
	}
	packetRequest.AppendQuestionIncrementCount(questionBase)
	packetRequest.AppendQuestionIncrementCount(question2)
	packetRequest.AppendQuestionIncrementCount(question3)

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
	if headerQuestionCount != 3 || questionsLength != 3 {
		t.Error("Question count mismatch")
	}

	responseQuestion2 := packetResponse.Questions[1]
	isResponseQuestionCorrect2 := responseQuestion2.DomainName[0].String() != "q2" ||
		responseQuestion2.DomainName[1].String() != "example" ||
		responseQuestion2.DomainName[2].String() != "com"

	responseQuestion3 := packetResponse.Questions[2]
	isResponseQuestionCorrect3 := responseQuestion3.DomainName[0].String() != "q3" ||
		responseQuestion2.DomainName[1].String() != "example" ||
		responseQuestion2.DomainName[2].String() != "com"

	if isResponseQuestionCorrect2 && isResponseQuestionCorrect3 {
		t.Errorf("Incorrect question in response")
	}
}

func TestTransitivePointer(t *testing.T) {
	port := common.InitializeDnsServer()
	connection, err := common.CreateConnection(port)
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	// Setup
	questionBase := &dnssection.DnsQuestion{
		QueryType:  1,
		QueryClass: 1,
	}
	questionBase.DomainName.SetDomainName("example.com")
	questionBaseAbsoluteIndex := byte(12)

	question2 := &dnssection.DnsQuestion{
		DomainName: []dnsname.NamePart{
			[]byte{2, 'q', '2'},
			[]byte{0b11 << 6, questionBaseAbsoluteIndex},
		},
		QueryType:  1,
		QueryClass: 1,
	}
	question2AbsoluteIndex := questionBaseAbsoluteIndex + (2 + 2 + (7 + 1) + (3 + 1) + 1)

	question3 := &dnssection.DnsQuestion{
		DomainName: []dnsname.NamePart{
			[]byte{2, 'a', 'b'},
			[]byte{0b11 << 6, question2AbsoluteIndex},
		},
		QueryType:  1,
		QueryClass: 1,
	}

	packetRequest := packet.DnsPacket{
		Header: &dnssection.DnsHeader{},
	}
	packetRequest.AppendQuestionIncrementCount(questionBase)
	packetRequest.AppendQuestionIncrementCount(question2)
	packetRequest.AppendQuestionIncrementCount(question3)

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
	if headerQuestionCount != 3 || questionsLength != 3 {
		t.Error("Question count mismatch")
	}

	responseQuestion3 := packetResponse.Questions[2]
	isResponseQuestionCorrect := responseQuestion3.DomainName[0].String() != "ab" ||
		responseQuestion3.DomainName[1].String() != "q2" ||
		responseQuestion3.DomainName[2].String() != "example" ||
		responseQuestion3.DomainName[3].String() != "com"

	if isResponseQuestionCorrect {
		t.Errorf("Incorrect question in response")
	}
}
