package document_test

import (
	"image"
	"testing"

	"github.com/ubavic/bas-celik/document"
)

var documentId1 = document.IdDocument{}
var documentId2 = document.IdDocument{
	GivenName:       "Петар",
	ParentGivenName: "Арсеније",
	Surname:         "Петровић",
	Street:          "Његошева",
	HouseNumber:     "9",
	HouseLetter:     "Б",
	Place:           "Подгорица",
	PlaceOfBirth:    "Његуши",
	StateOfBirth:    "Црна Гора",
}
var documentId3 = document.IdDocument{
	GivenName:        "Pablo Diego",
	Surname:          "Ruiz Picasso",
	HouseNumber:      "7",
	HouseLetter:      "A",
	Street:           "Rue des Grands-Augustins",
	ApartmentNumber:  "21",
	Floor:            "6",
	Community:        "Saint-Germain-des-Prés",
	Place:            "Paris",
	PlaceOfBirth:     "Málaga",
	CommunityOfBirth: "Andalucía",
	StateOfBirth:     "Reino de España",
}

func Test_GetFullName_ID(t *testing.T) {
	testCases := []struct {
		value    document.IdDocument
		expected string
	}{
		{
			value:    documentId1,
			expected: "",
		},
		{
			value:    documentId2,
			expected: "Петар, Арсеније, Петровић",
		},
		{
			value:    documentId3,
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

func Test_GetFullAddress_ID(t *testing.T) {
	testCases := []struct {
		value            document.IdDocument
		expected         string
		expectedReversed string
	}{
		{
			value:            documentId1,
			expected:         "",
			expectedReversed: "",
		},
		{
			value:            documentId2,
			expected:         "Његошева 9Б, Подгорица",
			expectedReversed: "Подгорица, Његошева 9Б",
		},
		{
			value:            documentId3,
			expected:         "Rue des Grands-Augustins 7A/6/21, Saint-Germain-des-Prés, Paris",
			expectedReversed: "Paris, Saint-Germain-des-Prés, Rue des Grands-Augustins 7A/6/21",
		},
	}

	for _, testCase := range testCases {
		result := testCase.value.GetFullAddress(false)
		if result != testCase.expected {
			t.Errorf("Expected '%s' but got '%s'", testCase.expected, result)
		}

		resultReversed := testCase.value.GetFullAddress(true)
		if resultReversed != testCase.expectedReversed {
			t.Errorf("Expected '%s' but got '%s'", testCase.expectedReversed, resultReversed)
		}
	}
}

func Test_GetFullPlaceOfBirth_ID(t *testing.T) {
	testCases := []struct {
		value    document.IdDocument
		expected string
	}{
		{
			value:    documentId1,
			expected: "",
		},
		{
			value:    documentId2,
			expected: "Његуши, Црна Гора",
		},
		{
			value:    documentId3,
			expected: "Málaga, Andalucía, Reino de España",
		},
	}

	for _, testCase := range testCases {
		result := testCase.value.GetFullPlaceOfBirth()
		if result != testCase.expected {
			t.Errorf("Expected '%s' but got '%s'", testCase.expected, result)
		}
	}
}

func Test_BuildPdfID(t *testing.T) {
	document.UnsetData(t)

	_, _, err := documentId1.BuildPdf()
	if err == nil {
		t.Errorf("Expected error but got %v", err)
	}

	document.SetDataFromLocalFiles(t)

	_, _, err = documentId1.BuildPdf()
	if err == nil {
		t.Errorf("Expected error but got %v", err)
	}

	rect := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(rect)
	documentId1.Portrait = img
	documentId2.Portrait = img

	_, _, err = documentId1.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	_, _, err = documentId2.BuildPdf()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
