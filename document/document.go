package document

import (
	"embed"
	"fmt"
	"strings"
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
func SetData(fontRegularFS, fontBoldFS, rfzoLogoFS embed.FS) error {
	fontFile, err := fontRegularFS.ReadFile("assets/liberationSansRegular.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	fontRegular = fontFile

	fontFile, err = fontBoldFS.ReadFile("assets/liberationSansBold.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	fontBold = fontFile

	rfzoLogoFile, err := rfzoLogoFS.ReadFile("assets/rfzo.png")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}
	rfzoLogo = rfzoLogoFile

	return nil
}

// Expects a pointer to a date in the format DDMMYYYY.
// Modifies, in place, date to format DD.MM.YYYY.
func FormatDate(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	if chars[4] == "0" {
		*in = "Nije dostupan"
		return
	}
	*in = chars[0] + chars[1] + "." + chars[2] + chars[3] + "." + chars[4] + chars[5] + chars[6] + chars[7] + "."
}

// Expects a pointer to a date in the format YYYYMMDD.
// Modifies, in place, date to format DD.MM.YYYY.
func FormatDate2(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	*in = chars[6] + chars[7] + "." + chars[4] + chars[5] + "." + chars[0] + chars[1] + chars[2] + chars[3]
}
