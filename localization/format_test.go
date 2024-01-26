package localization_test

import (
	"testing"

	"github.com/ubavic/bas-celik/localization"
)

func Test_FormatDate(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"23051987", "23.05.1987."},
		{"01010000", "Nije dostupan"},
		{"123", "123"},
		{"", ""},
	}

	for _, testCase := range testCases {
		value := testCase.value
		localization.FormatDate(&value)
		if value != testCase.expected {
			t.Errorf("Got %s but expected %s", value, testCase.expected)
		}
	}
}

func Test_FormatDateYMD(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"19870523", "23.05.1987"},
		{"123", "123"},
		{"", ""},
	}

	for _, testCase := range testCases {
		value := testCase.value
		localization.FormatDateYMD(&value)
		if value != testCase.expected {
			t.Errorf("Got %s but expected %s", value, testCase.expected)
		}
	}
}

func Test_FormFormatYesNo(t *testing.T) {
	testCases := []struct {
		value    bool
		script   localization.Script
		expected string
	}{
		{true, localization.Latin, "Da"},
		{true, localization.Cyrillic, "Да"},
		{true, 10, "Да"},
		{false, localization.Latin, "Ne"},
		{false, localization.Cyrillic, "Не"},
		{false, 155, "Не"},
	}

	for _, testCase := range testCases {
		result := localization.FormatYesNo(testCase.value, testCase.script)
		if result != testCase.expected {
			t.Errorf("Got %s for value %t and script %d, but expected %s", result, testCase.value, testCase.script, testCase.expected)
		}
	}
}
