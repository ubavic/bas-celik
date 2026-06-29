package card

import (
	"encoding/binary"
	"testing"

	"github.com/ubavic/bas-celik/v2/document"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Encodes a single record of the simple TLV format (two-byte little-endian tag,
// two-byte little-endian length, value) used by medical cards.
func medRecord(tag uint16, value []byte) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint16(out[0:], tag)
	binary.LittleEndian.PutUint16(out[2:], uint16(len(value)))
	return append(out, value...)
}

func medConcat(records ...[]byte) []byte {
	out := []byte{}
	for _, r := range records {
		out = append(out, r...)
	}
	return out
}

// Encodes a string the way text fields are stored on medical cards: UTF-16
// little-endian with a byte order mark. card.descramble reverses this.
func u16(s string) []byte {
	out, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder(), []byte(s))
	if err != nil {
		panic(err)
	}
	return out
}

// Builds synthetic medical card files whose TLV tag layout mirrors a real
// Serbian medical insurance card (ATR MEDICAL_ATR_1). All values are
// fabricated; only the tag layout and the UTF-16 encoding of text fields
// reproduce a genuine card so the tag-to-field mapping in
// MedicalCard.GetDocument is exercised without exposing any personal data.
func syntheticMedicalCard() MedicalCard {
	documentFile := medConcat(
		medRecord(1553, u16("Републички фонд за здравствено осигурање")),
		medRecord(1554, []byte("06110101000000")),
		medRecord(1555, []byte("12345678901")),
		medRecord(1557, []byte("01012020")),
		medRecord(1558, []byte("01012030")),
		medRecord(1559, []byte("1234567890123456")),
		medRecord(1560, []byte("SR")),
	)

	fixedPersonalFile := medConcat(
		medRecord(1569, []byte("12345678901")),
		medRecord(1570, u16("Перић")),
		medRecord(1571, u16("Perić")),
		medRecord(1572, u16("Петар")),
		medRecord(1573, u16("Petar")),
		medRecord(1574, []byte("01011990")),
	)

	variablePersonalFile := medConcat(
		medRecord(1586, []byte("31122025")),
		medRecord(1587, []byte("1")),
	)

	variableAdminFile := medConcat(
		medRecord(1601, u16("Марко")),
		medRecord(1602, u16("Marko")),
		medRecord(1603, []byte("02")),
		medRecord(1604, []byte("0101990710020")),
		medRecord(1605, u16("Кнеза Милоша")),
		medRecord(1606, []byte("11000")),
		medRecord(1607, u16("Врачар")),
		medRecord(1608, u16("Београд")),
		medRecord(1609, []byte("00123")),
		medRecord(1610, u16("5")),
		medRecord(1614, []byte("123")),
		medRecord(1615, u16("Самостална делатност")),
		medRecord(1624, []byte("15062021")),
		medRecord(1626, u16("Србија")),
		medRecord(1630, u16("Петар Перић")),
		medRecord(1631, u16("Б")),
		medRecord(1633, []byte("0101990710020")),
	)

	return MedicalCard{
		medicalDocumentFile:  documentFile,
		fixedPersonalFile:    fixedPersonalFile,
		variablePersonalFile: variablePersonalFile,
		variableAdminFile:    variableAdminFile,
	}
}

func Test_MedicalCard_GetDocument(t *testing.T) {
	card := syntheticMedicalCard()

	doc, err := card.GetDocument()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	medical, ok := doc.(*document.MedicalDocument)
	if !ok {
		t.Fatalf("expected *document.MedicalDocument, got %T", doc)
	}

	cases := []struct {
		field    string
		got      string
		expected string
	}{
		{"InsurerName", medical.InsurerName, "Републички фонд за здравствено осигурање"},
		{"InsurerID", medical.InsurerID, "06110101000000"},
		{"CardId", medical.CardId, "12345678901"},
		{"DateOfIssue", medical.DateOfIssue, "01.01.2020."},
		{"DateOfExpiry", medical.DateOfExpiry, "01.01.2030."},
		{"ChipSerialNumber", medical.ChipSerialNumber, "1234567890123456"},
		{"PrintLanguage", medical.PrintLanguage, "SR"},
		{"InsurantNumber", medical.InsurantNumber, "12345678901"},
		{"FamilyName", medical.FamilyName, "Перић"},
		{"FamilyNameLatin", medical.FamilyNameLatin, "Perić"},
		{"GivenName", medical.GivenName, "Петар"},
		{"GivenNameLatin", medical.GivenNameLatin, "Petar"},
		{"DateOfBirth", medical.DateOfBirth, "01.01.1990."},
		{"ValidUntil", medical.ValidUntil, "31.12.2025."},
		{"ParentName", medical.ParentName, "Марко"},
		{"ParentNameLatin", medical.ParentNameLatin, "Marko"},
		{"Gender", medical.Gender, "Женско"},
		{"PersonalNumber", medical.PersonalNumber, "0101990710020"},
		{"Street", medical.Street, "Кнеза Милоша"},
		{"PostNumber", medical.PostNumber, "11000"},
		{"Municipality", medical.Municipality, "Врачар"},
		{"Place", medical.Place, "Београд"},
		{"StreetCode", medical.StreetCode, "00123"},
		{"Number", medical.Number, "5"},
		{"InsuranceBasisRZZO", medical.InsuranceBasisRZZO, "123"},
		{"InsuranceDescription", medical.InsuranceDescription, "Самостална делатност"},
		{"InsuranceStartDate", medical.InsuranceStartDate, "15.06.2021."},
		{"Country", medical.Country, "Србија"},
		{"TaxpayerName", medical.TaxpayerName, "Петар Перић"},
		{"TaxpayerResidence", medical.TaxpayerResidence, "Б"},
		{"TaxpayerIdNumber", medical.TaxpayerIdNumber, "0101990710020"},
	}

	for _, c := range cases {
		if c.got != c.expected {
			t.Errorf("%s: got %q, expected %q", c.field, c.got, c.expected)
		}
	}

	if !medical.PermanentlyValid {
		t.Errorf("PermanentlyValid: got false, expected true")
	}
}
