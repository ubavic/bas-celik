package card

import "github.com/telesma-app/iso7816"

// Constructs an ISO 7816-4 command APDU from the expected response length.
func buildAPDU(cla, ins, p1, p2 byte, data []byte, ne uint) []byte {
	if ne > 0x10000 {
		panic("APDU expected response length too large")
	}

	command := iso7816.Command{
		CLA:  cla,
		INS:  ins,
		P1:   p1,
		P2:   p2,
		Data: data,
		Le:   uint32(ne),
	}
	apdu, err := command.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return apdu
}
