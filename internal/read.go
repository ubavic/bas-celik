package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
)

type LaunchConfig struct {
	PdfPath               string
	JsonPath              string
	Verbose               bool
	GetValidUntilFromRfzo bool
	Reader                uint
}

func readAndSave(cfg LaunchConfig) error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

	if len(cfg.PdfPath) > 0 {
		if _, err := os.Stat(cfg.PdfPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("opening file %s: %w", cfg.PdfPath, err)
		}
	}

	if len(cfg.JsonPath) > 0 {
		if _, err := os.Stat(cfg.JsonPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("opening file %s: %w", cfg.JsonPath, err)
		}
	}

	readersNames, err := ctx.ListReaders()
	if err != nil {
		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		return fmt.Errorf("no reader found")
	}

	if cfg.Reader >= uint(len(readersNames)) {
		return fmt.Errorf("only %d readers found", len(readersNames))
	}

	sCard, err := ctx.Connect(readersNames[cfg.Reader], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return fmt.Errorf("connecting reader %s: %w", readersNames[cfg.Reader], err)
	}

	defer sCard.Disconnect(scard.LeaveCard)

	cardDoc, err := card.DetectCardDocument(sCard)
	if err != nil {
		return fmt.Errorf("detecting card type: %w", err)
	}

	doc, err := card.ReadCard(cardDoc)
	if err != nil {
		return fmt.Errorf("reading card: %w", err)
	}

	switch doc := doc.(type) {
	case *document.MedicalDocument:
		if cfg.GetValidUntilFromRfzo {
			err := doc.UpdateValidUntilDateFromRfzo()
			if err != nil {
				return fmt.Errorf("updating `ValidUntil` date: %w", err)
			}
		}
	}

	if len(cfg.PdfPath) > 0 {
		pdf, _, err := doc.BuildPdf()
		if err != nil {
			return fmt.Errorf("generating pdf: %w", err)
		}

		err = os.WriteFile(cfg.PdfPath, pdf, 0600)
		if err != nil {
			return fmt.Errorf("writing file %s: %w", cfg.PdfPath, err)
		}
	}

	if len(cfg.JsonPath) > 0 {
		json, err := doc.BuildJson()
		if err != nil {
			return fmt.Errorf("generating json: %w", err)
		}

		err = os.WriteFile(cfg.JsonPath, json, 0600)
		if err != nil {
			return fmt.Errorf("writing file %s: %w", cfg.JsonPath, err)
		}
	}

	return nil
}
