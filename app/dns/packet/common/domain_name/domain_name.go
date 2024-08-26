package dnsname

import (
	"fmt"
	"strings"
)

type DomainName []NamePart

func (name *DomainName) SerializedParts() []string {
	result := make([]string, len(*name))
	for i, part := range *name {
		result[i] = part.Serialize()
	}

	return result
}

func (name DomainName) Bytes() []byte {
	result := []byte{}

	for _, part := range name {
		result = append(result, part...)
	}

	lastPart := name[len(name)-1]
	if lastPart.IsPointer() == false {
		result = append(result, 0)
	}

	return result
}

func (name *DomainName) SetDomainName(domainName string) {
	domainParts := strings.Split(domainName, ".")
	result := make([]NamePart, len(domainParts))

	for i, domainPart := range domainParts {
		result[i] = NewPart(domainPart)
	}

	*name = result
}

func (name DomainName) DerefDomainName(dnsPacketBytes []byte) DomainName {
	result := []NamePart{}

	if len(name) == 0 {
		return result
	}

	lastPartIndex := len(name) - 1
	lastNamePart := name[lastPartIndex]

	iterationParts := name
	for {
		result = append(result, iterationParts[:lastPartIndex]...)

		if lastNamePart.IsPointer() == false {
			break
		}

		derefedMemoryRegion := dnsPacketBytes[lastNamePart.PointerAbsoluteOffset():]
		derefedParts, _, err := ParseDomainName(derefedMemoryRegion)
		iterationParts = derefedParts

		if err != nil {
			panic(err)
		}

		lastPartIndex = len(derefedParts) - 1
		lastNamePart = derefedParts[lastPartIndex]
	}

	result = append(result, lastNamePart)

	return result
}

func ParseDomainName(domainNameBytes []byte) ([]NamePart, int, error) {
	domainParts, currentPartBaseIndex, totalLength := []NamePart{}, 0, 0

	for domainNameBytes[currentPartBaseIndex] != 0 {
		if isBytePointer(domainNameBytes[currentPartBaseIndex]) {
			pointerPart := domainNameBytes[currentPartBaseIndex:(currentPartBaseIndex + 2)]

			totalLength += 2
			domainParts = append(domainParts, pointerPart)

			break
		}

		currentPartNameLength := int(domainNameBytes[currentPartBaseIndex])
		nextPartBaseIndex := (currentPartBaseIndex + 1) + currentPartNameLength
		currentPartBytes := domainNameBytes[currentPartBaseIndex:nextPartBaseIndex]

		totalLength += currentPartNameLength + 1
		domainParts = append(domainParts, currentPartBytes)

		if nextPartBaseIndex >= len(domainNameBytes) {
			return nil, 0, fmt.Errorf("Unable to parse domain name due to overflow")
		}

		currentPartBaseIndex = nextPartBaseIndex
	}
	if domainNameBytes[currentPartBaseIndex] == 0 {
		totalLength += 1
	}

	return domainParts, totalLength, nil
}
