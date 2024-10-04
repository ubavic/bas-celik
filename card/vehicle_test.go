package card

import (
	"errors"
	"testing"

	"github.com/ubavic/bas-celik/card/cardErrors"
)

func Test_parseVehicleCardFileSize(t *testing.T) {
	testCases := []struct {
		data           []byte
		expectedLength uint
		expectedOffset uint
		expectedError  error
	}{
		{
			data:          []byte{},
			expectedError: cardErrors.ErrInvalidLength,
		},
		{
			data:          []byte{0x01, 0x02, 0x03, 0x04},
			expectedError: cardErrors.ErrInvalidLength,
		},
		{
			data: []byte{
				0x78, 0x0E, 0x4F, 0x0C, 0xA0, 0x00, 0x00, 0x00,
				0x18, 0x65, 0x56, 0x4C, 0x2D, 0x30, 0x30, 0x31,
				0x72,
			},
			expectedError: cardErrors.ErrInvalidLength,
		},
		{
			data: []byte{
				0x78, 0x0E, 0x4F, 0x0C, 0xA0, 0x00, 0x00, 0x00,
				0x18, 0x65, 0x56, 0x4C, 0x2D, 0x30, 0x30, 0x31,
				0x72, 0x27,
			},
			expectedLength: 41,
			expectedOffset: 16,
			expectedError:  nil,
		},
		{
			data:          []byte{0x01, 0x01, 0x01, 0x00, 0x80},
			expectedError: cardErrors.ErrInvalidFormat,
		},
	}

	for _, testCase := range testCases {
		length, offset, err := parseVehicleCardFileSize(testCase.data)
		if err == nil && testCase.expectedError == nil {
			if length != testCase.expectedLength {
				t.Errorf("Expected length to be %d, but it is %d", testCase.expectedLength, length)
			}
			if offset != testCase.expectedOffset {
				t.Errorf("Expected offset to be %d, but it is %d", testCase.expectedOffset, offset)
			}
		} else {
			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("Expected error to be '%v', but it is '%v'", testCase.expectedError, err)
			}
		}
	}

}
