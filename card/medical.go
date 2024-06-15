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
	smartCard Card
}

// Possibly the first version of the medical card. Newer version has the GEMALTO_ATR_2 for the ATR.
var MEDICAL_ATR = []byte{
	0x3B, 0xF4, 0x13, 0x00, 0x00, 0x81, 0x31, 0xFE,
	0x45, 0x52, 0x46, 0x5A, 0x4F, 0xED,
}

func readMedicalCard(card MedicalCard) (*document.MedicalDocument, error) {
	s1 := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x56, 0x53, 0x5A, 0x4B, 0x01}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, s1, 0)

	_, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, err
	}

	doc := document.MedicalDocument{}

	rsp, err := card.readFile([]byte{0x0D, 0x01}, false)
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
	assignField(fields, 1557, &doc.CardIssueDate)
	localization.FormatDate(&doc.CardIssueDate)
	assignField(fields, 1558, &doc.CardExpiryDate)
	localization.FormatDate(&doc.CardExpiryDate)
	assignField(fields, 1560, &doc.Language)

	rsp, err = card.readFile([]byte{0x0D, 0x02}, false)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err = parseTLV(rsp)
	if err != nil {
		return nil, err
	}
	descramble(fields, 1570)
	assignField(fields, 1570, &doc.SurnameCyrl)
	descramble(fields, 1571)
	assignField(fields, 1571, &doc.Surname)
	descramble(fields, 1572)
	assignField(fields, 1572, &doc.GivenNameCyrl)
	descramble(fields, 1573)
	assignField(fields, 1573, &doc.GivenName)
	assignField(fields, 1574, &doc.DateOfBirth)
	localization.FormatDate(&doc.DateOfBirth)
	assignField(fields, 1569, &doc.InsuranceNumber)

	rsp, err = card.readFile([]byte{0x0D, 0x03}, false)
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

	rsp, err = card.readFile([]byte{0x0D, 0x04}, false)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	fields, err = parseTLV(rsp)
	if err != nil {
		return nil, err
	}
	descramble(fields, 1601)
	assignField(fields, 1601, &doc.ParentNameCyrl)
	descramble(fields, 1602)
	assignField(fields, 1602, &doc.ParentName)
	if string(fields[1603]) == "01" {
		doc.Sex = "Mушко"
	} else {
		doc.Sex = "Женско"
	}
	assignField(fields, 1604, &doc.PersonalNumber)
	descramble(fields, 1605)
	assignField(fields, 1605, &doc.AddressStreet)
	descramble(fields, 1607)
	assignField(fields, 1607, &doc.AddressMunicipality)
	descramble(fields, 1608)
	assignField(fields, 1608, &doc.AddressTown)
	descramble(fields, 1610)
	assignField(fields, 1610, &doc.AddressNumber)
	descramble(fields, 1612)
	assignField(fields, 1612, &doc.AddressApartmentNumber)
	assignField(fields, 1614, &doc.InsuranceReason)
	descramble(fields, 1615)
	assignField(fields, 1615, &doc.InsuranceDescription)
	descramble(fields, 1616)
	assignField(fields, 1616, &doc.InsuranceHolderRelation)
	assignBoolField(fields, 1617, &doc.InsuranceHolderIsFamilyMember)
	assignField(fields, 1618, &doc.InsuranceHolderPersonalNumber)
	assignField(fields, 1619, &doc.InsuranceHolderInsuranceNumber)
	descramble(fields, 1620)
	assignField(fields, 1620, &doc.InsuranceHolderSurnameCyrl)
	descramble(fields, 1621)
	assignField(fields, 1621, &doc.InsuranceHolderSurname)
	descramble(fields, 1622)
	assignField(fields, 1622, &doc.InsuranceHolderNameCyrl)
	descramble(fields, 1623)
	assignField(fields, 1623, &doc.InsuranceHolderName)
	assignField(fields, 1624, &doc.InsuranceStartDate)
	localization.FormatDate(&doc.InsuranceStartDate)
	descramble(fields, 1626)
	assignField(fields, 1626, &doc.AddressState)
	descramble(fields, 1630)
	assignField(fields, 1630, &doc.ObligeeName)
	descramble(fields, 1631)
	assignField(fields, 1631, &doc.ObligeePlace)
	assignField(fields, 1632, &doc.ObligeeIdNumber)
	if len(doc.ObligeeIdNumber) == 0 {
		assignField(fields, 1633, &doc.ObligeeIdNumber)
	}
	descramble(fields, 1634)
	assignField(fields, 1634, &doc.ObligeeActivity)

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

func (card MedicalCard) readFile(name []byte, _ bool) ([]byte, error) {
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

	rsp, err := card.readFile([]byte{0x0D, 0x01}, false)
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
