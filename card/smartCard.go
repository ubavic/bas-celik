package card

import "github.com/telesma-app/pcsc"

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

func (card *VirtualCard) Status() (*pcsc.CardStatus, error) {
	status := pcsc.CardStatus{ATR: card.atr, ReaderNames: []string{"Virtual"}, State: pcsc.CardStatePowered}
	return &status, nil
}

func Transmit(cmd []byte) ([]byte, error) {
	return []byte{0x90, 0x00}, nil
}
