package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func TestHandleGetTerminals(t *testing.T) {
	mock := &mockPkcsModuleSession{
		slots:     []uint{0, 1},
		slotNames: []string{"Reader 0", "Reader 1"},
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)

	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})

	var buf bytes.Buffer
	err := srv.handleGetTerminals(&session, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetTerminalsPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != operationGetTerminals {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetTerminals)
	}

	if len(rsp.Payload.Terminals) != 2 {
		t.Fatalf("expected 2 terminals, got %d", len(rsp.Payload.Terminals))
	}

	if rsp.Payload.Terminals[0].Id != "0" {
		t.Errorf("terminal[0].Id = %q, want %q", rsp.Payload.Terminals[0].Id, "0")
	}
	if rsp.Payload.Terminals[0].Name != "Reader 0" {
		t.Errorf("terminal[0].Name = %q, want %q", rsp.Payload.Terminals[0].Name, "Reader 0")
	}

	if session.vendor != pkcs11.CardVendorPosta {
		t.Errorf("session.vendor = %d, want %d", session.vendor, pkcs11.CardVendorPosta)
	}

	if session.moduleSession == nil {
		t.Error("session.moduleSession should not be nil after GET_TERMINALS")
	}
}

func TestHandleGetTerminalsInvalidProvider(t *testing.T) {
	mock := &mockPkcsModuleSession{}
	srv := newTestServer(mock, pkcs11.CardVendorPosta)

	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(99)},
	})

	var buf bytes.Buffer
	err := srv.handleGetTerminals(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error for invalid provider id")
	}
}

func TestHandleGetTerminalsSessionProviderError(t *testing.T) {
	srv := &SmartBoxServer{
		sessions:      make(map[string]SmartboxSession),
		loadedVendors: []pkcs11.CardVendor{pkcs11.CardVendorPosta},
		sessionProvider: func(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
			return nil, fmt.Errorf("module not available")
		},
	}

	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})

	var buf bytes.Buffer
	err := srv.handleGetTerminals(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error from session provider")
	}
}

func TestHandleGetTerminalsListSlotsError(t *testing.T) {
	mock := &mockPkcsModuleSession{
		listSlotsErr: fmt.Errorf("no slots"),
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})

	var buf bytes.Buffer
	err := srv.handleGetTerminals(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error from ListSlots")
	}
}

func TestHandleGetTerminalsNoSlots(t *testing.T) {
	mock := &mockPkcsModuleSession{
		slots:     []uint{},
		slotNames: []string{},
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})

	var buf bytes.Buffer
	err := srv.handleGetTerminals(&session, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetTerminalsPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(rsp.Payload.Terminals) != 0 {
		t.Errorf("expected 0 terminals, got %d", len(rsp.Payload.Terminals))
	}
}
