package document

import (
	"fmt"
	"image"
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

	err := pdf.AddTTFFontData("liberationsans", font)
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

	move := func() {
		pdf.SetXY(textLeftMargin, pdf.GetY()+24.67)
	}

	tab := func() {
		pdf.SetXY(textLeftMargin+128, pdf.GetY())
	}

	line := func(y float64) {
		pdf.Line(leftMargin, y, rightMargin, y)
	}

	pdf.SetLineWidth(0.83)
	pdf.SetLineType("solid")
	line(59.041)

	pdf.SetXY(textLeftMargin+1.0, 68.5)
	pdf.SetCharSpacing(-0.2)
	pdf.Cell(nil, "ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")
	pdf.SetCharSpacing(-0.1)

	line(88)

	pdf.ImageFrom(doc.Photo, leftMargin, 102.8, &gopdf.Rect{W: 119.9, H: 159})
	pdf.SetLineWidth(0.48)
	pdf.SetFillColor(255, 255, 255)
	pdf.Rectangle(leftMargin, 102.8, 179, 262, "D", 0, 0)
	pdf.SetFillColor(0, 0, 0)

	pdf.SetLineWidth(1.08)
	line(276)
	pdf.SetXY(textLeftMargin, 284)
	pdf.SetFontSize(11.1)
	pdf.Cell(nil, "Podaci o građaninu")

	line(300)
	pdf.SetXY(textLeftMargin, 309)
	pdf.Cell(nil, "Prezime:")
	tab()
	pdf.Cell(nil, doc.Surname)

	move()
	pdf.Cell(nil, "Ime:")
	tab()
	pdf.Cell(nil, doc.GivenName)

	move()
	pdf.Cell(nil, "Ime jednog roditelja:")
	tab()
	pdf.Cell(nil, doc.ParentName)

	move()
	pdf.Cell(nil, "Datum rođenja:")
	tab()
	pdf.Cell(nil, doc.DateOfBirth)

	move()
	pdf.Cell(nil, "Mesto rođenja, opština i")
	tab()
	pdf.Cell(nil, doc.formatPlaceOfBirth())
	pdf.SetXY(textLeftMargin, pdf.GetY()+12)
	pdf.Cell(nil, "država:")

	move()
	pdf.Cell(nil, "Prebivalište:")
	tab()
	pdf.Cell(nil, doc.formatAddress())

	move()
	pdf.Cell(nil, "Datum promene adrese:")
	tab()
	pdf.Cell(nil, doc.AddressDate)

	move()
	pdf.Cell(nil, "JMBG:")
	tab()
	pdf.Cell(nil, doc.PersonalNumber)

	move()
	pdf.Cell(nil, "Pol:")
	tab()
	pdf.Cell(nil, doc.Sex)

	line(534)
	pdf.SetXY(textLeftMargin, 543)
	pdf.Cell(nil, "Podaci o dokumentu")

	line(559)

	pdf.SetXY(textLeftMargin, 567)
	pdf.Cell(nil, "Dokument izdaje:")
	tab()
	pdf.Cell(nil, doc.IssuingAuthority)

	move()
	pdf.Cell(nil, "Broj dokumenta:")
	tab()
	pdf.Cell(nil, doc.DocumentNumber)

	move()
	pdf.Cell(nil, "Datum izdavanja:")
	tab()
	pdf.Cell(nil, doc.IssuingDate)

	move()
	pdf.Cell(nil, "Važi do:")
	tab()
	pdf.Cell(nil, doc.ExpiryDate)

	line(657)
	line(660)

	pdf.SetXY(textLeftMargin, 669.2)
	pdf.Cell(nil, "Datum štampe: "+time.Now().Format("02.01.2006."))

	pdf.SetLineWidth(0.83)
	line(730.6)

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

	line(794.5)

	fileName := strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}
