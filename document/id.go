package document

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/localization"
)

const ID_TYPE_APOLLO = ""
const ID_TYPE_ID = "ID"
const ID_TYPE_IDENTITY_FOREIGNER = "IF"
const ID_TYPE_RESIDENCE_PERMIT = "RP"

// Represents a document stored on a Serbian ID card.
type IdDocument struct {
	Portrait             image.Image
	DocRegNo             string
	DocumentType         string
	IssuingDate          string
	ExpiryDate           string
	IssuingAuthority     string
	DocumentSerialNumber string
	ChipSerialNumber     string
	DocumentName         string
	PersonalNumber       string
	Surname              string
	GivenName            string
	ParentGivenName      string
	Sex                  string
	PlaceOfBirth         string
	CommunityOfBirth     string
	StateOfBirth         string
	StateOfBirthCode     string
	DateOfBirth          string
	StatusOfForeigner    string
	NationalityFull      string
	PurposeOfStay        string
	ENote                string
	State                string
	Community            string
	Place                string
	Street               string
	HouseNumber          string
	HouseLetter          string
	Entrance             string
	Floor                string
	ApartmentNumber      string
	AddressDate          string
	AddressLabel         string
}

func (doc *IdDocument) GetFullName() string {
	return localization.JoinWithComma(doc.GivenName, doc.ParentGivenName, doc.Surname)
}

func (doc *IdDocument) GetFullAddress(reverse bool) string {
	var streetAndNumber = doc.Street

	if doc.HouseNumber != "" || doc.HouseLetter != "" || doc.Entrance != "" {
		streetAndNumber += " " + doc.HouseNumber + doc.HouseLetter

		if doc.Floor != "" {
			streetAndNumber += "/" + doc.Floor
		}

		if doc.ApartmentNumber != "" {
			streetAndNumber += "/" + doc.ApartmentNumber
		}
	}

	if reverse {
		return localization.JoinWithComma(doc.Place, doc.Community, streetAndNumber)
	} else {
		return localization.JoinWithComma(streetAndNumber, doc.Community, doc.Place)
	}
}

func (doc *IdDocument) GetFullPlaceOfBirth() string {
	return localization.JoinWithComma(doc.PlaceOfBirth, doc.CommunityOfBirth, doc.StateOfBirth)
}

func (doc *IdDocument) BuildPdf() (data []byte, fileName string, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				retErr = x
			default:
				retErr = errors.New("unknown panic")
			}
		}
	}()

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	err := pdf.AddTTFFontData("liberationsans", fontRegular)
	if err != nil {
		panic(fmt.Errorf("loading font: %w", err))
	}

	err = pdf.SetFont("liberationsans", "", 13.5)
	if err != nil {
		panic(fmt.Errorf("setting font: %w", err))
	}

	ipw := IdPdfWriter{
		pdf:            &pdf,
		leftMargin:     58.8,
		rightMargin:    535,
		textLeftMargin: 67.3,
		doc:            doc,
	}

	if doc.DocumentType == ID_TYPE_APOLLO || doc.DocumentType == ID_TYPE_ID {
		ipw.printRegularId()
	} else if doc.DocumentType == ID_TYPE_IDENTITY_FOREIGNER {
		ipw.printForeignerId()
	} else if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		ipw.printResidencePermit()
	}

	fileName = doc.formatFilename() + ".pdf"

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}

func (doc *IdDocument) BuildJson() ([]byte, error) {
	var bs bytes.Buffer
	err := jpeg.Encode(&bs, doc.Portrait, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, fmt.Errorf("creating json: %w", err)
	}

	type Alias IdDocument
	return json.Marshal(&struct {
		Portrait string
		*Alias
	}{
		Portrait: base64.StdEncoding.EncodeToString(bs.Bytes()),
		Alias:    (*Alias)(doc),
	})
}

func (doc *IdDocument) BuildExcel() ([]byte, string, error) {
	xlsx, err := CreateExcel(*doc)
	filename := doc.formatFilename() + ".xlsx"
	return xlsx, filename, err
}

func (doc *IdDocument) formatFilename() string {
	return strings.ToLower(doc.GivenName + "_" + doc.Surname)
}
