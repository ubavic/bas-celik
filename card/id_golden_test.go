package card

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"testing"

	"github.com/ubavic/bas-celik/v2/document"
)

// Encodes a single record of the simple TLV format used by Serbian ID cards:
// a two-byte little-endian tag, a two-byte little-endian length and the value.
func tlvRecord(tag uint16, value []byte) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint16(out[0:], tag)
	binary.LittleEndian.PutUint16(out[2:], uint16(len(value)))
	return append(out, value...)
}

func tlvConcat(records ...[]byte) []byte {
	out := []byte{}
	for _, r := range records {
		out = append(out, r...)
	}
	return out
}

// Produces a minimal valid JPEG so the portrait decoding step succeeds.
func syntheticPortrait() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Builds synthetic ID card files whose TLV tag layout mirrors a real Serbian
// Gemalto ID card (ATR GEMALTO_ATR_2). All values are fabricated; only the tag
// layout reproduces a genuine card so the tag-to-field mapping in
// Gemalto.GetDocument is exercised without exposing any personal data.
func syntheticIdCard() Gemalto {
	documentFile := tlvConcat(
		tlvRecord(1546, []byte("123456789")),
		tlvRecord(1547, []byte("ID")),
		tlvRecord(1548, []byte("RS12345678")),
		tlvRecord(1549, []byte("01012020")),
		tlvRecord(1550, []byte("01012030")),
		tlvRecord(1551, []byte("MUP Republike Srbije")),
	)

	personalFile := tlvConcat(
		tlvRecord(1558, []byte("0101990710020")),
		tlvRecord(1559, []byte("Petrović")),
		tlvRecord(1560, []byte("Petar")),
		tlvRecord(1561, []byte("Marko")),
		tlvRecord(1562, []byte("M")),
		tlvRecord(1563, []byte("Beograd")),
		tlvRecord(1565, []byte("SRB")),
		tlvRecord(1566, []byte("01011990")),
		tlvRecord(1567, []byte("SRB")),
	)

	residenceFile := tlvConcat(
		tlvRecord(1568, []byte("SRB")),
		tlvRecord(1569, []byte("Vračar")),
		tlvRecord(1570, []byte("Beograd")),
		tlvRecord(1571, []byte("Kneza Miloša")),
		tlvRecord(1572, []byte("12")),
		tlvRecord(1575, []byte("3")),
		tlvRecord(1578, []byte("5")),
		tlvRecord(1580, []byte("15062021")),
		tlvRecord(1581, []byte("stan")),
	)

	return Gemalto{
		documentFile:  documentFile,
		personalFile:  personalFile,
		residenceFile: residenceFile,
		photoFile:     syntheticPortrait(),
	}
}

func Test_Gemalto_GetDocument(t *testing.T) {
	card := syntheticIdCard()

	doc, err := card.GetDocument()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	id, ok := doc.(*document.IdDocument)
	if !ok {
		t.Fatalf("expected *document.IdDocument, got %T", doc)
	}

	cases := []struct {
		field    string
		got      string
		expected string
	}{
		{"DocRegNo", id.DocRegNo, "123456789"},
		{"DocumentType", id.DocumentType, "ID"},
		{"DocumentSerialNumber", id.DocumentSerialNumber, "RS12345678"},
		{"IssuingDate", id.IssuingDate, "01.01.2020."},
		{"ExpiryDate", id.ExpiryDate, "01.01.2030."},
		{"IssuingAuthority", id.IssuingAuthority, "MUP Republike Srbije"},
		{"PersonalNumber", id.PersonalNumber, "0101990710020"},
		{"Surname", id.Surname, "Petrović"},
		{"GivenName", id.GivenName, "Petar"},
		{"ParentGivenName", id.ParentGivenName, "Marko"},
		{"Sex", id.Sex, "M"},
		{"PlaceOfBirth", id.PlaceOfBirth, "Beograd"},
		{"StateOfBirth", id.StateOfBirth, "SRB"},
		{"DateOfBirth", id.DateOfBirth, "01.01.1990."},
		{"StateOfBirthCode", id.StateOfBirthCode, "SRB"},
		{"State", id.State, "SRB"},
		{"Community", id.Community, "Vračar"},
		{"Place", id.Place, "Beograd"},
		{"Street", id.Street, "Kneza Miloša"},
		{"HouseNumber", id.HouseNumber, "12"},
		{"Floor", id.Floor, "3"},
		{"ApartmentNumber", id.ApartmentNumber, "5"},
		{"AddressDate", id.AddressDate, "15.06.2021."},
		{"AddressLabel", id.AddressLabel, "stan"},
	}

	for _, c := range cases {
		if c.got != c.expected {
			t.Errorf("%s: got %q, expected %q", c.field, c.got, c.expected)
		}
	}

	if id.Portrait == nil {
		t.Errorf("Portrait: expected decoded image, got nil")
	}
}
