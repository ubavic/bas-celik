package document_test

import (
	"testing"

	"github.com/ubavic/bas-celik/document"
)

var documentMedical1 = document.MedicalDocument{}
var documentMedical2 = document.MedicalDocument{
	GivenName:                  "Петар",
	ParentName:                 "Арсеније",
	Surname:                    "Петровић",
	AddressStreet:              "Његошева",
	AddressNumber:              "9",
	AddressApartmentNumber:     "5",
	AddressTown:                "Подгорица",
	AddressState:               "Црна Гора",
	PersonalNumber:             "01019950900200",
	DateOfBirth:                "01.01.1995",
	InsuranceHolderName:        "Petar",
	InsuranceHolderSurname:     "Petrović",
	InsuranceHolderNameCyrl:    "Петар",
	InsuranceHolderSurnameCyrl: "Петровић",
	InsuranceNumber:            "12345678",
	InsuranceStartDate:         "29.03.2014",
}
var documentMedical3 = document.MedicalDocument{
	GivenName:     "Pablo Diego",
	Surname:       "Ruiz Picasso",
	AddressStreet: "Rue des Grands-Augustins",
	AddressTown:   "Paris",
}

func Test_GetFullName_Medical(t *testing.T) {
	testCases := []struct {
		value    document.MedicalDocument
		expected string
	}{
		{
			value:    documentMedical1,
			expected: "",
		},
		{
			value:    documentMedical2,
			expected: "Петар, Арсеније, Петровић",
		},
		{
			value:    documentMedical3,
			expected: "Pablo Diego, Ruiz Picasso",
		},
	}

	for _, testCase := range testCases {
		result := testCase.value.GetFullName()
		if result != testCase.expected {
			t.Errorf("Expected '%s' but got '%s'", testCase.expected, result)
		}
	}
}

func Test_GetFullStreetAddress_Medical(t *testing.T) {
	testCases := []struct {
		value    document.MedicalDocument
		expected string
	}{
		{
			value:    documentMedical1,
			expected: "",
		},
		{
			value:    documentMedical2,
			expected: "Његошева, Број: 9 Стан: 5",
		},
		{
			value:    documentMedical3,
			expected: "Rue des Grands-Augustins",
		},
	}

	for _, testCase := range testCases {
		result := testCase.value.GetFullStreetAddress()
		if result != testCase.expected {
			t.Errorf("Expected '%s' but got '%s'", testCase.expected, result)
		}
	}
}

func Test_GetFullPlaceAddress_Medical(t *testing.T) {
	testCases := []struct {
		value    document.MedicalDocument
		expected string
	}{
		{
			value:    documentMedical1,
			expected: "",
		},
		{
			value:    documentMedical2,
			expected: "Подгорица, Црна Гора",
		},
		{
			value:    documentMedical3,
			expected: "Paris",
		},
	}

	for _, testCase := range testCases {
		result := testCase.value.GetFullPlaceAddress()
		if result != testCase.expected {
			t.Errorf("Expected '%s' but got '%s'", testCase.expected, result)
		}
	}
}

func Test_BuildPdfMedical(t *testing.T) {
	document.SetDataFromLocalFiles(t)

	_, _, err := documentMedical1.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	_, _, err = documentMedical2.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
