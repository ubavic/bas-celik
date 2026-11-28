// Package card provides functions for communication with smart cards.
// It includes implementations for handling different types of smart cards
// and reading associated documents.
package card

import (
	"context"
	"errors"
	"fmt"

	"github.com/telesma-app/iso7816"
	"github.com/telesma-app/pcsc"
	doc "github.com/ubavic/bas-celik/v2/document"
)

// Card represents a PC/SC connection or a virtual card used in tests.
type Card interface {
	Status() (*pcsc.CardStatus, error)
	Transmit(context.Context, []byte) ([]byte, error)
	BeginTransaction(context.Context) error
	EndTransaction(pcsc.Disposition) error
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

	atr := Atr(smartCardStatus.ATR)

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

	return nil, fmt.Errorf("unexpected card type (ATR: %s)", atr)
}

// Reads binary data from the card starting from the specified offset and with the specified length.
func read(card Card, offset, length uint) ([]byte, error) {
	if length == 0 {
		return nil, nil
	}

	if offset >= 0xFFFF {
		return nil, fmt.Errorf("reading binary: offset too large")
	}

	readSize := min(length, 0xFF)
	apu := buildAPDU(0x00, 0xB0, byte((0xFF00&offset)>>8), byte(offset&0xFF), nil, readSize)
	rsp, err := card.Transmit(context.Background(), apu)
	if err != nil {
		return nil, fmt.Errorf("reading binary: %w", err)
	}

	response, err := iso7816.ParseResponse(rsp)
	if err != nil {
		return nil, fmt.Errorf("reading binary: bad status code")
	}

	return response.Data, nil
}

// Checks if the card response indicates no error.
func responseOK(rsp []byte) bool {
	if len(rsp) < 2 {
		return false
	}

	return iso7816.NewStatusWord(rsp[len(rsp)-2], rsp[len(rsp)-1]).IsSuccess()
}

// Trim four bytes from the start of the slice
func trim4b(data []byte) []byte {
	if len(data) > 4 {
		return data[4:]
	}

	return data
}
