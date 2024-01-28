package card

import (
	"fmt"
	"slices"
	"testing"
)

func Test_buildAPDU(t *testing.T) {
	longData := make([]byte, 0x100)
	for i := range longData {
		longData[i] = byte(i)
	}

	testCases := []struct {
		parameters   [4]byte
		data         []byte
		ne           uint
		expectedAPDU []byte
	}{
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{},
			ne:           0,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{},
			ne:           0xFF,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0xFF},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{},
			ne:           0x01FF,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0x01, 0xFF},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{},
			ne:           0x10000,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0x00, 0x00},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{0x00, 0x00, 0x00},
			ne:           0,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0x03, 0x00, 0x00, 0x00},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{0x00, 0x00, 0x00},
			ne:           0x01,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0x03, 0x00, 0x00, 0x00, 0x01},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         []byte{0x00, 0x00, 0x00},
			ne:           0x100,
			expectedAPDU: []byte{0x00, 0xA4, 0x04, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00},
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         longData,
			ne:           0,
			expectedAPDU: append([]byte{0x00, 0xA4, 0x04, 0x01, 0x00, 0x01, 0x00}, longData...),
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         longData,
			ne:           0x01FF,
			expectedAPDU: append(append([]byte{0x00, 0xA4, 0x04, 0x01, 0x00, 0x01, 0x00}, longData...), 0x01, 0xff),
		},
		{
			parameters:   [4]byte{0x00, 0xA4, 0x04, 0x01},
			data:         longData,
			ne:           0x01FF,
			expectedAPDU: append(append([]byte{0x00, 0xA4, 0x04, 0x01, 0x00, 0x01, 0x00}, longData...), 0x01, 0xff),
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				apdu := buildAPDU(testCase.parameters[0], testCase.parameters[1], testCase.parameters[2], testCase.parameters[3], testCase.data, testCase.ne)
				if !slices.Equal[[]byte](apdu, testCase.expectedAPDU) {
					t.Errorf("Expected APDU %v, but got %v", testCase.expectedAPDU, apdu)
				}
			},
		)
	}
}
