package reader

import (
	"fmt"
	"slices"

	"github.com/ebfe/scard"
)

var created = false

type ReaderPoller struct {
	readerListerContext *scard.Context
	singleReaderContext *scard.Context
	readerLister        ReaderLister
	currentReader       string
	readers             []string
	readerStateChange   chan<- string
}

type ReaderLister interface {
	SetReaders([]string, string)
	HookReaderChange(func(string))
}

func NewPoller(readerLister ReaderLister, readerStateChange chan<- string) (*ReaderPoller, error) {
	if created {
		panic("you can create only single instance of ReaderPoller")
	}
	created = true

	readerListerContext, err := scard.EstablishContext()
	if err != nil {
		return nil, fmt.Errorf("creating reader poller: %w", err)
	}

	singleReaderContext, err := scard.EstablishContext()
	if err != nil {
		return nil, fmt.Errorf("creating reader poller: %w", err)
	}

	poller := ReaderPoller{
		readerListerContext: readerListerContext,
		singleReaderContext: singleReaderContext,
		readerLister:        readerLister,
		readerStateChange:   readerStateChange,
	}

	readerLister.HookReaderChange(poller.SetReader)

	go poller.pollReaders()

	return &poller, nil
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
	}
}

func (rp *ReaderPoller) SetReader(newReader string) {
	if rp.currentReader == newReader {
		return
	}

	rp.currentReader = newReader
	rp.readerStateChange <- newReader
	rp.singleReaderContext.Cancel()
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

		err := rp.singleReaderContext.GetStatusChange(states, -1)
		if err != nil {
			return
		}

		rp.readerStateChange <- selectedReader
	}
}
