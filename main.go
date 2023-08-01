package main

import (
	"errors"
	"flag"
	"os"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/gui"

	"embed"
	"fmt"
)

//go:embed assets/free-sans-regular.ttf
var font embed.FS

//go:embed assets/defaultPhoto.png
var defaultPhoto embed.FS

func main() {
	verboseFlag := flag.Bool("verbose", true, "Provide additional details in the terminal. Useful for debugging GUI")
	pdfPath := flag.String("pdf", "", "Set PDF export path. This command suppresses GUI")
	flag.Parse()

	ctx, err := scard.EstablishContext()
	if err != nil {
		fmt.Printf("Error establishing context: %s", err)
		return
	}

	defer ctx.Release()

	if len(*pdfPath) == 0 {
		gui.StartGui(ctx, *verboseFlag)
	} else {
		err := readAndSave(ctx, *pdfPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func readAndSave(ctx *scard.Context, pdfPath string) error {
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

	doc, err := card.ReadCard(sCard)
	if err != nil {
		return fmt.Errorf("reading card: %w", err)
	}

	pdf, _, err := doc.BuildPdf()
	if err != nil {
		return fmt.Errorf("generating pdf: %w", err)
	}

	err = os.WriteFile(pdfPath, pdf, 0600)
	if err != nil {
		return fmt.Errorf("writing file %s: %w", pdfPath, err)
	}

	return nil
}
