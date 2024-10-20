// Package card provides functions for communication with smart cards.
// It includes implementations for handling different types of smart cards
// and reading associated documents.
package card

import (
	"errors"
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
	Transmit([]byte) ([]byte, error)
}

// Represents a smart card with a document.
// All types of documents that Bas Celik can read should satisfy this interface
type CardDocument interface {
	ReadFile([]byte) ([]byte, error)
	InitCard() error
	ReadCard() error
	GetDocument() (doc.Document, error)
	Atr() Atr
}

var ErrUnknownCard = errors.New("unknown card")

// Detects Card Document from card's ATR
// Ambiguous cases are solved by reading specific card content
func DetectCardDocument(sc Card) (CardDocument, error) {
	var card CardDocument

	smartCardStatus, err := sc.Status()
	if err != nil {
		return nil, fmt.Errorf("reading card %w", err)
	}

	atr := Atr(smartCardStatus.Atr)

	if atr.Is(GEMALTO_ATR_1) {
		tempIdCard := Gemalto{smartCard: sc}
		if tempIdCard.testGemalto() {
			card = &Gemalto{atr: atr, smartCard: sc}
		} else {
			card = &VehicleCard{atr: atr, smartCard: sc}
		}
	} else if atr.Is(GEMALTO_ATR_2) || atr.Is(GEMALTO_ATR_3) {
		tempIdCard := Gemalto{smartCard: sc}
		tmpMedCard := MedicalCard{smartCard: sc}
		if tempIdCard.testGemalto() {
			card = &Gemalto{smartCard: sc}
		} else if tmpMedCard.testMedicalCard() {
			card = &MedicalCard{atr: atr, smartCard: sc}
		} else {
			card = &VehicleCard{atr: atr, smartCard: sc}
		}
	} else if atr.Is(GEMALTO_ATR_4) {
		card = &Gemalto{atr: atr, smartCard: sc}
	} else if atr.Is(APOLLO_ATR) {
		card = &Apollo{atr: atr, smartCard: sc}
	} else if atr.Is(MEDICAL_ATR_1) {
		card = &MedicalCard{atr: atr, smartCard: sc}
	} else if atr.Is(MEDICAL_ATR_2) {
		card = &MedicalCard{atr: atr, smartCard: sc}
	} else if atr.Is(VEHICLE_ATR_0) {
		card = &VehicleCard{atr: atr, smartCard: sc}
	} else if atr.Is(VEHICLE_ATR_2) {
		card = &VehicleCard{atr: atr, smartCard: sc}
	} else if atr.Is(VEHICLE_ATR_3) {
		card = &VehicleCard{atr: atr, smartCard: sc}
	} else if atr.Is(VEHICLE_ATR_4) {
		card = &VehicleCard{atr: atr, smartCard: sc}
	} else {
		card = &UnknownDocumentCard{atr: atr, smartCard: sc}
		return card, ErrUnknownCard
	}

	return card, nil
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

// Trim four bytes from the start of the slice
func trim4b(data []byte) []byte {
	if len(data) > 4 {
		return data[4:]
	}

	return data
}
