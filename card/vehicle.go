package card

import (
	"fmt"

	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/document"
)

type VehicleCard struct {
	smartCard *scard.Card
}

// possibly deprecated
var VEHICLE_ATR_0 = []byte{
	0x3B, 0xDB, 0x96, 0x00, 0x80, 0xB1, 0xFE, 0x45,
	0x1F, 0x83, 0x00, 0x31, 0xC0, 0x64, 0x1A, 0x18,
	0x01, 0x00, 0x0F, 0x90, 0x00, 0x52,
}

var VEHICLE_ATR_1 = []byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xb0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
}

var VEHICLE_ATR_2 = []byte{
	0x3B, 0x9D, 0x13, 0x81, 0x31, 0x60, 0x37, 0x80,
	0x31, 0xc0, 0x69, 0x4d, 0x54, 0x43, 0x4f, 0x53,
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
	data := emptyTree()

	for i := byte(1); i <= 3; i++ {
		rsp, err := card.readFile([]byte{0xD0, i*0x10 + 0x01}, false)
		if err != nil {
			return nil, fmt.Errorf("reading document %d file: %w", i, err)
		}

		parsed, err := parseBER(rsp)
		if err != nil {
			return nil, fmt.Errorf("parsing %d file: %w", i, err)
		}

		err = data.add(*parsed)
		if err != nil {
			return nil, fmt.Errorf("merging %d data: %w", i, err)
		}
	}

	data.assignToFrom(&doc.RegistrationNumberOfVehicle, 0x71, 0x81)
	data.assignToFrom(&doc.DateOfFirstRegistration, 0x71, 0x82)
	data.assignToFrom(&doc.VehicleIdNumber, 0x71, 0x8A)
	data.assignToFrom(&doc.VehicleMass, 0x71, 0x8C)
	data.assignToFrom(&doc.ExpiryDate, 0x71, 0x8D)
	data.assignToFrom(&doc.IssuingDate, 0x71, 0x8E)
	data.assignToFrom(&doc.TypeApprovalNumber, 0x71, 0x8F)
	data.assignToFrom(&doc.PowerWeightRatio, 0x71, 0x93)
	data.assignToFrom(&doc.VehicleMake, 0x71, 0xA3, 0x87)
	data.assignToFrom(&doc.VehicleType, 0x71, 0xA3, 0x88)
	data.assignToFrom(&doc.CommercialDescription, 0x71, 0xA3, 0x89)
	data.assignToFrom(&doc.MaximumPermissibleLadenMass, 0x71, 0xA4, 0x8B)
	data.assignToFrom(&doc.EngineCapacity, 0x71, 0xA5, 0x90)
	data.assignToFrom(&doc.MaximumNetPower, 0x71, 0xA5, 0x91)
	data.assignToFrom(&doc.TypeOfFuel, 0x71, 0xA5, 0x92)
	data.assignToFrom(&doc.NumberOfSeats, 0x71, 0xA6, 0x94)
	data.assignToFrom(&doc.NumberOfStandingPlaces, 0x71, 0xA6, 0x95)
	data.assignToFrom(&doc.SerialNumber, 0x71, 0xC9)
	data.assignToFrom(&doc.StateIssuing, 0x71, 0x9F33)
	data.assignToFrom(&doc.CompetentAuthority, 0x71, 0x9F35)
	data.assignToFrom(&doc.AuthorityIssuing, 0x71, 0x9F36)
	data.assignToFrom(&doc.UnambiguousNumber, 0x71, 0x9F38)
	data.assignToFrom(&doc.VehicleCategory, 0x72, 0x98)
	data.assignToFrom(&doc.NumberOfAxles, 0x72, 0x99)
	data.assignToFrom(&doc.VehicleLoad, 0x72, 0xC4)
	data.assignToFrom(&doc.YearOfProduction, 0x72, 0xC5)
	data.assignToFrom(&doc.EngineIdNumber, 0x72, 0xA5, 0x9E)
	data.assignToFrom(&doc.OwnersSurnameOrBusinessName, 0x72, 0xA1, 0xA2, 0x83)
	data.assignToFrom(&doc.OwnerName, 0x72, 0xA1, 0xA2, 0x84)
	data.assignToFrom(&doc.OwnerAddress, 0x72, 0xA1, 0xA2, 0x85)
	data.assignToFrom(&doc.UsersSurnameOrBusinessName, 0x72, 0xA1, 0xA9, 0x83)
	data.assignToFrom(&doc.UsersName, 0x72, 0xA1, 0xA9, 0x84)
	data.assignToFrom(&doc.UsersAddress, 0x72, 0xA1, 0xA9, 0x85)
	data.assignToFrom(&doc.ColourOfVehicle, 0x72, 0x9F24)

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

	len32, _, err := parseBerLength(header[offset+3:])
	if err != nil {
		return nil, fmt.Errorf("parsing size: %w", err)
	}

	length := uint(len32)

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
