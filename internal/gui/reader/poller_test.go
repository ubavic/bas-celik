package reader

import (
	"slices"
	"testing"

	"github.com/telesma-app/pcsc"
)

func TestPollReadersTracksSelectedCard(t *testing.T) {
	first := &pcsc.ReaderInfo{Name: "First", State: pcsc.ReaderStatePresent}
	second := &pcsc.ReaderInfo{Name: "Second", State: pcsc.ReaderStateEmpty}
	third := &pcsc.ReaderInfo{Name: "Third", State: pcsc.ReaderStatePresent}
	events := []pcsc.DeviceEvent{
		{Type: pcsc.DeviceEventReaderConnected, ReaderInfo: third},
		{Type: pcsc.DeviceEventCardInserted, ReaderInfo: third},
		{Type: pcsc.DeviceEventCardRemoved, ReaderInfo: first},
		{Type: pcsc.DeviceEventReaderDisconnected, ReaderInfo: first},
		{Type: pcsc.DeviceEventCardInserted, ReaderInfo: second},
		{Type: pcsc.DeviceEventCardRemoved, ReaderInfo: second},
		{Type: pcsc.DeviceEventReaderDisconnected, ReaderInfo: second},
	}
	listen := make(chan pcsc.DeviceEvent, len(events))
	for _, event := range events {
		listen <- event
	}
	close(listen)

	type cardEvent struct {
		name    string
		present bool
	}
	var received []cardEvent
	lister := &testReaderLister{}
	poller := ReaderPoller{
		watcher:      testWatcher{readers: []*pcsc.ReaderInfo{second, first}, events: listen},
		readerLister: lister,
		onCardEvent: func(name string, present bool) {
			received = append(received, cardEvent{name, present})
		},
	}
	poller.pollReaders()

	wantEvents := []cardEvent{
		{"First", true}, {"First", false}, {"Second", false},
		{"Second", true}, {"Second", false}, {"Third", true},
	}
	if !slices.Equal(received, wantEvents) {
		t.Fatalf("card events = %v, want %v", received, wantEvents)
	}
	wantSelections := []string{"First", "First", "Second", "Third"}
	if !slices.Equal(lister.selections, wantSelections) {
		t.Fatalf("selected readers = %v, want %v", lister.selections, wantSelections)
	}
}

type testWatcher struct {
	readers []*pcsc.ReaderInfo
	events  <-chan pcsc.DeviceEvent
}

func (w testWatcher) Snapshot() pcsc.Snapshot {
	return pcsc.Snapshot{Readers: w.readers}
}

func (w testWatcher) Listen() <-chan pcsc.DeviceEvent {
	return w.events
}

func (w testWatcher) Close() error { return nil }

type testReaderLister struct {
	selections []string
}

func (l *testReaderLister) SetReaders(_ []string, selected string) {
	l.selections = append(l.selections, selected)
}

func (l *testReaderLister) HookReaderChange(func(string)) {}
