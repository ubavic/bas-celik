//go:build windows

package reader

import (
	"github.com/ebfe/scard"
)

func (rp *ReaderPoller) waitForReaderChange(readersCount int) {
	state := scard.ReaderState{
		Reader:       `\\?PnP?\Notification`,
		CurrentState: scard.StateFlag(readersCount << 16),
		UserData:     nil,
		Atr:          nil,
	}

	states := []scard.ReaderState{state}
	rp.readerListerContext.GetStatusChange(states, -1)
}
