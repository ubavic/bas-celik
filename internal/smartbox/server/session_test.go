package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func TestSmartboxSessionCertificate(t *testing.T) {
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

	session := SmartboxSession{
		moduleSession: mock,
		certificateId: certId,
	}

	got := session.Certificate()
	if got == nil {
		t.Fatal("expected certificate, got nil")
	}
	if got.Subject.CommonName != "Test User" {
		t.Errorf("CommonName = %q, want %q", got.Subject.CommonName, "Test User")
	}
}

func TestSmartboxSessionCertificateNotFound(t *testing.T) {
	cert, _ := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)

	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: []byte{0x01}, Certificate: cert},
		},
	}

	session := SmartboxSession{
		moduleSession: mock,
		certificateId: []byte{0xFF},
	}

	got := session.Certificate()
	if got != nil {
		t.Error("expected nil certificate for non-matching id")
	}
}

func TestSmartboxSessionCertificateNilModule(t *testing.T) {
	session := SmartboxSession{}
	got := session.Certificate()
	if got != nil {
		t.Error("expected nil certificate when moduleSession is nil")
	}
}

func TestSmartboxSessionPublic(t *testing.T) {
	cert, _ := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)
	certId := []byte{0x01}

	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: certId, Certificate: cert},
		},
	}

	session := SmartboxSession{
		moduleSession: mock,
		certificateId: certId,
	}

	pub := session.Public()
	if pub == nil {
		t.Fatal("expected public key, got nil")
	}

	if _, ok := pub.(*rsa.PublicKey); !ok {
		t.Error("expected *rsa.PublicKey")
	}
}

func TestSmartboxSessionPublicNilModule(t *testing.T) {
	session := SmartboxSession{}
	pub := session.Public()
	if pub != nil {
		t.Error("expected nil public key")
	}
}

func TestSmartboxSessionSign(t *testing.T) {
	cert, key := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)
	certId := []byte{0x01}

	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: certId, Certificate: cert},
		},
		signFunc: func(cid, digest []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
		},
	}

	session := SmartboxSession{
		moduleSession: mock,
		certificateId: certId,
	}

	digest := make([]byte, 32)
	rand.Read(digest)

	sig, err := session.Sign(rand.Reader, digest, crypto.SHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig) == 0 {
		t.Fatal("expected non-empty signature")
	}

	err = rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, digest, sig)
	if err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestSmartboxSessionSignNilModule(t *testing.T) {
	session := SmartboxSession{certificateId: []byte{0x01}}
	_, err := session.Sign(rand.Reader, make([]byte, 32), crypto.SHA256)
	if err == nil {
		t.Fatal("expected error for nil module")
	}
}

func TestSmartboxSessionSignNoCertId(t *testing.T) {
	session := SmartboxSession{moduleSession: &mockPkcsModuleSession{}}
	_, err := session.Sign(rand.Reader, make([]byte, 32), crypto.SHA256)
	if err == nil {
		t.Fatal("expected error for nil certificate ID")
	}
}

func TestSmartboxSessionSignWrongHash(t *testing.T) {
	session := SmartboxSession{
		moduleSession: &mockPkcsModuleSession{},
		certificateId: []byte{0x01},
	}
	_, err := session.Sign(rand.Reader, make([]byte, 20), crypto.SHA1)
	if err == nil {
		t.Fatal("expected error for non-SHA256 hash")
	}
}
