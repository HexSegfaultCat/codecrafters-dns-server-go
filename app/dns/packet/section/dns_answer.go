package section

import (
	. "app/dns/packet/common/dns_class"
	. "app/dns/packet/common/dns_type"
	. "app/dns/packet/common/domain_name"
	"encoding/binary"
)

type DnsAnswer struct {
	DomainName  DomainName  // QNAME
	RecordType  RecordType  // QTYPE
	RecordClass RecordClass // QCLASS
	TimeToLive  uint32      // TTL
	Length      uint16      // RDLENGTH
	Data        []byte      // RDATA
}

func (answer *DnsAnswer) Bytes() []byte {
	result := []byte(answer.DomainName)

	result = binary.BigEndian.AppendUint16(result, uint16(answer.RecordType))
	result = binary.BigEndian.AppendUint16(result, uint16(answer.RecordClass))
	result = binary.BigEndian.AppendUint32(result, answer.TimeToLive)
	result = binary.BigEndian.AppendUint16(result, answer.Length)
	result = append(result, answer.Data...)

	return result
}

func (answer *DnsAnswer) SerializeToMap() map[string]interface{} {
	return map[string]interface{}{
		"NAME":     answer.DomainName.Parts(),
		"TYPE":     answer.RecordType,
		"CLASS":    answer.RecordClass,
		"TTL":      answer.TimeToLive,
		"RDLENGTH": answer.Length,
		"RDATA":    answer.Data,
	}
}
