package section

import (
	. "app/dns/packet/common/dns_class"
	. "app/dns/packet/common/dns_type"
	. "app/dns/packet/common/domain_name"
	"encoding/binary"
)

type DnsQuestion struct {
	DomainName DomainName // QNAME
	QueryType  QueryType  // QTYPE
	QueryClass QueryClass // QCLASS
}

func (question *DnsQuestion) Bytes() []byte {
	result := []byte(question.DomainName)

	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryType))
	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryClass))

	return result
}

func (question *DnsQuestion) SerializeToMap() map[string]interface{} {
	return map[string]interface{}{
		"QNAME":  question.DomainName.Parts(),
		"QTYPE":  question.QueryType,
		"QCLASS": question.QueryClass,
	}
}
