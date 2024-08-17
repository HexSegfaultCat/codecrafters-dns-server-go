package dname

import "strings"

type DomainName []byte

func (name *DomainName) Parts() []string {
	domainNameBytes := []byte(*name)

	finalDomainParts, currentPartBaseIndex := []string{}, 0
	for currentPartBaseIndex < len(domainNameBytes) {
		currentPartTextLength := int(domainNameBytes[currentPartBaseIndex])
		currentPartTextBaseIndex := currentPartBaseIndex + 1

		nextPartBaseIndex := currentPartTextBaseIndex + currentPartTextLength

		currentPartText := domainNameBytes[currentPartTextBaseIndex:nextPartBaseIndex]
		finalDomainParts = append(finalDomainParts, string(currentPartText))

		currentPartBaseIndex = nextPartBaseIndex
	}

	return finalDomainParts
}

func (name *DomainName) SetDomainName(domainName string) {
	*name = []byte{}

	domainParts := strings.Split(domainName, ".")
	for _, domainPart := range domainParts {
		partBytes := []byte(domainPart)
		partLength := byte(len(domainPart))

		fullPart := append([]byte{partLength}, partBytes...)

		*name = append(*name, fullPart...)
	}

	*name = append(*name, 0)
}

func NewByDomainName(domainName string) DomainName {
	name := DomainName{}
	name.SetDomainName(domainName)

	return name
}
