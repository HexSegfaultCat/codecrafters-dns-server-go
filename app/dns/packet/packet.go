package packet

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"app/dns/packet/common/dns_class"
	"app/dns/packet/common/dns_type"
	"app/dns/packet/common/domain_name"
	"app/dns/packet/section"
)

type DnsSerializer interface {
	SerializeToMap() map[string]interface{}
}

type DnsPacket struct {
	Header    *dnssection.DnsHeader
	Questions []*dnssection.DnsQuestion
	Answers   []*dnssection.DnsAnswer
}

func (packet *DnsPacket) Bytes() []byte {
	result := packet.Header[:]

	for _, question := range packet.Questions {
		result = append(result, question.Bytes()...)
	}

	for _, answer := range packet.Answers {
		result = append(result, answer.Bytes()...)
	}

	return result
}

func (packet *DnsPacket) AppendQuestionIncrementCount(question *dnssection.DnsQuestion) {
	packet.Questions = append(packet.Questions, question)

	newCount := packet.Header.QuestionCount() + 1
	packet.Header.SetQuestionCount(newCount)
}

func (packet *DnsPacket) AppendAnswerIncrementCount(answer *dnssection.DnsAnswer) {
	packet.Answers = append(packet.Answers, answer)

	newCount := packet.Header.AnswerRecordCount() + 1
	packet.Header.SetAnswerRecordCount(newCount)
}

func ParsePacketFromBytes(packetBytes []byte) (*DnsPacket, error) {
	dnsHeader := dnssection.DnsHeader(packetBytes[:12])

	result := &DnsPacket{
		Header:    &dnsHeader,
		Questions: make([]*dnssection.DnsQuestion, dnsHeader.QuestionCount()),
		Answers:   make([]*dnssection.DnsAnswer, dnsHeader.AnswerRecordCount()),
	}

	leftSkipBytes := packetBytes[12:]
	for i := range result.Questions {
		domainName, leftOffset, err := dnsname.ParseDomainName(leftSkipBytes)
		if err != nil {
			panic(err)
		}

		queryType := binary.BigEndian.Uint16(
			leftSkipBytes[leftOffset:(leftOffset + 2)],
		)
		queryClass := binary.BigEndian.Uint16(
			leftSkipBytes[(leftOffset + 2):(leftOffset + 4)],
		)

		result.Questions[i] = &dnssection.DnsQuestion{
			DomainName: domainName,
			QueryType:  dnstype.QueryType(queryType),
			QueryClass: dnsclass.QueryClass(queryClass),
		}

		leftSkipBytes = leftSkipBytes[(leftOffset + 4):]
	}
	for i := range result.Answers {
		domainName, leftOffset, err := dnsname.ParseDomainName(leftSkipBytes)
		if err != nil {
			panic(err)
		}

		recordType := binary.BigEndian.Uint16(
			leftSkipBytes[leftOffset:(leftOffset + 2)],
		)
		recordClass := binary.BigEndian.Uint16(
			leftSkipBytes[(leftOffset + 2):(leftOffset + 4)],
		)
		timeToLive := binary.BigEndian.Uint32(
			leftSkipBytes[(leftOffset + 4):(leftOffset + 8)],
		)
		length := binary.BigEndian.Uint16(
			leftSkipBytes[(leftOffset + 8):(leftOffset + 10)],
		)
		data := leftSkipBytes[(leftOffset + 10):(leftOffset + 10 + int(length))]

		result.Answers[i] = &dnssection.DnsAnswer{
			DomainName:  domainName,
			RecordType:  dnstype.RecordType(recordType),
			RecordClass: dnsclass.RecordClass(recordClass),
			TimeToLive:  timeToLive,
			Length:      length,
			Data:        data,
		}

		leftSkipBytes = leftSkipBytes[(leftOffset + 10 + int(length)):]
	}

	return result, nil
}

func (packet *DnsPacket) SerializeToMap() map[string]interface{} {
	result := make(map[string]interface{})

	structKeys := reflect.TypeOf(*packet)
	structValues := reflect.ValueOf(*packet)

	for fieldIndex := 0; fieldIndex < structValues.NumField(); fieldIndex++ {
		fieldKey := structKeys.Field(fieldIndex)
		fieldValue := structValues.Field(fieldIndex)

		switch fieldValue.Kind() {
		case reflect.Slice:
			serializedSlice := make([]map[string]interface{}, fieldValue.Len())
			for i := range serializedSlice {
				serializedSlice[i] = fieldValue.Index(i).Interface().(DnsSerializer).SerializeToMap()
			}
			result[fieldKey.Name] = serializedSlice
		case reflect.Ptr:
			result[fieldKey.Name] = fieldValue.Interface().(DnsSerializer).SerializeToMap()
		default:
			panic("Unknown type")
		}
	}

	return result
}

func (packet *DnsPacket) DumpPacket(label string) string {
	prefix := fmt.Sprintf("[%d][%s] ", packet.Header.PacketIdentifier(), label)

	marshalled, _ := json.MarshalIndent(packet.SerializeToMap(), prefix, "  ")
	return prefix + string(marshalled) + "\n"
}
