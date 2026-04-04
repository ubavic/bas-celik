package server

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type PkcsModuleSession interface {
	ListSlots() ([]uint, []string, error)
	OpenSessionAndLogin(pin string, terminalIndex int) error
	GetCertificates() ([]pkcs11.NamedCert, error)
	Sign(certId []byte, message []byte) ([]byte, error)
	CloseSession() error
}

type SmartboxSession struct {
	id            string
	module        PkcsModuleSession
	terminalId    int
	vendor        pkcs11.CardVendor
	certificateId []byte
}

func (ss *SmartboxSession) Certificate() *x509.Certificate {
	if ss.module == nil {
		return nil
	}

	namedCerts, err := ss.module.GetCertificates()
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
	if ss.module == nil {
		return nil, fmt.Errorf("nil module")
	}

	if ss.certificateId == nil {
		return nil, fmt.Errorf("certificate ID not set")
	}

	if opts.HashFunc() != crypto.SHA256 {
		return nil, fmt.Errorf("only sha256 supported")
	}

	signed, err := ss.module.Sign(ss.certificateId, digest)

	return signed, err
}

type ModulePath struct {
	Vendor pkcs11.CardVendor
	Path   string
}

func (s *SmartBoxServer) setModulePaths(paths []ModulePath) int {
	s.modulePaths = make([]ModulePath, 0, len(paths))

	for _, mPath := range paths {
		if mPath.Path == "" {
			continue
		}

		fInfo, err := os.Stat(mPath.Path)
		if err != nil {
			continue
		}

		if fInfo.IsDir() {
			continue
		}

		s.modulePaths = append(s.modulePaths, mPath)
	}

	return len(s.modulePaths)
}
