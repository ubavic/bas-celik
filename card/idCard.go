package card

import (
	"bytes"
	"fmt"
	"image"

	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/localization"
)

// Location of the file with document data.
var ID_DOCUMENT_FILE_LOC = []byte{0x0F, 0x02}

// Location of the file with personal data.
var ID_PERSONAL_FILE_LOC = []byte{0x0F, 0x03}

// Location of the file with residence data.
var ID_RESIDENCE_FILE_LOC = []byte{0x0F, 0x04}

// Location of the the portrait. Portrait is encoded as JPEG.
var ID_PHOTO_FILE_LOC = []byte{0x0F, 0x06}

func parseIdDocumentFile(data []byte, doc *document.IdDocument) error {
	fields, err := parseTLV(data)
	if err != nil {
		return err
	}
	assignField(fields, 1546, &doc.DocumentNumber)
	assignField(fields, 1547, &doc.DocumentType)
	assignField(fields, 1548, &doc.DocumentSerialNumber)
	assignField(fields, 1549, &doc.IssuingDate)
	assignField(fields, 1550, &doc.ExpiryDate)
	assignField(fields, 1551, &doc.IssuingAuthority)
	localization.FormatDate(&doc.IssuingDate)
	localization.FormatDate(&doc.ExpiryDate)

	return nil
}

func parseIdPersonalFile(data []byte, doc *document.IdDocument) error {
	fields, err := parseTLV(data)
	if err != nil {
		return err
	}

	assignField(fields, 1558, &doc.PersonalNumber)
	assignField(fields, 1559, &doc.Surname)
	assignField(fields, 1560, &doc.GivenName)
	assignField(fields, 1561, &doc.ParentGivenName)
	assignField(fields, 1562, &doc.Sex)
	assignField(fields, 1563, &doc.PlaceOfBirth)
	assignField(fields, 1564, &doc.CommunityOfBirth)
	assignField(fields, 1565, &doc.StateOfBirth)
	assignField(fields, 1566, &doc.DateOfBirth)
	localization.FormatDate(&doc.DateOfBirth)

	return nil
}

func parseIdResidenceFile(data []byte, doc *document.IdDocument) error {
	fields, err := parseTLV(data)
	if err != nil {
		return err
	}
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
	localization.FormatDate(&doc.AddressDate)

	return nil
}

func parseAndAssignIdPhotoFile(data []byte, doc *document.IdDocument) error {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decoding photo file: %w", err)
	}

	doc.Portrait = img

	return nil
}
