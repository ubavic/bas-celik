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
	Test() bool
	Atr() Atr
}

// Represents a different types of smart card documents.
// Each value of `CardDocumentType` is represented with a struct
// that satisfies `CardDocument` interface.
type CardDocumentType uint8

const (
	UnknownDocumentCardType = CardDocumentType(iota)
	ApolloIdDocumentCardType
	GemaltoIdDocumentCardType
	MedicalDocumentCardType
	VehicleDocumentCardType
)

var ErrUnknownCard = errors.New("unknown card")

// Detects Card Document from card's ATR
// Ambiguous cases are solved by reading specific card content
func DetectCardDocument(sc Card) (CardDocument, error) {
	smartCardStatus, err := sc.Status()
	if err != nil {
		return nil, fmt.Errorf("reading card status %w", err)
	}

	atr := Atr(smartCardStatus.Atr)

	possibleCardTypes := DetectCardDocumentByAtr(atr)

	for _, cardType := range possibleCardTypes {
		switch cardType {
		case ApolloIdDocumentCardType:
			card := &Apollo{atr: atr, smartCard: sc}
			return card, nil
		case GemaltoIdDocumentCardType:
			card := Gemalto{atr: atr, smartCard: sc}
			if card.Test() {
				return &card, nil
			}
		case VehicleDocumentCardType:
			card := VehicleCard{atr: atr, smartCard: sc}
			if card.Test() {
				return &card, nil
			}
		case MedicalDocumentCardType:
			card := MedicalCard{atr: atr, smartCard: sc}
			if card.Test() {
				return &card, nil
			}
		default:
			card := &UnknownDocumentCard{atr: atr, smartCard: sc}
			return card, ErrUnknownCard
		}
	}

	return nil, errors.New("unexpected card type")
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
