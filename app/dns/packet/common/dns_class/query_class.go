package dnsclass

type QueryClass uint16

const (
	AnyClass = QueryClass(255) // *
)

func (recordValue RecordClass) QueryValue() QueryClass {
	return QueryClass(recordValue)
}
