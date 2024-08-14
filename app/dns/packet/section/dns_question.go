package section

import (
	. "app/dns/packet/query_class"
	. "app/dns/packet/query_type"
	"encoding/binary"
	"strings"
)

type DnsQuestion struct {
	DomainName []byte     // QNAME
	QueryType  QueryType  // QTYPE
	QueryClass QueryClass // QCLASS
}

func (question *DnsQuestion) SetDomainName(domainName string) {
	var result []byte

	domainParts := strings.Split(domainName, ".")
	for _, domainPart := range domainParts {
		partBytes := []byte(domainPart)
		partLength := byte(len(domainPart))
		fullPart := append([]byte{partLength}, partBytes...)

		result = append(result, fullPart...)
	}

	result = append(result, 0)

	question.DomainName = result
}

func (question *DnsQuestion) Bytes() []byte {
	result := question.DomainName

	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryType))
	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryClass))

	return result
}

func (question *DnsQuestion) SerializeToMap() map[string]interface{} {
	domainPart, partIndex := []string{}, 0
	for partIndex < len(question.DomainName) {
		partLength := int(question.DomainName[partIndex])
		nextPartIndex := partIndex + 1 + partLength

		partText := question.DomainName[partIndex+1 : nextPartIndex]
		domainPart = append(domainPart, string(partText))

		partIndex = nextPartIndex
	}

	return map[string]interface{}{
		"QNAME":  domainPart,
		"QTYPE":  question.QueryType,
		"QCLASS": question.QueryClass,
	}
}
