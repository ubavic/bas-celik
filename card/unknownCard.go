package card

import "errors"

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

func (card UnknownDocumentCard) initCard() error {
	return nil
}
