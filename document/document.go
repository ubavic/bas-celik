package document

import (
	"embed"
	"fmt"
)

// Represents any document handled by Bas Celik
type Document interface {
	BuildPdf() ([]byte, string, error) // Renders document to pdf
	BuildJson() ([]byte, error)        // Renders document to json
}

var (
	fontRegular []byte // regular font used for PDF render
	fontBold    []byte // bold font used for PDF render
	rfzoLogo    []byte // logo used in PDF render of medical cards
)

// Sets fonts and graphics used for rendering PDF
func SetData(embedFS embed.FS) error {
	fontFile, err := embedFS.ReadFile("embed/liberationSansRegular.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	fontRegular = fontFile

	fontFile, err = embedFS.ReadFile("embed/liberationSansBold.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	fontBold = fontFile

	rfzoLogoFile, err := embedFS.ReadFile("embed/rfzo.png")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	rfzoLogo = rfzoLogoFile

	return nil
}
