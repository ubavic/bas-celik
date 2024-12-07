//go:build darwin

package reader

import "time"

func (rp *ReaderPoller) waitForReaderChange(_ int) {
	time.Sleep(2 * time.Second)
}
