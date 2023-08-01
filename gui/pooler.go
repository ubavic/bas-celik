package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
)

func pooler(ctx *scard.Context) {
	loaded := false

	setStartPage("Konekcija sa čitačem...", "", nil)

	readersNames, err := ctx.ListReaders()
	if err != nil {
		loaded = false
		setStartPage(
			"Greška pri povezivanju sa čitačem",
			"Da li je čitač povezan za računar? Ponovo pokrenite aplikaciju nakon povezivanja.",
			fmt.Errorf("listing readers: %w", err))
		return
	}

	if len(readersNames) == 0 {
		setStartPage(
			"Nijedan čitač nije detektovan",
			"Da li je čitač povezan za računar? Ponovo pokrenite aplikaciju nakon povezivanja.",
			fmt.Errorf("no reader found"))
		return
	}

	for {
		sCard, err := ctx.Connect(readersNames[0], scard.ShareShared, scard.ProtocolAny)
		if err == nil {
			if !loaded {
				setStartPage("Čitam sa kartice...", "", nil)
				doc, err := card.ReadCard(sCard)
				if err != nil {
					setStartPage(
						"Greška pri očitavanju kartice",
						"",
						fmt.Errorf("reading from card: %w", err))
				} else {
					setStatus("Dokument uspešno pročitan", nil)
					setUI(doc)
					loaded = true
				}
			}
			sCard.Disconnect(scard.LeaveCard)
		} else {
			loaded = false
			setStartPage(
				"Greška pri čitanju kartice",
				"Da li je kartica prisutna?",
				fmt.Errorf("connecting reader %s: %w", readersNames[0], err))
		}

		time.Sleep(500 * time.Millisecond)
	}
}
