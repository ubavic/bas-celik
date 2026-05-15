package smartbox

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

type GetCertificatesInput struct {
	TerminalId int    `json:"terminalId"`
	Pin        string `json:"pin"`
}

type GetCertificatesPayload struct {
	Certificates []CertificateAlias `json:"certificates"`
}

type CertificateAlias struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

func (s *SmartBoxServer) handleGetCertificates(session *SmartboxSession, data []byte, w io.Writer) error {
	msg := Message[GetCertificatesInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if session.moduleSession == nil {
		return fmt.Errorf("pkcs11 module not loaded")
	}

	err := session.moduleSession.OpenSessionAndLogin(msg.Input.Pin, msg.Input.TerminalId)
	if err != nil {
		return err
	}

	certs, err := session.moduleSession.GetCertificates()
	if err != nil {
		return err
	}

	rsp := Response[GetCertificatesPayload]{
		Operation: operationGetCertificates,
		Payload: GetCertificatesPayload{
			Certificates: GetCertificateAliases(GetValidCertificates(certs, session.vendor)),
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}

func GetValidCertificates(namedCerts []pkcs11.NamedCert, cardVendor pkcs11.CardVendor) []pkcs11.NamedCert {
	now := time.Now()

	validNamedCertificates := make([]pkcs11.NamedCert, 0, len(namedCerts))
	for _, namedCert := range namedCerts {
		if now.Before(namedCert.Certificate.NotBefore) {
			continue
		}

		if now.After(namedCert.Certificate.NotAfter) {
			continue
		}

		if cardVendor == pkcs11.CardVendorMup {
			if namedCert.Certificate.KeyUsage&x509.KeyUsageContentCommitment == 0 {
				continue
			}
		}

		validNamedCertificates = append(validNamedCertificates, namedCert)
	}

	return validNamedCertificates
}

func GetCertificateAliases(namedCerts []pkcs11.NamedCert) []CertificateAlias {
	aliases := make([]CertificateAlias, 0, len(namedCerts))
	for _, namedCert := range namedCerts {
		alias := CertificateAlias{
			Alias: hex.EncodeToString(namedCert.Id),
			Name:  namedCert.Certificate.Subject.CommonName,
		}
		aliases = append(aliases, alias)
	}

	return aliases
}
