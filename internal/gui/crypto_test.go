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

func Test_SanitizeFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "cert",
		},
		{
			input:    "    ",
			expected: "cert",
		},
		{
			input:    "!@#",
			expected: "cert",
		},
		{
			input:    "Pera Miloš Perić",
			expected: "Pera_Miloš_Perić",
		},
		{
			input:    "Пера Милош Перић",
			expected: "Пера_Милош_Перић",
		},
		{
			input:    "ЖАРКО ЉУБА ПЕРИЋ",
			expected: "ЖАРКО_ЉУБА_ПЕРИЋ",
		},
		{
			input:    "John    Doe !!!#",
			expected: "John_Doe",
		},
		{
			input:    "John    Doe !!!#",
			expected: "John_Doe",
		},
		{
			input:    "John    Doe !!!#",
			expected: "John_Doe",
		},
		{
			input:    "Аркадий Иванович Свидригайлов",
			expected: "Аркадий_Иванович_Свидригайлов",
		},
	}

	for _, testCase := range testCases {
		result := sanitizeFilename(testCase.input)
		if testCase.expected != result {
			t.Errorf("Expected `%s` but got `%s`", testCase.expected, result)
		}
	}
}
