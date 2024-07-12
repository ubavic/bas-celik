// Package card provides functions for communication with smart cards.
// It includes implementations for handling different types of smart cards
// and reading associated documents.
package card

import (
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/ebfe/scard"
	doc "github.com/ubavic/bas-celik/document"
)

// Represents a physical or virtual smart card.
// Essentially it is just a wrapper for the scard.Card type,
// but it also allows virtual cards which can be useful for testing.
type Card interface {
	Status() (*scard.CardStatus, error)
	Transmit(cmd []byte) ([]byte, error)
}

// Represents a smart card with a document.
// All types of documents that Bas Celik can read should satisfy this interface
type CardDocument interface {
	readFile([]byte, bool) ([]byte, error)
}

// Reads a smart card and returns the associated document.
// It determines the type of the card based on its ATR value and initializes
// the appropriate card implementation for further reading.
func ReadCard(sc Card) (doc.Document, error) {
	var card CardDocument

	smartCardStatus, err := sc.Status()
	if err != nil {
		return nil, fmt.Errorf("reading card %w", err)
	}

	if slices.Equal(smartCardStatus.Atr, GEMALTO_ATR_1) {
		tempIdCard := Gemalto{smartCard: sc}
		if tempIdCard.testGemalto() {
			card = Gemalto{smartCard: sc}
		} else {
			card = VehicleCard{smartCard: sc}
		}
	} else if slices.Equal(smartCardStatus.Atr, GEMALTO_ATR_2) {
		tempIdCard := Gemalto{smartCard: sc}
		tmpMedCard := MedicalCard{smartCard: sc}
		if tempIdCard.testGemalto() {
			card = Gemalto{smartCard: sc}
		} else if tmpMedCard.testMedicalCard() {
			card = MedicalCard{smartCard: sc}
		} else {
			card = VehicleCard{smartCard: sc}
		}
	} else if slices.Equal(smartCardStatus.Atr, GEMALTO_ATR_3) {
		card = Gemalto{smartCard: sc}
	} else if slices.Equal(smartCardStatus.Atr, APOLLO_ATR) {
		card = Apollo{smartCard: sc}
	} else if slices.Equal(smartCardStatus.Atr, MEDICAL_ATR) {
		card = MedicalCard{smartCard: sc}
	} else if slices.Equal(smartCardStatus.Atr, VEHICLE_ATR_0) {
		card = VehicleCard{smartCard: sc}
	} else if slices.Equal(smartCardStatus.Atr, VEHICLE_ATR_2) {
		card = VehicleCard{smartCard: sc}
	} else if slices.Equal(smartCardStatus.Atr, VEHICLE_ATR_3) {
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
		d, err = readIdCard(card)
	case Gemalto:
		d, err = readIdCard(card)
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

// Reads binary data from the card starting from the specified offset and with the specified length.
func read(card Card, offset, length uint) ([]byte, error) {
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

// Checks if the card response indicates no error.
func responseOK(rsp []byte) bool {
	if len(rsp) < 2 {
		return false
	}

	return slices.Equal(rsp[len(rsp)-2:], []byte{0x90, 0x00})
}
