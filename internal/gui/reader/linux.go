//go:build linux

package reader

import (
	"github.com/ebfe/scard"
)

func (rp *ReaderPoller) waitForReaderChange(_ int) {
	state := scard.ReaderState{
		Reader:       `\\?PnP?\Notification`,
		CurrentState: scard.StateUnaware,
	}

	states := []scard.ReaderState{state}
	rp.readerListerContext.GetStatusChange(states, -1)
}
