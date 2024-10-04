package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
	"github.com/ubavic/bas-celik/localization"
)

func pageID(doc *document.IdDocument) *fyne.Container {
	var personalInformationGroupObjects, docGroupObjects []fyne.CanvasObject

	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.GetFullName(), 350)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 100)
	sexF := widgets.NewField("Pol", doc.Sex, 50)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF := widgets.NewField("Mesto rođenja, opština i država", doc.GetFullPlaceOfBirth(), 350)
	addressF := widgets.NewField("Prebivalište i adresa stana", doc.GetFullAddress(), 350)
	addressDateF := widgets.NewField("Datum promene adrese", doc.AddressDate, 10)

	personalInformationGroupObjects = []fyne.CanvasObject{nameF, birthRow, birthPlaceF, addressF, addressDateF}

	if doc.DocumentType == document.ID_TYPE_IDENTITY_FOREIGNER {
		personalInformationGroupObjects = append(personalInformationGroupObjects,
			widgets.NewField("Nacionalnost", doc.NationalityFull, 200),
			widgets.NewField("Status stranca", doc.StatusOfForeigner, 200))
	} else if doc.DocumentType == document.ID_TYPE_RESIDENCE_PERMIT {
		personalInformationGroupObjects = append(personalInformationGroupObjects,
			widgets.NewField("Nacionalnost", doc.NationalityFull, 200),
			widgets.NewField("Osnov boravka", doc.PurposeOfStay, 200),
			widgets.NewField("Napomena", doc.ENote, 200))
	}

	personInformationGroup := widgets.NewGroup("Podaci o građaninu", personalInformationGroupObjects...)

	if doc.DocumentType == document.ID_TYPE_RESIDENCE_PERMIT {
		docGroupObjects = append(docGroupObjects, widgets.NewField("Naziv dokumenta", doc.DocumentName, 100))
	}
	docGroupObjects = append(docGroupObjects, widgets.NewField("Dokument izdaje", doc.IssuingAuthority, 10))
	documentNumberF := widgets.NewField("Broj dokumenta", doc.DocRegNo, 100)
	issueDateF := widgets.NewField("Datum izdavanja", doc.IssuingDate, 100)
	expiryDateF := widgets.NewField("Važi do", doc.ExpiryDate, 100)
	docRow := container.New(layout.NewHBoxLayout(), documentNumberF, issueDateF, expiryDateF)
	docGroupObjects = append(docGroupObjects, docRow)
	docGroup := widgets.NewGroup("Podaci o dokumentu", docGroupObjects...)

	colRight := container.New(layout.NewVBoxLayout(), personInformationGroup, docGroup)

	imgWidget := canvas.NewImageFromImage(doc.Portrait)
	imgWidget.SetMinSize(fyne.Size{Width: 200, Height: 250})
	imgWidget.FillMode = canvas.ImageFillContain
	colLeft := container.New(layout.NewVBoxLayout(), imgWidget)

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}

func pageMedical(doc *document.MedicalDocument) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.GetFullName(), 350)
	genderF := widgets.NewField("Pol", doc.Gender, 170)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 170)
	birthRow := container.New(layout.NewHBoxLayout(), genderF, birthDateF)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 170)
	insurantNumberF := widgets.NewField("LBO", doc.InsurantNumber, 170)
	idsRow := container.New(layout.NewHBoxLayout(), personalNumberF, insurantNumberF)
	countryF := widgets.NewField("Država", doc.Country, 170)
	municipalityF := widgets.NewField("Opština", doc.Municipality, 170)
	address1Row := container.New(layout.NewHBoxLayout(), countryF, municipalityF)
	placeF := widgets.NewField("Mesto", doc.Place, 170)
	street := widgets.NewField("Ulica", doc.Street, 170)
	address2Row := container.New(layout.NewHBoxLayout(), placeF, street)
	addressNumber := widgets.NewField("Broj", doc.Number, 76)
	addressEntrance := widgets.NewField("Ulaz", doc.Apartment, 76)
	addressApartmentNumber := widgets.NewField("Stan", doc.Apartment, 80)
	address3Row := container.New(layout.NewHBoxLayout(), addressNumber, addressEntrance, addressApartmentNumber)

	generalGroup := widgets.NewGroup("Opšti podaci", nameF, birthRow, idsRow, address1Row, address2Row, address3Row)

	insuranceBasisF := widgets.NewField("Osnov osiguranja", doc.InsuranceBasisRZZO, 170)
	insuranceStartDateF := widgets.NewField("Datum početka osiguranja", doc.InsuranceStartDate, 170)
	insuranceRow := container.New(layout.NewHBoxLayout(), insuranceBasisF, insuranceStartDateF)
	insuranceDescriptionF := widgets.NewField("Opis", doc.InsuranceDescription, 350)
	insuranceGroup := widgets.NewGroup("Podaci o osiguranju", insuranceRow, insuranceDescriptionF)

	colLeft := container.New(layout.NewVBoxLayout(), generalGroup, insuranceGroup)

	carrierNameF := widgets.NewField("Nosilac osiguranja", doc.CarrierGivenNameLatin+" "+doc.CarrierFamilyNameLatin, 350)
	carrierInsurantNumberF := widgets.NewField("LBO", doc.CarrierInsurantNumber, 170)
	carrierIdNumberF := widgets.NewField("JMBG", doc.CarrierIdNumber, 170)
	carrierRow1 := container.New(layout.NewHBoxLayout(), carrierInsurantNumberF, carrierIdNumberF)
	carrierFamilyMemberF := widgets.NewField("Član porodice", localization.FormatYesNo(doc.CarrierFamilyMember, localization.Latin), 170)
	carrierRelationshipF := widgets.NewField("Srodstvo", doc.CarrierRelationship, 170)
	carrierRow2 := container.New(layout.NewHBoxLayout(), carrierFamilyMemberF, carrierRelationshipF)
	carrierGroup := widgets.NewGroup("Podaci o nosiocu osiguranja", carrierNameF, carrierRow1, carrierRow2)

	cardNumber := widgets.NewField("Broj zdravstvene isprave", doc.CardId, 270)
	cardRow0 := container.New(layout.NewHBoxLayout(), cardNumber)
	dateOfIssueF := widgets.NewField("Datum izdavanja", doc.DateOfIssue, 170)
	dateOfExpiryF := widgets.NewField("Datum važenja", doc.DateOfExpiry, 170)
	cardRow1 := container.New(layout.NewHBoxLayout(), dateOfIssueF, dateOfExpiryF)
	validUntilF := widgets.NewField("Overena do*", doc.ValidUntil, 170)
	permanentlyValidF := widgets.NewField("Trajna overa", localization.FormatYesNo(doc.PermanentlyValid, localization.Latin), 170)
	cardRow2 := container.New(layout.NewHBoxLayout(), validUntilF, permanentlyValidF)
	cardGroup := widgets.NewGroup("Podaci o kartici", cardRow0, cardRow1, cardRow2)

	taxpayerNameF := widgets.NewField("Naziv / ime i prezime", doc.TaxpayerName, 350)
	taxpayerActivityCodeF := widgets.NewField("Oznaka delatnosti", doc.TaxpayerActivityCode, 170)
	taxpayerPlaceF := widgets.NewField("Sedište", doc.TaxpayerResidence, 170)
	taxpayerRow1 := container.New(layout.NewHBoxLayout(), taxpayerActivityCodeF, taxpayerPlaceF)
	taxpayerNumberF := widgets.NewField("Registarski broj", doc.TaxpayerNumber, 170)
	taxpayerIdNumberF := widgets.NewField("PIB/JMBG", doc.TaxpayerIdNumber, 170)
	taxpayerRow2 := container.New(layout.NewHBoxLayout(), taxpayerNumberF, taxpayerIdNumberF)
	taxpayerGroup := widgets.NewGroup("Podaci o obavezniku", taxpayerNameF, taxpayerRow1, taxpayerRow2)

	colRight := container.New(layout.NewVBoxLayout(), carrierGroup, cardGroup, taxpayerGroup)

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
