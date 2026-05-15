package smartbox

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestHandleGetInfoNewSession(t *testing.T) {
	srv := newTestServer(nil)

	data, _ := json.Marshal(Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input: GetInfoInput{
			SbSession: "test-session-123",
			Language:  "sr",
			Host:      "eporezi.purs.gov.rs",
		},
	})

	sessionId := ""
	var buf bytes.Buffer
	err := srv.handleGetInfo(&sessionId, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sessionId != "test-session-123" {
		t.Errorf("sessionId = %q, want %q", sessionId, "test-session-123")
	}

	if _, ok := srv.sessions["test-session-123"]; !ok {
		t.Error("session was not created in server sessions map")
	}

	var rsp Response[GetInfoPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != operationGetInfo {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetInfo)
	}
}

func TestHandleGetInfoExistingSession(t *testing.T) {
	srv := newTestServer(nil)
	srv.sessions["existing-session"] = SmartboxSession{
		id:            "existing-session",
		terminalId:    2,
		certificateId: []byte{0xAB},
	}

	data, _ := json.Marshal(Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input: GetInfoInput{
			SbSession: "existing-session",
		},
	})

	sessionId := ""
	var buf bytes.Buffer
	err := srv.handleGetInfo(&sessionId, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetInfoPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Payload.TerminalId != 2 {
		t.Errorf("terminalId = %d, want 2", rsp.Payload.TerminalId)
	}

	if rsp.Payload.CertificateId != "ab" {
		t.Errorf("certificateId = %q, want %q", rsp.Payload.CertificateId, "ab")
	}
}

func TestHandleGetInfoEmptySession(t *testing.T) {
	srv := newTestServer(nil)

	data, _ := json.Marshal(Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: ""},
	})

	sessionId := ""
	var buf bytes.Buffer
	err := srv.handleGetInfo(&sessionId, data, &buf)
	if err == nil {
		t.Fatal("expected error for empty session id")
	}
}
