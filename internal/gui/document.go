package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal/gui/translation"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
	"github.com/ubavic/bas-celik/localization"
)

func pageID(doc *document.IdDocument) *fyne.Container {
	var personalInformationGroupObjects, docGroupObjects []fyne.CanvasObject

	widthThird := (350 - 2*theme.Padding()) / 3

	nameF := widgets.NewField(t("id.name"), doc.GetFullName(), 350)
	birthDateF := widgets.NewField(t("id.birthDate"), doc.DateOfBirth, widthThird)
	sexF := widgets.NewField(t("id.sex"), doc.Sex, widthThird)
	personalNumberF := widgets.NewField(t("id.personalNumber"), doc.PersonalNumber, widthThird)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF := widgets.NewField(t("id.birthPlace"), doc.GetFullPlaceOfBirth(), 350)
	addressF := widgets.NewField(t("id.address"), doc.GetFullAddress(false), 350)
	addressDateF := widgets.NewField(t("id.addressDate"), doc.AddressDate, 10)

	personalInformationGroupObjects = []fyne.CanvasObject{nameF, birthRow, birthPlaceF, addressF, addressDateF}

	nationalityLabel := t("id.nationalityFull")
	foreignerStatusLabel := t("id.foreignerStatus")
	purposeOfStayLabel := t("id.purposeOfStay")
	eNoteLabel := t("id.eNote")

	if doc.DocumentType == document.ID_TYPE_IDENTITY_FOREIGNER {
		personalInformationGroupObjects = append(personalInformationGroupObjects,
			widgets.NewField(nationalityLabel, doc.NationalityFull, 200),
			widgets.NewField(foreignerStatusLabel, doc.StatusOfForeigner, 200))
	} else if doc.DocumentType == document.ID_TYPE_RESIDENCE_PERMIT {
		personalInformationGroupObjects = append(personalInformationGroupObjects,
			widgets.NewField(nationalityLabel, doc.NationalityFull, 200),
			widgets.NewField(purposeOfStayLabel, doc.PurposeOfStay, 200),
			widgets.NewField(eNoteLabel, doc.ENote, 200))
	}

	personInformationGroup := widgets.NewGroup(t("id.citizenInformation"), personalInformationGroupObjects...)

	if doc.DocumentType == document.ID_TYPE_RESIDENCE_PERMIT {
		docGroupObjects = append(docGroupObjects, widgets.NewField(t("id.documentName"), doc.DocumentName, 100))
	}

	docGroupObjects = append(docGroupObjects, widgets.NewField(t("id.issuingAuthority"), doc.IssuingAuthority, 350))
	documentNumberF := widgets.NewField(t("id.docRegNo"), doc.DocRegNo, widthThird)
	issueDateF := widgets.NewField(t("id.issuingDate"), doc.IssuingDate, widthThird)
	expiryDateF := widgets.NewField(t("id.expiryDate"), doc.ExpiryDate, widthThird)
	docRow := container.New(layout.NewHBoxLayout(), documentNumberF, issueDateF, expiryDateF)
	docGroupObjects = append(docGroupObjects, docRow)

	docGroup := widgets.NewGroup(t("id.documentInformation"), docGroupObjects...)

	colRight := container.New(layout.NewVBoxLayout(), personInformationGroup, docGroup)

	imgWidget := canvas.NewImageFromImage(doc.Portrait)
	imgWidget.SetMinSize(fyne.Size{Width: 200, Height: 250})
	imgWidget.FillMode = canvas.ImageFillContain
	colLeft := container.New(layout.NewVBoxLayout(), imgWidget)

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}

func pageMedical(doc *document.MedicalDocument) *fyne.Container {
	nameF := widgets.NewField(t("medical.fullName"), doc.GetFullName(), 350)
	genderF := widgets.NewField(t("medical.gender"), doc.Gender, 170)
	birthDateF := widgets.NewField(t("medical.dateOfBirth"), doc.DateOfBirth, 170)
	birthRow := container.New(layout.NewHBoxLayout(), genderF, birthDateF)

	personalNumberF := widgets.NewField(t("medical.personalNumber"), doc.PersonalNumber, 170)
	insurantNumberF := widgets.NewField(t("medical.insuranceNumber"), doc.InsurantNumber, 170)
	idsRow := container.New(layout.NewHBoxLayout(), personalNumberF, insurantNumberF)

	countryF := widgets.NewField(t("medical.country"), doc.Country, 170)
	municipalityF := widgets.NewField(t("medical.municipality"), doc.Municipality, 170)
	address1Row := container.New(layout.NewHBoxLayout(), countryF, municipalityF)

	placeF := widgets.NewField(t("medical.place"), doc.Place, 170)
	street := widgets.NewField(t("medical.street"), doc.Street, 170)
	address2Row := container.New(layout.NewHBoxLayout(), placeF, street)

	halfWidth := (170 - 3*theme.Padding()) / 2

	addressNumber := widgets.NewField(t("medical.number"), doc.Number, halfWidth)
	addressEntrance := widgets.NewField(t("medical.entrance"), doc.Apartment, halfWidth)
	addressApartmentNumber := widgets.NewField(t("medical.apartment"), doc.Apartment, 170)
	address3Row := container.New(layout.NewHBoxLayout(), addressNumber, addressEntrance, addressApartmentNumber)

	generalGroup := widgets.NewGroup(t("medical.generalInformation"), nameF, birthRow, idsRow, address1Row, address2Row, address3Row)

	insuranceBasisF := widgets.NewField(t("medical.insuranceBasis"), doc.InsuranceBasisRZZO, 170)
	insuranceStartDateF := widgets.NewField(t("medical.insuranceStartDate"), doc.InsuranceStartDate, 170)
	insuranceRow := container.New(layout.NewHBoxLayout(), insuranceBasisF, insuranceStartDateF)

	insuranceDescriptionF := widgets.NewField(t("medical.insuranceDescription"), doc.InsuranceDescription, 350)
	insuranceGroup := widgets.NewGroup(t("medical.insuranceInformation"), insuranceRow, insuranceDescriptionF)

	colLeft := container.New(layout.NewVBoxLayout(), generalGroup, insuranceGroup)

	carrierNameF := widgets.NewField(t("medical.carrier"), doc.CarrierGivenNameLatin+" "+doc.CarrierFamilyNameLatin, 350)
	carrierInsurantNumberF := widgets.NewField(t("medical.insuranceNumber"), doc.CarrierInsurantNumber, 170)
	carrierIdNumberF := widgets.NewField(t("medical.personalNumber"), doc.CarrierIdNumber, 170)
	carrierRow1 := container.New(layout.NewHBoxLayout(), carrierInsurantNumberF, carrierIdNumberF)

	carrierFamilyMemberF := widgets.NewField(t("medical.familyMember"), localization.FormatYesNo(doc.CarrierFamilyMember, translation.CurrentLanguage()), 170)
	carrierRelationshipF := widgets.NewField(t("medical.relationship"), doc.CarrierRelationship, 170)
	carrierRow2 := container.New(layout.NewHBoxLayout(), carrierFamilyMemberF, carrierRelationshipF)
	carrierGroup := widgets.NewGroup(t("medical.insuranceCarrierInformation"), carrierNameF, carrierRow1, carrierRow2)
	cardNumber := widgets.NewField(t("medical.cardId"), doc.CardId, 270)
	dateOfIssueF := widgets.NewField(t("medical.dateOfIssue"), doc.DateOfIssue, 170)
	dateOfExpiryF := widgets.NewField(t("medical.dateOfExpiry"), doc.DateOfExpiry, 170)
	cardRow1 := container.New(layout.NewHBoxLayout(), dateOfIssueF, dateOfExpiryF)

	validUntilF := widgets.NewField(t("medical.validUntil"), doc.ValidUntil, 170)
	permanentlyValidF := widgets.NewField(t("medical.permanentlyValid"), localization.FormatYesNo(doc.PermanentlyValid, translation.CurrentLanguage()), 170)
	cardRow2 := container.New(layout.NewHBoxLayout(), validUntilF, permanentlyValidF)

	cardGroup := widgets.NewGroup(t("medical.cardInformation"), cardNumber, cardRow1, cardRow2)

	taxpayerNameF := widgets.NewField(t("medical.taxpayerName"), doc.TaxpayerName, 350)
	taxpayerActivityCodeF := widgets.NewField(t("medical.taxpayerActivityCode"), doc.TaxpayerActivityCode, 170)
	taxpayerPlaceF := widgets.NewField(t("medical.taxpayerResidence"), doc.TaxpayerResidence, 170)
	taxpayerRow1 := container.New(layout.NewHBoxLayout(), taxpayerActivityCodeF, taxpayerPlaceF)
	taxpayerNumberF := widgets.NewField(t("medical.taxpayerNumber"), doc.TaxpayerNumber, 170)
	taxpayerIdNumberF := widgets.NewField(t("medical.taxpayerIdNumber"), doc.TaxpayerIdNumber, 170)
	taxpayerRow2 := container.New(layout.NewHBoxLayout(), taxpayerNumberF, taxpayerIdNumberF)

	taxpayerGroup := widgets.NewGroup(t("medical.taxpayerInformation"), taxpayerNameF, taxpayerRow1, taxpayerRow2)

	colRight := container.New(layout.NewVBoxLayout(), carrierGroup, cardGroup, taxpayerGroup)

	note := widget.NewLabel(t("medical.dateNote"))
	note.Wrapping = fyne.TextWrapWord

	return container.New(layout.NewVBoxLayout(), container.New(layout.NewHBoxLayout(), colLeft, colRight), note)
}

func pageVehicle(doc *document.VehicleDocument) *fyne.Container {
	issuingStateF := widgets.NewField(t("vehicle.stateIssuing"), doc.StateIssuing, 220)
	issuedByF := widgets.NewField(t("vehicle.authorityIssuing"), doc.AuthorityIssuing, 220)
	issueRow := container.New(layout.NewHBoxLayout(), issuingStateF, issuedByF)
	issuingDateF := widgets.NewField(t("vehicle.issuingDate"), doc.IssuingDate, 220)
	expiryDateF := widgets.NewField(t("vehicle.expiryDate"), doc.ExpiryDate, 220)
	dateRow := container.New(layout.NewHBoxLayout(), issuingDateF, expiryDateF)
	competentAuthorityF := widgets.NewField(t("vehicle.competentAuthority"), doc.CompetentAuthority, 350)
	docIdF := widgets.NewField(t("vehicle.unambiguousNumber"), doc.UnambiguousNumber, 220)
	serialNumberF := widgets.NewField(t("vehicle.serialNumber"), doc.SerialNumber, 220)
	idRow := container.New(layout.NewHBoxLayout(), docIdF, serialNumberF)
	documentGroup := widgets.NewGroup(t("vehicle.documentInformation"), issueRow, dateRow, competentAuthorityF, idRow)

	ownerNoLbl := ""
	if len(doc.OwnersPersonalNo) > 9 {
		ownerNoLbl = t("vehicle.ownersPersonalNo")
	} else {
		ownerNoLbl = t("vehicle.ownersCompanyNo")
	}

	ownerF := widgets.NewField(t("vehicle.ownersSurnameOrBusinessName"), doc.OwnersSurnameOrBusinessName, 220)
	ownerNameF := widgets.NewField(t("vehicle.ownerName"), doc.OwnerName, 220)
	ownerRow := container.New(layout.NewHBoxLayout(), ownerF, ownerNameF)
	ownerNumberF := widgets.NewField(ownerNoLbl, doc.OwnersPersonalNo, 350)
	ownerAddressF := widgets.NewField(t("vehicle.ownerAddress"), doc.OwnerAddress, 350)

	userWidgets := []fyne.CanvasObject{ownerRow, ownerNumberF, ownerAddressF}
	if doc.UsersSurnameOrBusinessName != "" || doc.UsersName != "" || doc.UsersAddress != "" {
		userNoLbl := ""
		if len(doc.OwnersPersonalNo) > 9 {
			userNoLbl = t("vehicle.usersPersonalNo")
		} else {
			userNoLbl = t("vehicle.usersCompanyNo")
		}

		userF := widgets.NewField(t("vehicle.usersSurnameOrBusinessName"), doc.UsersSurnameOrBusinessName, 220)
		userNameF := widgets.NewField(t("vehicle.usersName"), doc.UsersName, 220)
		userRow := container.New(layout.NewHBoxLayout(), userF, userNameF)
		userNumberF := widgets.NewField(userNoLbl, doc.UsersPersonalNo, 350)
		userAddressF := widgets.NewField(t("vehicle.usersAddress"), doc.UsersAddress, 350)
		userWidgets = append(userWidgets, userRow, userNumberF, userAddressF)
	}

	ownerGroup := widgets.NewGroup(t("vehicle.ownerInformation"), userWidgets...)

	colLeft := container.New(layout.NewVBoxLayout(), documentGroup, ownerGroup)

	registrationNumberF := widgets.NewField(t("vehicle.registrationNumberOfVehicle"), doc.RegistrationNumberOfVehicle, 220)
	dateOfFirstRegistrationF := widgets.NewField(t("vehicle.dateOfFirstRegistration"), doc.DateOfFirstRegistration, 220)
	vehicleRow0 := container.New(layout.NewHBoxLayout(), registrationNumberF, dateOfFirstRegistrationF)

	brandF := widgets.NewField(t("vehicle.vehicleMake"), doc.VehicleMake, 220)
	modelF := widgets.NewField(t("vehicle.commercialDescription"), doc.CommercialDescription, 220)
	vehicleRow1 := container.New(layout.NewHBoxLayout(), brandF, modelF)

	colorF := widgets.NewField(t("vehicle.colourOfVehicle"), doc.ColourOfVehicle, 220)
	yearOfProductionF := widgets.NewField(t("vehicle.yearOfProduction"), doc.YearOfProduction, 220)
	vehicleRow2 := container.New(layout.NewHBoxLayout(), colorF, yearOfProductionF)

	massF := widgets.NewField(t("vehicle.vehicleMass"), doc.VehicleMass, 220)
	maximalAllowedMassF := widgets.NewField(t("vehicle.maximumPermissibleLadenMass"), doc.MaximumPermissibleLadenMass, 220)
	vehicleRow3 := container.New(layout.NewHBoxLayout(), massF, maximalAllowedMassF)

	enginePowerF := widgets.NewField(t("vehicle.maximumNetPower"), doc.MaximumNetPower, 220)
	powerMassRatioF := widgets.NewField(t("vehicle.powerWeightRatio"), doc.PowerWeightRatio, 220)
	vehicleRow4 := container.New(layout.NewHBoxLayout(), enginePowerF, powerMassRatioF)

	engineNumberF := widgets.NewField(t("vehicle.engineIdNumber"), doc.EngineIdNumber, 220)
	engineCapacityF := widgets.NewField(t("vehicle.engineCapacity"), doc.EngineCapacity, 220)
	vehicleRow5 := container.New(layout.NewHBoxLayout(), engineNumberF, engineCapacityF)

	seatsF := widgets.NewField(t("vehicle.numberOfSeats"), doc.NumberOfSeats, 220)
	standingF := widgets.NewField(t("vehicle.numberOfStandingPlaces"), doc.NumberOfStandingPlaces, 220)
	vehicleRow6 := container.New(layout.NewHBoxLayout(), seatsF, standingF)

	insuranceHolderGroup := widgets.NewGroup(t("vehicle.vehicleInformation"),
		vehicleRow0, vehicleRow1, vehicleRow2, vehicleRow3,
		vehicleRow4, vehicleRow5, vehicleRow6,
	)

	colRight := container.New(layout.NewVBoxLayout(), insuranceHolderGroup)

	return container.New(layout.NewHBoxLayout(), colLeft, colRight)
}
