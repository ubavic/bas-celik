package card

import (
	"fmt"

	"github.com/ubavic/bas-celik/document"
)

// Represents a smart card that contains a Serbian vehicle document.
type VehicleCard struct {
	smartCard Card
}

// Possibly deprecated.
var VEHICLE_ATR_0 = []byte{
	0x3B, 0xDB, 0x96, 0x00, 0x80, 0xB1, 0xFE, 0x45,
	0x1F, 0x83, 0x00, 0x31, 0xC0, 0x64, 0x1A, 0x18,
	0x01, 0x00, 0x0F, 0x90, 0x00, 0x52,
}

var VEHICLE_ATR_1 = []byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xB0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
}

var VEHICLE_ATR_2 = []byte{
	0x3B, 0x9D, 0x13, 0x81, 0x31, 0x60, 0x37, 0x80,
	0x31, 0xC0, 0x69, 0x4D, 0x54, 0x43, 0x4F, 0x53,
	0x73, 0x02, 0x02, 0x04, 0x40,
}

func readVehicleCard(card VehicleCard) (*document.VehicleDocument, error) {
	s1 := []byte{0xA0, 0x00, 0x00, 0x01, 0x51, 0x00, 0x00}
	apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, s1, 0)
	_, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, err
	}

	doc := document.VehicleDocument{}
	data := BER{}

	for i := byte(0); i <= 3; i++ {
		rsp, err := card.readFile([]byte{0xD0, i*0x10 + 0x01}, false)
		if err != nil {
			return nil, fmt.Errorf("reading document %d file: %w", i, err)
		}

		parsed, err := ParseBER(rsp)
		if err != nil {
			return nil, fmt.Errorf("parsing %d file: %w", i, err)
		}

		err = data.merge(*parsed)
		if err != nil {
			return nil, fmt.Errorf("merging %d data: %w", i, err)
		}
	}

	data.assignFrom(&doc.RegistrationNumberOfVehicle, 0x71, 0x81)
	data.assignFrom(&doc.DateOfFirstRegistration, 0x71, 0x82)
	document.FormatDate2(&doc.DateOfFirstRegistration)
	data.assignFrom(&doc.VehicleIdNumber, 0x71, 0x8A)
	data.assignFrom(&doc.VehicleMass, 0x71, 0x8C)
	data.assignFrom(&doc.ExpiryDate, 0x71, 0x8D)
	document.FormatDate2(&doc.ExpiryDate)
	data.assignFrom(&doc.IssuingDate, 0x71, 0x8E)
	document.FormatDate2(&doc.IssuingDate)
	data.assignFrom(&doc.TypeApprovalNumber, 0x71, 0x8F)
	data.assignFrom(&doc.PowerWeightRatio, 0x71, 0x93)
	data.assignFrom(&doc.VehicleMake, 0x71, 0xA3, 0x87)
	data.assignFrom(&doc.VehicleType, 0x71, 0xA3, 0x88)
	data.assignFrom(&doc.CommercialDescription, 0x71, 0xA3, 0x89)
	data.assignFrom(&doc.MaximumPermissibleLadenMass, 0x71, 0xA4, 0x8B)
	data.assignFrom(&doc.EngineCapacity, 0x71, 0xA5, 0x90)
	data.assignFrom(&doc.MaximumNetPower, 0x71, 0xA5, 0x91)
	data.assignFrom(&doc.TypeOfFuel, 0x71, 0xA5, 0x92)
	data.assignFrom(&doc.NumberOfSeats, 0x71, 0xA6, 0x94)
	data.assignFrom(&doc.NumberOfStandingPlaces, 0x71, 0xA6, 0x95)
	data.assignFrom(&doc.StateIssuing, 0x71, 0x9F33)
	data.assignFrom(&doc.CompetentAuthority, 0x71, 0x9F35)
	data.assignFrom(&doc.AuthorityIssuing, 0x71, 0x9F36)
	data.assignFrom(&doc.UnambiguousNumber, 0x71, 0x9F38)
	data.assignFrom(&doc.VehicleCategory, 0x72, 0x98)
	data.assignFrom(&doc.NumberOfAxles, 0x72, 0x99)
	data.assignFrom(&doc.VehicleLoad, 0x72, 0xC4)
	data.assignFrom(&doc.YearOfProduction, 0x72, 0xC5)
	data.assignFrom(&doc.EngineIdNumber, 0x72, 0xA5, 0x9E)
	data.assignFrom(&doc.OwnersSurnameOrBusinessName, 0x71, 0xA1, 0xA2, 0x83)
	data.assignFrom(&doc.OwnerName, 0x71, 0xA1, 0xA2, 0x84)
	data.assignFrom(&doc.OwnerAddress, 0x71, 0xA1, 0xA2, 0x85)
	data.assignFrom(&doc.UsersSurnameOrBusinessName, 0x71, 0xA1, 0xA9, 0x83)
	data.assignFrom(&doc.UsersName, 0x71, 0xA1, 0xA9, 0x84)
	data.assignFrom(&doc.UsersAddress, 0x71, 0xA1, 0xA9, 0x85)
	data.assignFrom(&doc.OwnersPersonalNo, 0x72, 0xC2)
	data.assignFrom(&doc.UsersPersonalNo, 0x72, 0xC3)
	data.assignFrom(&doc.SerialNumber, 0x72, 0xC9)
	data.assignFrom(&doc.ColourOfVehicle, 0x72, 0x9F24)

	return &doc, nil
}

func (card VehicleCard) readFile(name []byte, _ bool) ([]byte, error) {
	output := make([]byte, 0)

	_, err := card.selectFile(name)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	const headerSize = uint(0x20)

	header, err := read(card.smartCard, 0, headerSize)
	if err != nil {
		return nil, fmt.Errorf("reading file header: %w", err)
	}

	offset := uint(header[1])

	len32, delta, err := parseBerLength(header[offset+3:])
	if err != nil {
		return nil, fmt.Errorf("parsing size: %w", err)
	}

	length := uint(len32 + delta)

	for length > 0 {
		toRead := min(length, 0x64)
		data, err := read(card.smartCard, offset, toRead)
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	return output, nil
}

// Initializes vehicle card by trying three different sets of commands.
// The procedure is reverse-engineered from the official binary.
func (card VehicleCard) initCard() error {
	tryToSelect := func(cmd1, cmd2, cmd3 []byte) error {
		apu := buildAPDU(0x00, 0xA4, 0x04, 0x00, cmd1, 0)
		rsp, err := card.smartCard.Transmit(apu)
		if err != nil {
			return fmt.Errorf("selecting file: %w", err)
		}

		if responseOK(rsp) {
			apu = buildAPDU(0x00, 0xA4, 0x04, 0x00, cmd2, 0)
			_, err = card.smartCard.Transmit(apu)
			if err != nil {
				return fmt.Errorf("selecting file: %w", err)
			}

			apu = buildAPDU(0x00, 0xA4, 0x04, 0x0C, cmd3, 0)
			_, err = card.smartCard.Transmit(apu)
			if err != nil {
				return fmt.Errorf("selecting file: %w", err)
			}

			return nil
		} else {
			return fmt.Errorf("selecting file: %w", err)
		}
	}

	err := tryToSelect(
		[]byte{0xA0, 0x00, 0x00, 0x01, 0x51, 0x00, 0x00},
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x77, 0x01, 0x08, 0x00, 0x07, 0x00, 0x00, 0xFE, 0x00, 0x00, 0x01, 0x00},
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x77, 0x01, 0x08, 0x00, 0x07, 0x00, 0x00, 0xFE, 0x00, 0x00, 0xAD, 0xF2})
	if err == nil {
		return nil
	}

	err = tryToSelect(
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00},
		[]byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x56, 0x4C, 0x04, 0x02, 0x01},
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x77, 0x01, 0x08, 0x00, 0x07, 0x00, 0x00, 0xFE, 0x00, 0x00, 0xAD, 0xF2})
	if err == nil {
		return nil
	}

	err = tryToSelect(
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x18, 0x43, 0x4D, 0x00},
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x18, 0x34, 0x14, 0x01, 0x00, 0x65, 0x56, 0x4C, 0x2D, 0x30, 0x30, 0x31},
		[]byte{0xA0, 0x00, 0x00, 0x00, 0x18, 0x65, 0x56, 0x4C, 0x2D, 0x30, 0x30, 0x31})
	if err == nil {
		return nil
	}

	return fmt.Errorf("card not responsive: %w", err)
}

func (card VehicleCard) selectFile(name []byte) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x02, 0x04, name, 0)

	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}
