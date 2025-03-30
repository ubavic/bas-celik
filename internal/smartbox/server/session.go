package server

import (
	"os"

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
