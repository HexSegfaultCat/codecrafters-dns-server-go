package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"app/dns/packet/query_class"
	"app/dns/packet/query_type"
	"app/dns/packet/section"
)

type DnsSerializer interface {
	SerializeToMap() map[string]interface{}
}

type DnsPacket struct {
	Header    *section.DnsHeader
	Questions []*section.DnsQuestion
}

func (packet *DnsPacket) Bytes() []byte {
	result := packet.Header[:]

	for _, question := range packet.Questions {
		result = append(result, question.Bytes()...)
	}

	return result
}

func (packet *DnsPacket) AppendQuestionIncrementCount(question *section.DnsQuestion) {
	packet.Questions = append(packet.Questions, question)

	newCount := packet.Header.QuestionCount() + 1
	packet.Header.SetQuestionCount(newCount)
}

func ParsePacketFromBytes(packetBytes []byte) (*DnsPacket, error) {
	dnsHeader := section.DnsHeader(packetBytes[:12])

	result := &DnsPacket{
		Header:    &dnsHeader,
		Questions: make([]*section.DnsQuestion, dnsHeader.QuestionCount()),
	}

	leftSkipBytes := packetBytes[12:]
	for i := range result.Questions {
		index := bytes.IndexByte(leftSkipBytes, 0)
		if index == -1 {
			return result, fmt.Errorf("Unable to parse question number %d\n", i+1)
		}

		domainName := leftSkipBytes[:(index + 1)]
		queryType := binary.BigEndian.Uint16(leftSkipBytes[(index + 1):(index + 3)])
		queryClass := binary.BigEndian.Uint16(leftSkipBytes[(index + 3):(index + 5)])

		result.Questions[i] = &section.DnsQuestion{
			DomainName: domainName,
			QueryType:  qtype.QueryType(queryType),
			QueryClass: qclass.QueryClass(queryClass),
		}

		leftSkipBytes = leftSkipBytes[(index + 5):]
	}

	return result, nil
}

func (packet *DnsPacket) SerializeToMap() map[string]interface{} {
	result := make(map[string]interface{})

	structKeys := reflect.TypeOf(*packet)
	structValues := reflect.ValueOf(*packet)

	for i := 0; i < structValues.NumField(); i++ {
		fieldKey := structKeys.Field(i)
		fieldValue := structValues.Field(i)

		switch fieldValue.Kind() {
		case reflect.Slice:
			serializedSlice := make([]map[string]interface{}, len(packet.Questions))
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

func (packet *DnsPacket) DumpPacket(isRequest bool) string {
	requestResponse := "RES"
	if isRequest {
		requestResponse = "REQ"
	}

	prefix := fmt.Sprintf("[%d][%s] ", packet.Header.PacketIdentifier(), requestResponse)

	marshalled, _ := json.MarshalIndent(packet.SerializeToMap(), prefix, "  ")
	return prefix + string(marshalled) + "\n"
}
