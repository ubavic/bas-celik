package card

import (
	"encoding/binary"
	"fmt"

	"github.com/ebfe/scard"
)

var GEMALTO_ATR_1 = []byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xB0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
}

var GEMALTO_ATR_2 = []byte{
	0x3B, 0xF9, 0x96, 0x00, 0x00, 0x80, 0x31, 0xFE,
	0x45, 0x53, 0x43, 0x45, 0x37, 0x20, 0x47, 0x43,
	0x4E, 0x33, 0x5E,
}

var GEMALTO_ATR_3 = []byte{
	0x3B, 0x9E, 0x96, 0x80, 0x31, 0xFE, 0x45, 0x53,
	0x43, 0x45, 0x20, 0x38, 0x2E, 0x30, 0x2D, 0x43,
	0x31, 0x56, 0x30, 0x0D, 0x0A, 0x6F,
}

type Gemalto struct {
	smartCard *scard.Card
}

func (card Gemalto) initCard() error {
	data := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x49, 0x44, 0x01}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("initializing card: %w", err)
	}

	if !responseOK(rsp) {
		return fmt.Errorf("initializing card: response not OK")
	}

	return nil
}

func (card Gemalto) readFile(name []byte, trim bool) ([]byte, error) {
	output := make([]byte, 0)

	_, err := card.selectFile(name, 4)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	data, err := read(card.smartCard, 0, 4)
	if err != nil {
		return nil, fmt.Errorf("reading file header: %w", err)
	}

	offset := uint(len(data))
	if offset < 3 {
		return nil, fmt.Errorf("file too short")
	}
	length := uint(binary.LittleEndian.Uint16(data[2:]))

	for length > 0 {
		data, err := read(card.smartCard, offset, length)
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	if trim {
		return output[4:], nil
	}

	return output, nil
}

func (card Gemalto) selectFile(name []byte, ne uint) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x08, 0x00, name, ne)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}

func (card Gemalto) testGemalto() bool {
	err := card.initCard()
	if err != nil {
		return false
	}

	_, err = card.readFile(DOCUMENT_FILE_LOC, false)
	return err == nil
}
