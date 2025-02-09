package icon

import (
	"embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var CertificateThemedResource *theme.ThemedResource
var PinThemedResource *theme.ThemedResource

func LoadIcons(embedFS embed.FS) error {
	certificateSvg, err := embedFS.ReadFile("embed/icons/certificate.svg")
	if err != nil {
		return fmt.Errorf("reading certificate icon: %w", err)
	}

	CertificateThemedResource = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "certificate.svg",
		StaticContent: certificateSvg,
	})

	pinSvg, err := embedFS.ReadFile("embed/icons/pin.svg")
	if err != nil {
		return fmt.Errorf("reading pin icon: %w", err)
	}

	PinThemedResource = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "pin.svg",
		StaticContent: pinSvg,
	})

	return nil
}
