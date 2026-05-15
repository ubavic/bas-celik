package server

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func TestHandleGetProviders(t *testing.T) {
	srv := newTestServer(nil, pkcs11.CardVendorPosta, pkcs11.CardVendorMup)

	data, _ := json.Marshal(Message[GetProvidersInput]{
		Operation: operationGetProviders,
	})

	var buf bytes.Buffer
	err := srv.handleGetProviders(data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetProvidersPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != operationGetProviders {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetProviders)
	}

	if len(rsp.Payload.Providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(rsp.Payload.Providers))
	}

	if rsp.Payload.Providers[0].Id != int(pkcs11.CardVendorPosta) {
		t.Errorf("provider[0].Id = %d, want %d", rsp.Payload.Providers[0].Id, int(pkcs11.CardVendorPosta))
	}

	if rsp.Payload.Providers[1].Id != int(pkcs11.CardVendorMup) {
		t.Errorf("provider[1].Id = %d, want %d", rsp.Payload.Providers[1].Id, int(pkcs11.CardVendorMup))
	}
}

func TestHandleGetProvidersEmpty(t *testing.T) {
	srv := newTestServer(nil)

	data, _ := json.Marshal(Message[GetProvidersInput]{
		Operation: operationGetProviders,
	})

	var buf bytes.Buffer
	err := srv.handleGetProviders(data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetProvidersPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(rsp.Payload.Providers) != 0 {
		t.Errorf("expected 0 providers, got %d", len(rsp.Payload.Providers))
	}
}
