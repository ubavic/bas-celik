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
)

type IdDocument struct {
	Portrait               image.Image
	DocumentNumber         string
	DocumentType           string
	DocumentSerialNumber   string
	IssuingDate            string
	ExpiryDate             string
	IssuingAuthority       string
	PersonalNumber         string
	Surname                string
	GivenName              string
	ParentGivenName        string
	Sex                    string
	PlaceOfBirth           string
	CommunityOfBirth       string
	StateOfBirth           string
	StateOfBirthCode       string
	DateOfBirth            string
	State                  string
	Community              string
	Place                  string
	Street                 string
	AddressNumber          string
	AddressLetter          string
	AddressEntrance        string
	AddressFloor           string
	AddressApartmentNumber string
	AddressDate            string
}

func (doc *IdDocument) FormatName() string {
	return doc.GivenName + ", " + doc.ParentGivenName + ", " + doc.Surname
}

func (doc *IdDocument) FormatAddress() string {
	var address strings.Builder

	address.WriteString(doc.Street)
	address.WriteString(" ")
	address.WriteString(doc.AddressNumber)
	address.WriteString(doc.AddressLetter)

	if len(doc.AddressApartmentNumber) != 0 {
		address.WriteString("/")
		address.WriteString(doc.AddressApartmentNumber)
	}

	if len(doc.Community) > 0 {
		address.WriteString(", ")
		address.WriteString(doc.Community)
	}

	address.WriteString(", ")
	address.WriteString(doc.Place)

	return address.String()
}

func (doc *IdDocument) FormatPlaceOfBirth() string {
	var placeOfBirth strings.Builder

	placeOfBirth.WriteString(doc.PlaceOfBirth)
	placeOfBirth.WriteString(", ")
	placeOfBirth.WriteString(doc.CommunityOfBirth)
	placeOfBirth.WriteString(", ")
	placeOfBirth.WriteString(doc.StateOfBirth)

	return placeOfBirth.String()
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
	line(0)

	err = pdf.ImageFrom(doc.Portrait, leftMargin, 102.8, &gopdf.Rect{W: 119.9, H: 159})
	if err != nil {
		panic(err)
	}

	pdf.SetLineWidth(0.48)
	pdf.SetFillColor(255, 255, 255)
	err = pdf.Rectangle(leftMargin, 102.8, 179, 262, "D", 0, 0)
	if err != nil {
		panic(err)
	}

	pdf.SetFillColor(0, 0, 0)

	pdf.SetY(276)
	line(1.08)
	moveY(8)
	pdf.SetXY(textLeftMargin, 284)
	err = pdf.SetFontSize(11.1)
	if err != nil {
		panic(err)
	}
	cell("Podaci o građaninu")
	moveY(16)
	line(0)
	moveY(9)

	putData("Prezime:", doc.Surname)
	putData("Ime:", doc.GivenName)
	putData("Ime jednog roditelja:", doc.ParentGivenName)
	putData("Datum rođenja:", doc.DateOfBirth)
	putData("Mesto rođenja, opština i država:", doc.FormatPlaceOfBirth())
	putData("Prebivalište:", doc.FormatAddress())
	putData("Datum promene adrese:", doc.AddressDate)
	putData("JMBG:", doc.PersonalNumber)
	putData("Pol:", doc.Sex)

	moveY(-8.67)
	line(0)
	moveY(9)
	cell("Podaci o dokumentu")
	moveY(16)

	line(0)
	moveY(9)
	putData("Dokument izdaje:", doc.IssuingAuthority)
	putData("Broj dokumenta:", doc.DocumentNumber)
	putData("Datum izdavanja:", doc.IssuingDate)
	putData("Važi do:", doc.ExpiryDate)

	moveY(-8.67)
	line(0)
	moveY(3)
	line(0)
	moveY(9)

	cell("Datum štampe: " + time.Now().Format("02.01.2006."))

	pdf.SetY(730.6)
	line(0.83)

	err = pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	pdf.SetXY(leftMargin, 739.7)
	cell("1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su")
	pdf.SetXY(leftMargin, 749.9)
	cell("ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.")
	pdf.SetXY(leftMargin, 759.7)
	cell("2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je")
	pdf.SetXY(leftMargin, 769.4)
	cell("ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu")
	pdf.SetXY(leftMargin, 779.1)
	cell("skraćuje se na prva dva karaktera")

	pdf.SetY(794.5)
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
