package document

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/localization"
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

func (doc *MedicalDocument) FormatName() string {
	return doc.GivenName + ", " + doc.ParentName + ", " + doc.Surname
}

func (doc *MedicalDocument) FormatStreetAddress() string {
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

func (doc *MedicalDocument) FormatPlaceAddress() string {
	var address strings.Builder

	address.WriteString(doc.AddressTown)

	address.WriteString(", ")
	address.WriteString(doc.AddressMunicipality)
	address.WriteString(", ")
	address.WriteString(doc.AddressState)

	return address.String()
}

func (doc *MedicalDocument) BuildPdf() (data []byte, fileName string, retErr error) {
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

	err = pdf.SetFont("liberationsans", "", 11)
	if err != nil {
		panic(fmt.Errorf("setting font: %w", err))
	}

	const leftMargin = 30.464
	const rightMargin = 535
	const textLeftMargin = 38.243

	cell := func(s string) {
		err := pdf.Cell(nil, s)
		if err != nil {
			panic(fmt.Errorf("putting text: %w", err))
		}
	}

	section := func(name string) {
		y := pdf.GetY() + 8
		pdf.Line(leftMargin, y, rightMargin, y)
		pdf.SetXY(textLeftMargin, y+12)
		cell(name)
		pdf.Line(leftMargin, y+32, rightMargin, y+32)
		pdf.SetXY(textLeftMargin, y+41)
	}

	putData := func(label, data string) {
		cell(label)
		pdf.SetXY(textLeftMargin+144, pdf.GetY())

		texts, err := pdf.SplitTextWithWordWrap(data, 350)
		if err != nil && err != gopdf.ErrEmptyString {
			panic(fmt.Errorf("splitting text: %w", err))
		}

		for i, text := range texts {
			cell(text)
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
		panic(fmt.Errorf("decoding photo file: %w", err))
	}

	err = pdf.ImageFrom(rfzoLogoImage, 36.0, 14.0, &gopdf.Rect{W: 144, H: 51})
	if err != nil {
		panic(fmt.Errorf("inserting logo: %w", err))
	}

	pdf.SetXY(218.6, 49.6)
	cell("ПРЕГЛЕД КАРТИЦЕ ЗДРАВСТВЕНОГ ОСИГУРАЊА (КЗО)")

	pdf.SetY(68)
	section("Општи подаци о осигуранику")

	putData("Име:", doc.GivenNameCyrl+" ("+doc.GivenName+")")

	putData("Име једног родитеља:", doc.ParentNameCyrl+" ("+doc.ParentName+")")

	putData("Презиме:", doc.SurnameCyrl+" ("+doc.Surname+")")

	putData("Датум рођења:", doc.DateOfBirth)

	putData("Место, општина и држава:", doc.FormatPlaceAddress())

	putData("Улица:", doc.FormatStreetAddress())

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

	fileName = strings.ToLower(doc.GivenName + "_" + doc.Surname + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenName + " " + doc.Surname,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}

func (doc *MedicalDocument) BuildJson() ([]byte, error) {
	return json.Marshal(doc)
}
