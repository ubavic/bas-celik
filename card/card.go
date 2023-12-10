package card

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/ebfe/scard"
	doc "github.com/ubavic/bas-celik/document"
	"golang.org/x/exp/constraints"
)

type Card interface {
	readFile([]byte, bool) ([]byte, error)
}

func ReadCard(sc *scard.Card) (doc.Document, error) {
	var card Card

	smartCardStatus, err := sc.Status()
	if err != nil {
		return nil, fmt.Errorf("reading card %w", err)
	}

	if reflect.DeepEqual(smartCardStatus.Atr, GEMALTO_ATR_1) {
		tempIdCard := Gemalto{smartCard: sc}
		if tempIdCard.testGemalto() {
			card = Gemalto{smartCard: sc}
		} else {
			card = VehicleCard{smartCard: sc}
		}
	} else if reflect.DeepEqual(smartCardStatus.Atr, GEMALTO_ATR_2) {
		tmpMedCard := MedicalCard{smartCard: sc}
		if tmpMedCard.testMedicalCard() {
			card = MedicalCard{smartCard: sc}
		} else {
			card = Gemalto{smartCard: sc}
		}
	} else if reflect.DeepEqual(smartCardStatus.Atr, GEMALTO_ATR_3) {
		card = Gemalto{smartCard: sc}
	} else if reflect.DeepEqual(smartCardStatus.Atr, APOLLO_ATR) {
		card = Apollo{smartCard: sc}
	} else if reflect.DeepEqual(smartCardStatus.Atr, MEDICAL_ATR) {
		card = MedicalCard{smartCard: sc}
	} else if reflect.DeepEqual(smartCardStatus.Atr, VEHICLE_ATR_0) {
		card = VehicleCard{smartCard: sc}
	} else if reflect.DeepEqual(smartCardStatus.Atr, VEHICLE_ATR_2) {
		card = VehicleCard{smartCard: sc}
	} else {
		return nil, fmt.Errorf("unknown card type: %s", hex.EncodeToString(smartCardStatus.Atr))
	}

	var d doc.Document

	switch card := card.(type) {
	case Gemalto:
		err = card.initCard()
	case VehicleCard:
		err = card.initCard()
	}

	if err != nil {
		return nil, fmt.Errorf("initializing card: %w", err)
	}

	switch card := card.(type) {
	case Apollo:
		d, err = readIDCard(card)
	case Gemalto:
		d, err = readIDCard(card)
	case MedicalCard:
		d, err = readMedicalCard(card)
	case VehicleCard:
		d, err = readVehicleCard(card)
	}

	if err != nil {
		return nil, fmt.Errorf("reading card: %w", err)
	}

	return d, nil
}

func assignField[T comparable](fields map[T][]byte, tag T, target *string) {
	val, ok := fields[tag]
	if ok {
		*target = string(val)
	} else {
		*target = ""
	}
}

func assignBoolField(fields map[uint][]byte, tag uint, target *bool) {
	val, ok := fields[tag]
	if ok && len(val) == 1 && val[0] == 0x31 {
		*target = true
	} else {
		*target = false
	}
}

func read(card *scard.Card, offset, length uint) ([]byte, error) {
	readSize := min(length, 0xFF)
	apu := buildAPDU(0x00, 0xB0, byte((0xFF00&offset)>>8), byte(offset&0xFF), nil, readSize)
	rsp, err := card.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("reading binary: %w", err)
	}

	if len(rsp) < 2 {
		return nil, fmt.Errorf("reading binary: bad status code")
	}

	return rsp[:len(rsp)-2], nil
}

func checkParseTLVData(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data %d", len(data))
	}
	return nil
}

// Parses simple TLV encoded data, where tag and length
// are encoded with two bytes
func parseTLV(data []byte) (map[uint][]byte, error) {
	err := checkParseTLVData(data)
	if err != nil {
		return nil, err
	}
	m := make(map[uint][]byte)
	offset := uint(0)

	for {
		tag := uint(binary.LittleEndian.Uint16(data[offset:]))
		length := uint(binary.LittleEndian.Uint16(data[offset+2:]))

		offset += 4
		value := data[offset : offset+length]
		m[tag] = value
		offset += length

		if offset >= uint(len(data)) {
			break
		}
	}

	return m, nil
}

func buildAPDU(cla, ins, p1, p2 byte, data []byte, ne uint) []byte {
	length := len(data)

	if length > 0xFFFF {
		panic(fmt.Errorf("APDU command length too large"))
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

	return apdu
}

func responseOK(rsp []byte) bool {
	return reflect.DeepEqual(rsp, []byte{0x90, 0x00})
}

func min[O constraints.Ordered](args ...O) O {
	if len(args) == 0 {
		return *new(O)
	}

	if args[0] != args[0] {
		return args[0]
	}

	min := args[0]
	for _, arg := range args[1:] {
		if arg != arg {
			return arg
		}

		if arg < min {
			min = arg
		}
	}
	return min
}
