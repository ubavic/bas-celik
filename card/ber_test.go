package card

import (
	"testing"
)

func Test_parseBerLength(t *testing.T) {
	testCases := []struct {
		data                []byte
		expectedLength      uint32
		expectedParsedBytes uint32
		expectedError       error
	}{
		{
			data:          []byte{},
			expectedError: ErrInvalidLength,
		},
		{
			data:                []byte{0x79},
			expectedLength:      0x79,
			expectedParsedBytes: 1,
			expectedError:       nil,
		},
		{
			data:          []byte{0x80, 0x91},
			expectedError: ErrInvalidFormat,
		},
		{
			data:                []byte{0x81, 0x01},
			expectedLength:      0x01,
			expectedParsedBytes: 2,
			expectedError:       nil,
		},
		{
			data:          []byte{0x81},
			expectedError: ErrInvalidLength,
		},
		{
			data:                []byte{0x82, 0x01, 0x02},
			expectedLength:      uint32(0x01)<<8 + uint32(0x02),
			expectedParsedBytes: 3,
			expectedError:       nil,
		},
		{
			data:          []byte{0x82},
			expectedError: ErrInvalidLength,
		},
		{
			data:                []byte{0x83, 0x01, 0x02, 0x03},
			expectedLength:      uint32(0x01)<<16 + uint32(0x02)<<8 + uint32(0x03),
			expectedParsedBytes: 4,
			expectedError:       nil,
		},
		{
			data:          []byte{0x83},
			expectedError: ErrInvalidLength,
		},
		{
			data:                []byte{0x84, 0x01, 0x02, 0x03, 0x04},
			expectedLength:      uint32(0x01)<<24 + uint32(0x02)<<16 + uint32(0x03)<<8 + uint32(0x04),
			expectedParsedBytes: 5,
			expectedError:       nil,
		},
		{
			data:          []byte{0x84},
			expectedError: ErrInvalidLength,
		},
	}

	for _, testCase := range testCases {
		length, parsedBytes, err := parseBerLength(testCase.data)

		if err == nil && testCase.expectedError == nil {
			if length != testCase.expectedLength {
				t.Errorf("Expected parsed length to be %d, but it is %d", testCase.expectedLength, length)
			}
			if parsedBytes != testCase.expectedParsedBytes {
				t.Errorf("Expected %d bytes to be parsed, but %d bytes was parsed", testCase.expectedParsedBytes, parsedBytes)
			}
		} else {
			if err != testCase.expectedError {
				t.Errorf("Expected error '%s', but error '%s'", testCase.expectedError.Error(), err.Error())
			}
		}
	}
}
