package internal

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	"github.com/telesma-app/pcsc"
)

var version string

func ProcessFlags() (LaunchConfig, bool) {
	launchCfg := LaunchConfig{}

	atrFlag := flag.Bool("atr", false, "Print the ATR form the card and exit")
	cyrillicLabelsFlag := flag.Bool("cyrillic-labels", false, "Print ID PDF with cyrillic labels. Ignored for other cards")
	excelPath := flag.String("excel", "", "Set Excel export path")
	jsonPath := flag.String("json", "", "Set JSON export path")
	listFlag := flag.Bool("list", false, "List connected readers and exit")
	pdfPath := flag.String("pdf", "", "Set PDF export path.")
	getValidUntilFromRfzo := flag.Bool("rfzoValidUntil", false, "Get the valid until date of medical card insurance from the RFZO API. Ignored for other cards")
	verboseFlag := flag.Bool("verbose", false, "Provide additional details in the terminal")
	versionFlag := flag.Bool("version", false, "Display version information and exit")
	readerIndex := flag.Uint("reader", 0, "Set reader")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return launchCfg, true
	}

	if *listFlag {
		err := listReaders()
		if err != nil {
			fmt.Println("Error listing readers:", err)
		}
		return launchCfg, true
	}

	if *atrFlag {
		err := printATR(*readerIndex)
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return launchCfg, true
	}

	launchCfg.CyrillicLabels = *cyrillicLabelsFlag
	launchCfg.JsonPath = *jsonPath
	launchCfg.PdfPath = *pdfPath
	launchCfg.ExcelPath = *excelPath
	launchCfg.Verbose = *verboseFlag
	launchCfg.Reader = *readerIndex
	launchCfg.GetValidUntilFromRfzo = *getValidUntilFromRfzo

	return launchCfg, false
}

func printATR(reader uint) error {
	sCard, err := openReader(reader)
	if err != nil {
		return err
	}
	defer sCard.Close()

	smartCardStatus, err := sCard.Status()
	if err != nil {
		return fmt.Errorf("reading card status: %w", err)
	}

	fmt.Println(hex.EncodeToString(smartCardStatus.ATR))
	return nil
}

func openReader(reader uint) (*pcsc.Card, error) {
	readersNames, err := readerNames()
	if err != nil {
		return nil, err
	}
	if len(readersNames) == 0 {
		return nil, fmt.Errorf("no reader found")
	}
	if reader >= uint(len(readersNames)) {
		return nil, fmt.Errorf("only %d readers found", len(readersNames))
	}

	sCard, err := pcsc.Open(readersNames[reader])
	if err != nil {
		return nil, fmt.Errorf("connecting reader %s: %w", readersNames[reader], err)
	}
	return sCard, nil
}

func readerNames() ([]string, error) {
	var names []string
	for reader, err := range pcsc.Enumerate() {
		if err != nil {
			return nil, fmt.Errorf("listing readers: %w", err)
		}
		names = append(names, reader.Name)
	}
	return names, nil
}

func listReaders() error {
	readersNames, err := readerNames()
	if err != nil {
		return err
	}

	if len(readersNames) == 0 {
		fmt.Println("No readers found.")
		return nil
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

func SetVersion(v string) {
	version = strings.TrimSpace(v)
}
