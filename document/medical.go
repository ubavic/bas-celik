package document

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/localization"
)

const rfzoServiceUrl = "https://www.rfzo.rs/proveraUplateDoprinosa2.php"

// Card number doesn't have exactly 11 digits.
var ErrInvalidCardNo = errors.New("invalid card number length")

// Insurance number doesn't have exactly 11 digits.
var ErrInvalidInsuranceNo = errors.New("invalid insurance number length")

// Date `ValidUntil` could not be extracted from RFZO response.
var ErrNoSubmatchFound = errors.New("no submatch found")

// Represents a document stored on a Serbian public medical insurance card.
type MedicalDocument struct {
	InsurerName            string
	InsurerID              string
	CardId                 string
	DateOfIssue            string
	DateOfExpiry           string
	ChipSerialNumber       string
	PrintLanguage          string
	PersonalNumber         string
	FamilyNameLatin        string
	GivenNameLatin         string
	ParentNameLatin        string
	FamilyName             string
	GivenName              string
	ParentName             string
	Gender                 string
	InsurantNumber         string
	DateOfBirth            string
	Apartment              string
	Number                 string
	Street                 string
	Place                  string
	Municipality           string
	Country                string
	ValidUntil             string
	PermanentlyValid       bool
	CarrierGivenNameLatin  string
	CarrierFamilyNameLatin string
	CarrierGivenName       string
	CarrierFamilyName      string
	CarrierIdNumber        string
	CarrierInsurantNumber  string
	CarrierFamilyMember    bool
	CarrierRelationship    string
	InsuranceBasisRZZO     string
	InsuranceStartDate     string
	InsuranceDescription   string
	TaxpayerName           string
	TaxpayerResidence      string
	TaxpayerNumber         string
	TaxpayerIdNumber       string
	TaxpayerActivityCode   string
}

func (doc *MedicalDocument) GetFullName() string {
	return localization.JoinWithComma(doc.GivenNameLatin, doc.ParentNameLatin, doc.FamilyNameLatin)
}

func (doc *MedicalDocument) GetFullStreetAddress() string {
	var address strings.Builder

	address.WriteString(doc.Street)
	if len(doc.Number) > 0 {
		address.WriteString(", Број: ")
		address.WriteString(doc.Number)
	}

	if len(doc.Apartment) > 0 {
		address.WriteString(" Стан: ")
		address.WriteString(doc.Apartment)
	}

	return address.String()
}

func (doc *MedicalDocument) GetFullPlaceAddress() string {
	return localization.JoinWithComma(doc.Place, doc.Municipality, doc.Country)
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

	putData("Име:", doc.GivenName+" ("+doc.GivenNameLatin+")")

	putData("Име једног родитеља:", doc.ParentName+" ("+doc.ParentNameLatin+")")

	putData("Презиме:", doc.FamilyName+" ("+doc.FamilyNameLatin+")")

	putData("Датум рођења:", doc.DateOfBirth)

	putData("Место, општина и држава:", doc.GetFullPlaceAddress())

	putData("Улица:", doc.GetFullStreetAddress())

	putData("Пол:", doc.Gender)

	putData("Језик:", doc.PrintLanguage)

	putData("ЛБО:", doc.InsurantNumber)

	putData("ЈМБГ:", doc.PersonalNumber)

	section("Подаци о картици здравственог осигурања")

	putData("Датум издавања:", doc.DateOfIssue)

	putData("Датум важења:", doc.DateOfExpiry)

	putData("Оверена до:", doc.ValidUntil)

	putData("Трајно оверена:", localization.FormatYesNo(doc.PermanentlyValid, localization.Cyrillic))

	section("Подаци о носиоцу осигурања")

	putData("Име:", doc.CarrierGivenName+" ("+doc.CarrierGivenNameLatin+")")

	putData("Презиме:", doc.CarrierFamilyName+" ("+doc.CarrierFamilyName+")")

	putData("ЛБО:", doc.CarrierInsurantNumber)

	putData("ЈМБГ:", doc.CarrierIdNumber)

	putData("Члан породице:", localization.FormatYesNo(doc.CarrierFamilyMember, localization.Cyrillic))

	putData("Сродство:", doc.CarrierRelationship)

	section("Подаци о осигурању")

	putData("Основ осигурања:", doc.InsuranceBasisRZZO)

	putData("Датум почетка осигурања:", doc.InsuranceStartDate)

	putData("Опис:", doc.InsuranceDescription)

	section("Подаци о обвезнику плаћања доприноса")

	putData("Назив:", doc.TaxpayerName)

	putData("Седиште:", doc.TaxpayerResidence)

	putData("Регистарски број:", doc.TaxpayerNumber)

	putData("ПИБ/ЈМБГ:", doc.TaxpayerIdNumber)

	putData("Делатност:", doc.TaxpayerActivityCode)

	fileName = strings.ToLower(doc.GivenNameLatin + "_" + doc.FamilyNameLatin + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.GivenNameLatin + " " + doc.FamilyNameLatin,
		Author:       "Baš Čelik",
		Subject:      "Lična karta",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}

func (doc *MedicalDocument) BuildJson() ([]byte, error) {
	return json.Marshal(doc)
}

func (doc *MedicalDocument) UpdateValidUntilDateFromRfzo() error {
	if len([]rune(doc.CardId)) != 11 {
		return ErrInvalidCardNo
	}

	if len([]rune(doc.InsurantNumber)) != 11 {
		return ErrInvalidInsuranceNo
	}

	resp, err := http.PostForm(rfzoServiceUrl, url.Values{"zk": {doc.CardId}, "lbo": {doc.InsurantNumber}})
	if err != nil {
		return fmt.Errorf("posting: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	date, err := ParseValidUntilDateFromRfzoResponse(string(body))
	if err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	doc.ValidUntil = date

	return nil
}

func ParseValidUntilDateFromRfzoResponse(response string) (string, error) {
	regex, err := regexp.Compile(`оверена до: <strong>(\d+\.\d+\.\d+\.)</strong>`)
	if err != nil {
		return "", fmt.Errorf("compiling regex: %w", err)
	}

	matches := regex.FindStringSubmatch(response)
	if len(matches) < 2 {
		return "", ErrNoSubmatchFound
	}

	return matches[1], nil
}
