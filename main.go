package main

import (
	"errors"
	"flag"
	"os"
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
	verboseFlag := flag.Bool("verbose", false, "Provide additional details in the terminal. Useful for debugging GUI")
	pdfPath := flag.String("pdf", "", "Set PDF export path. This command suppresses GUI")
	flag.Parse()

	doc, err := document.NewDocument(font, defaultPhoto)
	if err != nil {
		fmt.Printf("Error creating document: %s", err)
		return
	}

	ctx, err := scard.EstablishContext()
	if err != nil {
		fmt.Printf("Error establishing context: %s", err)
		return
	}

	defer ctx.Release()

	if len(*pdfPath) == 0 {
		gui(ctx, doc, *verboseFlag)
	} else {
		err := readAndSave(ctx, doc, *pdfPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func readAndSave(ctx *scard.Context, doc *document.Document, pdfPath string) error {
	if _, err := os.Stat(pdfPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("opening file %s: %w", pdfPath, err)
	}

	readersNames, err := ctx.ListReaders()
	if err != nil {
		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		return fmt.Errorf("no reader found")
	}

	sCard, err := ctx.Connect(readersNames[0], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return fmt.Errorf("connecting reader %s: %w", readersNames[0], err)
	}

	defer sCard.Disconnect(scard.LeaveCard)

	err = card.ReadCard(sCard, doc)
	if err != nil {
		return fmt.Errorf("reading card: %w", err)
	}

	pdf, _, err := doc.Pdf()
	if err != nil {
		return fmt.Errorf("generating pdf: %w", err)
	}

	err = os.WriteFile(pdfPath, pdf, 0600)
	if err != nil {
		return fmt.Errorf("writing file %s: %w", pdfPath, err)
	}

	return nil
}

func gui(ctx *scard.Context, doc *document.Document, verbose bool) {
	app := app.New()
	window := app.NewWindow("Baš Čelik")
	app.Settings().SetTheme(ui.MyTheme{})
	ui := ui.SetUI(&window, doc, verbose)

	go statusPoler(ctx, doc)

	window.SetContent(container.New(layout.NewPaddedLayout(), ui))
	window.ShowAndRun()
}

func statusPoler(ctx *scard.Context, doc *document.Document) {
	ui.SetStatus("Konekcija sa čitačem...", nil)

	readersNames, err := ctx.ListReaders()
	if err != nil {
		ui.SetStatus(
			"Greška pri konekciji sa čitačem.",
			fmt.Errorf("listing readers: %w", err))
		return
	}

	if len(readersNames) == 0 {
		ui.SetStatus(
			"Nijedan čitač nije detektovan.",
			fmt.Errorf("no reader found"))
		return
	}

	for {
		sCard, err := ctx.Connect(readersNames[0], scard.ShareShared, scard.ProtocolAny)
		if err == nil {
			if !doc.Loaded {
				ui.SetStatus("Čitam sa kartice...", nil)
				err := card.ReadCard(sCard, doc)
				if err != nil {
					ui.SetStatus(
						"Greška pri očitavanju kartice",
						fmt.Errorf("reading from card: %w", err))
					doc.Clear()
				}

				ui.SetStatus("LK uspešno pročitana", nil)
				ui.UpdateUI()
			}
			sCard.Disconnect(scard.LeaveCard)

		} else {
			ui.SetStatus(
				"Greška pri čitanju kartice. Da li je kartica prisutna?",
				fmt.Errorf("connecting reader %s: %w", readersNames[0], err))
			doc.Clear()
			ui.UpdateUI()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
