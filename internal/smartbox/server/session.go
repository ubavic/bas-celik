package server

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"io"
	"slices"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type PkcsModuleSession interface {
	ListSlots() ([]uint, []string, error)
	OpenSessionAndLogin(pin string, terminalIndex int) error
	GetCertificates() ([]pkcs11.NamedCert, error)
	SignDigest(certId []byte, sha256digest []byte) ([]byte, error)
}

type SmartboxSession struct {
	id            string
	moduleSession PkcsModuleSession
	terminalId    int
	vendor        pkcs11.CardVendor
	certificateId []byte
}

func (ss *SmartboxSession) Certificate() *x509.Certificate {
	if ss.moduleSession == nil {
		return nil
	}

	namedCerts, err := ss.moduleSession.GetCertificates()
	if err != nil {
		return nil
	}

	var cert *x509.Certificate
	for _, namedCert := range namedCerts {
		if slices.Equal(namedCert.Id, ss.certificateId) {
			cert = namedCert.Certificate
		}
	}

	return cert
}

func (ss *SmartboxSession) Public() crypto.PublicKey {
	certificate := ss.Certificate()
	if certificate == nil {
		return nil
	}

	return certificate.PublicKey
}

func (ss *SmartboxSession) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if ss.moduleSession == nil {
		return nil, fmt.Errorf("nil module")
	}

	if ss.certificateId == nil {
		return nil, fmt.Errorf("certificate ID not set")
	}

	if opts.HashFunc() != crypto.SHA256 {
		return nil, fmt.Errorf("only sha256 supported")
	}

	signed, err := ss.moduleSession.SignDigest(ss.certificateId, digest)

	return signed, err
}
