package gui

import "testing"

func Test_ExtractCertShortName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "CERT #",
		},
		{
			input:    "AB",
			expected: "CERT #AB",
		},
		{
			input:    "ABCDEF",
			expected: "CERT #ABCDEF",
		},
		{
			input:    "1234567890",
			expected: "CERT #123456",
		},
	}

	for _, testCase := range testCases {
		result := extractCertShortName(testCase.input)
		if testCase.expected != result {
			t.Errorf("Expected `%s` but got `%s`", testCase.expected, result)
		}
	}

}
