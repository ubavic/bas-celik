package card

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"reflect"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/document"
	doc "github.com/ubavic/bas-celik/document"
)

var DOCUMENT_FILE_LOC = []byte{0x0F, 0x02}
var PERSONAL_FILE_LOC = []byte{0x0F, 0x03}
var RESIDENCE_FILE_LOC = []byte{0x0F, 0x04}
var PHOTO_FILE_LOC = []byte{0x0F, 0x06}

type Card interface {
	readFile([]byte, bool) ([]byte, error)
}

func ReadCard(sc *scard.Card, doc *doc.Document) error {
	var card Card

	smartCardStatus, err := sc.Status()
	if err != nil {
		return fmt.Errorf("error reading card")
	}

	if reflect.DeepEqual(smartCardStatus.Atr, GEMALTO_ATR_1) || reflect.DeepEqual(smartCardStatus.Atr, GEMALTO_ATR_2) {
		card = Gemalto{smartCard: sc}
		connectGemalto(sc)
	} else if reflect.DeepEqual(smartCardStatus.Atr, APOLLO_ATR) {
		card = Apollo{smartCard: sc}
	} else {
		return fmt.Errorf("unknown card type")
	}

	var fields map[uint]string

	assignField := func(tag uint, target *string) {
		val, ok := fields[tag]
		if ok {
			*target = val
		} else {
			*target = ""
		}
	}

	rsp, err := card.readFile(DOCUMENT_FILE_LOC, false)
	if err != nil {
		return fmt.Errorf("error reading document file %w", err)
	}

	fields = parseResponse(rsp)
	assignField(1546, &doc.DocumentNumber)
	assignField(1549, &doc.IssuingDate)
	assignField(1550, &doc.ExpiryDate)
	assignField(1551, &doc.IssuingAuthority)
	doc.IssuingDate = document.FormatDate(doc.IssuingDate)
	doc.ExpiryDate = document.FormatDate(doc.ExpiryDate)

	rsp, err = card.readFile(PERSONAL_FILE_LOC, false)
	if err != nil {
		return fmt.Errorf("error reading personal file %w", err)
	}

	fields = parseResponse(rsp)
	assignField(1558, &doc.PersonalNumber)
	assignField(1559, &doc.Surname)
	assignField(1560, &doc.GivenName)
	assignField(1561, &doc.ParentName)
	assignField(1562, &doc.Sex)
	assignField(1563, &doc.PlaceOfBirth)
	assignField(1564, &doc.CommunityOfBirth)
	assignField(1565, &doc.StateOfBirth)
	assignField(1566, &doc.DateOfBirth)
	doc.DateOfBirth = document.FormatDate(doc.DateOfBirth)

	rsp, err = card.readFile(RESIDENCE_FILE_LOC, false)
	if err != nil {
		return fmt.Errorf("error reading residence file %w", err)
	}

	fields = parseResponse(rsp)
	assignField(1568, &doc.State)
	assignField(1569, &doc.Community)
	assignField(1570, &doc.Place)
	assignField(1571, &doc.Street)
	assignField(1572, &doc.AddressNumber)
	assignField(1573, &doc.AddressLetter)
	assignField(1574, &doc.AddressEntrance)
	assignField(1575, &doc.AddressFloor)
	assignField(1578, &doc.AddressApartmentNumber)
	assignField(1580, &doc.AddressDate)
	doc.AddressDate = document.FormatDate(doc.AddressDate)

	rsp, err = card.readFile(PHOTO_FILE_LOC, true)
	if err != nil {
		return fmt.Errorf("error reading photo file %w", err)
	}

	doc.Photo, _, err = image.Decode(bytes.NewReader(rsp))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error decoding photo file %w", err)
	}

	sc.Disconnect(scard.LeaveCard)

	return nil
}

func selectFile(card *scard.Card, name []byte, ne uint) ([]byte, error) {
	apu, err := buildAPDU(0x00, 0xA4, 0x08, 0x00, name, ne)

	if err != nil {
		return nil, fmt.Errorf("error selecting file: %w", err)
	}

	rsp, err := card.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("error selecting file: %w", err)
	}

	return rsp, nil
}

func read(card *scard.Card, offset, length uint) ([]byte, error) {
	readSize := length
	if readSize >= 0xFE {
		readSize = 0xFE
	}

	apu, err := buildAPDU(0x00, 0xB0, byte((0xFF00&offset)>>8), byte(offset&0xFF), nil, readSize)
	if err != nil {
		return nil, fmt.Errorf("error reading binary: %w", err)
	}

	rsp, err := card.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("error reading binary: %w", err)
	}

	if len(rsp) < 2 {
		return nil, fmt.Errorf("error reading binary: bad status code")
	}

	return rsp[:len(rsp)-2], nil
}

func parseResponse(data []byte) map[uint]string {
	m := make(map[uint]string)
	offset := uint(0)

	for {
		tag := uint(binary.LittleEndian.Uint16(data[offset:]))
		length := uint(binary.LittleEndian.Uint16(data[offset+2:]))

		offset += 4
		value := data[offset : offset+length]
		m[tag] = string(value)
		offset += length

		if offset >= uint(len(data)) {
			break
		}
	}

	return m
}

func buildAPDU(cla, ins, p1, p2 byte, data []byte, ne uint) ([]byte, error) {
	length := len(data)

	if length > 65535 {
		return nil, errors.New("length is too large")
	}

	apdu := make([]byte, 4)
	apdu[0] = cla
	apdu[1] = ins
	apdu[2] = p1
	apdu[3] = p2

	if length == 0 {
		if ne != 0 {
			if ne <= 256 {
				l := byte(0x00)
				if ne != 256 {
					l = byte(ne)
				}
				apdu = append(apdu, l)
			} else {
				var l1, l2 byte
				if ne == 65536 {
					l1 = 0
					l2 = 0
				} else {
					l1 = byte(ne >> 8)
					l2 = byte(ne)
				}
				apdu = append(apdu, []byte{l1, l2}...)
			}
		}
	} else {
		if ne == 0 {
			if length <= 255 {
				apdu = append(apdu, byte(length))
				apdu = append(apdu, data...)
			} else {
				l := []byte{0x0, byte(length >> 8), byte(length)}
				apdu = append(apdu, l...)
				apdu = append(apdu, data...)
			}
		} else {
			if length <= 255 && ne <= 256 {
				apdu = append(apdu, byte(length))
				apdu = append(apdu, data...)
				if ne != 256 {
					apdu = append(apdu, byte(ne))
				} else {
					apdu = append(apdu, 0x00)
				}
			} else {
				l := []byte{0x00, byte(length >> 8), byte(length)}
				apdu = append(apdu, l...)
				apdu = append(apdu, data...)
				if ne != 65536 {
					neB := []byte{byte(ne >> 8), byte(ne)}
					apdu = append(apdu, neB...)
				}
			}

		}
	}

	return apdu, nil
}

func responseOK(rsp []byte) bool {
	return reflect.DeepEqual(rsp, []byte{0x90, 0x00})
}
