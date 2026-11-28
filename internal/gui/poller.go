package gui

import (
	"context"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/telesma-app/pcsc"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/document"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func releaseReaderCard() {
	state.mu.Lock()
	sCard := state.readerCard
	state.readerCard = nil
	state.cardDocument = nil
	state.mu.Unlock()

	if sCard != nil {
		_ = sCard.Close()
	}
}

func connectToCard(selectedReader string, present bool) {
	releaseReaderCard()
	fyne.Do(func() {
		state.mu.Lock()
		state.cryptoUi = nil
		state.certs = nil
		state.selectedCert = -1
		state.cryptoUiContainer.Hide()
		state.cryptoUiContainer.RemoveAll()
		state.mu.Unlock()

		switch {
		case selectedReader == "":
			setStartPage("error.noReader", t("error.noReaderExplanation"), nil)
		case !present:
			setStartPage("error.readingCard", t("error.isCardPresent"), nil)
		default:
			setStartPage("poller.connectingReader", "", nil)
		}
	})
	if !present {
		return
	}

	sCard, err := pcsc.Open(selectedReader)
	if err == nil {
		err = sCard.BeginTransaction(context.Background())
	}
	if err != nil {
		if sCard != nil {
			_ = sCard.Close()
		}
		fyne.Do(func() {
			setStartPage("error.readingCard", t("error.isCardPresent"),
				fmt.Errorf("connecting reader %s: %w", selectedReader, err))
		})
		return
	}

	fyne.Do(func() {
		setStartPage("poller.readingFromCard", "", nil)
	})

	cardDoc, err := card.DetectCardDocument(sCard)
	if cardDoc != nil {
		logger.Info("ATR read: " + cardDoc.Atr().String())
	}
	var doc document.Document
	if err == nil {
		doc, err = initCardAndReadDoc(cardDoc)
	}

	endErr := sCard.EndTransaction(pcsc.DispositionLeaveCard)
	if err == nil {
		err = endErr
	}
	if err != nil {
		_ = sCard.Close()
		fyne.Do(func() {
			readErr := fmt.Errorf("reading from card: %w", err)
			if errors.Is(err, card.ErrUnknownCard) && cardDoc != nil {
				setUnknownCardPage(cardDoc.Atr().String(), readErr)
			} else {
				setStartPage("error.readingCard", "", readErr)
			}
		})
		return
	}

	state.mu.Lock()
	if state.readerStopped {
		state.mu.Unlock()
		_ = sCard.Close()
		return
	}
	state.readerCard = sCard
	state.mu.Unlock()

	fyne.Do(func() {
		state.mu.Lock()
		if state.readerCard != sCard {
			state.mu.Unlock()
			return
		}
		state.cardDocument = cardDoc
		state.mu.Unlock()

		setStatus("poller.documentRead", nil)
		setUI(doc)
		go autoSave(doc)
		showWindowFromTray()
	})
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
