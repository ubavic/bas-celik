package card

import (
	"strings"

	"github.com/telesma-app/pcsc"
)

func FormatState(state pcsc.ReaderState) string {
	states := []string{}

	if state&pcsc.ReaderStateUnaware != 0 {
		states = append(states, "StateUnaware")
	}
	if state&pcsc.ReaderStateIgnore != 0 {
		states = append(states, "StateIgnore")
	}
	if state&pcsc.ReaderStateChanged != 0 {
		states = append(states, "StateChanged")
	}
	if state&pcsc.ReaderStateUnknown != 0 {
		states = append(states, "StateUnknown")
	}
	if state&pcsc.ReaderStatePresent != 0 {
		states = append(states, "StatePresent")
	}
	if state&pcsc.ReaderStateATRMatch != 0 {
		states = append(states, "StateAtrmatch")
	}
	if state&pcsc.ReaderStateExclusive != 0 {
		states = append(states, "StateExclusive")
	}
	if state&pcsc.ReaderStateMute != 0 {
		states = append(states, "StateMute")
	}
	if state&pcsc.ReaderStateUnpowered != 0 {
		states = append(states, "StateUnpowered")
	}

	return strings.Join(states, " ")
}
