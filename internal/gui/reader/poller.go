package reader

import (
	"fmt"
	"slices"
	"sync/atomic"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

var created = false
var createdPoller *ReaderPoller

type ReaderPoller struct {
	readerListerContext *scard.Context
	singleReaderContext *scard.Context
	readerLister        ReaderLister
	currentReader       string
	readers             []string
	readerPollerStarted atomic.Bool
	onCardEvent         func(string, *scard.Context)
}

type ReaderLister interface {
	SetReaders([]string, string)
	HookReaderChange(func(string))
}

func NewPoller(readerLister ReaderLister, onCardEvent func(string, *scard.Context)) (*ReaderPoller, error) {
	if created {
		logger.Info("NewPoller called more than once")
		return createdPoller, nil
	}

	readerListerContext, err := scard.EstablishContext()
	if err != nil {
		return nil, fmt.Errorf("creating reader list context: %w", err)
	}

	singleReaderContext, err := scard.EstablishContext()
	if err != nil {
		return nil, fmt.Errorf("creating single reader context: %w", err)
	}

	poller := ReaderPoller{
		readerListerContext: readerListerContext,
		singleReaderContext: singleReaderContext,
		readerLister:        readerLister,
		onCardEvent:         onCardEvent,
	}

	readerLister.HookReaderChange(poller.SetReader)

	created = true
	createdPoller = &poller

	return createdPoller, nil
}

func (rp *ReaderPoller) StartPoller() {
	rp.onCardEvent("", rp.singleReaderContext)
	go rp.pollReaders()
}

func (rp *ReaderPoller) pollReaders() {
	newReaders, _ := rp.readerListerContext.ListReaders()
	newReader := ""
	if len(rp.readers) > 0 {
		newReader = rp.readers[0]
	}

	rp.readerLister.SetReaders(newReaders, newReader)

	readerToSelect := ""

	for {
		newReaders, _ = rp.readerListerContext.ListReaders()

		readersCount := len(newReaders)

		if slices.Compare(rp.readers, newReaders) != 0 {
			logger.Debug("New readers list: " + fmt.Sprintf("%v", newReaders))

			rp.readers = newReaders

			if slices.Contains(rp.readers, rp.currentReader) {
				readerToSelect = rp.currentReader
			} else {
				if len(rp.readers) == 0 {
					readerToSelect = ""
				} else {
					readerToSelect = rp.readers[0]
				}
			}

			rp.readerLister.SetReaders(rp.readers, readerToSelect)
			rp.SetReader(readerToSelect)
		}

		rp.waitForReaderChange(readersCount)
		logger.Debug("Readers changed...")
	}
}

func (rp *ReaderPoller) SetReader(newReader string) {
	if rp.currentReader == newReader {
		return
	}

	logger.Debug("Setting reader: " + newReader)

	rp.currentReader = newReader
	CancelReaderPoler()
	rp.onCardEvent(newReader, rp.singleReaderContext)
	go rp.readerPoller(newReader)
}

func (rp *ReaderPoller) readerPoller(selectedReader string) {
	for {
		if selectedReader == "" {
			return
		}

		state := scard.ReaderState{
			Reader:       selectedReader,
			CurrentState: scard.StateUnaware,
		}

		states := []scard.ReaderState{state}

		rp.singleReaderContext.GetStatusChange(states, 0)
		for i := range states {
			states[i].CurrentState = states[i].EventState
		}

		logger.Debug("Reader poller event 1: " + card.FormatState(states[0].CurrentState))

		rp.readerPollerStarted.Store(true)
		err := rp.singleReaderContext.GetStatusChange(states, -1)
		if err != nil {
			return
		}

		logger.Debug("Reader poller event 2: " + card.FormatState(states[0].CurrentState))

		rp.onCardEvent(selectedReader, rp.singleReaderContext)
	}
}

func CancelReaderPoler() {
	logger.Debug("Canceling reader poller...")
	if createdPoller.readerPollerStarted.Load() {
		createdPoller.readerPollerStarted.Store(false)
		createdPoller.singleReaderContext.Cancel()
	}
}

func RestartReaderPoler() {
	logger.Debug("Restarting reader poller...")
	go createdPoller.readerPoller(createdPoller.currentReader)
}
