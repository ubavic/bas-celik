package document_test

import (
	"testing"

	"github.com/ubavic/bas-celik/document"
)

var documentMedical1 = document.MedicalDocument{}
var documentMedical2 = document.MedicalDocument{
	GivenNameLatin:         "Петар",
	ParentNameLatin:        "Арсеније",
	FamilyNameLatin:        "Петровић",
	Street:                 "Његошева",
	Number:                 "9",
	Apartment:              "5",
	Place:                  "Подгорица",
	Country:                "Црна Гора",
	PersonalNumber:         "01019950900200",
	DateOfBirth:            "01.01.1995",
	CarrierGivenNameLatin:  "Petar",
	CarrierFamilyNameLatin: "Petrović",
	CarrierGivenName:       "Петар",
	CarrierFamilyName:      "Петровић",
	InsurantNumber:         "12345678",
	InsuranceStartDate:     "29.03.2014",
	CardId:                 "12345678901",
}
var documentMedical3 = document.MedicalDocument{
	GivenNameLatin:  "Pablo Diego",
	FamilyNameLatin: "Ruiz Picasso",
	Street:          "Rue des Grands-Augustins",
	Place:           "Paris",
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

func Test_GetExpiryDateFromRfzo(t *testing.T) {
	err := documentMedical1.UpdateValidUntilDateFromRfzo()
	if err != document.ErrInvalidCardNo {
		t.Errorf("Expected the InvalidCardNo error but got %v", err)
	}

	err = documentMedical2.UpdateValidUntilDateFromRfzo()
	if err != document.ErrInvalidInsuranceNo {
		t.Errorf("Expected the InvalidInsuranceNo error but got %v", err)
	}

}

func Test_parseDateFromRfzoResponse(t *testing.T) {
	_, err := document.ParseValidUntilDateFromRfzoResponse("")
	if err != document.ErrNoSubmatchFound {
		t.Errorf("Expected the NoSubmatchFound error but got %v", err)
	}

	date, err := document.ParseValidUntilDateFromRfzoResponse("Ваши иницијали су <strong>Н.Н.</strong> (ЛБО: 123456789)<br />Матична филијала: <strong>Београд</strong>.<br/>Ваша здравствена књижица је оверена до: <strong>3.4.2025.</strong>")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if date != "3.4.2025." {
		t.Errorf("Expected date `3.4.2025.` but got `%s`", date)
	}

	date, err = document.ParseValidUntilDateFromRfzoResponse("Ваша здравствена књижица је оверена до: <strong>31.12.2023.</strong>")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if date != "31.12.2023." {
		t.Errorf("Expected date `31.12.2023.` but got `%s`", date)
	}
}
