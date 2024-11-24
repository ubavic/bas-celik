//go:build darwin

package reader

import "time"

func (rp *ReaderPoller) waitForReaderChange() {
	time.Sleep(2 * time.Second)
}
