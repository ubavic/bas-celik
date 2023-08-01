package document

import (
	"bytes"
	"fmt"
	"image"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/localization"
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

func (doc *MedicalDocument) formatStreetAddress() string {
	var address strings.Builder

	address.WriteString(doc.AddressStreet)
	if len(doc.AddressNumber) > 0 {
		address.WriteString(", Број: ")
		address.WriteString(doc.AddressNumber)
	}

	if len(doc.AddressApartmentNumber) > 0 {
		address.WriteString(" Стан: ")
		address.WriteString(doc.AddressApartmentNumber)
	}

	return address.String()
}

func (doc *MedicalDocument) formatPlaceAddress() string {
	var address strings.Builder

	address.WriteString(doc.AddressTown)

	address.WriteString(", ")
	address.WriteString(doc.AddressMunicipality)
	address.WriteString(", ")
	address.WriteString(doc.AddressState)

	return address.String()
}

func (doc MedicalDocument) BuildUI(pdfHandler func(), statusBar *widgets.StatusBar) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.formatName(), 350)
	sexF := widgets.NewField("Pol", doc.Sex, 170)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 170)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 170)
	insuranceNumberF := widgets.NewField("LBO", doc.InsuranceNumber, 170)
	idsRow := container.New(layout.NewHBoxLayout(), personalNumberF, insuranceNumberF)
	stateF := widgets.NewField("Država", doc.AddressState, 170)
	municipalityF := widgets.NewField("Opština", doc.AddressMunicipality, 170)
	address1Row := container.New(layout.NewHBoxLayout(), stateF, municipalityF)
	townF := widgets.NewField("Mesto", doc.AddressTown, 170)
	street := widgets.NewField("Ulica", doc.AddressStreet, 170)
	address2Row := container.New(layout.NewHBoxLayout(), townF, street)
	addressNumber := widgets.NewField("Broj", doc.AddressNumber, 76)
	addressEntrance := widgets.NewField("Ulaz", doc.AddressApartmentNumber, 76)
	addressApartmentNumber := widgets.NewField("Stan", doc.AddressApartmentNumber, 80)
	address3Row := container.New(layout.NewHBoxLayout(), addressNumber, addressEntrance, addressApartmentNumber)

	generalGroup := widgets.NewGroup("Opšti podaci", nameF, birthRow, idsRow, address1Row, address2Row, address3Row)

	insuranceReasonF := widgets.NewField("Osnov osiguranja", doc.InsuranceReason, 170)
	insuracaStartDateF := widgets.NewField("Datum početka osiguranja", doc.InsuranceStartDate, 170)
	insuranceRow := container.New(layout.NewHBoxLayout(), insuranceReasonF, insuracaStartDateF)
	insuranceDescriptionF := widgets.NewField("Opis", doc.InsuranceDescription, 350)
	insuranceGroup := widgets.NewGroup("Podaci o osiguranju", insuranceRow, insuranceDescriptionF)

	colLeft := container.New(layout.NewVBoxLayout(), generalGroup, insuranceGroup)

	insuranceHolderNameF := widgets.NewField("Nosilac osiguranja", doc.InsuranceHolderName+" "+doc.InsuranceHolderSurname, 350)
	insuranceHolderInsuranceNumberF := widgets.NewField("LBO", doc.InsuranceHolderInsuranceNumber, 170)
	insuranceHolderPersonalNumberF := widgets.NewField("JMBG", doc.InsuranceHolderPersonalNumber, 170)
	insuranceHolderRow1 := container.New(layout.NewHBoxLayout(), insuranceHolderInsuranceNumberF, insuranceHolderPersonalNumberF)
	insuranceHolderIsFamilyMemberF := widgets.NewField("Član porodice", localization.FormatYesNo(doc.InsuranceHolderIsFamilyMember, localization.Latin), 170)
	insuranceHolderRelationF := widgets.NewField("Srodstvo", doc.InsuranceHolderRelation, 170)
	insuranceHolderRow2 := container.New(layout.NewHBoxLayout(), insuranceHolderIsFamilyMemberF, insuranceHolderRelationF)
	insuranceHolderGroup := widgets.NewGroup("Podaci o nosiocu osiguranja", insuranceHolderNameF, insuranceHolderRow1, insuranceHolderRow2)

	issueDateF := widgets.NewField("Datum izdavanja", doc.CardIssueDate, 170)
	expiryDateF := widgets.NewField("Datum važenja", doc.CardExpiryDate, 170)
	cardRow1 := container.New(layout.NewHBoxLayout(), issueDateF, expiryDateF)
	validUntilF := widgets.NewField("Overena do", doc.ValidUntil, 170)
	permanentlyValidF := widgets.NewField("Trajna overa", localization.FormatYesNo(doc.PermanentlyValid, localization.Latin), 170)
	cardRow2 := container.New(layout.NewHBoxLayout(), validUntilF, permanentlyValidF)
	cardGroup := widgets.NewGroup("Podaci o kartici", cardRow1, cardRow2)

	obligeeNameF := widgets.NewField("Naziv / ime i prezime", doc.ObligeeName, 350)
	obligeeActivityF := widgets.NewField("Oznaka delatnosti", doc.ObligeeActivity, 170)
	obligeePlaceF := widgets.NewField("Sedište", doc.ObligeePlace, 170)
	obligeeRow1 := container.New(layout.NewHBoxLayout(), obligeeActivityF, obligeePlaceF)
	obligeeRegistrationNumberF := widgets.NewField("Registarski broj", doc.ObligeeRegistrationNumber, 170)
	obligeeIdNumberF := widgets.NewField("PIB/JMBG", doc.ObligeeIdNumber, 170)
	obligeeRow2 := container.New(layout.NewHBoxLayout(), obligeeRegistrationNumberF, obligeeIdNumberF)
	obligeeGroup := widgets.NewGroup("Podaci o obavezniku", obligeeNameF, obligeeRow1, obligeeRow2)

	colRight := container.New(layout.NewVBoxLayout(), insuranceHolderGroup, cardGroup, obligeeGroup)

	cols := container.New(layout.NewHBoxLayout(), colLeft, colRight)

	saveButton := widget.NewButton("Sačuvaj PDF", pdfHandler)
	buttonBar := container.New(layout.NewHBoxLayout(), statusBar, layout.NewSpacer(), saveButton)

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

	err = pdf.SetFont("liberationsans", "", 11)
	if err != nil {
		return nil, "", fmt.Errorf("setting font: %w", err)
	}

	const leftMargin = 30.464
	const rightMargin = 535
	const textLeftMargin = 38.243

	section := func(name string) {
		y := pdf.GetY() + 8
		pdf.Line(leftMargin, y, rightMargin, y)
		pdf.SetXY(textLeftMargin, y+12)
		pdf.Cell(nil, name)
		pdf.Line(leftMargin, y+32, rightMargin, y+32)
		pdf.SetXY(textLeftMargin, y+41)
	}

	putData := func(label, data string) {
		pdf.Cell(nil, label)
		pdf.SetXY(textLeftMargin+144, pdf.GetY())

		texts, _ := pdf.SplitText(data, 350)
		for i, text := range texts {
			_ = pdf.Cell(nil, text)
			if i < len(texts)-1 {
				pdf.SetXY(textLeftMargin+144, pdf.GetY()+11)
			}
		}

		pdf.SetXY(textLeftMargin, pdf.GetY()+14)
	}

	pdf.SetLineWidth(0.58)
	pdf.SetLineType("solid")

	rfzoLogoImage, _, err := image.Decode(bytes.NewReader(rfzoLogo))
	if err != nil {
		fmt.Println(err)
		return nil, "", fmt.Errorf("decoding photo file: %w", err)
	}

	pdf.ImageFrom(rfzoLogoImage, 36.0, 14.0, &gopdf.Rect{W: 144, H: 51})

	pdf.SetXY(218.6, 49.6)
	pdf.Cell(nil, "ПРЕГЛЕД КАРТИЦЕ ЗДРАВСТВЕНОГ ОСИГУРАЊА (КЗО)")

	pdf.SetY(68)
	section("Општи подаци о осигуранику")

	putData("Име:", doc.GivenNameCyrl+" ("+doc.GivenName+")")

	putData("Име једног родитеља:", doc.ParentNameCyrl+" ("+doc.ParentName+")")

	putData("Презиме:", doc.SurnameCyrl+" ("+doc.Surname+")")

	putData("Датум рођења:", doc.DateOfBirth)

	putData("Место, општина и држава:", doc.formatPlaceAddress())

	putData("Улица:", doc.formatStreetAddress())

	putData("Пол:", doc.Sex)

	putData("Језик:", doc.Language)

	putData("ЛБО:", doc.InsuranceNumber)

	putData("ЈМБГ:", doc.PersonalNumber)

	section("Подаци о картици здравственог осигурања")

	putData("Датум издавања:", doc.CardIssueDate)

	putData("Датум важења:", doc.CardExpiryDate)

	putData("Оверена до:", doc.ValidUntil)

	putData("Трајно оверена:", localization.FormatYesNo(doc.PermanentlyValid, localization.Cyrillic))

	section("Подаци о носиоцу осигурања")

	putData("Име:", doc.InsuranceHolderNameCyrl+" ("+doc.InsuranceHolderName+")")

	putData("Презиме:", doc.InsuranceHolderSurnameCyrl+" ("+doc.InsuranceHolderSurnameCyrl+")")

	putData("ЛБО:", doc.InsuranceHolderInsuranceNumber)

	putData("ЈМБГ:", doc.InsuranceHolderPersonalNumber)

	putData("Члан породице:", localization.FormatYesNo(doc.InsuranceHolderIsFamilyMember, localization.Cyrillic))

	putData("Сродство:", doc.DateOfBirth)

	section("Подаци о осигурању")

	putData("Основ осигурања:", doc.InsuranceReason)

	putData("Датум почетка осигурања:", doc.InsuranceStartDate)

	putData("Опис:", doc.InsuranceDescription)

	section("Подаци о обвезнику плаћања доприноса")

	putData("Назив:", doc.ObligeeName)

	putData("Седиште:", doc.ObligeePlace)

	putData("Регистарски број:", doc.ObligeeRegistrationNumber)

	putData("ПИБ/ЈМБГ:", doc.ObligeeIdNumber)

	putData("Делатност:", doc.ObligeeActivity)

	fileName := strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}
