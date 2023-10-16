package document

import (
	"embed"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/ubavic/bas-celik/widgets"
)

type Document interface {
	BuildPdf() ([]byte, string, error)
	BuildUI(func(), *widgets.StatusBar) *fyne.Container
}

var font, rfzoLogo []byte

func SetData(fontFS, rfzoLogoFS embed.FS) error {
	fontFile, err := fontFS.ReadFile("assets/free-sans-regular.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	font = fontFile

	rfzoLogoFile, err := rfzoLogoFS.ReadFile("assets/rfzo.png")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	rfzoLogo = rfzoLogoFile

	return nil
}

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
