package card

import (
	"fmt"
	"testing"
)

func Test_responseOK(t *testing.T) {
	testCases := []struct {
		value  []byte
		result bool
	}{
		{[]byte{0x0F, 0x0F}, false},
		{[]byte{0x90, 0x00}, true},
		{[]byte{0x01, 0xFF, 0x90, 0x00}, true},
		{[]byte{0x01, 0xFF, 0x00, 0x00}, false},
		{[]byte{0xA1}, false},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				res := responseOK(testCase.value)

				if res != testCase.result {
					t.Errorf("Expected %t, but got %t", testCase.result, res)
				}
			},
		)
	}
}
