package section

import "encoding/binary"

type DnsHeader [12]byte

func setOrClearBitAtIndex(byte *byte, bitOffset uint8, value bool) {
	if value {
		*byte |= (1 << bitOffset)
	} else {
		*byte &= 0xff - (1 << bitOffset)
	}
}

func setOrClearBitsWithMask(byte *byte, mask byte, bitOffset uint8, value uint8) bool {
	if value > mask {
		return false
	}

	var sectionInverseMask uint8 = 0xff - (mask << bitOffset)
	maskOverlay := value << bitOffset

	*byte = (*byte & sectionInverseMask) | maskOverlay
	return true
}

// ID
func (header *DnsHeader) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(header[0:2])
}
func (header *DnsHeader) SetPacketIdentifier(value uint16) {
	binary.BigEndian.PutUint16(header[0:2], value)
}

// QR: 0b10000000
func (header *DnsHeader) QueryResponseIndicator() bool {
	return (header[2] >> 7) == 1
}
func (header *DnsHeader) SetQueryResponseIndicator(value bool) {
	setOrClearBitAtIndex(&header[2], 7, value)
}

// OPCODE: 0b01111000
func (header *DnsHeader) OperationCode() uint8 {
	return (header[2] >> 3) & 0b1111
}
func (header *DnsHeader) SetOperationCode(value uint8) {
	ok := setOrClearBitsWithMask(&header[2], 0b1111, 3, value)
	if !ok {
		panic("Value may only set 4 lower bits")
	}
}

// AA: 0b00000100
func (header *DnsHeader) AuthoritativeAnswer() bool {
	return (header[2] >> 2) == 1
}
func (header *DnsHeader) SetAuthoritativeAnswer(value bool) {
	setOrClearBitAtIndex(&header[2], 2, value)
}

// TC: 0b00000010
func (header *DnsHeader) Truncation() bool {
	return (header[2] >> 1) == 1
}
func (header *DnsHeader) SetTruncation(value bool) {
	setOrClearBitAtIndex(&header[2], 1, value)
}

// RD: 0b00000001
func (header *DnsHeader) RecursionDesired() bool {
	return (header[2] & 0b1) == 1
}
func (header *DnsHeader) SetRecursionDesired(value bool) {
	setOrClearBitAtIndex(&header[2], 0, value)
}

// RA: 0b10000000
func (header *DnsHeader) RecursionAvailable() bool {
	return (header[3] >> 7) == 1
}
func (header *DnsHeader) SetRecursionAvailable(value bool) {
	setOrClearBitAtIndex(&header[3], 7, value)
}

// Z: 0b01110000
func (header *DnsHeader) Reserved() uint8 {
	return (header[3] >> 4) & 0b111
}
func (header *DnsHeader) SetReserved(value uint8) {
	ok := setOrClearBitsWithMask(&header[3], 0b111, 4, value)
	if !ok {
		panic("Value may only set 3 lower bits")
	}
}

// RCODE: 0b00001111
func (header *DnsHeader) ResponseCode() uint8 {
	return header[3] & 0b1111
}
func (header *DnsHeader) SetResponseCode(value uint8) {
	ok := setOrClearBitsWithMask(&header[3], 0b1111, 0, value)
	if !ok {
		panic("Value may only set 4 lower bits")
	}
}

// QDCOUNT
func (header *DnsHeader) QuestionCount() uint16 {
	return binary.BigEndian.Uint16(header[4:6])
}
func (header *DnsHeader) SetQuestionCount(value uint16) {
	binary.BigEndian.PutUint16(header[4:6], value)
}

// ANCOUNT
func (header *DnsHeader) AnswerRecordCount() uint16 {
	return binary.BigEndian.Uint16(header[6:8])
}
func (header *DnsHeader) SetAnswerRecordCount(value uint16) {
	binary.BigEndian.PutUint16(header[6:8], value)
}

// NSCOUNT
func (header *DnsHeader) AuthorityRecordCount() uint16 {
	return binary.BigEndian.Uint16(header[8:10])
}
func (header *DnsHeader) SetAuthorityRecordCount(value uint16) {
	binary.BigEndian.PutUint16(header[8:10], value)
}

// ARCOUNT
func (header *DnsHeader) AdditionalRecordCount() uint16 {
	return binary.BigEndian.Uint16(header[10:12])
}
func (header *DnsHeader) SetAdditionalRecordCount(value uint16) {
	binary.BigEndian.PutUint16(header[10:12], value)
}

func (header *DnsHeader) SerializeToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":      header.PacketIdentifier(),
		"QR":      header.QueryResponseIndicator(),
		"OPCODE":  header.OperationCode(),
		"AA":      header.AuthoritativeAnswer(),
		"TC":      header.Truncation(),
		"RD":      header.RecursionDesired(),
		"RA":      header.RecursionAvailable(),
		"Z":       header.Reserved(),
		"RCODE":   header.ResponseCode(),
		"QDCOUNT": header.QuestionCount(),
		"ANCOUNT": header.AnswerRecordCount(),
		"NSCOUNT": header.AuthorityRecordCount(),
		"ARCOUNT": header.AdditionalRecordCount(),
	}
}
