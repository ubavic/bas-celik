package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
)

func readAndSave(pdfPath, jsonPath string, reader uint) error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

	if len(pdfPath) > 0 {
		if _, err := os.Stat(pdfPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("opening file %s: %w", pdfPath, err)
		}
	}

	if len(jsonPath) > 0 {
		if _, err := os.Stat(jsonPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("opening file %s: %w", jsonPath, err)
		}
	}

	readersNames, err := ctx.ListReaders()
	if err != nil {
		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		return fmt.Errorf("no reader found")
	}

	if reader >= uint(len(readersNames)) {
		return fmt.Errorf("only %d readers found", len(readersNames))
	}

	sCard, err := ctx.Connect(readersNames[reader], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return fmt.Errorf("connecting reader %s: %w", readersNames[0], err)
	}

	defer sCard.Disconnect(scard.LeaveCard)

	doc, err := card.ReadCard(sCard)
	if err != nil {
		return fmt.Errorf("reading card: %w", err)
	}

	if len(pdfPath) > 0 {
		pdf, _, err := doc.BuildPdf()
		if err != nil {
			return fmt.Errorf("generating pdf: %w", err)
		}

		err = os.WriteFile(pdfPath, pdf, 0600)
		if err != nil {
			return fmt.Errorf("writing file %s: %w", pdfPath, err)
		}
	}

	if len(jsonPath) > 0 {
		json, err := doc.BuildJson()
		if err != nil {
			return fmt.Errorf("generating json: %w", err)
		}

		err = os.WriteFile(jsonPath, json, 0600)
		if err != nil {
			fmt.Println(fmt.Errorf("writing file %s: %w", jsonPath, err))
		}
	}

	return nil
}

func printATR() error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

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

	smartCardStatus, err := sCard.Status()
	if err != nil {
		return fmt.Errorf("reading card %w", err)
	}

	fmt.Println(hex.EncodeToString(smartCardStatus.Atr))

	return nil
}

func listReaders() error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

	readersNames, err := ctx.ListReaders()
	if err != nil {
		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		return errors.New("no reader found")
	}

	for i, name := range readersNames {
		fmt.Println(i, "|", name)
	}

	return nil
}

func printVersion() {
	ver := strings.TrimSpace(version)
	fmt.Println("bas-celik", ver)
	fmt.Println("https://github.com/ubavic/bas-celik")
}
