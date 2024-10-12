package card

import (
	"errors"

	"github.com/ubavic/bas-celik/document"
)

type UnknownDocumentCard struct {
	atr       Atr
	smartCard Card
}

func (card UnknownDocumentCard) Atr() Atr {
	return card.atr
}

func (card UnknownDocumentCard) readFile(_ []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (card UnknownDocumentCard) InitCard() error {
	return nil
}

func (card UnknownDocumentCard) ReadCard() error {
	return nil
}

func (card UnknownDocumentCard) GetDocument() (document.Document, error) {
	return nil, nil
}
