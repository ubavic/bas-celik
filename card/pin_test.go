package card_test

import (
	"slices"
	"testing"

	"github.com/ubavic/bas-celik/card"
)

func Test_ValidatePin(t *testing.T) {
	testCases := []struct {
		pin            string
		expectedResult bool
	}{
		{
			pin:            "",
			expectedResult: false,
		},
		{
			pin:            "1",
			expectedResult: false,
		},
		{
			pin:            "123",
			expectedResult: false,
		},
		{
			pin:            "1234",
			expectedResult: true,
		},
		{
			pin:            "dddd",
			expectedResult: false,
		},
		{
			pin:            "123a",
			expectedResult: false,
		},
		{
			pin:            "123456789",
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		pinValid := card.ValidatePin(testCase.pin)
		if pinValid != testCase.expectedResult {
			t.Errorf("Expected %t but got %t.", testCase.expectedResult, pinValid)
		}
	}
}

func Test_PadPin(t *testing.T) {
	testCases := []struct {
		pin            string
		expectedResult []byte
	}{
		{
			pin:            "1234",
			expectedResult: []byte{0x31, 0x32, 0x33, 0x34, 0, 0, 0, 0},
		},
		{
			pin:            "",
			expectedResult: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			pin:            "12345678",
			expectedResult: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38},
		},
		{
			pin:            "123456789",
			expectedResult: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38},
		},
	}

	for _, testCase := range testCases {
		paddedPin := card.PadPin(testCase.pin)
		if !slices.Equal[[]byte](paddedPin, testCase.expectedResult) {
			t.Errorf("Expected %v but got %v.", testCase.expectedResult, paddedPin)
		}
	}
}
