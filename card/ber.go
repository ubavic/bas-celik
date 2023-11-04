package card

import (
	"encoding/binary"
	"errors"
)

type BER struct {
	tag       uint32
	primitive bool
	data      []byte
	children  []BER
}

func parseBER(data []byte) (*BER, error) {
	primitive, constructed, err := parseBERLayer(data)

	if err != nil {
		return nil, err
	}

	ber := BER{
		tag:       0,
		primitive: false,
		data:      nil,
		children:  []BER{},
	}

	for t, v := range primitive {
		prim := BER{
			tag:       t,
			primitive: true,
			data:      v,
			children:  nil,
		}
		ber.add(prim)
	}

	for t, v := range constructed {
		children, err := parseBER(v)
		if err != nil {
			return nil, err
		}
		prim := BER{
			tag:       t,
			primitive: true,
			data:      nil,
			children:  children.children,
		}
		ber.add(prim)
	}

	return &ber, nil

}

func (tree BER) access(address ...uint32) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("bad address")
	} else if len(address) == 1 {
		return tree.data, nil
	} else {
		var found *BER = nil
		for i := range tree.children {
			if tree.children[i].tag == address[0] {
				found = &tree.children[i]
				break
			}
		}
		if found != nil {
			return found.access(address[1:]...)
		} else {
			return nil, errors.New("tag not found")
		}
	}
}

// Insert a new node into BER tree
func (into *BER) add(new BER) error {
	if into.primitive {
		return errors.New("can't add a value into primitive value")
	}

	var targetField *BER

	alradyExists := false
	for _, v := range into.children {
		if v.tag == new.tag {
			alradyExists = true
			targetField = &v
		}
	}

	if !alradyExists {
		into.children = append(into.children, new)
		return nil
	} else {
		if targetField.primitive == new.primitive {
			if targetField.primitive {
				*targetField = new
			} else {
				for _, vv := range new.children {
					err := targetField.add(vv)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return errors.New("types don't match")
		}
	}

	return nil
}

func emptyTree() BER {
	return BER{
		tag:       0,
		primitive: false,
		data:      []byte{},
		children:  []BER{},
	}
}

// Parses one level of BER-TLV encoded data
// according to ISO/IEC 7816-4 (2005)
// Returns map of primitive and constructed fields
func parseBERLayer(data []byte) (map[uint32][]byte, map[uint32][]byte, error) {
	primF := make(map[uint32][]byte)
	consF := make(map[uint32][]byte)
	var primitive bool
	offset := uint32(0)

	for {
		tag := uint32(data[offset])

		if tag&0x20 != 0 {
			primitive = false
		} else {
			primitive = true
		}

		if 0x1F&tag != 0x1F {
			offset += 1
		} else {
			if data[offset+1]&0x80 == 0x00 {
				tag = uint32(data[offset])<<8 + uint32(data[offset+1])
				offset += 2
			} else {
				tag = uint32(data[offset])<<16 +
					uint32(data[offset+1])<<8 +
					uint32(data[offset+2])
				offset += 3
			}
		}

		length, offsetDelta, err := parseBerLength(data)
		if err != nil {
			return nil, nil, errors.New("invalid length")
		}

		offset += offsetDelta
		value := data[offset : offset+length]

		if primitive {
			primF[tag] = value
		} else {
			consF[tag] = value
		}

		offset += length

		if offset == uint32(len(data)) {
			break
		} else if offset > uint32(len(data)) {
			return nil, nil, errors.New("invalid length")
		}
	}

	return primF, consF, nil
}

func (tree *BER) assignToFrom(target *string, address ...uint32) {
	bytes, error := tree.access(address...)
	if error == nil {
		*target = string(bytes)
	}
}

func parseBerLength(data []byte) (uint32, uint32, error) {
	length := uint32(data[0])
	offset := uint32(0)
	if length < 0x80 {
		offset = 1
	} else if length == 0x81 {
		length = uint32(data[1])
		offset = 2
	} else if length == 0x82 {
		length = uint32(binary.BigEndian.Uint16(data[1:]))
		offset = 3
	} else if length == 0x83 {
		length = 0x00FFFFFF & binary.BigEndian.Uint32(data)
		offset = 4
	} else if length == 0x84 {
		length = binary.BigEndian.Uint32(data[1:])
		offset = 5
	} else {
		return 0, 0, errors.New("invalid length")
	}

	return length, uint32(offset), nil
}
