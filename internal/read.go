package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
)

func readAndSave(pdfPath, jsonPath string, reader uint, getMedicalExpiryDateFromRfzo bool) error {
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
		return fmt.Errorf("connecting reader %s: %w", readersNames[reader], err)
	}

	defer sCard.Disconnect(scard.LeaveCard)

	doc, err := card.ReadCard(sCard)
	if err != nil {
		return fmt.Errorf("reading card: %w", err)
	}

	switch doc := doc.(type) {
	case *document.MedicalDocument:
		if getMedicalExpiryDateFromRfzo {
			err := doc.UpdateValidUntilDateFromRfzo()
			if err != nil {
				return fmt.Errorf("updating `ValidUntil` date: %w", err)
			}
		}
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
			return fmt.Errorf("writing file %s: %w", jsonPath, err)
		}
	}

	return nil
}
