package dnsclass

type RecordClass uint16

const (
	_        = RecordClass(iota)
	Internet // IN
	Csnet    // CS (obsolete)
	Chaos    // CH
	Hesiod   // HS
)
