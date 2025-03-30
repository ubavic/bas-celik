package server

import (
	"encoding/json"
	"fmt"
	"io"
)

type GetInfoInput struct {
	SbSession string `json:"sbSession"`
	Language  string `json:"language"`
	Host      string `json:"host"`
}

type GetInfoPayload struct {
	TerminalId    int    `json:"terminalId"`
	ProviderId    int    `json:"providerId"`
	CertificateId string `json:"certificateId"`
}

func (s *SmartBoxServer) handleGetInfo(sessionId *string, data []byte, w io.Writer) error {
	msg := Message[GetInfoInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if msg.Input.SbSession == "" {
		return fmt.Errorf("invalid session id")
	}

	session, ok := s.sessions[msg.Input.SbSession]
	if !ok {
		s.sessions[msg.Input.SbSession] = SmartboxSession{
			id: msg.Input.SbSession,
		}
	}

	*sessionId = msg.Input.SbSession

	rsp := Response[GetInfoPayload]{
		Operation: operationGetInfo,
		Payload: GetInfoPayload{
			TerminalId:    session.terminalId,
			ProviderId:    int(session.vendor),
			CertificateId: session.certificateId,
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}
