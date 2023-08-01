package document

import (
	"embed"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
)

type Document interface {
	BuildPdf() ([]byte, string, error)
	BuildUI(func()) *fyne.Container
}

var font []byte

func SetFont(fontFS embed.FS) error {
	fontFile, err := fontFS.ReadFile("assets/free-sans-regular.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	font = fontFile
	return nil
}

func FormatDate(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	*in = chars[0] + chars[1] + "." + chars[2] + chars[3] + "." + chars[4] + chars[5] + chars[6] + chars[7]
}
