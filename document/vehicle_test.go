package document_test

import (
	"testing"

	"github.com/ubavic/bas-celik/document"
)

var documentVehicle1 = document.VehicleDocument{}
var documentVehicle2 = document.VehicleDocument{
	ColourOfVehicle:             "Bela",
	EngineCapacity:              "2000",
	NumberOfSeats:               "4",
	NumberOfAxles:               "2",
	NumberOfStandingPlaces:      "0",
	OwnerName:                   "Petar",
	OwnersSurnameOrBusinessName: "Petrović",
	OwnerAddress:                "Kralja Aleksandra Karađorđevića, Beograd",
	OwnersPersonalNo:            "1234567890",
	SerialNumber:                "122333444455555",
	VehicleMake:                 "Opel",
	UsersName:                   "Ivana",
	UsersSurnameOrBusinessName:  "Ivanović",
	UsersPersonalNo:             "987654321",
	UsersAddress:                "Kneza Miloša, Beograd",
	YearOfProduction:            "2005",
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
