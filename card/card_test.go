package card

import (
	"testing"

	"github.com/ebfe/scard"
)

func TestReadCardInit(t *testing.T) {
	var card scard.Card
	_, err := ReadCard(&card)

	if err == nil {
		t.Errorf("Expected error here!")
	}
}

func Test_responseOK(t *testing.T) {
	byteSliceTest := []struct {
		value  []byte
		result bool
	}{{[]byte{0x0F, 0x0F}, false}, {[]byte{0x90, 0x00}, true}}

	for _, testRes := range byteSliceTest {
		res := responseOK(testRes.value)

		if res != testRes.result {
			t.Errorf("Expected res to be %t and it is %t", res, testRes.result)
		}
	}
}
