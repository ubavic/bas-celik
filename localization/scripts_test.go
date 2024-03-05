package localization_test

import (
	"testing"

	"github.com/ubavic/bas-celik/localization"
)

func Test_CyrillicToLatin(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"брза вижљаста лија хоће да ђипи преко њушке флегматичног џукца", "brza vižljasta lija hoće da đipi preko njuške flegmatičnog džukca"},
		{"БРЗА ВИЖЉАСТА ЛИЈА ХОЋЕ ДА ЂИПИ ПРЕКО ЊУШКЕ ФЛЕГМАТИЧНОГ ЏУКЦА", "BRZA VIŽLjASTA LIJA HOĆE DA ĐIPI PREKO NjUŠKE FLEGMATIČNOG DžUKCA"},
		{"0123456789", "0123456789"},
		{"", ""},
	}

	for _, testCase := range testCases {
		latin := localization.CyrillicToLatin(testCase.value)
		if latin != testCase.expected {
			t.Errorf("Got '%s' but expected '%s'", latin, testCase.expected)
		}
	}
}
