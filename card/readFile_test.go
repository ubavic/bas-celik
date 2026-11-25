package card

import (
	"testing"

	"github.com/ebfe/scard"
)

// Answers every READ BINARY with the same data and every other command with 9000.
type shortHeaderCard struct {
	header []byte
}

func (c *shortHeaderCard) Status() (*scard.CardStatus, error) {
	return &scard.CardStatus{}, nil
}

func (c *shortHeaderCard) Transmit(apdu []byte) ([]byte, error) {
	if len(apdu) > 1 && apdu[1] == 0xB0 {
		return append(append([]byte{}, c.header...), 0x90, 0x00), nil
	}

	return []byte{0x90, 0x00}, nil
}

func (c *shortHeaderCard) BeginTransaction() error {
	return nil
}

func (c *shortHeaderCard) EndTransaction(scard.Disposition) error {
	return nil
}

func Test_ReadFile_ShortHeader(t *testing.T) {
	testCases := []struct {
		name string
		read func() ([]byte, error)
	}{
		{
			"Gemalto",
			func() ([]byte, error) {
				card := Gemalto{smartCard: &shortHeaderCard{header: []byte{0x02, 0x0F, 0x10}}}
				return card.ReadFile(ID_DOCUMENT_FILE_LOC)
			}},
		{
			"Gemalto certificate",
			func() ([]byte, error) {
				card := Gemalto{smartCard: &shortHeaderCard{header: []byte{0x10}}}
				return card.readCertificateFile([]byte{0x71, 0x02})
			}},
		{
			"Medical",
			func() ([]byte, error) {
				card := MedicalCard{smartCard: &shortHeaderCard{header: []byte{0x01, 0x0D, 0x10}}}
				return card.ReadFile(MED_DOCUMENT_FILE_LOC)
			}},
		{
			"Apollo",
			func() ([]byte, error) {
				card := Apollo{smartCard: &shortHeaderCard{header: []byte{0x00, 0x00, 0x00, 0x00, 0x10}}}
				return card.ReadFile(ID_DOCUMENT_FILE_LOC)
			}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.read()
			if err == nil {
				t.Errorf("Expected an error for a short header, got nil")
			}
		})
	}
}
