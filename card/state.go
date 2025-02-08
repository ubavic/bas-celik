package card

import (
	"strings"

	"github.com/ebfe/scard"
)

func FormatState(state scard.StateFlag) string {
	states := []string{}

	if state&scard.StateUnaware != 0 {
		states = append(states, "StateUnaware")
	}
	if state&scard.StateIgnore != 0 {
		states = append(states, "StateIgnore")
	}
	if state&scard.StateChanged != 0 {
		states = append(states, "StateChanged")
	}
	if state&scard.StateUnknown != 0 {
		states = append(states, "StateUnknown")
	}
	if state&scard.StatePresent != 0 {
		states = append(states, "StatePresent")
	}
	if state&scard.StateAtrmatch != 0 {
		states = append(states, "StateAtrmatch")
	}
	if state&scard.StateExclusive != 0 {
		states = append(states, "StateExclusive")
	}
	if state&scard.StateMute != 0 {
		states = append(states, "StateMute")
	}
	if state&scard.StateUnpowered != 0 {
		states = append(states, "StateUnpowered")
	}

	return strings.Join(states, " ")
}
