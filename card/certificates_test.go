package card

import (
	"bytes"
	"compress/zlib"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ebfe/scard"
)

// All certificate data is generated for tests; no real card contents are stored.
type certificateTestCard struct {
	atr         Atr
	files       map[uint16][]byte
	selected    uint16
	readResult  []byte
	readErr     error
	statusErr   error
	commands    int
	vehicleFile bool
}

func (c *certificateTestCard) Status() (*scard.CardStatus, error) {
	return &scard.CardStatus{Atr: c.atr}, c.statusErr
}
func (c *certificateTestCard) BeginTransaction() error                { return nil }
func (c *certificateTestCard) EndTransaction(scard.Disposition) error { return nil }

func (c *certificateTestCard) Transmit(apdu []byte) ([]byte, error) {
	c.commands++
	switch apdu[1] {
	case 0xA4:
		if c.vehicleFile && bytes.Equal(apdu, buildAPDU(0, 0xA4, 2, 4, []byte{0xD0, 1}, 0)) {
			return []byte{0x90, 0}, nil
		}
		// The real PKS card accepts a card-manager AID also used during vehicle
		// detection, but rejects the vehicle application and its data files.
		if bytes.Equal(apdu, buildAPDU(0, 0xA4, 4, 0, []byte{0xA0, 0, 0, 1, 0x51, 0, 0}, 0)) {
			return []byte{0x90, 0}, nil
		}
		if bytes.Equal(apdu, buildAPDU(0, 0xA4, 4, 0, []byte("\xa0\x00\x00\x00\x63PKCS-15"), 0)) {
			return []byte{0x90, 0}, nil
		}
		if len(apdu) == 7 && apdu[2] == 0 && apdu[4] == 2 {
			c.selected = binary.BigEndian.Uint16(apdu[5:])
			if _, ok := c.files[c.selected]; ok {
				return []byte{0x90, 0}, nil
			}
		}
		return []byte{0x6A, 0x82}, nil
	case 0xB0:
		if c.readErr != nil || c.readResult != nil {
			return c.readResult, c.readErr
		}
		offset := int(apdu[2])<<8 | int(apdu[3])
		end := offset + int(apdu[4])
		data := c.files[c.selected]
		if end > len(data) {
			return []byte{0x62, 0x82}, nil
		}
		return append(bytes.Clone(data[offset:end]), 0x90, 0), nil
	default:
		return nil, fmt.Errorf("unexpected instruction %X", apdu[1])
	}
}

func compressedCertificate(t *testing.T, der []byte) []byte {
	t.Helper()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(der); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	file := append(make([]byte, 6), compressed.Bytes()...)
	binary.LittleEndian.PutUint16(file, uint16(len(file)-2))
	return file
}

func newCertificateTestCard(t *testing.T) *certificateTestCard {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	c := &certificateTestCard{atr: GEMALTO_ATR_3, files: make(map[uint16][]byte)}
	for _, id := range []uint16{0x7102, 0x7103} {
		template := &x509.Certificate{
			SerialNumber: big.NewInt(int64(id)),
			Subject:      pkix.Name{CommonName: "Synthetic certificate", Organization: []string{"Test organisation"}},
			NotBefore:    time.Unix(0, 0), NotAfter: time.Unix(2000000000, 0),
		}
		der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
		if err != nil {
			t.Fatal(err)
		}
		c.files[id] = compressedCertificate(t, der)
	}
	return c
}

func TestReadGemaltoCertificatesWithoutIdentityApplication(t *testing.T) {
	c := newCertificateTestCard(t)
	doc, err := DetectCardDocument(c)
	if !errors.Is(err, ErrUnknownCard) || doc == nil || !doc.Atr().Is(GEMALTO_ATR_3) {
		t.Fatalf("expected unknown document with shared ATR; got %T, %v", doc, err)
	}
	certs, err := ReadGemaltoCertificates(c)
	if err != nil || len(certs) != 2 {
		t.Fatalf("got %d certificates, error %v", len(certs), err)
	}
	for i, cert := range certs {
		if cert.SerialNumber.Int64() != int64(0x7102+i) || cert.Subject.CommonName != "Synthetic certificate" {
			t.Fatalf("unexpected certificate %d", i)
		}
	}
}

func TestReadGemaltoCertificatesRejectsUnknownATR(t *testing.T) {
	c := newCertificateTestCard(t)
	c.atr = Atr{1, 2, 3}
	certs, err := ReadGemaltoCertificates(c)
	if !errors.Is(err, ErrUnknownCard) || len(certs) != 0 || c.commands != 0 {
		t.Fatalf("unexpected result: %d certificates, %d APDUs, %v", len(certs), c.commands, err)
	}
	c.statusErr = errors.New("card removed")
	_, err = ReadGemaltoCertificates(c)
	if !errors.Is(err, c.statusErr) {
		t.Fatalf("lost status error: %v", err)
	}
}

func TestVehicleDocumentTakesPriorityOverCertificates(t *testing.T) {
	c := newCertificateTestCard(t)
	c.vehicleFile = true
	doc, err := DetectCardDocument(c)
	if _, ok := doc.(*VehicleCard); !ok || err != nil {
		t.Fatalf("vehicle document with its mandatory file was not recognized: %T, %v", doc, err)
	}
}

func TestCertificateReadErrors(t *testing.T) {
	for _, tc := range []struct {
		name     string
		response []byte
		err      error
	}{
		{"short status", []byte{0x90}, nil},
		{"denied", []byte{0x69, 0x82}, nil},
		{"empty data", []byte{0x90, 0}, nil},
		{"short header", []byte{1, 0x90, 0}, nil},
		{"too small", []byte{0, 0, 0x90, 0}, nil},
		{"too large", []byte{0xff, 0xff, 0x90, 0}, nil},
		{"removed", nil, scard.ErrRemovedCard},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newCertificateTestCard(t)
			c.readResult, c.readErr = tc.response, tc.err
			certs, err := ReadGemaltoCertificates(c)
			if err == nil || len(certs) != 0 {
				t.Fatalf("got %d certificates, error %v", len(certs), err)
			}
			if tc.err != nil && !errors.Is(err, tc.err) {
				t.Fatalf("lost transport error: %v", err)
			}
		})
	}
}

func TestCertificatePartialReadAndRetry(t *testing.T) {
	c := newCertificateTestCard(t)
	second := c.files[0x7103]
	delete(c.files, 0x7103)
	g := &Gemalto{smartCard: c}
	if err := g.LoadCertificates(); err == nil || len(g.GetCertificates()) != 1 {
		t.Fatalf("expected a partial result and error: %v", err)
	}
	c.files[0x7103] = second
	if err := g.LoadCertificates(); err != nil || len(g.GetCertificates()) != 2 {
		t.Fatalf("retry failed: %v", err)
	}
	// Callers must not be able to mutate the cached certificates.
	certs := g.GetCertificates()
	certs[0].Raw[0] = 0
	certs[0].Subject.Organization[0] = "changed"
	if g.GetCertificates()[0].Subject.Organization[0] != "Test organisation" {
		t.Fatal("GetCertificates exposes cached data")
	}
	commands := c.commands
	if err := g.LoadCertificates(); err != nil || commands != c.commands {
		t.Fatalf("successful read was not cached: %v", err)
	}
}

func TestCertificateMalformedContent(t *testing.T) {
	for _, tc := range []struct {
		name string
		file []byte
	}{
		{"invalid zlib", []byte{6, 0, 0, 0, 0, 0, 0, 0}},
		{"invalid DER", compressedCertificate(t, []byte("not X.509"))},
		{"decompression limit", compressedCertificate(t, bytes.Repeat([]byte{0}, 65537))},
		{"truncated file", []byte{0x20, 0, 0, 0, 0, 0, 0, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newCertificateTestCard(t)
			c.files[0x7102] = tc.file
			certs, err := ReadGemaltoCertificates(c)
			if err == nil || len(certs) != 1 || certs[0].SerialNumber.Int64() != 0x7103 {
				t.Fatalf("expected the other certificate to remain readable: %d, %v", len(certs), err)
			}
		})
	}
}
