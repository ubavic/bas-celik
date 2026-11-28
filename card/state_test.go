package card_test

import (
	"testing"

	"github.com/telesma-app/pcsc"
	"github.com/ubavic/bas-celik/v2/card"
)

func TestFormatState(t *testing.T) {
	tests := []struct {
		name  string
		state pcsc.ReaderState
		want  string
	}{
		{
			name:  "no flags",
			state: 0,
			want:  "",
		},
		{
			name:  "single flag",
			state: pcsc.ReaderStatePresent,
			want:  "StatePresent",
		},
		{
			name:  "state unaware is indistinguishable from zero",
			state: pcsc.ReaderStateUnaware,
			want:  "",
		},
		{
			name:  "multiple flags keep declaration order",
			state: pcsc.ReaderStateMute | pcsc.ReaderStateChanged,
			want:  "StateChanged StateMute",
		},
		{
			name: "all known flags",
			state: pcsc.ReaderStateIgnore |
				pcsc.ReaderStateChanged |
				pcsc.ReaderStateUnknown |
				pcsc.ReaderStatePresent |
				pcsc.ReaderStateATRMatch |
				pcsc.ReaderStateExclusive |
				pcsc.ReaderStateMute |
				pcsc.ReaderStateUnpowered,
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
