package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/ui"

	"embed"
	"fmt"
)

//go:embed assets/free-sans-regular.ttf
var font embed.FS

//go:embed assets/defaultPhoto.png
var defaultPhoto embed.FS

func main() {
	doc, err := document.NewDocument(font, defaultPhoto)
	if err != nil {
		fmt.Println("Error creating document:", err)
		return
	}

	ctx, err := scard.EstablishContext()
	if err != nil {
		fmt.Println("Error EstablishContext:", err)
		return
	}

	defer ctx.Release()

	app := app.New()
	window := app.NewWindow("Baš Čelik")
	app.Settings().SetTheme(ui.MyTheme{})
	ui := ui.BuildUI(&window, doc)

	go statusPoler(card.Readers{Context: ctx}, doc)

	window.SetContent(container.New(layout.NewPaddedLayout(), ui))
	window.ShowAndRun()
}

func statusPoler(readers card.Readers, doc *document.Document) {
	doc.Loaded = false
	ui.SetStatus("Konekcija sa čitačem...", false)

	readersNames, err := readers.List()
	if err != nil {
		ui.SetStatus("Greška pri konekciji sa čitačem", true)
	} else {
		if len(readersNames) > 0 {
			for {
				scard, err := readers.ConnectReader(readersNames[0])
				if err != nil {
					fmt.Println(err)
					if strings.Contains(err.Error(), "No smart card inserted") {
						ui.SetStatus("Nije detektovana kartica u čitaču", true)
					} else {
						ui.SetStatus("Greška pri očitavanju kartice", true)
					}
					doc.Loaded = false
					doc.Clear()
				} else if !doc.Loaded {
					err := card.ReadCard(scard, doc)
					if err != nil {
						doc.Loaded = false
						ui.SetStatus("Greška pri očitavanju kartice", true)
					}

					doc.Loaded = true
					ui.SetStatus("LK uspešno pročitana", false)
				}
				ui.UpdateUI()
				time.Sleep(500 * time.Millisecond)
			}
		} else {
			ui.SetStatus("Greška pri konekciji sa čitačem", true)
		}
	}

}
