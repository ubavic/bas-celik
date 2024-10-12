package card

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ubavic/bas-celik/card/tlv"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/localization"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Represents a smart card that holds a Serbian medical insurance document.
type MedicalCard struct {
	atr                  Atr
	smartCard            Card
	medicalDocumentFile  []byte
	fixedPersonalFile    []byte
	variablePersonalFile []byte
	variableAdminFile    []byte
}

// Possibly the first version of the medical card. Newer version has the GEMALTO_ATR_2 for the ATR.
var MEDICAL_ATR_1 = Atr([]byte{
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

func (card MedicalCard) ReadCard() error {
	var err error

	card.medicalDocumentFile, err = card.readFile(MED_DOCUMENT_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading document file: %w", err)
	}

	card.fixedPersonalFile, err = card.readFile(MED_FIXED_PERSONAL_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading fixed personal file: %w", err)
	}

	card.variablePersonalFile, err = card.readFile(MED_VARIABLE_PERSONAL_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading variable personal file: %w", err)
	}

	card.variableAdminFile, err = card.readFile(MED_VARIABLE_ADMIN_FILE_LOC)
	if err != nil {
		return fmt.Errorf("reading variable administrative file: %w", err)
	}

	return nil
}

func (card MedicalCard) GetDocument() (document.Document, error) {
	doc := document.MedicalDocument{}

	err := parseMedicalDocumentFile(card.medicalDocumentFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing document file: %w", err)
	}

	err = parseMedicalFixedPersonalFile(card.fixedPersonalFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing fixed personal file: %w", err)
	}

	err = parseMedicalVariablePersonalFile(card.variablePersonalFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing variable personal file: %w", err)
	}

	err = parseMedicalVariableAdminFile(card.variableAdminFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing variable administrative file: %w", err)
	}

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

	fields, err := tlv.ParseTLV(rsp)
	if err != nil {
		return false
	}
	descramble(fields, 1553)

	return strings.Compare(string(fields[1553]), "Републички фонд за здравствено осигурање") == 0
}

func (card MedicalCard) Atr() Atr {
	return card.atr
}

func (card MedicalCard) InitCard() error {
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

func parseMedicalDocumentFile(data []byte, doc *document.MedicalDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}

	descramble(fields, 1553)
	tlv.AssignField(fields, 1553, &doc.InsurerName)
	tlv.AssignField(fields, 1554, &doc.InsurerID)
	tlv.AssignField(fields, 1555, &doc.CardId)
	tlv.AssignField(fields, 1557, &doc.DateOfIssue)
	localization.FormatDate(&doc.DateOfIssue)
	tlv.AssignField(fields, 1558, &doc.DateOfExpiry)
	localization.FormatDate(&doc.DateOfExpiry)
	tlv.AssignField(fields, 1560, &doc.PrintLanguage)

	return nil
}

func parseMedicalFixedPersonalFile(data []byte, doc *document.MedicalDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}
	descramble(fields, 1570)
	tlv.AssignField(fields, 1570, &doc.FamilyName)
	descramble(fields, 1571)
	tlv.AssignField(fields, 1571, &doc.FamilyNameLatin)
	descramble(fields, 1572)
	tlv.AssignField(fields, 1572, &doc.GivenName)
	descramble(fields, 1573)
	tlv.AssignField(fields, 1573, &doc.GivenNameLatin)
	tlv.AssignField(fields, 1574, &doc.DateOfBirth)
	localization.FormatDate(&doc.DateOfBirth)
	tlv.AssignField(fields, 1569, &doc.InsurantNumber)

	return nil
}

func parseMedicalVariablePersonalFile(data []byte, doc *document.MedicalDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}
	tlv.AssignField(fields, 1586, &doc.ValidUntil)
	localization.FormatDate(&doc.ValidUntil)
	tlv.AssignBoolField(fields, 1587, &doc.PermanentlyValid)

	return nil
}

func parseMedicalVariableAdminFile(data []byte, doc *document.MedicalDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}
	descramble(fields, 1601)
	tlv.AssignField(fields, 1601, &doc.ParentName)
	descramble(fields, 1602)
	tlv.AssignField(fields, 1602, &doc.ParentNameLatin)
	if string(fields[1603]) == "01" {
		doc.Gender = "Mушко"
	} else {
		doc.Gender = "Женско"
	}
	tlv.AssignField(fields, 1604, &doc.PersonalNumber)
	descramble(fields, 1605)
	tlv.AssignField(fields, 1605, &doc.Street)
	descramble(fields, 1607)
	tlv.AssignField(fields, 1607, &doc.Municipality)
	descramble(fields, 1608)
	tlv.AssignField(fields, 1608, &doc.Place)
	descramble(fields, 1610)
	tlv.AssignField(fields, 1610, &doc.Number)
	descramble(fields, 1612)
	tlv.AssignField(fields, 1612, &doc.Apartment)
	tlv.AssignField(fields, 1614, &doc.InsuranceBasisRZZO)
	descramble(fields, 1615)
	tlv.AssignField(fields, 1615, &doc.InsuranceDescription)
	descramble(fields, 1616)
	tlv.AssignField(fields, 1616, &doc.CarrierRelationship)
	tlv.AssignBoolField(fields, 1617, &doc.CarrierFamilyMember)
	tlv.AssignField(fields, 1618, &doc.CarrierIdNumber)
	tlv.AssignField(fields, 1619, &doc.CarrierInsurantNumber)
	descramble(fields, 1620)
	tlv.AssignField(fields, 1620, &doc.CarrierFamilyName)
	descramble(fields, 1621)
	tlv.AssignField(fields, 1621, &doc.CarrierFamilyNameLatin)
	descramble(fields, 1622)
	tlv.AssignField(fields, 1622, &doc.CarrierGivenName)
	descramble(fields, 1623)
	tlv.AssignField(fields, 1623, &doc.CarrierGivenNameLatin)
	tlv.AssignField(fields, 1624, &doc.InsuranceStartDate)
	localization.FormatDate(&doc.InsuranceStartDate)
	descramble(fields, 1626)
	tlv.AssignField(fields, 1626, &doc.Country)
	descramble(fields, 1630)
	tlv.AssignField(fields, 1630, &doc.TaxpayerName)
	descramble(fields, 1631)
	tlv.AssignField(fields, 1631, &doc.TaxpayerResidence)
	tlv.AssignField(fields, 1632, &doc.TaxpayerIdNumber)
	if len(doc.TaxpayerIdNumber) == 0 {
		tlv.AssignField(fields, 1633, &doc.TaxpayerIdNumber)
	}
	tlv.AssignField(fields, 1634, &doc.TaxpayerActivityCode)

	return nil
}
