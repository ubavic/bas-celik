package card

import (
	"bytes"
	"fmt"
	"image"

	"github.com/ubavic/bas-celik/card/tlv"
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
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}
	tlv.AssignField(fields, 1546, &doc.DocumentNumber)
	tlv.AssignField(fields, 1547, &doc.DocumentType)
	tlv.AssignField(fields, 1548, &doc.DocumentSerialNumber)
	tlv.AssignField(fields, 1549, &doc.IssuingDate)
	tlv.AssignField(fields, 1550, &doc.ExpiryDate)
	tlv.AssignField(fields, 1551, &doc.IssuingAuthority)
	localization.FormatDate(&doc.IssuingDate)
	localization.FormatDate(&doc.ExpiryDate)

	return nil
}

func parseIdPersonalFile(data []byte, doc *document.IdDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}

	tlv.AssignField(fields, 1558, &doc.PersonalNumber)
	tlv.AssignField(fields, 1559, &doc.Surname)
	tlv.AssignField(fields, 1560, &doc.GivenName)
	tlv.AssignField(fields, 1561, &doc.ParentGivenName)
	tlv.AssignField(fields, 1562, &doc.Sex)
	tlv.AssignField(fields, 1563, &doc.PlaceOfBirth)
	tlv.AssignField(fields, 1564, &doc.CommunityOfBirth)
	tlv.AssignField(fields, 1565, &doc.StateOfBirth)
	tlv.AssignField(fields, 1566, &doc.DateOfBirth)
	localization.FormatDate(&doc.DateOfBirth)

	return nil
}

func parseIdResidenceFile(data []byte, doc *document.IdDocument) error {
	fields, err := tlv.ParseTLV(data)
	if err != nil {
		return err
	}
	tlv.AssignField(fields, 1568, &doc.State)
	tlv.AssignField(fields, 1569, &doc.Community)
	tlv.AssignField(fields, 1570, &doc.Place)
	tlv.AssignField(fields, 1571, &doc.Street)
	tlv.AssignField(fields, 1572, &doc.AddressNumber)
	tlv.AssignField(fields, 1573, &doc.AddressLetter)
	tlv.AssignField(fields, 1574, &doc.AddressEntrance)
	tlv.AssignField(fields, 1575, &doc.AddressFloor)
	tlv.AssignField(fields, 1578, &doc.AddressApartmentNumber)
	tlv.AssignField(fields, 1580, &doc.AddressDate)
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
