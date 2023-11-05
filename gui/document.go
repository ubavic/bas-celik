package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/gui/widgets"
	"github.com/ubavic/bas-celik/localization"
)

func pageID(doc *document.IdDocument) *fyne.Container {
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.FormatName(), 350)
	birthDateF := widgets.NewField("Datum rođenja", doc.DateOfBirth, 100)
	sexF := widgets.NewField("Pol", doc.Sex, 50)
	personalNumberF := widgets.NewField("JMBG", doc.PersonalNumber, 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF := widgets.NewField("Mesto rođenja, opština i država", doc.FormatPlaceOfBirth(), 350)
	addressF := widgets.NewField("Prebivalište i adresa stana", doc.FormatAddress(), 350)
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
	nameF := widgets.NewField("Ime, ime roditelja, prezime", doc.FormatName(), 350)
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

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}

func pageVehicle(doc *document.VehicleDocument) *fyne.Container {
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

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}
