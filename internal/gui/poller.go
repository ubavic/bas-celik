package gui

import (
	"fmt"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal/logger"
)

func connectToCard(selectedReader string, ctx *scard.Context) {
	readers, _ := ctx.ListReaders()
	if selectedReader == "" || len(readers) == 0 {
		setStartPage("error.noReader", "error.noReaderExplanation", nil)
		return
	}

	setStartPage("poller.connectingReader", "", nil)

	sCard, err := ctx.Connect(selectedReader, scard.ShareShared, scard.ProtocolAny)
	if err == nil {
		tryToProcessCard(sCard)
	} else {
		state.mu.Lock()
		state.cardDocument = nil
		state.toolbar.DisablePinChange()
		state.mu.Unlock()

		setStartPage(
			"error.readingCard",
			"error.isCardPresent",
			fmt.Errorf("connecting reader %s: %w", selectedReader, err))
	}
}

func tryToProcessCard(sCard *scard.Card) bool {
	loaded := false

	setStartPage("poller.readingFromCard", "", nil)

	cardDoc, err := card.DetectCardDocument(sCard)
	if cardDoc != nil {
		logger.Info("ATR read: " + cardDoc.Atr().String())
	}
	if err != nil {
		message := ""
		if err == card.ErrUnknownCard {
			message = "error.unknownCard"
		}
		setStartPage(
			"error.readingCard",
			message,
			fmt.Errorf("reading from card: %w", err))
	} else {
		state.mu.Lock()
		state.cardDocument = cardDoc
		state.mu.Unlock()

		doc, err := initCardAndReadDoc(cardDoc)
		if err != nil {
			setStartPage(
				"error.readingCard",
				"",
				fmt.Errorf("reading from card: %w", err))
		} else {
			setStatus("poller.documentRead", nil)
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
