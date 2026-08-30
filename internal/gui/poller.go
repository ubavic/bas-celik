package gui

import (
	"fmt"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/document"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func connectToCard(selectedReader string, ctx *scard.Context) {
	state.mu.Lock()
	state.cardDocument = nil
	state.cryptoUi = nil
	state.certs = nil
	state.selectedCert = -1
	state.mu.Unlock()

	state.cryptoUiContainer.Hide()
	state.cryptoUiContainer.RemoveAll()

	readers, _ := ctx.ListReaders()
	if selectedReader == "" || len(readers) == 0 {
		setStartPage("error.noReader", t("error.noReaderExplanation"), nil)
		return
	}

	setStartPage("poller.connectingReader", "", nil)

	sCard, err := ctx.Connect(selectedReader, scard.ShareShared, scard.ProtocolAny)
	if err == nil {
		err = sCard.BeginTransaction()
		if err == nil {
			tryToProcessCard(sCard)
			sCard.EndTransaction(scard.LeaveCard)
			return
		}
	}

	setStartPage(
		"error.readingCard",
		t("error.isCardPresent"),
		fmt.Errorf("connecting reader %s: %w", selectedReader, err))
}

func tryToProcessCard(sCard *scard.Card) bool {
	loaded := false

	setStartPage("poller.readingFromCard", "", nil)

	cardDoc, err := card.DetectCardDocument(sCard)
	if cardDoc != nil {
		logger.Info("ATR read: " + cardDoc.Atr().String())
	}

	if err != nil {
		if err == card.ErrUnknownCard && cardDoc != nil {
			setUnknownCardPage(cardDoc.Atr().String(), fmt.Errorf("reading from card: %w", err))
		} else {
			setStartPage(
				"error.readingCard",
				"",
				fmt.Errorf("reading from card: %w", err))
		}
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

		go autoSave(doc)

		showWindowFromTray()
	}

	return loaded
}

func showWindowFromTray() {
	if state.runInBackground && state.autoSaveMode != AutoSaveSaveAndOpen {
		state.window.Show()
		state.window.RequestFocus()
	}
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
