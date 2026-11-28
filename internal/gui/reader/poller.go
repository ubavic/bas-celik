package reader

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/telesma-app/pcsc"
)

type ReaderPoller struct {
	watcher      pcsc.Watcher
	readerLister ReaderLister
	onCardEvent  func(string, bool)
	selections   chan string
}

type ReaderLister interface {
	SetReaders([]string, string)
	HookReaderChange(func(string))
}

func NewPoller(readerLister ReaderLister, onCardEvent func(string, bool)) (*ReaderPoller, error) {
	watcher, err := pcsc.Watch()
	if err != nil {
		return nil, fmt.Errorf("watching smart card readers: %w", err)
	}

	poller := &ReaderPoller{
		watcher:      watcher,
		readerLister: readerLister,
		onCardEvent:  onCardEvent,
		selections:   make(chan string, 1),
	}
	readerLister.HookReaderChange(poller.selectReader)

	return poller, nil
}

func (rp *ReaderPoller) StartPoller() {
	go rp.pollReaders()
}

func (rp *ReaderPoller) Close() error {
	return rp.watcher.Close()
}

func (rp *ReaderPoller) selectReader(name string) {
	// Fyne delivers selections serially. Keep the latest if a card is still being read.
	select {
	case <-rp.selections:
	default:
	}
	rp.selections <- name
}

func (rp *ReaderPoller) pollReaders() {
	readers := make(map[string]bool)
	for _, info := range rp.watcher.Snapshot().Readers {
		if !strings.HasPrefix(info.Name, "Windows Hello for Business") {
			readers[info.Name] = info.State&pcsc.ReaderStatePresent != 0
		}
	}

	selected := rp.updateReaders(readers, "")
	rp.onCardEvent(selected, readers[selected])

	events := rp.watcher.Listen()
	for {
		select {
		case name := <-rp.selections:
			if _, exists := readers[name]; exists && name != selected {
				selected = name
				rp.onCardEvent(selected, readers[selected])
			}
		case event, ok := <-events:
			if !ok {
				return
			}

			info := event.ReaderInfo
			if strings.HasPrefix(info.Name, "Windows Hello for Business") {
				continue
			}
			name := info.Name
			switch event.Type {
			case pcsc.DeviceEventReaderConnected, pcsc.DeviceEventReaderDisconnected:
				if event.Type == pcsc.DeviceEventReaderConnected {
					readers[name] = info.State&pcsc.ReaderStatePresent != 0
				} else {
					delete(readers, name)
				}
				newSelection := rp.updateReaders(readers, selected)
				if newSelection != selected {
					selected = newSelection
					rp.onCardEvent(selected, readers[selected])
				}
			case pcsc.DeviceEventCardInserted, pcsc.DeviceEventCardRemoved:
				present := event.Type == pcsc.DeviceEventCardInserted
				if previous, exists := readers[name]; exists && previous != present {
					readers[name] = present
					if name == selected {
						rp.onCardEvent(name, present)
					}
				}
			}
		}
	}
}

func (rp *ReaderPoller) updateReaders(readers map[string]bool, selected string) string {
	names := slices.Sorted(maps.Keys(readers))
	if _, exists := readers[selected]; !exists {
		selected = ""
		if len(names) > 0 {
			selected = names[0]
		}
	}
	rp.readerLister.SetReaders(names, selected)
	return selected
}
