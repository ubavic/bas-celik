package tlv

import (
	"encoding/binary"

	"github.com/ubavic/bas-celik/card/cardErrors"
)

// Parses simple TLV-encoded data and returns a map of tags to values.
// It assumes that tag and length are encoded with two bytes each.
func ParseTLV(data []byte) (map[uint][]byte, error) {
	if len(data) == 0 {
		return nil, cardErrors.ErrInvalidLength
	}

	m := make(map[uint][]byte)
	offset := uint(0)

	for {
		tag := uint(binary.LittleEndian.Uint16(data[offset:]))
		length := uint(binary.LittleEndian.Uint16(data[offset+2:]))

		offset += 4

		if offset+length > uint(len(data)) {
			return nil, cardErrors.ErrInvalidLength
		}

		value := data[offset : offset+length]
		m[tag] = value
		offset += length

		if offset >= uint(len(data)) {
			break
		}
	}

	return m, nil
}

// Assigns the value from the provided fields map to the target string, based on the specified tag.
// If the tag is not present in the map, the target is set to an empty string.
func AssignField[T comparable](fields map[T][]byte, tag T, target *string) {
	val, ok := fields[tag]
	if ok {
		*target = string(val)
	} else {
		*target = ""
	}
}

// Assigns a boolean value from the provided fields map to the target, based on the specified tag.
// If the tag is not present in the map or the value is not 0x31, the target is set to false.
func AssignBoolField(fields map[uint][]byte, tag uint, target *bool) {
	val, ok := fields[tag]
	if ok && len(val) == 1 && val[0] == 0x31 {
		*target = true
	} else {
		*target = false
	}
}
