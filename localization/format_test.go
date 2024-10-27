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
			t.Errorf("Got '%s' but expected '%s'", value, testCase.expected)
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
			t.Errorf("Got '%s', but expected '%s'", value, testCase.expected)
		}
	}
}

func Test_FormFormatYesNo(t *testing.T) {
	testCases := []struct {
		value    bool
		script   localization.Language
		expected string
	}{
		{true, localization.SrLatin, "Da"},
		{true, localization.SrCyrillic, "Да"},
		{true, localization.En, "Yes"},
		{false, localization.SrLatin, "Ne"},
		{false, localization.SrCyrillic, "Не"},
		{false, localization.En, "No"},
	}

	for _, testCase := range testCases {
		result := localization.FormatYesNo(testCase.value, testCase.script)
		if result != testCase.expected {
			t.Errorf("Got '%s' for value %t and script %s, but expected '%s'", result, testCase.value, testCase.script, testCase.expected)
		}
	}
}

func Test_JoinWithComma(t *testing.T) {
	testCases := []struct {
		value    []string
		expected string
	}{
		{
			value:    []string{""},
			expected: "",
		},
		{
			value:    []string{"aaa", "", "bb"},
			expected: "aaa, bb",
		},
		{
			value:    []string{"", "", "a", "b", "c"},
			expected: "a, b, c",
		},
	}

	for _, testCase := range testCases {
		result := localization.JoinWithComma(testCase.value...)
		if result != testCase.expected {
			t.Errorf("Got '%s', but expected '%s'", result, testCase.expected)
		}
	}
}
