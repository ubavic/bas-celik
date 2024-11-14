package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
)

func establishContextAndStartPollers() {
	setStartPage(t("poller.connectingReader"), "", nil)

	// if we don't sequentially establish contexts,
	// they will not be established successfully

	ctx1, err := scard.EstablishContext()
	if err != nil {
		setStartPage(
			t("error.driver"),
			t("error.driverExplanation"),
			fmt.Errorf("establishing context: %w", err))
	}

	ctx2, err := scard.EstablishContext()
	if err != nil {
		setStartPage(
			t("error.driver"),
			t("error.driverExplanation"),
			fmt.Errorf("establishing context: %w", err))
	}

	go pollReaders(ctx1)
	go poller(ctx2)
}

func pollReaders(ctx *scard.Context) {
	for {
		readersNames, _ := ctx.ListReaders()

		if len(readersNames) > 0 {
			states := make([]scard.ReaderState, 0, len(readersNames))
			for _, readerName := range readersNames {
				states = append(states, scard.ReaderState{
					Reader:       readerName,
					CurrentState: scard.StateUnaware,
				})
			}

			state.mu.Lock()
			state.toolbar.SetReaders(readersNames)
			state.mu.Unlock()

			ctx.GetStatusChange(states, 0)
			for i := range states {
				states[i].CurrentState = states[i].EventState
			}
			ctx.GetStatusChange(states, 3*time.Second)

		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func poller(ctx *scard.Context) {
	loaded := false
	breakCardLoop := false
	selectedReader := ""

	for {
		breakCardLoop = false

		state.mu.Lock()
		selectedReader = state.toolbar.GetReaderName()
		state.mu.Unlock()

		for !breakCardLoop {
			sCard, err := ctx.Connect(selectedReader, scard.ShareShared, scard.ProtocolAny)
			if err == nil {
				if !loaded {
					loaded = tryToProcessCard(sCard)
				}
			} else {
				state.mu.Lock()
				state.cardDocument = nil
				state.toolbar.DisablePinChange()
				state.mu.Unlock()

				loaded = false
				setStartPage(
					t("error.readingCard"),
					t("error.isCardPresent"),
					fmt.Errorf("connecting reader %s: %w", selectedReader, err))
			}

			state.mu.Lock()
			breakCardLoop = state.toolbar.ReaderChanged()
			state.mu.Unlock()

			if !breakCardLoop {
				time.Sleep(500 * time.Millisecond)
			}
		}

		setStartPage(
			t("poller.connectingReader"),
			"",
			nil)
	}
}

func tryToProcessCard(sCard *scard.Card) bool {
	loaded := false

	setStartPage(t("poller.readingFromCard"), "", nil)

	cardDoc, err := card.DetectCardDocument(sCard)
	if err != nil {
		message := ""
		if err == card.ErrUnknownCard {
			message = t("error.unknownCard") + ": " + cardDoc.Atr().String()
		}
		setStartPage(
			t("error.readingCard"),
			message,
			fmt.Errorf("reading from card: %w", err))
	} else {
		state.mu.Lock()
		state.cardDocument = cardDoc
		state.mu.Unlock()

		doc, err := initCardAndReadDoc(cardDoc)
		if err != nil {
			setStartPage(
				t("error.readingCard"),
				"",
				fmt.Errorf("reading from card: %w", err))
		} else {
			setStatus(t("poller.documentRead"), nil)
			setUI(doc)
			loaded = true
		}

		switch cardDoc.(type) {
		case *card.Gemalto:
			state.toolbar.EnablePinChange()
		}
	}

	return loaded
}

func initCardAndReadDoc(cardDoc card.CardDocument) (document.Document, error) {
	err := cardDoc.InitCard()
	if err != nil {
		return nil, err
	}

	err = cardDoc.ReadCard()
	if err != nil {
		return nil, err
	}

	doc, err := cardDoc.GetDocument()
	if err != nil {
		return nil, err
	}

	return doc, nil
}
