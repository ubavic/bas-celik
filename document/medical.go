package document

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/widgets"
)

type MedicalDocument struct {
	PersonalNumber                 string
	Surname                        string
	GivenName                      string
	ParentName                     string
	SurnameCyrl                    string
	GivenNameCyrl                  string
	ParentNameCyrl                 string
	Sex                            string
	InsuranceNumber                string
	Language                       string
	DateOfBirth                    string
	AddressApartmentNumber         string
	AddressNumber                  string
	AddressStreet                  string
	AddressTown                    string
	AddressMunicipality            string
	AddressState                   string
	CardIssueDate                  string
	CardExpiryDate                 string
	ValidUntil                     string
	PermanentlyValid               bool
	InsuranceHolderName            string
	InsuranceHolderSurname         string
	InsuranceHolderNameCyrl        string
	InsuranceHolderSurnameCyrl     string
	InsuranceHolderPersonalNumber  string
	InsuranceHolderInsuranceNumber string
	InsuranceHolderIsFamilyMember  bool
	InsuranceHolderRelation        string
	InsuranceReason                string
	InsuranceStartDate             string
	InsuranceDescription           string
	ObligeeName                    string
	ObligeePlace                   string
	ObligeeRegistrationNumber      string
	ObligeeIdNumber                string
	ObligeeActivity                string
}

func (doc *MedicalDocument) formatName() string {
	return doc.GivenName + ", " + doc.ParentName + ", " + doc.Surname
}

func (doc *MedicalDocument) formatAddress() string {
	return doc.GivenName + ", " + doc.ParentName + ", " + doc.Surname
}

func (doc MedicalDocument) BuildUI(pdfHandler func()) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.formatName(), 350)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 100)
	sexF := widgets.NewField("Pol", doc.Sex, 80)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	adressF := widgets.NewField("Mesto, opština i država", doc.formatAddress(), 350)
	insuranceNumberF := widgets.NewField("LBO", doc.InsuranceNumber, 10)
	colRight := container.New(layout.NewVBoxLayout(), nameF, birthRow, adressF, insuranceNumberF)

	cols := container.New(layout.NewHBoxLayout(), colRight)

	saveButton := widget.NewButton("Sačuvaj PDF", pdfHandler)
	buttonBar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), saveButton)

	return container.New(layout.NewVBoxLayout(), cols, buttonBar)
}

func (doc *MedicalDocument) BuildPdf() ([]byte, string, error) {
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

	fileName := strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}
