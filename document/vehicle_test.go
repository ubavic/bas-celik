package document_test

import (
	"testing"

	"github.com/ubavic/bas-celik/v2/document"
)

var documentVehicle1 = document.VehicleDocument{}
var documentVehicle2 = document.VehicleDocument{
	ColourOfVehicle:             "Bela",
	EngineCapacity:              "1910",
	NumberOfSeats:               "4",
	NumberOfAxles:               "2",
	NumberOfStandingPlaces:      "0",
	OwnerName:                   "Petar",
	OwnersSurnameOrBusinessName: "Petrović",
	OwnerAddress:                "Kralja Aleksandra Karađorđevića, Subotica",
	OwnersPersonalNo:            "1234567890",
	SerialNumber:                "122333444455555",
	VehicleMake:                 "Opel",
	UsersName:                   "Ivana",
	UsersSurnameOrBusinessName:  "Ivanović",
	UsersPersonalNo:             "987654321",
	UsersAddress:                "Kneza Miloša 6, Subotica",
	YearOfProduction:            "2005",
	VehicleMass:                 "980",
	VehicleLoad:                 "0",
	PowerWeightRatio:            "0",
	HomologationMark:            "-",
	VehicleCategory:             "PUTNICKO VOZILO",
	TypeOfFuel:                  "DIZEL",
	StateIssuing:                "Srbija",
	UnambiguousNumber:           "123456789",
	AuthorityIssuing:            "PS SUBOTICA",
	DateOfFirstRegistration:     "09.12.2005",
	EngineIdNumber:              "937A555555555",
	MaximumNetPower:             "85",
	RegistrationNumberOfVehicle: "SU BČ 123",
	MaximumPermissibleLadenMass: "2800",
	IssuingDate:                 "10.02.2020",
	ExpiryDate:                  "10.02.2021",
	CommercialDescription:       "Corsa",
}

func Test_BuildPdfVehicle(t *testing.T) {
	unsetDocumentConfig()

	_, _, err := documentVehicle1.BuildPdf()
	if err == nil {
		t.Errorf("Expected error but got %v", err)
	}

	setDocumentConfigFromLocalFiles(t)

	_, _, err = documentVehicle1.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	_, _, err = documentVehicle2.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
