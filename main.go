package main

import (
	"flag"

	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/gui"

	"embed"
	"fmt"
)

//go:embed assets/version
var version string

//go:embed assets/liberationSansRegular.ttf
var fontRegular embed.FS

//go:embed assets/liberationSansBold.ttf
var fontBold embed.FS

//go:embed assets/rfzo.png
var rfzoLogo embed.FS

func main() {
	atrFlag := flag.Bool("atr", false, "Print the ATR form the card and exit. Useful for debugging")
	jsonPath := flag.String("json", "", "Set JSON export path. This command suppresses GUI")
	listFlag := flag.Bool("list", false, "List connected readers and exit")
	pdfPath := flag.String("pdf", "", "Set PDF export path. This command suppresses GUI")
	verboseFlag := flag.Bool("verbose", false, "Provide additional details in the terminal. Useful for debugging")
	versionFlag := flag.Bool("version", false, "Display version information and exit")
	readerIndex := flag.Uint("reader", 0, "Set reader")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return
	}

	if *listFlag {
		err := listReaders()
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return
	}

	if *atrFlag {
		err := printATR()
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return
	}

	err := document.SetData(fontRegular, fontBold, rfzoLogo)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}

	if len(*pdfPath) == 0 && len(*jsonPath) == 0 {
		gui.StartGui(*verboseFlag, version)
		gui.StartNativeMessaging()
	} else {
		err := readAndSave(*pdfPath, *jsonPath, *readerIndex)
		if err != nil {
			fmt.Println("Error saving document:", err)
		}
	}
}
