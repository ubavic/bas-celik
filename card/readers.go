package card

import (
	"fmt"

	"github.com/ebfe/scard"
)

type Readers struct {
	Context *scard.Context
}

func (rs Readers) List() ([]string, error) {
	return rs.Context.ListReaders()
}

func (rs Readers) ConnectReader(reader string) (*scard.Card, error) {
	if rs.Context == nil {
		return nil, fmt.Errorf("invalid context")
	}

	card, err := rs.Context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, fmt.Errorf("error connecting to reader: %w", err)
	}

	status, _ := card.Status()
	if status.State == scard.Absent {
		return nil, fmt.Errorf("card absent")
	}

	return card, nil
}
