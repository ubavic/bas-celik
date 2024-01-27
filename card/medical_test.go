package card

import (
	"encoding/hex"
	"fmt"
	"slices"
	"testing"
)

func Test_descramble(t *testing.T) {
	// Some raw bytes from different cards
	testCases := []struct {
		data, expectedData string
	}{
		{
			"200435043f04430431043b04380447043a043804200044043e043d04340420003704300420003704340440043004320441044204320435043d043e0420003e044104380433044304400430045a043504",
			"Републички фонд за здравствено осигурање",
		},
		{
			"210440043104380458043004",
			"Србија",
		},
		{
			"110415041e041304200410041404",
			"БЕОГРАД",
		},
		{
			"170430043f043e0441043b0435043d0438042000430420003f044004380432044004350434043d043e043c04200034044004430448044204320443042c00200034044004430433043e043c0420003f044004300432043d043e043c0420003b043804460443042c0020003a043e04340420003f0440043504340443043704350442043d0438043a0430042c00200046043804320438043b043d04300420003b0438044604300420003d043004200041043b04430436043104380420004304200032043e045804410446043804",
			"Запослени у привредном друштву, другом правном лицу, код предузетника, цивилна лица на служби у војсци",
		},
		{
			"110443045f04350442042000200435043f04430431043b0438043a0435042000210440043104380458043504",
			"Буџет Републике Србије",
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				decoded, err := hex.DecodeString(testCase.data)
				if err != nil {
					t.Errorf("Unexpected error %v", err)
				}

				fields := make(map[uint][]byte, 0)
				fields[uint(i)] = decoded

				descramble(fields, uint(i))
				if !slices.Equal(fields[uint(i)], []byte(testCase.expectedData)) {
					t.Errorf("Got %s, but expected %s", string(fields[uint(i)]), testCase.expectedData)
				}
			},
		)
	}

	t.Run(
		"Empty case",
		func(t *testing.T) {
			fields := make(map[uint][]byte, 0)
			descramble(fields, 1)
			if !slices.Equal(fields[1], []byte{}) {
				t.Errorf("Got %v, but expected empty slice", fields[1])
			}
		},
	)
}
