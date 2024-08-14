package section

import "testing"

func TestOperationCodeWhenByteHasAllBitsSet(t *testing.T) {
	dnsHeader := &DnsHeader{}
	dnsHeader[2] = 0xff

	want := uint8(0b1111)
	actual := dnsHeader.OperationCode()

	if actual != want {
		t.Fatal()
	}
}

func TestSetOperationCodeForCorrectValue(t *testing.T) {
	dnsHeader := &DnsHeader{}
	dnsHeader.SetOperationCode(0b1101)

	want := uint8(0b01101000)
	actual := dnsHeader[2]

	if actual != want {
		t.Fatal()
	}
}
