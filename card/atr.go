package card

import (
	"encoding/hex"
	"slices"
)

type Atr []byte

func (atr Atr) String() string {
	return hex.EncodeToString([]byte(atr))
}

func (atr Atr) Is(otherAtr Atr) bool {
	return slices.Equal(atr, otherAtr)
}
