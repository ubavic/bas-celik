package gui

import (
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
)

func pooler(ctx *scard.Context) {
	loaded := false

	setStatus("Konekcija sa čitačem...", nil)

	readersNames, err := ctx.ListReaders()
	if err != nil {
		loaded = false
		setStatus(
			"Greška pri konekciji sa čitačem.",
			fmt.Errorf("listing readers: %w", err))
		return
	}

	if len(readersNames) == 0 {
		setStatus(
			"Nijedan čitač nije detektovan.",
			fmt.Errorf("no reader found"))
		return
	}

	for {
		sCard, err := ctx.Connect(readersNames[0], scard.ShareShared, scard.ProtocolAny)
		if err == nil {
			if !loaded {
				setStatus("Čitam sa kartice...", nil)
				doc, err := card.ReadCard(sCard)
				if err != nil {
					setStatus(
						"Greška pri očitavanju kartice",
						fmt.Errorf("reading from card: %w", err))
				} else {
					setStatus("LK uspešno pročitana", nil)
					setUI(doc)
					loaded = true
				}
			}
			sCard.Disconnect(scard.LeaveCard)
		} else {
			loaded = false
			setStatus(
				"Greška pri čitanju kartice. Da li je kartica prisutna?",
				fmt.Errorf("connecting reader %s: %w", readersNames[0], err))
		}

		time.Sleep(500 * time.Millisecond)
	}
}
