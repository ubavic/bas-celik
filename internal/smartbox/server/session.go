package server

import (
	"crypto/x509"
	"os"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type PkcsModule interface {
	ListSlots() ([]uint, []string, error)
	OpenSessionAndLogin(pin string, terminalIndex int) error
	GetCertificates(pin string, terminalIndex int) ([][]byte, []*x509.Certificate, error)
	Sign(certId []byte, message []byte) ([]byte, error)
	CloseSession(terminalIndex int) error
}

type Session struct {
	id            string
	module        PkcsModule
	terminalId    int
	vendor        pkcs11.CardVendor
	certificateId string
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
