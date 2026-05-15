package server

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func TestGetValidCertificates(t *testing.T) {
	now := time.Now()
	validCert, _ := generateTestCert(now.Add(-time.Hour), now.Add(time.Hour), x509.KeyUsageDigitalSignature)
	expiredCert, _ := generateTestCert(now.Add(-48*time.Hour), now.Add(-24*time.Hour), x509.KeyUsageDigitalSignature)
	futureCert, _ := generateTestCert(now.Add(24*time.Hour), now.Add(48*time.Hour), x509.KeyUsageDigitalSignature)

	certs := []pkcs11.NamedCert{
		{Id: []byte{0x01}, Certificate: validCert},
		{Id: []byte{0x02}, Certificate: expiredCert},
		{Id: []byte{0x03}, Certificate: futureCert},
	}

	valid := GetValidCertificates(certs, pkcs11.CardVendorPosta)
	if len(valid) != 1 {
		t.Fatalf("expected 1 valid cert, got %d", len(valid))
	}
	if !bytes.Equal(valid[0].Id, []byte{0x01}) {
		t.Errorf("expected cert id 0x01, got %x", valid[0].Id)
	}
}

func TestGetValidCertificatesMupFilter(t *testing.T) {
	now := time.Now()
	withCommitment, _ := generateTestCert(now.Add(-time.Hour), now.Add(time.Hour),
		x509.KeyUsageDigitalSignature|x509.KeyUsageContentCommitment)
	withoutCommitment, _ := generateTestCert(now.Add(-time.Hour), now.Add(time.Hour),
		x509.KeyUsageDigitalSignature)

	certs := []pkcs11.NamedCert{
		{Id: []byte{0x01}, Certificate: withCommitment},
		{Id: []byte{0x02}, Certificate: withoutCommitment},
	}

	valid := GetValidCertificates(certs, pkcs11.CardVendorMup)
	if len(valid) != 1 {
		t.Fatalf("expected 1 valid cert for MUP, got %d", len(valid))
	}
	if !bytes.Equal(valid[0].Id, []byte{0x01}) {
		t.Errorf("expected cert with ContentCommitment, got id %x", valid[0].Id)
	}
}

func TestGetValidCertificatesMupAllowsCommitment(t *testing.T) {
	now := time.Now()
	cert, _ := generateTestCert(now.Add(-time.Hour), now.Add(time.Hour),
		x509.KeyUsageContentCommitment)

	certs := []pkcs11.NamedCert{
		{Id: []byte{0x01}, Certificate: cert},
	}

	valid := GetValidCertificates(certs, pkcs11.CardVendorMup)
	if len(valid) != 1 {
		t.Fatalf("expected 1 valid cert for MUP with ContentCommitment, got %d", len(valid))
	}
}

func TestGetValidCertificatesEmpty(t *testing.T) {
	valid := GetValidCertificates(nil, pkcs11.CardVendorPosta)
	if len(valid) != 0 {
		t.Errorf("expected 0 valid certs, got %d", len(valid))
	}
}

func TestGetCertificateAliases(t *testing.T) {
	cert, _ := generateTestCert(time.Now().Add(-time.Hour), time.Now().Add(time.Hour), x509.KeyUsageDigitalSignature)

	certs := []pkcs11.NamedCert{
		{Id: []byte{0xAB, 0xCD}, Certificate: cert},
	}

	aliases := GetCertificateAliases(certs)
	if len(aliases) != 1 {
		t.Fatalf("expected 1 alias, got %d", len(aliases))
	}
	if aliases[0].Alias != "abcd" {
		t.Errorf("alias = %q, want %q", aliases[0].Alias, "abcd")
	}
	if aliases[0].Name != "Test User" {
		t.Errorf("name = %q, want %q", aliases[0].Name, "Test User")
	}
}

func TestGetCertificateAliasesEmpty(t *testing.T) {
	aliases := GetCertificateAliases(nil)
	if len(aliases) != 0 {
		t.Errorf("expected 0 aliases, got %d", len(aliases))
	}
}

func TestHandleGetCertificates(t *testing.T) {
	cert, _ := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)
	certId := []byte{0x01, 0x02}

	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: certId, Certificate: cert},
		},
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)

	session := SmartboxSession{
		moduleSession: mock,
		vendor:        pkcs11.CardVendorPosta,
	}

	data, _ := json.Marshal(Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input: GetCertificatesInput{
			TerminalId: 0,
			Pin:        "1234",
		},
	})

	var buf bytes.Buffer
	err := srv.handleGetCertificates(&session, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetCertificatesPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != operationGetCertificates {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetCertificates)
	}

	if len(rsp.Payload.Certificates) != 1 {
		t.Fatalf("expected 1 certificate, got %d", len(rsp.Payload.Certificates))
	}

	if rsp.Payload.Certificates[0].Alias != hex.EncodeToString(certId) {
		t.Errorf("alias = %q, want %q", rsp.Payload.Certificates[0].Alias, hex.EncodeToString(certId))
	}

	if mock.loginPin != "1234" {
		t.Errorf("login pin = %q, want %q", mock.loginPin, "1234")
	}
}

func TestHandleGetCertificatesNilModule(t *testing.T) {
	srv := newTestServer(nil, pkcs11.CardVendorPosta)
	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{Pin: "1234"},
	})

	var buf bytes.Buffer
	err := srv.handleGetCertificates(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error for nil module session")
	}
}

func TestHandleGetCertificatesLoginError(t *testing.T) {
	mock := &mockPkcsModuleSession{
		loginErr: fmt.Errorf("wrong pin"),
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{moduleSession: mock}

	data, _ := json.Marshal(Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{Pin: "0000"},
	})

	var buf bytes.Buffer
	err := srv.handleGetCertificates(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error for login failure")
	}
}
