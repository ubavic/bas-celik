package server

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func setupIntegrationServer(t *testing.T) (*httptest.Server, *x509.Certificate, *rsa.PrivateKey, []byte) {
	t.Helper()

	cert, key := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature|x509.KeyUsageContentCommitment,
	)
	certId := []byte{0xAA, 0xBB}

	mock := &mockPkcsModuleSession{
		slots:     []uint{0},
		slotNames: []string{"Virtual Reader"},
		certs: []pkcs11.NamedCert{
			{Id: certId, Certificate: cert},
		},
		signFunc: func(cid, digest []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
		},
	}

	srv := &SmartBoxServer{
		sessions:      make(map[string]SmartboxSession),
		loadedVendors: []pkcs11.CardVendor{pkcs11.CardVendorPosta},
		sessionProvider: func(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
			return mock, nil
		},
	}

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	return ts, cert, key, certId
}

func dialWS(t *testing.T, ts *httptest.Server) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("websocket dial failed: %v", err)
	}
	t.Cleanup(func() { conn.CloseNow() })

	return conn
}

func readJSON[T any](t *testing.T, conn *websocket.Conn) T {
	t.Helper()

	ctx := context.Background()
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read websocket message: %v", err)
	}

	var msg T
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("failed to unmarshal message: %v\nraw: %s", err, string(data))
	}

	return msg
}

func writeJSON(t *testing.T, conn *websocket.Conn, v any) {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	ctx := context.Background()
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("failed to write websocket message: %v", err)
	}
}

func TestIntegrationOnOpen(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)
	conn := dialWS(t, ts)

	rsp := readJSON[Response[OnOpenPayload]](t, conn)
	if rsp.Operation != onOpenEvent {
		t.Errorf("operation = %q, want %q", rsp.Operation, onOpenEvent)
	}
	if rsp.Payload.AppName != "SmartBox" {
		t.Errorf("appName = %q, want %q", rsp.Payload.AppName, "SmartBox")
	}
}

func TestIntegrationGetInfo(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)
	conn := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, conn)

	writeJSON(t, conn, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input: GetInfoInput{
			SbSession: "session-abc",
			Language:  "sr",
		},
	})

	rsp := readJSON[Response[GetInfoPayload]](t, conn)
	if rsp.Operation != operationGetInfo {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetInfo)
	}
}

func TestIntegrationGetProviders(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)
	conn := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, conn)

	writeJSON(t, conn, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "s1"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, conn)

	writeJSON(t, conn, Message[GetProvidersInput]{
		Operation: operationGetProviders,
	})

	rsp := readJSON[Response[GetProvidersPayload]](t, conn)
	if rsp.Operation != operationGetProviders {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetProviders)
	}
	if len(rsp.Payload.Providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(rsp.Payload.Providers))
	}
	if rsp.Payload.Providers[0].Name != pkcs11.CardVendorPosta.String() {
		t.Errorf("provider name = %q, want %q", rsp.Payload.Providers[0].Name, pkcs11.CardVendorPosta.String())
	}
}

func TestIntegrationGetTerminals(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)
	conn := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, conn)

	writeJSON(t, conn, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "s1"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, conn)

	writeJSON(t, conn, Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})

	rsp := readJSON[Response[GetTerminalsPayload]](t, conn)
	if rsp.Operation != operationGetTerminals {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetTerminals)
	}
	if len(rsp.Payload.Terminals) != 1 {
		t.Fatalf("expected 1 terminal, got %d", len(rsp.Payload.Terminals))
	}
	if rsp.Payload.Terminals[0].Name != "Virtual Reader" {
		t.Errorf("terminal name = %q, want %q", rsp.Payload.Terminals[0].Name, "Virtual Reader")
	}
}

func TestIntegrationFullFlow(t *testing.T) {
	ts, _, _, certId := setupIntegrationServer(t)
	conn := dialWS(t, ts)

	onOpen := readJSON[Response[OnOpenPayload]](t, conn)
	if onOpen.Operation != onOpenEvent {
		t.Fatalf("expected ON_OPEN_EVENT, got %q", onOpen.Operation)
	}

	writeJSON(t, conn, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "flow-session"},
	})
	infoRsp := readJSON[Response[GetInfoPayload]](t, conn)
	if infoRsp.Operation != operationGetInfo {
		t.Fatalf("expected GET_INFO, got %q", infoRsp.Operation)
	}

	writeJSON(t, conn, Message[GetProvidersInput]{
		Operation: operationGetProviders,
	})
	provRsp := readJSON[Response[GetProvidersPayload]](t, conn)
	if len(provRsp.Payload.Providers) == 0 {
		t.Fatal("expected at least 1 provider")
	}

	writeJSON(t, conn, Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(provRsp.Payload.Providers[0].Id)},
	})
	termRsp := readJSON[Response[GetTerminalsPayload]](t, conn)
	if len(termRsp.Payload.Terminals) == 0 {
		t.Fatal("expected at least 1 terminal")
	}

	writeJSON(t, conn, Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{TerminalId: 0, Pin: "1234"},
	})
	certRsp := readJSON[Response[GetCertificatesPayload]](t, conn)
	if len(certRsp.Payload.Certificates) == 0 {
		t.Fatal("expected at least 1 certificate")
	}
	if certRsp.Payload.Certificates[0].Alias != hex.EncodeToString(certId) {
		t.Errorf("cert alias = %q, want %q", certRsp.Payload.Certificates[0].Alias, hex.EncodeToString(certId))
	}

	xmlDoc := `<Document><Data>test payload</Data></Document>`
	writeJSON(t, conn, Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: certRsp.Payload.Certificates[0],
			Pin:         "1234",
			Xml:         base64.StdEncoding.EncodeToString([]byte(xmlDoc)),
		},
	})
	signRsp := readJSON[Response[GetSignedXmlPayload]](t, conn)
	if signRsp.Operation != operationGetSignedXml {
		t.Fatalf("expected GET_SIGNED_XML, got %q", signRsp.Operation)
	}

	signedBytes, err := base64.StdEncoding.DecodeString(signRsp.Payload.Xml)
	if err != nil {
		t.Fatalf("failed to decode signed XML: %v", err)
	}
	signedXml := string(signedBytes)
	if !strings.Contains(signedXml, "<SignatureValue>") {
		t.Error("signed XML missing SignatureValue")
	}
	if !strings.Contains(signedXml, "test payload") {
		t.Error("signed XML missing original content")
	}
}

func TestIntegrationSessionNotFoundBeforeGetInfo(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("websocket dial failed: %v", err)
	}
	defer conn.CloseNow()

	_ = readJSON[Response[OnOpenPayload]](t, conn)

	writeJSON(t, conn, Message[GetProvidersInput]{
		Operation: operationGetProviders,
	})

	// Server returns an error internally for missing session, which breaks
	// the handler loop. The deferred writer close may flush an empty frame
	// before the connection is closed. Drain any remaining messages until
	// we get a read error (connection closed).
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		readCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		_, _, err = conn.Read(readCtx)
		cancel()
		if err != nil {
			return
		}
	}
	t.Error("expected connection to eventually close after session-not-found error")
}

func TestTwoSessionsOneFailsLogin(t *testing.T) {
	cert, key := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature|x509.KeyUsageContentCommitment,
	)
	certId := []byte{0xAA, 0xBB}

	goodMock := &mockPkcsModuleSession{
		slots:     []uint{0},
		slotNames: []string{"Reader A"},
		certs:     []pkcs11.NamedCert{{Id: certId, Certificate: cert}},
		signFunc: func(cid, digest []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
		},
	}

	failMock := &mockPkcsModuleSession{
		slots:     []uint{0},
		slotNames: []string{"Reader B"},
		loginErr:  fmt.Errorf("incorrect PIN"),
	}

	var mu sync.Mutex
	providerCalls := 0

	srv := &SmartBoxServer{
		sessions:      make(map[string]SmartboxSession),
		loadedVendors: []pkcs11.CardVendor{pkcs11.CardVendorPosta},
		sessionProvider: func(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
			mu.Lock()
			defer mu.Unlock()
			providerCalls++
			if providerCalls == 1 {
				return goodMock, nil
			}
			return failMock, nil
		},
	}

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	connA := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, connA)
	writeJSON(t, connA, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "session-a"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, connA)
	writeJSON(t, connA, Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})
	_ = readJSON[Response[GetTerminalsPayload]](t, connA)

	connB := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, connB)
	writeJSON(t, connB, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "session-b"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, connB)
	writeJSON(t, connB, Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})
	_ = readJSON[Response[GetTerminalsPayload]](t, connB)

	writeJSON(t, connB, Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{TerminalId: 0, Pin: "0000"},
	})

	ctx := context.Background()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		readCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		_, _, err := connB.Read(readCtx)
		cancel()
		if err != nil {
			break
		}
	}

	writeJSON(t, connA, Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{TerminalId: 0, Pin: "1234"},
	})
	certRsp := readJSON[Response[GetCertificatesPayload]](t, connA)
	if len(certRsp.Payload.Certificates) == 0 {
		t.Fatal("session A: expected certificates after session B failed login")
	}
	if certRsp.Payload.Certificates[0].Alias != hex.EncodeToString(certId) {
		t.Errorf("session A: cert alias = %q, want %q",
			certRsp.Payload.Certificates[0].Alias, hex.EncodeToString(certId))
	}
}

func TestTwoSessionsOneClosesConnection(t *testing.T) {
	cert, key := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature|x509.KeyUsageContentCommitment,
	)
	certId := []byte{0xAA, 0xBB}

	mock := &mockPkcsModuleSession{
		slots:     []uint{0},
		slotNames: []string{"Virtual Reader"},
		certs:     []pkcs11.NamedCert{{Id: certId, Certificate: cert}},
		signFunc: func(cid, digest []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
		},
	}

	srv := &SmartBoxServer{
		sessions:      make(map[string]SmartboxSession),
		loadedVendors: []pkcs11.CardVendor{pkcs11.CardVendorPosta},
		sessionProvider: func(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
			return mock, nil
		},
	}

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ctx := context.Background()
	connA, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("session A: websocket dial failed: %v", err)
	}

	_ = readJSON[Response[OnOpenPayload]](t, connA)
	writeJSON(t, connA, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "session-a"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, connA)

	connB := dialWS(t, ts)
	_ = readJSON[Response[OnOpenPayload]](t, connB)
	writeJSON(t, connB, Message[GetInfoInput]{
		Operation: operationGetInfo,
		Input:     GetInfoInput{SbSession: "session-b"},
	})
	_ = readJSON[Response[GetInfoPayload]](t, connB)

	connA.Close(websocket.StatusGoingAway, "client closing")
	time.Sleep(100 * time.Millisecond)

	writeJSON(t, connB, Message[GetTerminalsInput]{
		Operation: operationGetTerminals,
		Input:     GetTerminalsInput{ProviderId: stringOrInt(pkcs11.CardVendorPosta)},
	})
	termRsp := readJSON[Response[GetTerminalsPayload]](t, connB)
	if len(termRsp.Payload.Terminals) == 0 {
		t.Fatal("session B: expected terminals after session A disconnected")
	}

	writeJSON(t, connB, Message[GetCertificatesInput]{
		Operation: operationGetCertificates,
		Input:     GetCertificatesInput{TerminalId: 0, Pin: "1234"},
	})
	certRsp := readJSON[Response[GetCertificatesPayload]](t, connB)
	if len(certRsp.Payload.Certificates) == 0 {
		t.Fatal("session B: expected certificates after session A disconnected")
	}

	xmlDoc := `<Document><Data>session B payload</Data></Document>`
	writeJSON(t, connB, Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: certRsp.Payload.Certificates[0],
			Pin:         "1234",
			Xml:         base64.StdEncoding.EncodeToString([]byte(xmlDoc)),
		},
	})
	signRsp := readJSON[Response[GetSignedXmlPayload]](t, connB)
	if signRsp.Operation != operationGetSignedXml {
		t.Fatalf("session B: expected GET_SIGNED_XML, got %q", signRsp.Operation)
	}

	signedBytes, err := base64.StdEncoding.DecodeString(signRsp.Payload.Xml)
	if err != nil {
		t.Fatalf("session B: failed to decode signed XML: %v", err)
	}
	if !strings.Contains(string(signedBytes), "session B payload") {
		t.Error("session B: signed XML missing original content")
	}
}

func TestServeHTTPOriginRejection(t *testing.T) {
	ts, _, _, _ := setupIntegrationServer(t)

	ctx := context.Background()
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	_, resp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Origin": []string{"https://evil.example.com"},
		},
	})
	if err == nil {
		t.Fatal("expected error for rejected origin")
	}
	if resp != nil && resp.StatusCode != http.StatusForbidden {
		t.Logf("status = %d (origin rejection may produce different status codes)", resp.StatusCode)
	}
}
