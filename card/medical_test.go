package card

import (
	"context"
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

func TestMedicalCardAcceptsSelectWithResponseAvailable(t *testing.T) {
	document := medRecord(1553, u16("Републички фонд за здравствено осигурање"))
	chip := &scriptedMedicalCard{
		t: t,
		steps: []medicalExchange{
			{[]byte{0x00, 0xA4, 0x04, 0x00, 0x0D, 0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x56, 0x53, 0x5A, 0x4B, 0x01}, []byte{0x61, 0x24}},
			{[]byte{0x00, 0xC0, 0x00, 0x00, 0x24}, []byte{0x61, 0x00}},
			{[]byte{0x00, 0xC0, 0x00, 0x00, 0x00}, []byte{0x90, 0x00}},
			{[]byte{0x00, 0xA4, 0x00, 0x00, 0x02, 0x0D, 0x01}, []byte{0x61, 0x0A}},
			{[]byte{0x00, 0xC0, 0x00, 0x00, 0x0A}, []byte{0x90, 0x00}},
			{[]byte{0x00, 0xB0, 0x00, 0x00, 0x04}, []byte{0x00, 0x00, byte(len(document)), 0x00, 0x90, 0x00}},
		},
	}
	for offset := 0; offset < len(document); offset += 32 {
		end := min(offset+32, len(document))
		chip.steps = append(chip.steps, medicalExchange{
			command:  []byte{0x00, 0xB0, 0x00, byte(offset + 4), byte(end - offset)},
			response: append(slices.Clone(document[offset:end]), 0x90, 0x00),
		})
	}
	card := MedicalCard{smartCard: chip}
	if !card.Test() {
		t.Fatal("failed to detect the medical card")
	}
	if len(chip.steps) != 0 {
		t.Fatalf("%d APDU exchanges were skipped", len(chip.steps))
	}
}

func TestMedicalCardLimitsPendingSelectResponses(t *testing.T) {
	chip := &scriptedMedicalCard{
		t: t,
		steps: []medicalExchange{
			{[]byte{0x00, 0xA4, 0x00, 0x00, 0x02, 0x0D, 0x01}, []byte{0x61, 0x00}},
		},
	}
	for range 16 {
		chip.steps = append(chip.steps, medicalExchange{
			command:  []byte{0x00, 0xC0, 0x00, 0x00, 0x00},
			response: []byte{0x61, 0x00},
		})
	}
	medical := MedicalCard{smartCard: chip}
	err := medical.selectWithGetResponse(0x00, MED_DOCUMENT_FILE_LOC)
	if err == nil || err.Error() != "too many GET RESPONSE commands after SELECT" {
		t.Fatalf("unexpected response to pending SELECT: %v", err)
	}
	if len(chip.steps) != 0 {
		t.Fatalf("%d APDU exchanges were skipped", len(chip.steps))
	}
}

type medicalExchange struct {
	command  []byte
	response []byte
}

type scriptedMedicalCard struct {
	Card
	t     *testing.T
	steps []medicalExchange
}

func (card *scriptedMedicalCard) Transmit(_ context.Context, command []byte) ([]byte, error) {
	card.t.Helper()
	if len(card.steps) == 0 {
		card.t.Fatalf("unexpected APDU %x", command)
	}
	step := card.steps[0]
	card.steps = card.steps[1:]
	if !slices.Equal(command, step.command) {
		card.t.Fatalf("APDU %x, want %x", command, step.command)
	}
	return step.response, nil
}

func Test_parseMedicalGender(t *testing.T) {
	testCases := []struct {
		value    []byte
		expected string
	}{
		{[]byte("010"), ""},
		{[]byte("01"), "Мушко"},
		{[]byte("02"), "Женско"},
		{[]byte("1"), ""},
		{[]byte{}, ""},
		{nil, ""},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				gender := parseMedicalGender(testCase.value)
				if gender != testCase.expected {
					t.Errorf("Got %q, but expected %q", gender, testCase.expected)
				}
			},
		)
	}
}
