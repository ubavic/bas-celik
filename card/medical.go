package card

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/localization"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Represents a smart card that holds a Serbian medical insurance document.
type MedicalCard struct {
	atr       Atr
	smartCard Card
}

// Possibly the first version of the medical card. Newer version has the GEMALTO_ATR_2 for the ATR.
var MEDICAL_ATR = Atr([]byte{
	0x3B, 0xF4, 0x13, 0x00, 0x00, 0x81, 0x31, 0xFE,
	0x45, 0x52, 0x46, 0x5A, 0x4F, 0xED,
})

// Available since March 2023?
var MEDICAL_ATR_2 = []byte{
	0x3B, 0x9E, 0x97, 0x80, 0x31, 0xFE, 0x45, 0x53,
	0x43, 0x45, 0x20, 0x38, 0x2E, 0x30, 0x2D, 0x43,
	0x31, 0x56, 0x30, 0x0D, 0x0A, 0x6E,
}

// Location of the file with document data.
var MED_DOCUMENT_FILE_LOC = []byte{0x0D, 0x01}

// Location of the file with fixed personal data.
var MED_FIXED_PERSONAL_FILE_LOC = []byte{0x0D, 0x02}

// Location of the file with variable personal data.
var MED_VARIABLE_PERSONAL_FILE_LOC = []byte{0x0D, 0x03}

// Location of the file with variable administrative data.
var MED_VARIABLE_ADMIN_FILE_LOC = []byte{0x0D, 0x04}

func readMedicalCard(card MedicalCard) (*document.MedicalDocument, error) {
	doc := document.MedicalDocument{}

	rsp, err := card.readFile(MED_DOCUMENT_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err := parseTLV(rsp)
	if err != nil {
		return nil, err
	}

	descramble(fields, 1553)
	assignField(fields, 1553, &doc.InsurerName)
	assignField(fields, 1554, &doc.InsurerID)
	assignField(fields, 1555, &doc.CardId)
	assignField(fields, 1557, &doc.DateOfIssue)
	localization.FormatDate(&doc.DateOfIssue)
	assignField(fields, 1558, &doc.DateOfExpiry)
	localization.FormatDate(&doc.DateOfExpiry)
	assignField(fields, 1560, &doc.PrintLanguage)

	rsp, err = card.readFile(MED_FIXED_PERSONAL_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err = parseTLV(rsp)
	if err != nil {
		return nil, err
	}
	descramble(fields, 1570)
	assignField(fields, 1570, &doc.FamilyName)
	descramble(fields, 1571)
	assignField(fields, 1571, &doc.FamilyNameLatin)
	descramble(fields, 1572)
	assignField(fields, 1572, &doc.GivenName)
	descramble(fields, 1573)
	assignField(fields, 1573, &doc.GivenNameLatin)
	assignField(fields, 1574, &doc.DateOfBirth)
	localization.FormatDate(&doc.DateOfBirth)
	assignField(fields, 1569, &doc.InsurantNumber)

	rsp, err = card.readFile(MED_VARIABLE_PERSONAL_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err = parseTLV(rsp)
	if err != nil {
		return nil, err
	}
	assignField(fields, 1586, &doc.ValidUntil)
	localization.FormatDate(&doc.ValidUntil)
	assignBoolField(fields, 1587, &doc.PermanentlyValid)

	rsp, err = card.readFile(MED_VARIABLE_ADMIN_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err = parseTLV(rsp)
	if err != nil {
		return nil, err
	}
	descramble(fields, 1601)
	assignField(fields, 1601, &doc.ParentName)
	descramble(fields, 1602)
	assignField(fields, 1602, &doc.ParentNameLatin)
	if string(fields[1603]) == "01" {
		doc.Gender = "Mушко"
	} else {
		doc.Gender = "Женско"
	}
	assignField(fields, 1604, &doc.PersonalNumber)
	descramble(fields, 1605)
	assignField(fields, 1605, &doc.Street)
	descramble(fields, 1607)
	assignField(fields, 1607, &doc.Municipality)
	descramble(fields, 1608)
	assignField(fields, 1608, &doc.Place)
	descramble(fields, 1610)
	assignField(fields, 1610, &doc.Number)
	descramble(fields, 1612)
	assignField(fields, 1612, &doc.Apartment)
	assignField(fields, 1614, &doc.InsuranceBasisRZZO)
	descramble(fields, 1615)
	assignField(fields, 1615, &doc.InsuranceDescription)
	descramble(fields, 1616)
	assignField(fields, 1616, &doc.CarrierRelationship)
	assignBoolField(fields, 1617, &doc.CarrierFamilyMember)
	assignField(fields, 1618, &doc.CarrierIdNumber)
	assignField(fields, 1619, &doc.CarrierInsurantNumber)
	descramble(fields, 1620)
	assignField(fields, 1620, &doc.CarrierFamilyName)
	descramble(fields, 1621)
	assignField(fields, 1621, &doc.CarrierFamilyNameLatin)
	descramble(fields, 1622)
	assignField(fields, 1622, &doc.CarrierGivenName)
	descramble(fields, 1623)
	assignField(fields, 1623, &doc.CarrierGivenNameLatin)
	assignField(fields, 1624, &doc.InsuranceStartDate)
	localization.FormatDate(&doc.InsuranceStartDate)
	descramble(fields, 1626)
	assignField(fields, 1626, &doc.Country)
	descramble(fields, 1630)
	assignField(fields, 1630, &doc.TaxpayerName)
	descramble(fields, 1631)
	assignField(fields, 1631, &doc.TaxpayerResidence)
	assignField(fields, 1632, &doc.TaxpayerIdNumber)
	if len(doc.TaxpayerIdNumber) == 0 {
		assignField(fields, 1633, &doc.TaxpayerIdNumber)
	}
	assignField(fields, 1634, &doc.TaxpayerActivityCode)

	return &doc, nil
}

// Decodes UTF16 encoded data on medical cards.
func descramble(fields map[uint][]byte, tag uint) {
	bs, ok := fields[tag]
	if ok {
		utf8, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder(), bs)
		if err == nil {
			fields[tag] = utf8
			return
		}
	}

	fields[tag] = []byte{}
}

func (card MedicalCard) readFile(name []byte) ([]byte, error) {
	output := make([]byte, 0)

	_, err := card.selectFile(name)
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

func (card MedicalCard) selectFile(name []byte) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x00, 0x00, name, 0)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}

// Newer medical cards share ATR with the ID cards (GEMALTO_ATR_2)
// Function testMedicalCard tests if s smart card is a medical card.
func (card MedicalCard) testMedicalCard() bool {
	s1 := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x56, 0x53, 0x5A, 0x4B, 0x01}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, s1, 0)
	_, err := card.smartCard.Transmit(apu)
	if err != nil {
		return false
	}

	rsp, err := card.readFile([]byte{0x0D, 0x01})
	if err != nil {
		return false
	}

	fields, err := parseTLV(rsp)
	if err != nil {
		return false
	}
	descramble(fields, 1553)

	return strings.Compare(string(fields[1553]), "Републички фонд за здравствено осигурање") == 0
}

func (card MedicalCard) Atr() Atr {
	return card.atr
}

func (card MedicalCard) initCard() error {
	s1 := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x56, 0x53, 0x5A, 0x4B, 0x01}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, s1, 0)

	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return err
	}

	if !responseOK(rsp) {
		return fmt.Errorf("initializing card: response not OK")
	}

	return nil
}
