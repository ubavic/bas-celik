package card

import "github.com/ebfe/scard"

type VirtualCard struct {
	atr   []byte
	files map[uint32][]byte
}

func MakeVirtualCard(atr []byte, fs map[uint32][]byte) *VirtualCard {
	vc := VirtualCard{
		atr:   atr,
		files: fs,
	}

	return &vc
}

func (card *VirtualCard) Status() (*scard.CardStatus, error) {
	status := scard.CardStatus{Atr: card.atr, Reader: "Virtual", State: scard.Powered}
	return &status, nil
}

func Transmit(cmd []byte) ([]byte, error) {
	return []byte{0x90, 0x00}, nil
}
