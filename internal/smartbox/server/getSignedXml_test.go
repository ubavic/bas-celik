package server

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/beevik/etree"
	"github.com/moov-io/signedxml"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func TestFormatSubject(t *testing.T) {
	subject := pkix.Name{
		CommonName: "Petar Petrović",
		Country:    []string{"RS"},
		Names: []pkix.AttributeTypeAndValue{
			{Type: asn1.ObjectIdentifier{2, 5, 4, 5}, Value: "12345678"},
			{Type: asn1.ObjectIdentifier{2, 5, 4, 5}, Value: "PNORS-12345678"},
			{Type: asn1.ObjectIdentifier{2, 5, 4, 42}, Value: "Petar"},
			{Type: asn1.ObjectIdentifier{2, 5, 4, 4}, Value: "Petrović"},
		},
	}

	result := formatSubject(subject)

	if !strings.Contains(result, "CN=Petar Petrović") {
		t.Errorf("expected CN in result, got %q", result)
	}
	if !strings.Contains(result, "GIVENNAME=Petar") {
		t.Errorf("expected GIVENNAME in result, got %q", result)
	}
	if !strings.Contains(result, "SURNAME=Petrović") {
		t.Errorf("expected SURNAME in result, got %q", result)
	}
	if !strings.Contains(result, "C=RS") {
		t.Errorf("expected C=RS in result, got %q", result)
	}

	pnorsIdx := strings.Index(result, "PNORS-12345678")
	firstSerIdx := strings.Index(result, "SERIALNUMBER=")
	if pnorsIdx < 0 || firstSerIdx < 0 {
		t.Fatalf("expected PNORS serial number first, got %q", result)
	}
	if pnorsIdx > strings.Index(result[firstSerIdx+13:], "SERIALNUMBER=")+firstSerIdx+13 {
		t.Error("PNORS serial should come first")
	}
}

func TestFormatSubjectNoCountry(t *testing.T) {
	subject := pkix.Name{
		CommonName: "Test",
	}

	result := formatSubject(subject)
	if !strings.HasSuffix(result, "C=") {
		t.Errorf("expected empty country, got %q", result)
	}
}

func TestPopulateSignature(t *testing.T) {
	cert, _ := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)

	doc := etree.NewDocument()
	doc.ReadFromString(`<Root><Data>test</Data></Root>`)

	err := populateSignature(doc.Root(), cert)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sig := doc.Root().FindElement("./Signature")
	if sig == nil {
		t.Fatal("Signature element not found")
	}

	xmlns := sig.SelectAttrValue("xmlns", "")
	if xmlns != "http://www.w3.org/2000/09/xmldsig#" {
		t.Errorf("xmlns = %q, want xmldsig namespace", xmlns)
	}

	signedInfo := sig.FindElement("./SignedInfo")
	if signedInfo == nil {
		t.Fatal("SignedInfo not found")
	}

	canon := signedInfo.FindElement("./CanonicalizationMethod")
	if canon == nil {
		t.Fatal("CanonicalizationMethod not found")
	}
	if algo := canon.SelectAttrValue("Algorithm", ""); algo != "http://www.w3.org/TR/2001/REC-xml-c14n-20010315" {
		t.Errorf("CanonicalizationMethod Algorithm = %q", algo)
	}

	sigMethod := signedInfo.FindElement("./SignatureMethod")
	if sigMethod == nil {
		t.Fatal("SignatureMethod not found")
	}
	if algo := sigMethod.SelectAttrValue("Algorithm", ""); algo != "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256" {
		t.Errorf("SignatureMethod Algorithm = %q", algo)
	}

	ref := signedInfo.FindElement("./Reference")
	if ref == nil {
		t.Fatal("Reference not found")
	}
	if uri := ref.SelectAttrValue("URI", "x"); uri != "" {
		t.Errorf("Reference URI = %q, want empty", uri)
	}

	transform := ref.FindElement("./Transforms/Transform")
	if transform == nil {
		t.Fatal("Transform not found")
	}
	if algo := transform.SelectAttrValue("Algorithm", ""); algo != "http://www.w3.org/2000/09/xmldsig#enveloped-signature" {
		t.Errorf("Transform Algorithm = %q", algo)
	}

	digestValue := ref.FindElement("./DigestValue")
	if digestValue == nil {
		t.Fatal("DigestValue not found")
	}

	x509Cert := sig.FindElement("./KeyInfo/X509Data/X509Certificate")
	if x509Cert == nil {
		t.Fatal("X509Certificate not found")
	}
	if x509Cert.Text() == "" {
		t.Error("X509Certificate text should not be empty")
	}

	serialNumber := sig.FindElement("./KeyInfo/X509Data/X509IssuerSerial/X509SerialNumber")
	if serialNumber == nil {
		t.Fatal("X509SerialNumber not found")
	}
	if serialNumber.Text() != cert.SerialNumber.String() {
		t.Errorf("serial = %q, want %q", serialNumber.Text(), cert.SerialNumber.String())
	}

	subjectName := sig.FindElement("./KeyInfo/X509Data/X509SubjectName")
	if subjectName == nil {
		t.Fatal("X509SubjectName not found")
	}
}

func TestPopulateSignatureAlreadyExists(t *testing.T) {
	cert, _ := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)

	doc := etree.NewDocument()
	doc.ReadFromString(`<Root><Signature></Signature></Root>`)

	err := populateSignature(doc.Root(), cert)
	if err == nil {
		t.Fatal("expected error when Signature already exists")
	}
}

func TestHandleGetSignedXml(t *testing.T) {
	cert, key := generateTestCert(
		time.Now().Add(-time.Hour),
		time.Now().Add(time.Hour),
		x509.KeyUsageDigitalSignature,
	)
	certId := []byte{0x01, 0x02}

	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: certId, Certificate: cert},
		},
		signFunc: func(cid, digest []byte) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
		},
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{
		moduleSession: mock,
		vendor:        pkcs11.CardVendorPosta,
	}

	xmlDoc := `<Document><Data>test data</Data></Document>`
	xmlB64 := base64.StdEncoding.EncodeToString([]byte(xmlDoc))

	data, _ := json.Marshal(Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: CertificateAlias{
				Alias: hex.EncodeToString(certId),
				Name:  "Test User",
			},
			Pin: "1234",
			Xml: xmlB64,
		},
	})

	var buf bytes.Buffer
	err := srv.handleGetSignedXml(&session, data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[GetSignedXmlPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != operationGetSignedXml {
		t.Errorf("operation = %q, want %q", rsp.Operation, operationGetSignedXml)
	}

	signedXmlBytes, err := base64.StdEncoding.DecodeString(rsp.Payload.Xml)
	if err != nil {
		t.Fatalf("failed to decode signed XML: %v", err)
	}

	signedXml := string(signedXmlBytes)
	if !strings.Contains(signedXml, "test data") {
		t.Error("signed XML should still contain original data")
	}

	validator, err := signedxml.NewValidator(signedXml)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	validator.Certificates = []x509.Certificate{*cert}

	refs, err := validator.ValidateReferences()
	if err != nil {
		t.Fatalf("XML signature validation failed: %v", err)
	}
	if len(refs) == 0 {
		t.Error("expected at least one validated reference")
	}
}

func TestHandleGetSignedXmlNilModule(t *testing.T) {
	srv := newTestServer(nil, pkcs11.CardVendorPosta)
	session := SmartboxSession{}

	data, _ := json.Marshal(Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: CertificateAlias{Alias: "0102"},
			Xml:         base64.StdEncoding.EncodeToString([]byte("<Root/>")),
		},
	})

	var buf bytes.Buffer
	err := srv.handleGetSignedXml(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error for nil module session")
	}
}

func TestHandleGetSignedXmlInvalidBase64(t *testing.T) {
	mock := &mockPkcsModuleSession{
		certs: []pkcs11.NamedCert{
			{Id: []byte{0x01}, Certificate: func() *x509.Certificate {
				c, _ := generateTestCert(time.Now().Add(-time.Hour), time.Now().Add(time.Hour), x509.KeyUsageDigitalSignature)
				return c
			}()},
		},
	}

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{moduleSession: mock}

	data, _ := json.Marshal(Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: CertificateAlias{Alias: "01"},
			Xml:         "not-valid-base64!@#$",
		},
	})

	var buf bytes.Buffer
	err := srv.handleGetSignedXml(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestHandleGetSignedXmlCertNotFound(t *testing.T) {
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

	srv := newTestServer(mock, pkcs11.CardVendorPosta)
	session := SmartboxSession{moduleSession: mock}

	data, _ := json.Marshal(Message[GetSignedXmlInput]{
		Operation: operationGetSignedXml,
		Input: GetSignedXmlInput{
			Certificate: CertificateAlias{Alias: "ff"},
			Xml:         base64.StdEncoding.EncodeToString([]byte("<Root/>")),
		},
	})

	var buf bytes.Buffer
	err := srv.handleGetSignedXml(&session, data, &buf)
	if err == nil {
		t.Fatal("expected error when certificate not found")
	}
}
