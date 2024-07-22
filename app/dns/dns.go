package dns

import "encoding/binary"

type HeaderSection struct {
	PacketIdentifier uint16 // ID: 16 bit

	QueryResponseIndicator uint8 // QR: 1 bit
	OperationCode          uint8 // OPCODE: 4 bits
	AuthoritativeAnswer    uint8 // AA: 1 bit
	Truncation             uint8 // TC: 1 bit
	RecursionDesired       uint8 // RD: 1 bit

	RecursionAvailable uint8 // RA: 1 bit
	Reserved           uint8 // Z: 3 bits
	ResponseCode       uint8 // RCODE: 4 bits

	QuestionCount         uint16 // QDCOUNT: 16 bits
	AnswerRecordCount     uint16 // ANCOUNT: 16 bits
	AuthorityRecordCount  uint16 // NSCOUNT: 16 bits
	AdditionalRecordCount uint16 // ARCOUNT: 16 bits
}

func (hs *HeaderSection) ToBytes() []byte {
	result := make([]byte, 12)

	var byteAtIndex2 byte = 0
	byteAtIndex2 |= hs.QueryResponseIndicator << 7 // 0b10000000
	byteAtIndex2 |= hs.OperationCode << 3          // 0b01111000
	byteAtIndex2 |= hs.AuthoritativeAnswer << 2    // 0b00000100
	byteAtIndex2 |= hs.Truncation << 1             // 0b00000010
	byteAtIndex2 |= hs.RecursionDesired            // 0b00000001

	var byteAtIndex3 byte = 0
	byteAtIndex3 |= hs.RecursionAvailable << 7 // 0b10000000
	byteAtIndex3 |= hs.Reserved << 4           // 0b01110000
	byteAtIndex3 |= hs.ResponseCode            // 0b00001111

	binary.BigEndian.PutUint16(result[0:2], hs.PacketIdentifier)
	result[2] = byteAtIndex2
	result[3] = byteAtIndex3
	binary.BigEndian.PutUint16(result[4:6], hs.QuestionCount)
	binary.BigEndian.PutUint16(result[6:8], hs.AnswerRecordCount)
	binary.BigEndian.PutUint16(result[8:10], hs.AuthorityRecordCount)
	binary.BigEndian.PutUint16(result[10:12], hs.AdditionalRecordCount)

	return result
}
