//go:build !darwin

package reader

import (
	"github.com/ebfe/scard"
)

func (rp *ReaderPoller) waitForReaderChange() {
	state := scard.ReaderState{
		Reader:       `\\?PnP?\Notification`,
		CurrentState: scard.StateUnaware,
	}

	states := []scard.ReaderState{state}
	rp.readerListerContext.GetStatusChange(states, -1)
}
