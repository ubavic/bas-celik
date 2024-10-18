package document

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"math"
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

	const leftMargin = 58.8
	const rightMargin = 535
	const textLeftMargin = 67.3

	line := func(width float64) {
		if width > 0 {
			pdf.SetLineWidth(width)
		}

		y := pdf.GetY()
		pdf.Line(leftMargin, y, rightMargin, y)
	}

	moveY := func(y float64) {
		pdf.SetXY(pdf.GetX(), pdf.GetY()+y)
	}

	cell := func(s string) {
		err := pdf.Cell(nil, s)
		if err != nil {
			panic(fmt.Errorf("putting text: %w", err))
		}
	}

	putData := func(label, data string) {
		y := pdf.GetY()

		pdf.SetX(textLeftMargin)
		texts, err := pdf.SplitTextWithWordWrap(label, 120)
		if err != nil && err != gopdf.ErrEmptyString {
			panic(err)
		}

		for i, text := range texts {
			cell(text)
			if i < len(texts)-1 {
				pdf.SetXY(textLeftMargin, pdf.GetY()+12)
			}
		}

		y1 := pdf.GetY()

		pdf.SetXY(textLeftMargin+128, y)
		texts, err = pdf.SplitTextWithWordWrap(data, 350)
		if err != nil && err != gopdf.ErrEmptyString {
			panic(err)
		}

		for i, text := range texts {
			cell(text)
			if i < len(texts)-1 {
				pdf.SetXY(textLeftMargin+128, pdf.GetY()+12)
			}
		}

		y2 := pdf.GetY()

		pdf.SetXY(textLeftMargin, math.Max(y1, y2)+24.67)
	}

	pdf.SetLineType("solid")
	pdf.SetY(59.041)
	line(0.83)

	pdf.SetXY(textLeftMargin+1.0, 68.5)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		pdf.SetY(64.95)
	}

	err = pdf.SetCharSpacing(-0.2)
	if err != nil {
		panic(err)
	}
	cell("ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")

	err = pdf.SetCharSpacing(-0.1)
	if err != nil {
		panic(err)
	}

	pdf.SetY(88)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		pdf.SetY(79.8)
	}

	line(0)

	imageY := 102.8
	imageHeight := 159.0
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		imageY = 86
	}

	err = pdf.ImageFrom(doc.Portrait, leftMargin, imageY, &gopdf.Rect{W: 119.9, H: imageHeight})
	if err != nil {
		panic(err)
	}

	pdf.SetLineWidth(0.48)
	pdf.SetFillColor(255, 255, 255)
	err = pdf.Rectangle(leftMargin, imageY, 179, imageY+imageHeight, "D", 0, 0)
	if err != nil {
		panic(err)
	}

	pdf.SetFillColor(0, 0, 0)

	pdf.SetY(276)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		pdf.SetY(250)
	}
	line(1.08)
	moveY(8)
	pdf.SetX(textLeftMargin)
	err = pdf.SetFontSize(11.1)
	if err != nil {
		panic(err)
	}

	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		cell("Podaci o strancu")
	} else {
		cell("Podaci o građaninu")
	}

	moveY(16)
	line(0)
	moveY(9)

	putData("Prezime:", doc.Surname)
	putData("Ime:", doc.GivenName)

	if doc.DocumentType != ID_TYPE_RESIDENCE_PERMIT {
		putData("Ime jednog roditelja:", doc.ParentGivenName)
	} else {
		putData("Državljanstvo:", doc.NationalityFull)
	}
	putData("Datum rođenja:", doc.DateOfBirth)
	putData("Mesto rođenja,\nopština i država:", doc.GetFullPlaceOfBirth())
	putData("Prebivalište:", doc.GetFullAddress(true))
	putData("Datum promene adrese:", doc.AddressDate)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		putData("Evidencijski broj\nstranca:", doc.PersonalNumber)
	} else {
		putData("JMBG:", doc.PersonalNumber)
	}
	putData("Pol:", doc.Sex)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		putData("Osnov boravka:", doc.PurposeOfStay)
		putData("Napomena:", doc.ENote)
	}

	moveY(-8.67)
	line(0)
	moveY(9)
	cell("Podaci o dokumentu")
	moveY(16)

	line(0)
	moveY(9)
	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		putData("Naziv dokumenta:", doc.DocumentName)
	}
	putData("Dokument izdaje:", doc.IssuingAuthority)
	putData("Broj dokumenta:", doc.DocRegNo)
	putData("Datum izdavanja:", doc.IssuingDate)
	putData("Važi do:", doc.ExpiryDate)

	moveY(-8.67)
	line(0)
	moveY(3)
	line(0)
	moveY(9)

	cell("Datum štampe: " + time.Now().Format("02.01.2006."))

	moveY(19)

	if doc.DocumentType == ID_TYPE_APOLLO || doc.DocumentType == ID_TYPE_ID {
		if pdf.GetY() < 700 {
			pdf.SetY(730.6)
		}
	}

	line(0.83)

	err = pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	moveY(4)
	if doc.DocumentType == ID_TYPE_APOLLO || doc.DocumentType == ID_TYPE_ID {
		moveY(6)
	}

	pdf.SetX(leftMargin)

	if doc.DocumentType == ID_TYPE_RESIDENCE_PERMIT {
		cell("1. U čipu dozvole za privremeni boravak i rad, podaci o imenu i prezimenu imaoca dozvole ispisani su onako")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("kako su ispisani na samom obrascu dozvole za privremeni boravak latiničnim pismom.")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("2. Ako se ime ili prezime stranca sastoji od dve ili više reči čija dužina prelazi 30 karaktera za ime,")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("odnosno 36 karaktera za prezime u čip se upisuje puno ime stranca, a na obrascu dozvole za privremeni boravak")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("se upisuje do 30 karaktera za ime, odnosno 36 karaktera za prezime.")
	} else {
		cell("1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu")
		pdf.SetX(leftMargin)
		moveY(9.7)
		cell("skraćuje se na prva dva karaktera")
	}

	moveY(9.7)

	if doc.DocumentType == ID_TYPE_APOLLO || doc.DocumentType == ID_TYPE_ID {
		moveY(6)
	}

	line(0)

	fileName = strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

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

func (doc *IdDocument) BuildExcel() ([]byte, error) {
	return CreateExcel(*doc)
}
