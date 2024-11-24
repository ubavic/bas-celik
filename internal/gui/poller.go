package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
)

func cardLoop(readerSelection <-chan string) {
	selectedReader := ""

	ctx, err := scard.EstablishContext()
	if err != nil {
		setStartPage(
			t("error.driver"),
			t("error.driverExplanation"),
			fmt.Errorf("establishing context: %w", err))
		select {}
	}

	for {
		setStartPage(t("error.noReader"), t("error.noReaderExplanation"), nil)

		for selectedReader == "" {
			selectedReader = <-readerSelection
		}

		setStartPage(t("poller.connectingReader"), "", nil)

		for selectedReader != "" {
			sCard, err := ctx.Connect(selectedReader, scard.ShareShared, scard.ProtocolAny)
			if err == nil {
				time.Sleep(50 * time.Microsecond)
				tryToProcessCard(sCard)
			} else {
				state.mu.Lock()
				state.cardDocument = nil
				state.toolbar.DisablePinChange()
				state.mu.Unlock()

				setStartPage(
					t("error.readingCard"),
					t("error.isCardPresent"),
					fmt.Errorf("connecting reader %s: %w", selectedReader, err))
			}

			selectedReader = <-readerSelection
			state.mu.Lock()
			state.cardDocument = nil
			state.toolbar.DisablePinChange()
			state.mu.Unlock()

			if selectedReader == "" {
				break
			}
		}
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
			state.mu.Lock()
			state.toolbar.EnablePinChange()
			state.mu.Unlock()
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
