package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type mockPkcsModuleSession struct {
	slots        []uint
	slotNames    []string
	listSlotsErr error

	loginPin    string
	loginSlotId int
	loginErr    error

	certs       []pkcs11.NamedCert
	getCertsErr error

	signFunc func(certId, digest []byte) ([]byte, error)
}

func (m *mockPkcsModuleSession) ListSlots() ([]uint, []string, error) {
	return m.slots, m.slotNames, m.listSlotsErr
}

func (m *mockPkcsModuleSession) OpenSessionAndLogin(pin string, terminalIndex int) error {
	m.loginPin = pin
	m.loginSlotId = terminalIndex
	return m.loginErr
}

func (m *mockPkcsModuleSession) GetCertificates() ([]pkcs11.NamedCert, error) {
	return m.certs, m.getCertsErr
}

func (m *mockPkcsModuleSession) SignDigest(certId []byte, sha256digest []byte) ([]byte, error) {
	if m.signFunc != nil {
		return m.signFunc(certId, sha256digest)
	}
	return nil, fmt.Errorf("sign not configured")
}

func generateTestCert(notBefore, notAfter time.Time, keyUsage x509.KeyUsage) (*x509.Certificate, *rsa.PrivateKey) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject: pkix.Name{
			CommonName:   "Test User",
			Country:      []string{"RS"},
			SerialNumber: "12345678",
			ExtraNames: []pkix.AttributeTypeAndValue{
				{Type: []int{2, 5, 4, 5}, Value: "PNORS-12345678"},
				{Type: []int{2, 5, 4, 42}, Value: "Petar"},
				{Type: []int{2, 5, 4, 4}, Value: "Petrović"},
			},
		},
		Issuer: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Test Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  keyUsage,
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	cert, _ := x509.ParseCertificate(certDER)

	return cert, key
}

func newTestServer(mock *mockPkcsModuleSession, vendors ...pkcs11.CardVendor) *SmartBoxServer {
	return &SmartBoxServer{
		sessions:      make(map[string]SmartboxSession),
		loadedVendors: vendors,
		sessionProvider: func(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
			return mock, nil
		},
	}
}
