package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/gui/widgets"
	"github.com/ubavic/bas-celik/localization"
)

func pageID(doc *document.IdDocument) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.GetFullName(), 350)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 100)
	sexF := widgets.NewField("Pol", doc.Sex, 50)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF := widgets.NewField("Mesto rođenja, opština i država", doc.GetFullPlaceOfBirth(), 350)
	addressF := widgets.NewField("Prebivalište i adresa stana", doc.GetFullAddress(), 350)
	addressDateF := widgets.NewField("Datum promene adrese", doc.AddressDate, 10)
	personInformationGroup := widgets.NewGroup("Podaci o građaninu", nameF, birthRow, birthPlaceF, addressF, addressDateF)

	issuedByF := widgets.NewField("Dokument izdaje", doc.IssuingAuthority, 10)
	documentNumberF := widgets.NewField("Broj dokumenta", doc.DocumentNumber, 100)
	issueDateF := widgets.NewField("Datum izdavanja", doc.IssuingDate, 100)
	expiryDateF := widgets.NewField("Važi do", doc.ExpiryDate, 100)
	docRow := container.New(layout.NewHBoxLayout(), documentNumberF, issueDateF, expiryDateF)
	docGroup := widgets.NewGroup("Podaci o dokumentu", issuedByF, docRow)
	colRight := container.New(layout.NewVBoxLayout(), personInformationGroup, docGroup)

	imgWidget := canvas.NewImageFromImage(doc.Portrait)
	imgWidget.SetMinSize(fyne.Size{Width: 200, Height: 250})
	imgWidget.FillMode = canvas.ImageFillContain
	colLeft := container.New(layout.NewVBoxLayout(), imgWidget)

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}

func pageMedical(doc *document.MedicalDocument) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.GetFullName(), 350)
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
	validUntilF := widgets.NewField("Overena do*", doc.ValidUntil, 170)
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

	note := widget.NewLabel("* Datum isteka overe ne mora da bude ažuran na kartici. " +
		"Pritiskom na dugme Ažuriraj, ovaj podatak će se ažurirati sa podatkom dostupnim na web servisu RFZO-a. " +
		"Ova akcija zahteva konekciju sa internetom.")
	note.Wrapping = fyne.TextWrapWord

	return container.New(layout.NewVBoxLayout(), container.New(layout.NewHBoxLayout(), colLeft, colRight), note)
}

func pageVehicle(doc *document.VehicleDocument) *fyne.Container {
	issuingStateF := widgets.NewField("Država izdavanja", doc.StateIssuing, 220)
	issuedByF := widgets.NewField("Dokument izdao", doc.AuthorityIssuing, 220)
	issueRow := container.New(layout.NewHBoxLayout(), issuingStateF, issuedByF)
	issuingDateF := widgets.NewField("Datum izdavanja", doc.IssuingDate, 220)
	expiryDateF := widgets.NewField("Važi do", doc.ExpiryDate, 220)
	dateRow := container.New(layout.NewHBoxLayout(), issuingDateF, expiryDateF)
	competentAuthorityF := widgets.NewField("Zabrana otuđenja", doc.CompetentAuthority, 350)
	docIdF := widgets.NewField("Broj saobraćajne", doc.UnambiguousNumber, 220)
	serialNumberF := widgets.NewField("Serijski broj", doc.SerialNumber, 220)
	idRow := container.New(layout.NewHBoxLayout(), docIdF, serialNumberF)
	documentGroup := widgets.NewGroup("Podaci o dokumentu", issueRow, dateRow, competentAuthorityF, idRow)

	ownerF := widgets.NewField("Vlasnik", doc.OwnersSurnameOrBusinessName, 220)
	ownerNameF := widgets.NewField("Ime vlasnika", doc.OwnerName, 220)
	ownerRow := container.New(layout.NewHBoxLayout(), ownerF, ownerNameF)
	ownerNumberF := widgets.NewField("JMBG vlasnika", doc.OwnersPersonalNo, 350)
	ownerAddressF := widgets.NewField("Adresa vlasnika", doc.OwnerAddress, 350)

	userWidgets := []fyne.CanvasObject{ownerRow, ownerNumberF, ownerAddressF}
	if doc.UsersSurnameOrBusinessName != "" || doc.UsersName != "" || doc.UsersAddress != "" {
		userF := widgets.NewField("Korisnik", doc.UsersSurnameOrBusinessName, 220)
		userNameF := widgets.NewField("Ime korisnika", doc.UsersName, 220)
		userRow := container.New(layout.NewHBoxLayout(), userF, userNameF)
		userNumberF := widgets.NewField("JMBG korisnika", doc.UsersPersonalNo, 350)
		userAddressF := widgets.NewField("Adresa korisnika", doc.UsersAddress, 350)
		userWidgets = append(userWidgets, userRow, userNumberF, userAddressF)
	}

	ownerGroup := widgets.NewGroup("Podaci o vlasniku", userWidgets...)

	colLeft := container.New(layout.NewVBoxLayout(), documentGroup, ownerGroup)

	registrationNumberF := widgets.NewField("Registarski broj", doc.RegistrationNumberOfVehicle, 220)
	dateOfFirstRegistrationF := widgets.NewField("Datum prve registracije", doc.DateOfFirstRegistration, 220)
	vehicleRow0 := container.New(layout.NewHBoxLayout(), registrationNumberF, dateOfFirstRegistrationF)

	brandF := widgets.NewField("Marka", doc.VehicleMake, 220)
	modelF := widgets.NewField("Model", doc.CommercialDescription, 220)
	vehicleRow1 := container.New(layout.NewHBoxLayout(), brandF, modelF)

	colorF := widgets.NewField("Boja", doc.ColourOfVehicle, 220)
	yearOfProductionF := widgets.NewField("Godina proizvodnje", doc.YearOfProduction, 220)
	vehicleRow2 := container.New(layout.NewHBoxLayout(), colorF, yearOfProductionF)

	massF := widgets.NewField("Masa", doc.VehicleMass, 220)
	maximalAllowedMassF := widgets.NewField("Najveća dozvoljena masa", doc.MaximumPermissibleLadenMass, 220)
	vehicleRow3 := container.New(layout.NewHBoxLayout(), massF, maximalAllowedMassF)

	enginePowerF := widgets.NewField("Snaga motora", doc.MaximumNetPower, 220)
	powerMassRatioF := widgets.NewField("Specifična snaga", doc.PowerWeightRatio, 220)
	vehicleRow4 := container.New(layout.NewHBoxLayout(), enginePowerF, powerMassRatioF)

	engineNumberF := widgets.NewField("Broj motora", doc.EngineIdNumber, 220)
	engineCapacityF := widgets.NewField("Kapacitet motora", doc.EngineCapacity, 220)
	vehicleRow5 := container.New(layout.NewHBoxLayout(), engineNumberF, engineCapacityF)

	seatsF := widgets.NewField("Broj mesta za sedenje", doc.NumberOfSeats, 220)
	standingF := widgets.NewField("Broj mesta za stajanje", doc.NumberOfStandingPlaces, 220)
	vehicleRow6 := container.New(layout.NewHBoxLayout(), seatsF, standingF)

	insuranceHolderGroup := widgets.NewGroup("Podaci o vozilu",
		vehicleRow0, vehicleRow1, vehicleRow2, vehicleRow3,
		vehicleRow4, vehicleRow5, vehicleRow6,
	)

	colRight := container.New(layout.NewVBoxLayout(), insuranceHolderGroup)

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}
