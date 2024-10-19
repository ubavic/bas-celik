package card

import (
	"fmt"

	"github.com/ubavic/bas-celik/card/ber"
	"github.com/ubavic/bas-celik/card/cardErrors"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/localization"
)

// Represents a smart card that contains a Serbian vehicle document.
type VehicleCard struct {
	atr       Atr
	smartCard Card
	files     [4][]byte
}

// Possibly deprecated.
var VEHICLE_ATR_0 = Atr([]byte{
	0x3B, 0xDB, 0x96, 0x00, 0x80, 0xB1, 0xFE, 0x45,
	0x1F, 0x83, 0x00, 0x31, 0xC0, 0x64, 0x1A, 0x18,
	0x01, 0x00, 0x0F, 0x90, 0x00, 0x52,
})

var VEHICLE_ATR_1 = Atr([]byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xB0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
})

var VEHICLE_ATR_2 = Atr([]byte{
	0x3B, 0x9D, 0x13, 0x81, 0x31, 0x60, 0x37, 0x80,
	0x31, 0xC0, 0x69, 0x4D, 0x54, 0x43, 0x4F, 0x53,
	0x73, 0x02, 0x02, 0x04, 0x40,
})

var VEHICLE_ATR_3 = Atr([]byte{
	0x3B, 0x9D, 0x13, 0x81, 0x31, 0x60, 0x37, 0x80,
	0x31, 0xC0, 0x69, 0x4D, 0x54, 0x43, 0x4F, 0x53,
	0x73, 0x02, 0x05, 0x04, 0x47,
})

var VEHICLE_ATR_4 = Atr([]byte{
	0x3B, 0x9D, 0x18, 0x81, 0x31, 0xFC, 0x35, 0x80,
	0x31, 0xC0, 0x69, 0x4D, 0x54, 0x43, 0x4F, 0x53,
	0x73, 0x02, 0x05, 0x02, 0xD4,
})

// Initializes vehicle card by trying three different sets of commands.
// The procedure is reverse-engineered from the official binary.
func (card VehicleCard) InitCard() error {
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

func (card *VehicleCard) ReadCard() error {
	var err error

	for i := byte(0); i <= 3; i++ {
		card.files[int(i)], err = card.readFile([]byte{0xD0, i*0x10 + 0x01})
		if err != nil {
			return fmt.Errorf("reading document %d file: %w", i, err)
		}
	}

	return nil
}

func (card *VehicleCard) GetDocument() (document.Document, error) {
	doc := document.VehicleDocument{}
	data := ber.BER{}

	for i := byte(0); i <= 3; i++ {
		parsed, err := ber.ParseBER(card.files[int(i)])
		if err != nil {
			return nil, fmt.Errorf("parsing %d file: %w", i, err)
		}

		err = data.Merge(*parsed)
		if err != nil {
			return nil, fmt.Errorf("merging %d data: %w", i, err)
		}
	}

	data.AssignFrom(&doc.RegistrationNumberOfVehicle, 0x71, 0x81)
	data.AssignFrom(&doc.DateOfFirstRegistration, 0x71, 0x82)
	localization.FormatDateYMD(&doc.DateOfFirstRegistration)
	data.AssignFrom(&doc.VehicleIdNumber, 0x71, 0x8A)
	data.AssignFrom(&doc.VehicleMass, 0x71, 0x8C)
	data.AssignFrom(&doc.ExpiryDate, 0x71, 0x8D)
	localization.FormatDateYMD(&doc.ExpiryDate)
	data.AssignFrom(&doc.IssuingDate, 0x71, 0x8E)
	localization.FormatDateYMD(&doc.IssuingDate)
	data.AssignFrom(&doc.TypeApprovalNumber, 0x71, 0x8F)
	data.AssignFrom(&doc.PowerWeightRatio, 0x71, 0x93)
	data.AssignFrom(&doc.VehicleMake, 0x71, 0xA3, 0x87)
	data.AssignFrom(&doc.VehicleType, 0x71, 0xA3, 0x88)
	data.AssignFrom(&doc.CommercialDescription, 0x71, 0xA3, 0x89)
	data.AssignFrom(&doc.MaximumPermissibleLadenMass, 0x71, 0xA4, 0x8B)
	data.AssignFrom(&doc.EngineCapacity, 0x71, 0xA5, 0x90)
	data.AssignFrom(&doc.MaximumNetPower, 0x71, 0xA5, 0x91)
	data.AssignFrom(&doc.TypeOfFuel, 0x71, 0xA5, 0x92)
	data.AssignFrom(&doc.NumberOfSeats, 0x71, 0xA6, 0x94)
	data.AssignFrom(&doc.NumberOfStandingPlaces, 0x71, 0xA6, 0x95)
	data.AssignFrom(&doc.StateIssuing, 0x71, 0x9F33)
	data.AssignFrom(&doc.CompetentAuthority, 0x71, 0x9F35)
	data.AssignFrom(&doc.AuthorityIssuing, 0x71, 0x9F36)
	data.AssignFrom(&doc.UnambiguousNumber, 0x71, 0x9F38)
	data.AssignFrom(&doc.VehicleCategory, 0x72, 0x98)
	data.AssignFrom(&doc.NumberOfAxles, 0x72, 0x99)
	data.AssignFrom(&doc.VehicleLoad, 0x72, 0xC4)
	data.AssignFrom(&doc.YearOfProduction, 0x72, 0xC5)
	data.AssignFrom(&doc.EngineIdNumber, 0x72, 0xA5, 0x9E)
	data.AssignFrom(&doc.SerialNumber, 0x72, 0xC9)
	data.AssignFrom(&doc.ColourOfVehicle, 0x72, 0x9F24)
	data.AssignFrom(&doc.UsersPersonalNo, 0x72, 0xC3)
	data.AssignFrom(&doc.OwnersPersonalNo, 0x72, 0xC2)

	data.AssignFrom(&doc.OwnersSurnameOrBusinessName, 0x71, 0xA1, 0xA2, 0x83)
	data.AssignFrom(&doc.OwnerName, 0x71, 0xA1, 0xA2, 0x84)
	data.AssignFrom(&doc.OwnerAddress, 0x71, 0xA1, 0xA2, 0x85)

	data.AssignFrom(&doc.UsersSurnameOrBusinessName, 0x71, 0xA1, 0xA9, 0x83)
	data.AssignFrom(&doc.UsersName, 0x71, 0xA1, 0xA9, 0x84)
	data.AssignFrom(&doc.UsersAddress, 0x71, 0xA1, 0xA9, 0x85)
	if doc.UsersName == "" && doc.UsersSurnameOrBusinessName == "" && doc.UsersAddress == "" {
		data.AssignFrom(&doc.UsersSurnameOrBusinessName, 0x72, 0xA1, 0xA9, 0x83)
		data.AssignFrom(&doc.UsersName, 0x72, 0xA1, 0xA9, 0x84)
		data.AssignFrom(&doc.UsersAddress, 0x72, 0xA1, 0xA9, 0x85)
	}

	return &doc, nil
}

func (card *VehicleCard) Atr() Atr {
	return card.atr
}

func (card *VehicleCard) readFile(name []byte) ([]byte, error) {
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

	length, offset, err := parseVehicleCardFileSize(header)
	if err != nil {
		return nil, fmt.Errorf("parsing file header: %w", err)
	}

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

func parseVehicleCardFileSize(data []byte) (uint, uint, error) {
	if len(data) < 1 {
		return 0, 0, cardErrors.ErrInvalidLength
	}

	offset := uint(data[1]) + 2

	if offset >= uint(len(data)) {
		return 0, 0, cardErrors.ErrInvalidLength
	}

	_, _, offsetDelta1, err := ber.ParseTag(data[offset:])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing tag: %w", err)
	}

	if offset+uint(offsetDelta1) >= uint(len(data)) {
		return 0, 0, cardErrors.ErrInvalidLength
	}

	dataLength, offsetDelta2, err := ber.ParseLength(data[offset+uint(offsetDelta1):])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing size: %w", err)
	}

	length := uint(dataLength + offsetDelta1 + offsetDelta2)

	return length, offset, nil
}

func (card *VehicleCard) selectFile(name []byte) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x02, 0x04, name, 0)

	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}
