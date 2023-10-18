package document

import (
	"fmt"
	"image"
	"math"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/widgets"
)

type IdDocument struct {
	Loaded                 bool
	Photo                  image.Image
	DocumentNumber         string
	IssuingDate            string
	ExpiryDate             string
	IssuingAuthority       string
	PersonalNumber         string
	Surname                string
	GivenName              string
	ParentName             string
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

func (doc *IdDocument) formatName() string {
	return doc.GivenName + ", " + doc.ParentName + ", " + doc.Surname
}

func (doc *IdDocument) formatAddress() string {
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

func (doc *IdDocument) formatPlaceOfBirth() string {
	var placeOfBirth strings.Builder

	placeOfBirth.WriteString(doc.PlaceOfBirth)
	placeOfBirth.WriteString(", ")
	placeOfBirth.WriteString(doc.CommunityOfBirth)
	placeOfBirth.WriteString(", ")
	placeOfBirth.WriteString(doc.StateOfBirth)

	return placeOfBirth.String()
}

func (doc IdDocument) BuildUI(pdfHandler func(), statusBar *widgets.StatusBar) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.formatName(), 350)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 100)
	sexF := widgets.NewField("Pol", doc.Sex, 50)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF := widgets.NewField("Mesto rođenja, opština i država", doc.formatPlaceOfBirth(), 350)
	addressF := widgets.NewField("Prebivalište i adresa stana", doc.formatAddress(), 350)
	addressDateF := widgets.NewField("Datum promene adrese", doc.AddressDate, 10)
	personInformationGroup := widgets.NewGroup("Podaci o građaninu", nameF, birthRow, birthPlaceF, addressF, addressDateF)

	issuedByF := widgets.NewField("Dokument izdaje", doc.IssuingAuthority, 10)
	documentNumberF := widgets.NewField("Broj dokumenta", doc.DocumentNumber, 100)
	issueDateF := widgets.NewField("Datum izdavanja", doc.IssuingDate, 100)
	expiryDateF := widgets.NewField("Važi do", doc.ExpiryDate, 100)
	docRow := container.New(layout.NewHBoxLayout(), documentNumberF, issueDateF, expiryDateF)
	docGroup := widgets.NewGroup("Podaci o dokumentu", issuedByF, docRow)
	colRight := container.New(layout.NewVBoxLayout(), personInformationGroup, docGroup)

	imgWidget := canvas.NewImageFromImage(doc.Photo)
	imgWidget.SetMinSize(fyne.Size{Width: 200, Height: 250})
	imgWidget.FillMode = canvas.ImageFillContain
	colLeft := container.New(layout.NewVBoxLayout(), imgWidget)
	cols := container.New(layout.NewHBoxLayout(), colLeft, colRight)

	saveButton := widget.NewButton("Sačuvaj PDF", pdfHandler)
	buttonBar := container.New(layout.NewHBoxLayout(), statusBar, layout.NewSpacer(), saveButton)

	return container.New(layout.NewVBoxLayout(), cols, buttonBar)
}

func (doc *IdDocument) BuildPdf() ([]byte, string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	err := pdf.AddTTFFontData("liberationsans", fontRegular)
	if err != nil {
		return nil, "", fmt.Errorf("loading font: %w", err)
	}

	err = pdf.SetFont("liberationsans", "", 13.5)
	if err != nil {
		return nil, "", fmt.Errorf("setting font: %w", err)
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

	putData := func(label, data string) {
		y := pdf.GetY()

		pdf.SetX(textLeftMargin)
		texts, _ := pdf.SplitTextWithWordWrap(label, 120)
		for i, text := range texts {
			_ = pdf.Cell(nil, text)
			if i < len(texts)-1 {
				pdf.SetXY(textLeftMargin, pdf.GetY()+12)
			}
		}

		y1 := pdf.GetY()

		pdf.SetXY(textLeftMargin+128, y)
		texts, _ = pdf.SplitTextWithWordWrap(data, 350)
		for i, text := range texts {
			_ = pdf.Cell(nil, text)
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
	pdf.SetCharSpacing(-0.2)
	pdf.Cell(nil, "ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")
	pdf.SetCharSpacing(-0.1)

	pdf.SetY(88)
	line(0)

	pdf.ImageFrom(doc.Photo, leftMargin, 102.8, &gopdf.Rect{W: 119.9, H: 159})
	pdf.SetLineWidth(0.48)
	pdf.SetFillColor(255, 255, 255)
	pdf.Rectangle(leftMargin, 102.8, 179, 262, "D", 0, 0)
	pdf.SetFillColor(0, 0, 0)

	pdf.SetY(276)
	line(1.08)
	moveY(8)
	pdf.SetXY(textLeftMargin, 284)
	pdf.SetFontSize(11.1)
	pdf.Cell(nil, "Podaci o građaninu")
	moveY(16)
	line(0)
	moveY(9)

	putData("Prezime:", doc.Surname)
	putData("Ime:", doc.GivenName)
	putData("Ime jednog roditelja:", doc.ParentName)
	putData("Datum rođenja:", doc.DateOfBirth)
	putData("Mesto rođenja, opština i država:", doc.formatPlaceOfBirth())
	putData("Prebivalište:", doc.formatAddress())
	putData("Datum promene adrese:", doc.AddressDate)
	putData("JMBG:", doc.PersonalNumber)
	putData("Pol:", doc.Sex)

	moveY(-8.67)
	line(0)
	moveY(9)
	pdf.Cell(nil, "Podaci o dokumentu")
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

	pdf.Cell(nil, "Datum štampe: "+time.Now().Format("02.01.2006."))

	pdf.SetY(730.6)
	line(0.83)

	pdf.SetFontSize(9)

	pdf.SetXY(leftMargin, 739.7)
	pdf.Cell(nil, "1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su")
	pdf.SetXY(leftMargin, 749.9)
	pdf.Cell(nil, "ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.")
	pdf.SetXY(leftMargin, 759.7)
	pdf.Cell(nil, "2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je")
	pdf.SetXY(leftMargin, 769.4)
	pdf.Cell(nil, "ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu")
	pdf.SetXY(leftMargin, 779.1)
	pdf.Cell(nil, "skraćuje se na prva dva karaktera")

	pdf.SetY(794.5)
	line(0)

	fileName := strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}
