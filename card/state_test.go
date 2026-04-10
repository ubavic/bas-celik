package card_test

import (
	"testing"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/v2/card"
)

func TestFormatState(t *testing.T) {
	tests := []struct {
		name  string
		state scard.StateFlag
		want  string
	}{
		{
			name:  "no flags",
			state: 0,
			want:  "",
		},
		{
			name:  "single flag",
			state: scard.StatePresent,
			want:  "StatePresent",
		},
		{
			name:  "state unaware is indistinguishable from zero",
			state: scard.StateUnaware,
			want:  "",
		},
		{
			name:  "multiple flags keep declaration order",
			state: scard.StateMute | scard.StateChanged,
			want:  "StateChanged StateMute",
		},
		{
			name:  "all known flags",
			state: scard.StateIgnore |
				scard.StateChanged |
				scard.StateUnknown |
				scard.StatePresent |
				scard.StateAtrmatch |
				scard.StateExclusive |
				scard.StateMute |
				scard.StateUnpowered,
			want: "StateIgnore StateChanged StateUnknown StatePresent StateAtrmatch StateExclusive StateMute StateUnpowered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := card.FormatState(tt.state)
			if got != tt.want {
				t.Fatalf("FormatState(%v) = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}
