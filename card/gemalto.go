package card

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ubavic/bas-celik/document"
)

var GEMALTO_ATR_1 = Atr([]byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xB0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
})

// Available since January 2023 (maybe). Replaced very soon with an even newer version.
var GEMALTO_ATR_2 = Atr([]byte{
	0x3B, 0xF9, 0x96, 0x00, 0x00, 0x80, 0x31, 0xFE,
	0x45, 0x53, 0x43, 0x45, 0x37, 0x20, 0x47, 0x43,
	0x4E, 0x33, 0x5E,
})

// Available since July 2023.
var GEMALTO_ATR_3 = Atr([]byte{
	0x3B, 0x9E, 0x96, 0x80, 0x31, 0xFE, 0x45, 0x53,
	0x43, 0x45, 0x20, 0x38, 0x2E, 0x30, 0x2D, 0x43,
	0x31, 0x56, 0x30, 0x0D, 0x0A, 0x6F,
})

// Available since June 2024.
var GEMALTO_ATR_4 = Atr([]byte{
	0x3B, 0x9E, 0x96, 0x80, 0x31, 0xFE, 0x45, 0x53,
	0x43, 0x45, 0x20, 0x38, 0x2E, 0x30, 0x2D, 0x43,
	0x32, 0x56, 0x30, 0x0D, 0x0A, 0x6C,
})

// Gemalto represents ID cards based with Gemalto Java OS. Gemalto replaced Apollo cards around 2014.
type Gemalto struct {
	atr           Atr
	smartCard     Card
	documentFile  []byte
	personalFile  []byte
	residenceFile []byte
	photoFile     []byte
}

func (card *Gemalto) InitCard() error {
	data := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x49, 0x44, 0x01}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("initializing ID card: %w", err)
	}

	if responseOK(rsp) {
		return nil
	}

	data = []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x49, 0x46, 0x01}
	apu = buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)
	rsp, err = card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("initializing IF card: %w", err)
	}

	if responseOK(rsp) {
		return nil
	}

	data = []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x52, 0x50, 0x01}
	apu = buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)
	rsp, err = card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("initializing RP card: %w", err)
	}

	if responseOK(rsp) {
		return nil
	}

	return fmt.Errorf("initializing identity document card: unknown card type")
}

func (card *Gemalto) ReadCard() error {
	var err error

	card.documentFile, err = card.ReadFile(ID_DOCUMENT_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading document file: %w", err)
	}

	card.personalFile, err = card.ReadFile(ID_PERSONAL_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading personal file: %w", err)
	}

	card.residenceFile, err = card.ReadFile(ID_RESIDENCE_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading residence file: %w", err)
	}

	rsp, err := card.ReadFile(ID_PHOTO_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading photo file: %w", err)
	}

	card.photoFile = trim4b(rsp)

	return nil
}

func (card *Gemalto) GetDocument() (document.Document, error) {
	doc := document.IdDocument{}

	err := parseIdDocumentFile(card.documentFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing document file: %w", err)
	}

	err = parseIdPersonalFile(card.personalFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing personal file: %w", err)
	}

	err = parseIdResidenceFile(card.residenceFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing residence file: %w", err)
	}

	err = parseAndAssignIdPhotoFile(card.photoFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing photo file: %w", err)
	}

	return &doc, nil
}

func (card *Gemalto) Atr() Atr {
	return card.atr
}

func (card *Gemalto) ReadFile(name []byte) ([]byte, error) {
	output := make([]byte, 0)

	_, err := card.selectFile(name, 4)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	data, err := read(card.smartCard, 0, 4)
	if err != nil {
		return nil, fmt.Errorf("reading file header: %w", err)
	}

	offset := uint(len(data))
	if offset < 3 {
		return nil, fmt.Errorf("file too short")
	}
	length := uint(binary.LittleEndian.Uint16(data[2:]))

	for length > 0 {
		data, err := read(card.smartCard, offset, length)
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	return output, nil
}

func (card *Gemalto) selectFile(name []byte, ne uint) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x08, 0x00, name, ne)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}

func (card *Gemalto) Test() bool {
	err := card.InitCard()
	if err != nil {
		return false
	}

	_, err = card.ReadFile(ID_DOCUMENT_FILE_LOC)
	return err == nil
}

// Initialize card's cryptography application
func (card *Gemalto) InitCrypto() error {
	data := []byte{0xA0, 0x00, 0x00, 0x00, 0x63, 0x50, 0x4B, 0x43, 0x53, 0x2D, 0x31, 0x35}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)

	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("initializing cryptography application %w", err)
	}

	if !responseOK(rsp) {
		return errors.New("cryptography application not selected")
	}

	return nil
}

func (card *Gemalto) ChangePin(newPin, oldPin string) error {
	err := card.InitCrypto()
	if err != nil {
		return err
	}

	oldPinValid := ValidatePin(oldPin)
	if !oldPinValid {
		return errors.New("old pin not valid")
	}

	newPinValid := ValidatePin(newPin)
	if !newPinValid {
		return errors.New("new pin not valid")
	}

	apu := buildAPDU(0x00, 0x20, 0x00, 0x80, PadPin(oldPin), 0)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("verifying old pin %w", err)
	}

	if !responseOK(rsp) {
		return errors.New("verifying old pin")
	}

	data := make([]byte, 0, 8)
	data = append(data, PadPin(oldPin)...)
	data = append(data, PadPin(newPin)...)

	apu = buildAPDU(0x00, 0x24, 0x00, 0x80, data, 0)
	rsp, err = card.smartCard.Transmit(apu)
	if err != nil {
		return fmt.Errorf("changing pin %w", err)
	}

	if !responseOK(rsp) {
		return errors.New("changing pin")
	}

	return nil
}
