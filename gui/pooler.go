package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
)

func establishContextAndStartPooler() {
	setStartPage("Konekcija sa čitačem...", "", nil)

	ctx, err := scard.EstablishContext()
	if err != nil {
		setStartPage(
			"Greška pri upotrebi drajvera za pametne kartice.",
			"Da li program ima neophodne dozvole? Restartujte aplikaciju.",
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
				"Greška pri pretrazi dostupnih čitača",
				"Da li je čitač povezan za računar?",
				fmt.Errorf("listing readers: %w", err))

			time.Sleep(1000 * time.Millisecond)
			continue

		} else if len(readersNames) == 0 {
			loaded = false
			setStartPage(
				"Nijedan čitač nije detektovan",
				"Da li je čitač povezan za računar?",
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
					setStartPage("Čitam sa kartice...", "", nil)
					doc, err := card.ReadCard(sCard)
					if err != nil {
						setStartPage(
							"Greška pri čitanju kartice",
							"",
							fmt.Errorf("reading from card: %w", err))
					} else {
						setStatus("Dokument uspešno pročitan", nil)
						setUI(doc)
						loaded = true
					}
				}
				_ = sCard.Disconnect(scard.LeaveCard)
			} else {
				loaded = false
				setStartPage(
					"Greška pri čitanju kartice",
					"Da li je kartica prisutna?",
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
			"Povezivanje se čitačem u toku...",
			"",
			nil)
	}
}

func indexOf(element string, elements []string) int {
	for k, v := range elements {
		if element == v {
			return k
		}
	}
	return -1
}
