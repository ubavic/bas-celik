package internal

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/ebfe/scard"
)

var version string

func ProcessFlags(_jsonPath, _pdfPath *string, _verbose, _getMedicalExpiryDateFromRfzo *bool, _readerIndex *uint) bool {
	atrFlag := flag.Bool("atr", false, "Print the ATR form the card and exit")
	jsonPath := flag.String("json", "", "Set JSON export path")
	listFlag := flag.Bool("list", false, "List connected readers and exit")
	pdfPath := flag.String("pdf", "", "Set PDF export path.")
	getMedicalExpiryDateFromRfzo := flag.Bool("rfzoExpiryDate", false, "Get expiry date of medical card from RFZO API. Ignored for other cards")
	verboseFlag := flag.Bool("verbose", false, "Provide additional details in the terminal")
	versionFlag := flag.Bool("version", false, "Display version information and exit")
	readerIndex := flag.Uint("reader", 0, "Set reader")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return true
	}

	if *listFlag {
		err := listReaders()
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return true
	}

	if *atrFlag {
		err := printATR(*readerIndex)
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return true
	}

	*_jsonPath = *jsonPath
	*_pdfPath = *pdfPath
	*_verbose = *verboseFlag
	*_readerIndex = *readerIndex
	*_getMedicalExpiryDateFromRfzo = *getMedicalExpiryDateFromRfzo

	return false
}

func printATR(reader uint) error {
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

	if reader >= uint(len(readersNames)) {
		return fmt.Errorf("only %d readers found", len(readersNames))
	}

	sCard, err := ctx.Connect(readersNames[reader], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return fmt.Errorf("connecting reader %s: %w", readersNames[reader], err)
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

func SetVersion(v string) {
	version = v
}
