package qtype

type QueryType uint16

const (
	_                     = QueryType(iota)
	HostAddress           // A
	AuthorativeNameServer // NS
	MailDestination       // MD (obsolete - use `MailExchange`)
	MailForwarder         // MF (obsolete - use `MailExchange`)
	CanonicalName         // CNAME
	StartAuthorityZone    // SOA
	MailDomainName        // MB (experimental)
	MailGroupMember       // MG (experimental)
	MailRenameDomainName  // MR (experimental)
	Null                  // NULL (experimental)
	WellKnownService      // WKS
	DomainNamePointer     // PTR
	HostInformation       // HINFO
	MailListInformation   // MINFO
	MailExchange          // MX
	TextStrings           // TXT
)

const (
	_                  = QueryType(iota + 251)
	TransferEntireZone // AXFR
	MailboxRecords     // MAILB
	MailAgent          // MAILA
	AnyType            // *
)
