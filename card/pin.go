package card

import (
	"slices"
	"unicode"
)

// Checks if the PIN consists only of digits,
// and it's length is between 4 and 8.
func ValidatePin(pin string) bool {
	pinLength := len(pin)
	if pinLength < 4 || pinLength > 8 {
		return false
	}

	for _, r := range pin {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// Creates a 8 byte slice containing the PIN at the beginning.
func PadPin(pin string) []byte {
	data := make([]byte, 8)

	for i, r := range pin {
		if i < 8 {
			data[i] = byte(r)
		}
	}

	return data
}

func PinTriesLeft(rsp []byte) int {
	if slices.Equal(rsp, []byte{0x63, 0xC0}) {
		return 0
	}

	if slices.Equal(rsp, []byte{0x63, 0xC1}) {
		return 1
	}

	if slices.Equal(rsp, []byte{0x63, 0xC2}) {
		return 2
	}

	if slices.Equal(rsp, []byte{0x63, 0xC3}) {
		return 3
	}

	return -1
}
