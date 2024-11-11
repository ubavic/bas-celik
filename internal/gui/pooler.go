package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
)

func establishContextAndStartPooler() {
	setStartPage(t("pooler.connectingReader"), "", nil)

	ctx, err := scard.EstablishContext()
	if err != nil {
		setStartPage(
			t("error.driver"),
			t("error.driverExplanation"),
			fmt.Errorf("establishing context: %w", err))
	} else {
		pooler(ctx)
	}
}

func pooler(ctx *scard.Context) {
	loaded := false
	breakCardLoop := false
	selectedReader := ""
	selectedReaderIndex := 0
	var readersNames []string
	var err error

	for {
		breakCardLoop = false

		readersNames, err = ctx.ListReaders()
		if err != nil {
			loaded = false
			setStartPage(
				t("error.reader"),
				t("error.readerExplanation"),
				fmt.Errorf("listing readers: %w", err))

			time.Sleep(1000 * time.Millisecond)
			continue

		} else if len(readersNames) == 0 {
			loaded = false
			setStartPage(
				t("error.noReader"),
				t("error.noReaderExplanation"),
				fmt.Errorf("no reader found"))

			time.Sleep(1000 * time.Millisecond)
			continue
		}

		state.mu.Lock()
		selectedReader = state.toolbar.GetReaderName()

		selectedReaderIndex = indexOf(selectedReader, readersNames)
		if selectedReaderIndex < 0 {
			selectedReaderIndex = 0
		}

		state.toolbar.SetReaders(readersNames)
		state.mu.Unlock()

		for !breakCardLoop {
			sCard, err := ctx.Connect(readersNames[selectedReaderIndex], scard.ShareShared, scard.ProtocolAny)
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
					fmt.Errorf("connecting reader %s: %w", readersNames[selectedReaderIndex], err))
			}

			state.mu.Lock()
			breakCardLoop = state.toolbar.ReaderChanged()
			state.mu.Unlock()

			if !breakCardLoop {
				time.Sleep(500 * time.Millisecond)
			}
		}

		setStartPage(
			t("pooler.connectingReader"),
			"",
			nil)
	}
}

func tryToProcessCard(sCard *scard.Card) bool {
	loaded := false

	setStartPage(t("pooler.readingFromCard"), "", nil)

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
			setStatus(t("pooler.documentRead"), nil)
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

func indexOf(element string, elements []string) int {
	for k, v := range elements {
		if element == v {
			return k
		}
	}
	return -1
}
