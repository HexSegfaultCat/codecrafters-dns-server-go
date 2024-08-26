package dnsname

import "fmt"

type NamePart []byte

func isBytePointer(firstByte byte) bool {
	return (firstByte >> 6) == 0b11
}

func (part NamePart) IsPointer() bool {
	return isBytePointer(part[0])
}

func (part NamePart) PointerAbsoluteOffset() uint16 {
	leftByte := uint16(part[0] & 0b00111111)
	rightByte := uint16(part[1])

	return (leftByte << 8) | rightByte
}

func (part NamePart) Length() uint8 {
	return part[0]
}

func (part NamePart) String() string {
	return string(part[1:])
}

func NewPart(name string) NamePart {
	result := make([]byte, len(name)+1)

	result[0] = uint8(len(name))
	for i := range name {
		result[i+1] = name[i]
	}

	return result
}

func (part NamePart) Serialize() string {
	if part.IsPointer() {
		return fmt.Sprintf("P|%d", part.PointerAbsoluteOffset())
	}
	return fmt.Sprintf("%d|%s", part.Length(), part.String())
}
