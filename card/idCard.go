package card

import (
	"bytes"
	"fmt"
	"image"

	"github.com/ubavic/bas-celik/document"
	doc "github.com/ubavic/bas-celik/document"
)

var DOCUMENT_FILE_LOC = []byte{0x0F, 0x02}
var PERSONAL_FILE_LOC = []byte{0x0F, 0x03}
var RESIDENCE_FILE_LOC = []byte{0x0F, 0x04}
var PHOTO_FILE_LOC = []byte{0x0F, 0x06}

func readIDCard(card Card) (*doc.IdDocument, error) {
	rsp, err := card.readFile(DOCUMENT_FILE_LOC, false)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	doc := doc.IdDocument{}

	fields := parseResponse(rsp)
	assignField(fields, 1546, &doc.DocumentNumber)
	assignField(fields, 1549, &doc.IssuingDate)
	assignField(fields, 1550, &doc.ExpiryDate)
	assignField(fields, 1551, &doc.IssuingAuthority)
	document.FormatDate(&doc.IssuingDate)
	document.FormatDate(&doc.ExpiryDate)

	rsp, err = card.readFile(PERSONAL_FILE_LOC, false)
	if err != nil {
		return nil, fmt.Errorf("reading personal file: %w", err)
	}

	fields = parseResponse(rsp)
	assignField(fields, 1558, &doc.PersonalNumber)
	assignField(fields, 1559, &doc.Surname)
	assignField(fields, 1560, &doc.GivenName)
	assignField(fields, 1561, &doc.ParentName)
	assignField(fields, 1562, &doc.Sex)
	assignField(fields, 1563, &doc.PlaceOfBirth)
	assignField(fields, 1564, &doc.CommunityOfBirth)
	assignField(fields, 1565, &doc.StateOfBirth)
	assignField(fields, 1566, &doc.DateOfBirth)
	document.FormatDate(&doc.DateOfBirth)

	rsp, err = card.readFile(RESIDENCE_FILE_LOC, false)
	if err != nil {
		return nil, fmt.Errorf("reading residence file: %w", err)
	}

	fields = parseResponse(rsp)
	assignField(fields, 1568, &doc.State)
	assignField(fields, 1569, &doc.Community)
	assignField(fields, 1570, &doc.Place)
	assignField(fields, 1571, &doc.Street)
	assignField(fields, 1572, &doc.AddressNumber)
	assignField(fields, 1573, &doc.AddressLetter)
	assignField(fields, 1574, &doc.AddressEntrance)
	assignField(fields, 1575, &doc.AddressFloor)
	assignField(fields, 1578, &doc.AddressApartmentNumber)
	assignField(fields, 1580, &doc.AddressDate)
	document.FormatDate(&doc.AddressDate)

	rsp, err = card.readFile(PHOTO_FILE_LOC, true)
	if err != nil {
		return nil, fmt.Errorf("reading photo file: %w", err)
	}

	doc.Photo, _, err = image.Decode(bytes.NewReader(rsp))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("decoding photo file: %w", err)
	}

	doc.Loaded = true

	return &doc, nil
}
