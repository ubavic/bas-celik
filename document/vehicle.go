package document

import (
	"encoding/json"
	"errors"
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

type VehicleDocument struct {
	AuthorityIssuing            string
	ColourOfVehicle             string
	CommercialDescription       string
	CompetentAuthority          string
	DateOfFirstRegistration     string
	EngineCapacity              string
	EngineIdNumber              string
	EngineNumber                string
	EnginePower                 string
	EngineRatedSpeed            string
	ExpiryDate                  string
	HomologationMark            string
	IssuingDate                 string
	MaximumNetPower             string
	MaximumPermissibleLadenMass string
	NumberOfAxles               string
	NumberOfSeats               string
	NumberOfStandingPlaces      string
	OwnerAddress                string
	OwnerName                   string
	OwnersPersonalNo            string
	OwnersSurnameOrBusinessName string
	PowerWeightRatio            string
	RegistrationNumberOfVehicle string
	RestrictionToChangeOwner    string
	SerialNumber                string
	StateIssuing                string
	TypeApprovalNumber          string
	TypeOfFuel                  string
	UnambiguousNumber           string
	UsersAddress                string
	UsersName                   string
	UsersPersonalNo             string
	UsersSurnameOrBusinessName  string
	VehicleCategory             string
	VehicleIdNumber             string
	VehicleLoad                 string
	VehicleMake                 string
	VehicleMass                 string
	VehicleType                 string
	YearOfProduction            string
}

func (doc VehicleDocument) BuildUI(pdfHandler func(), statusBar *widgets.StatusBar) *fyne.Container {
	issuingStateF := widgets.NewField("Država izdavanja", doc.StateIssuing, 350)
	issuedByF := widgets.NewField("Dokument izdao", doc.AuthorityIssuing, 350)
	issuingDateF := widgets.NewField("Datum izdavanja", doc.IssuingDate, 170)
	expiryDateF := widgets.NewField("Važi do", doc.ExpiryDate, 170)
	dateRow := container.New(layout.NewHBoxLayout(), issuingDateF, expiryDateF)
	denialAuthorityF := widgets.NewField("Zabrana otuđenja", doc.RestrictionToChangeOwner, 350)
	docIdF := widgets.NewField("Broj dokumenta", doc.UnambiguousNumber, 350)
	serialNumberF := widgets.NewField("Serijski broj", doc.SerialNumber, 350)
	documentGroup := widgets.NewGroup("Podaci o dokumentu", issuingStateF, issuedByF, dateRow, denialAuthorityF, docIdF, serialNumberF)

	ownerF := widgets.NewField("Vlasnik", doc.OwnersSurnameOrBusinessName, 350)
	ownerNameF := widgets.NewField("Ime vlasnika", doc.OwnerName, 350)
	ownerNumberF := widgets.NewField("JMBG vlasnika", doc.OwnersPersonalNo, 350)
	ownerAddressF := widgets.NewField("Adresa vlasnika", doc.OwnerAddress, 350)
	userF := widgets.NewField("Korisnik", doc.UsersSurnameOrBusinessName, 350)
	userNameF := widgets.NewField("Ime korisnika", doc.UsersName, 350)
	userNumberF := widgets.NewField("JMBG korisnika", doc.UsersPersonalNo, 350)
	userAddressF := widgets.NewField("Adresa korisnika", doc.UsersAddress, 350)
	ownerGroup := widgets.NewGroup("Podaci o vlasniku", ownerF, ownerNameF, ownerNumberF, ownerAddressF, userF, userNameF, userNumberF, userAddressF)

	colLeft := container.New(layout.NewVBoxLayout(), documentGroup, ownerGroup)

	registrationNumberF := widgets.NewField("Registarski broj", doc.RegistrationNumberOfVehicle, 350)
	dateOfFirstRegistrationF := widgets.NewField("Datum prve registracije", doc.DateOfFirstRegistration, 170)

	brandF := widgets.NewField("Marka", doc.VehicleMake, 170)
	modelF := widgets.NewField("Model", "", 170)
	vehicleRow1 := container.New(layout.NewHBoxLayout(), brandF, modelF)

	colorF := widgets.NewField("Boja", doc.ColourOfVehicle, 170)
	yearOfProductionF := widgets.NewField("Godina proizvodnje", doc.YearOfProduction, 170)
	vehicleRow2 := container.New(layout.NewHBoxLayout(), colorF, yearOfProductionF)

	massF := widgets.NewField("Masa", doc.VehicleMass, 170)
	maximalAllowedMassF := widgets.NewField("Najveća dozvoljena masa", doc.MaximumPermissibleLadenMass, 170)
	vehicleRow3 := container.New(layout.NewHBoxLayout(), massF, maximalAllowedMassF)

	enginePowerF := widgets.NewField("Snaga motora", doc.EnginePower, 170)
	powerMassRatioF := widgets.NewField("Specifična snaga", doc.PowerWeightRatio, 170)
	vehicleRow4 := container.New(layout.NewHBoxLayout(), enginePowerF, powerMassRatioF)

	engineNumberF := widgets.NewField("Broj motora", doc.EngineNumber, 350)
	engineCapacityF := widgets.NewField("Kapacitet motora", doc.EngineCapacity, 350)
	vehicleRow5 := container.New(layout.NewHBoxLayout(), engineNumberF, engineCapacityF)

	seatsF := widgets.NewField("Broj mesta za sedenje", doc.NumberOfSeats, 170)
	standingF := widgets.NewField("Broj mesta za stajanje", doc.NumberOfStandingPlaces, 170)
	vehicleRow6 := container.New(layout.NewHBoxLayout(), seatsF, standingF)

	insuranceHolderGroup := widgets.NewGroup("Podaci o vozilu",
		registrationNumberF, dateOfFirstRegistrationF,
		vehicleRow1, vehicleRow2, vehicleRow3, vehicleRow4,
		vehicleRow5, vehicleRow6,
	)

	colRight := container.New(layout.NewVBoxLayout(), insuranceHolderGroup)

	cols := container.New(layout.NewHBoxLayout(), colLeft, colRight)

	saveButton := widget.NewButton("Sačuvaj PDF", pdfHandler)
	buttonBar := container.New(layout.NewHBoxLayout(), statusBar, layout.NewSpacer(), saveButton)

	return container.New(layout.NewVBoxLayout(), cols, buttonBar)
}

func (doc *VehicleDocument) BuildPdf() (data []byte, fileName string, retErr error) {
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

	err = pdf.AddTTFFontDataWithOption("liberationsans", fontBold, gopdf.TtfOption{Style: gopdf.Bold})
	if err != nil {
		panic(fmt.Errorf("loading font: %w", err))
	}

	const leftMargin = 28
	const rightMargin = 535
	const textLeftMargin = 38

	underlineOption := gopdf.CellOption{CoefUnderlineThickness: 1, CoefUnderlinePosition: -2}

	newLine := func() {
		pdf.SetXY(textLeftMargin, pdf.GetY()+20)
	}

	tab := func() {
		pdf.SetXY(343, pdf.GetY())
	}

	dashFormat := func(str string) string {
		if len(str) == 0 {
			return "-"
		} else {
			return str
		}
	}

	cell := func(s string) {
		err := pdf.Cell(nil, s)
		if err != nil {
			panic(fmt.Errorf("putting text: %w", err))
		}
	}

	putData := func(label, data string) {
		cell(label + ": " + data)
	}

	putUnderline := func(str string, size int) {
		err = pdf.SetFont("liberationsans", "U", size)
		if err != nil {
			panic(fmt.Errorf("setting font: %w", err))
		}
		err = pdf.CellWithOption(nil, str, underlineOption)
		if err != nil {
			panic(fmt.Errorf("cell: %w", err))
		}
		err = pdf.SetFont("liberationsans", "B", 12)
		if err != nil {
			panic(fmt.Errorf("setting font: %w", err))
		}
	}

	putParagraph := func(data string) {

		texts := strings.Split(data, ",")
		if len(texts) == 2 {
			cell(texts[0])
			pdf.SetXY(textLeftMargin, pdf.GetY()+14)
			cell(strings.TrimSpace(texts[1]))
			pdf.SetXY(textLeftMargin, pdf.GetY()+20)
			return
		}

		texts, err = pdf.SplitTextWithWordWrap(data, 500)
		if err != nil {
			panic(fmt.Errorf("splitting text: %w", err))
		}

		for i, text := range texts {
			cell(text)
			if i < len(texts)-1 {
				pdf.SetXY(textLeftMargin, pdf.GetY()+14)
			}
		}

		pdf.SetXY(textLeftMargin, pdf.GetY()+20)
	}

	err = pdf.SetFont("liberationsans", "B", 29)
	if err != nil {
		panic(fmt.Errorf("setting font: %w", err))
	}
	pdf.SetXY(textLeftMargin, 35)
	cell("Čitač saobraćajne dozvole")

	pdf.SetLineWidth(2.9)
	pdf.SetLineType("solid")
	pdf.Line(leftMargin, 72, rightMargin, 72)

	pdf.SetXY(textLeftMargin, 90)
	err = pdf.SetFontSize(21)
	if err != nil {
		panic(fmt.Errorf("setting font size: %w", err))
	}
	cell("Registarska oznaka: " + doc.RegistrationNumberOfVehicle)

	pdf.SetXY(textLeftMargin, 145)
	err = pdf.SetFontSize(12)
	if err != nil {
		panic(fmt.Errorf("setting font: %w", err))
	}

	putData("Datum izdavanja", doc.IssuingDate)
	tab()
	putUnderline("Važi do"+doc.ExpiryDate+": ", 12)
	newLine()

	putData("Saobraćajnu izdao", doc.AuthorityIssuing)
	tab()
	putData("Zabrana otuđenja", "")
	newLine()
	putParagraph(doc.RestrictionToChangeOwner)

	putData("Broj saobraćajne", doc.UnambiguousNumber)
	newLine()

	putData("Serijski broj", doc.SerialNumber)
	newLine()

	pdf.SetXY(textLeftMargin, 272)
	putUnderline("Podaci o vlasniku", 20)
	pdf.SetXY(textLeftMargin, pdf.GetY()+25)

	putData("Vlasnik", doc.OwnersSurnameOrBusinessName)
	newLine()

	putData("Ime vlasnika", doc.OwnerName)
	newLine()

	putData("Adresa vlasnika", doc.OwnerAddress)
	newLine()

	putData("Jmbg vlasnika", doc.OwnersPersonalNo)
	newLine()

	putData("Korisnik", doc.UsersSurnameOrBusinessName)
	newLine()

	putData("Ime korisnika", doc.UsersName)
	newLine()

	putData("Adresa korisnika", doc.UsersAddress)
	newLine()

	putData("Jmbg korisnika", doc.UsersAddress)
	newLine()

	pdf.SetXY(textLeftMargin, pdf.GetY()+6)
	putUnderline("Podaci o vozilu", 20)
	pdf.SetXY(textLeftMargin, pdf.GetY()+25)

	putData("Datum prve registracije", doc.DateOfFirstRegistration)
	tab()
	putData("Godina proizvodnje", doc.YearOfProduction)
	newLine()

	putData("Marka", doc.VehicleMake)
	tab()
	putData("Model", "")
	newLine()

	putData("Tip", dashFormat(doc.VehicleType))
	pdf.SetXY(pdf.GetX()+10, pdf.GetY())
	putData("Homologacijska oznaka", dashFormat(doc.HomologationMark))
	newLine()

	putData("Boja", doc.ColourOfVehicle)
	tab()
	putData("Broj osovina", doc.NumberOfAxles)
	newLine()

	putData("Broj šasije", doc.EngineNumber)
	tab()
	putData("Zapremina motora", doc.EngineCapacity)
	newLine()

	putData("Broj motora", doc.EngineNumber)
	tab()
	putData("Masa", doc.VehicleMass)
	newLine()

	putData("Snaga motora", doc.EnginePower)
	tab()
	putData("Nosivost", doc.VehicleLoad)
	newLine()

	putData("Odnos snaga/masa", doc.PowerWeightRatio)
	tab()
	cell("Najveća dozvoljena")
	newLine()

	putData("Kategorija", doc.VehicleCategory)
	tab()
	putData("masa", doc.MaximumPermissibleLadenMass)
	newLine()

	putData("Pogonsko gorivo", doc.TypeOfFuel)
	newLine()

	putData("Broj mesta za sedenje", doc.NumberOfSeats)
	tab()
	putData("Broj mesta za stajanje", doc.NumberOfStandingPlaces)
	newLine()

	fileName = strings.ToLower(doc.RegistrationNumberOfVehicle + "_" + doc.OwnersSurnameOrBusinessName + "_" + doc.OwnerName + ".pdf")

	pdf.SetInfo(gopdf.PdfInfo{
		Title:        doc.VehicleMake + " " + doc.CommercialDescription,
		Author:       "Baš Čelik",
		Subject:      "Saobraćajna dozvola",
		CreationDate: time.Now(),
	})

	return pdf.GetBytesPdf(), fileName, nil
}

func (doc *VehicleDocument) BuildJson() ([]byte, error) {
	return json.Marshal(doc)
}
