package dnssection

import (
	"encoding/binary"

	. "app/dns/packet/common/dns_class"
	. "app/dns/packet/common/dns_type"
	. "app/dns/packet/common/domain_name"
)

type DnsQuestion struct {
	DomainName DomainName // QNAME
	QueryType  QueryType  // QTYPE
	QueryClass QueryClass // QCLASS
}

func (question *DnsQuestion) Bytes() []byte {
	result := question.DomainName.Bytes()

	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryType))
	result = binary.BigEndian.AppendUint16(result, uint16(question.QueryClass))

	return result
}

func (question *DnsQuestion) SerializeToMap() map[string]interface{} {
	return map[string]interface{}{
		"QNAME":  question.DomainName.SerializedParts(),
		"QTYPE":  question.QueryType,
		"QCLASS": question.QueryClass,
	}
}
