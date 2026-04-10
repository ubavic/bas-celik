package server

import (
	"encoding/json"
	"fmt"
	"io"

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

	moduleSession, err := pkcs11.GetPkcsSession(pkcs11.CardVendor(providerID))
	if err != nil {
		return err
	}

	session.vendor = pkcs11.CardVendor(providerID)
	session.moduleSession = &moduleSession

	slotIds, slotNames, err := moduleSession.ListSlots()
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
