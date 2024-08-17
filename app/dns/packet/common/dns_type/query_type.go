package dnstype

type QueryType RecordType

const (
	_                  = QueryType(iota + 251)
	TransferEntireZone // AXFR
	MailboxRecords     // MAILB
	MailAgent          // MAILA
	AnyType            // *
)

func (recordValue RecordType) QueryValue() QueryType {
	return QueryType(recordValue)
}
