package qclass

type QueryClass uint16

const (
	_        = QueryClass(iota)
	Internet // IN
	Csnet    // CS (obsolete)
	Chaos    // CH
	Hesiod   // HS
)

const (
	AnyClass = QueryClass(255) // *
)
