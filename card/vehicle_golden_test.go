package card

import (
	"testing"

	"github.com/ubavic/bas-celik/v2/document"
)

// Encodes a single BER-TLV record (tag bytes + length + value).
func berRecord(tag, value []byte) []byte {
	out := append([]byte{}, tag...)
	n := len(value)
	switch {
	case n < 0x80:
		out = append(out, byte(n))
	case n < 0x100:
		out = append(out, 0x81, byte(n))
	default:
		out = append(out, 0x82, byte(n>>8), byte(n))
	}
	return append(out, value...)
}

func berConcat(records ...[]byte) []byte {
	out := []byte{}
	for _, r := range records {
		out = append(out, r...)
	}
	return out
}

// Builds synthetic vehicle card files that mirror the BER structure of a real
// Serbian vehicle registration card (ATR GEMALTO_ATR_4, "ESCE 8.0" variant).
// All values are fabricated; only the tag layout reproduces a genuine card so
// that the tag-to-field mapping in VehicleCard.GetDocument is exercised end to
// end (BER parsing, cross-file merge and date formatting).
func syntheticVehicleFiles() [4][]byte {
	file0 := berRecord([]byte{0x71}, berConcat(
		berRecord([]byte{0x80}, []byte{0x01, 0x00, 0x00}),
		berRecord([]byte{0x9F, 0x33}, []byte("SRB")),
		berRecord([]byte{0x9F, 0x35}, []byte("MUP Republike Srbije")),
		berRecord([]byte{0x9F, 0x36}, []byte("PU Beograd")),
		berRecord([]byte{0x9F, 0x37}, []byte("0")),
		berRecord([]byte{0x9F, 0x38}, []byte("SRB1234")),
		berRecord([]byte{0x81}, []byte("BG123AB")),
		berRecord([]byte{0x82}, []byte("20200115")),
		berRecord([]byte{0xA1}, berConcat(
			berRecord([]byte{0xA2}, berConcat(
				berRecord([]byte{0x83}, []byte("Petrović")),
				berRecord([]byte{0x84}, []byte("Petar")),
				berRecord([]byte{0x85}, []byte("Kneza Miloša 1, Beograd")),
			)),
			berRecord([]byte{0x86}, []byte("1")),
		)),
		berRecord([]byte{0xA3}, berConcat(
			berRecord([]byte{0x87}, []byte("ZASTAVA")),
			berRecord([]byte{0x88}, []byte("H")),
			berRecord([]byte{0x89}, []byte("101 SKALA")),
		)),
		berRecord([]byte{0x8A}, []byte("VX1234567890ABCDE")),
		berRecord([]byte{0xA4}, berConcat(
			berRecord([]byte{0x8B}, []byte("1500")),
		)),
		berRecord([]byte{0x8C}, []byte("1100")),
		berRecord([]byte{0x8D}, []byte("20251231")),
		berRecord([]byte{0x8E}, []byte("20240101")),
		berRecord([]byte{0x8F}, []byte("e1")),
		berRecord([]byte{0xA5}, berConcat(
			berRecord([]byte{0x90}, []byte("1100")),
			berRecord([]byte{0x91}, []byte("40")),
			berRecord([]byte{0x92}, []byte("BENZIN")),
		)),
		berRecord([]byte{0x93}, []byte("36")),
		berRecord([]byte{0xA6}, berConcat(
			berRecord([]byte{0x94}, []byte("5")),
			berRecord([]byte{0x95}, []byte("0")),
		)),
	))

	file1 := berRecord([]byte{0x72}, berConcat(
		berRecord([]byte{0x80}, []byte{0x01, 0x00, 0x00}),
		berRecord([]byte{0x98}, []byte("M1")),
		berRecord([]byte{0x99}, []byte("2")),
		berRecord([]byte{0xA5}, berConcat(
			berRecord([]byte{0x9E}, []byte("ENG123456")),
		)),
		berRecord([]byte{0x9F, 0x24}, []byte("PLAVA")),
	))

	file2 := berRecord([]byte{0x72}, berConcat(
		berRecord([]byte{0x80}, []byte{0x01, 0x00, 0x00}),
		berRecord([]byte{0xC2}, []byte("0101990710020")),
		berRecord([]byte{0xC4}, []byte("0")),
		berRecord([]byte{0xC5}, []byte("2020")),
	))

	file3 := berRecord([]byte{0x72}, berConcat(
		berRecord([]byte{0x80}, []byte{0x01, 0x00, 0x00}),
		berRecord([]byte{0xC9}, []byte("SRB1234567890X")),
	))

	return [4][]byte{file0, file1, file2, file3}
}

func Test_VehicleCard_GetDocument(t *testing.T) {
	card := VehicleCard{files: syntheticVehicleFiles()}

	doc, err := card.GetDocument()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vehicle, ok := doc.(*document.VehicleDocument)
	if !ok {
		t.Fatalf("expected *document.VehicleDocument, got %T", doc)
	}

	cases := []struct {
		field    string
		got      string
		expected string
	}{
		{"StateIssuing", vehicle.StateIssuing, "SRB"},
		{"CompetentAuthority", vehicle.CompetentAuthority, "MUP Republike Srbije"},
		{"AuthorityIssuing", vehicle.AuthorityIssuing, "PU Beograd"},
		{"UnambiguousNumber", vehicle.UnambiguousNumber, "SRB1234"},
		{"RegistrationNumberOfVehicle", vehicle.RegistrationNumberOfVehicle, "BG123AB"},
		{"DateOfFirstRegistration", vehicle.DateOfFirstRegistration, "15.01.2020"},
		{"OwnersSurnameOrBusinessName", vehicle.OwnersSurnameOrBusinessName, "Petrović"},
		{"OwnerName", vehicle.OwnerName, "Petar"},
		{"OwnerAddress", vehicle.OwnerAddress, "Kneza Miloša 1, Beograd"},
		{"VehicleMake", vehicle.VehicleMake, "ZASTAVA"},
		{"VehicleType", vehicle.VehicleType, "H"},
		{"CommercialDescription", vehicle.CommercialDescription, "101 SKALA"},
		{"VehicleIdNumber", vehicle.VehicleIdNumber, "VX1234567890ABCDE"},
		{"MaximumPermissibleLadenMass", vehicle.MaximumPermissibleLadenMass, "1500"},
		{"VehicleMass", vehicle.VehicleMass, "1100"},
		{"ExpiryDate", vehicle.ExpiryDate, "31.12.2025"},
		{"IssuingDate", vehicle.IssuingDate, "01.01.2024"},
		{"TypeApprovalNumber", vehicle.TypeApprovalNumber, "e1"},
		{"EngineCapacity", vehicle.EngineCapacity, "1100"},
		{"MaximumNetPower", vehicle.MaximumNetPower, "40"},
		{"TypeOfFuel", vehicle.TypeOfFuel, "BENZIN"},
		{"PowerWeightRatio", vehicle.PowerWeightRatio, "36"},
		{"NumberOfSeats", vehicle.NumberOfSeats, "5"},
		{"NumberOfStandingPlaces", vehicle.NumberOfStandingPlaces, "0"},
		{"VehicleCategory", vehicle.VehicleCategory, "M1"},
		{"NumberOfAxles", vehicle.NumberOfAxles, "2"},
		{"EngineIdNumber", vehicle.EngineIdNumber, "ENG123456"},
		{"ColourOfVehicle", vehicle.ColourOfVehicle, "PLAVA"},
		{"OwnersPersonalNo", vehicle.OwnersPersonalNo, "0101990710020"},
		{"VehicleLoad", vehicle.VehicleLoad, "0"},
		{"YearOfProduction", vehicle.YearOfProduction, "2020"},
		{"SerialNumber", vehicle.SerialNumber, "SRB1234567890X"},
	}

	for _, c := range cases {
		if c.got != c.expected {
			t.Errorf("%s: got %q, expected %q", c.field, c.got, c.expected)
		}
	}
}
