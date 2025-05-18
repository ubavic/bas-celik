package server

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type GetTerminalsInput struct {
	ProviderId stringOrInt `json:"providerId"`
}

type GetTerminalsPayload struct {
	Terminals []Terminal `json:"terminals"`
}

type Terminal struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *SmartBoxServer) handleGetTerminals(session *SmartboxSession, data []byte, w io.Writer) error {
	msg := Message[GetTerminalsInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	providerID := msg.Input.ProviderId

	if providerID < 0 || int(providerID) > int(pkcs11.CardVendorPks) {
		return fmt.Errorf("invalid provider id")
	}

	moduleIndex := slices.IndexFunc(s.modulePaths, func(m ModulePath) bool {
		return m.Vendor == pkcs11.CardVendor(providerID)
	})

	if moduleIndex < 0 {
		return fmt.Errorf("module not found")
	}

	module, err := pkcs11.NewPkcsExternalModule(s.modulePaths[moduleIndex].Path)
	if err != nil {
		return err
	}

	session.vendor = pkcs11.CardVendor(providerID)
	session.module = &module

	slotIds, slotNames, err := module.ListSlots()
	if err != nil {
		return err
	}

	terminals := make([]Terminal, 0, len(slotIds))
	for i, id := range slotIds {
		terminals = append(terminals, Terminal{
			Id:   fmt.Sprintf("%d", id),
			Name: slotNames[i],
		})
	}

	rsp := Response[GetTerminalsPayload]{
		Operation: operationGetTerminals,
		Payload: GetTerminalsPayload{
			Terminals: terminals,
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}
