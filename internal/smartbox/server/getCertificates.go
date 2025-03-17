package server

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"time"
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

func (s *SmartBoxServer) handleGetCertificates(session *Session, data []byte, w io.Writer) error {
	msg := Message[GetCertificatesInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if session.module == nil {
		return fmt.Errorf("pkcs11 module not loaded")
	}

	err := session.module.OpenSessionAndLogin(msg.Input.Pin, msg.Input.TerminalId)
	if err != nil {
		return err
	}

	_, certs, err := session.module.GetCertificates(msg.Input.Pin, msg.Input.TerminalId)
	if err != nil {
		return err
	}

	rsp := Response[GetCertificatesPayload]{
		Operation: operationGetCertificates,
		Payload: GetCertificatesPayload{
			Certificates: GetCertificateAliases(GetValidCertificates(certs)),
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}

func GetValidCertificates(certs []*x509.Certificate) []*x509.Certificate {
	now := time.Now()

	validCertificates := make([]*x509.Certificate, 0, len(certs))
	for _, cert := range certs {
		if now.Before(cert.NotBefore) {
			continue
		}

		if now.After(cert.NotAfter) {
			continue
		}

		validCertificates = append(validCertificates, cert)
	}

	return validCertificates
}

func GetCertificateAliases(certs []*x509.Certificate) []CertificateAlias {
	aliases := make([]CertificateAlias, 0, len(certs))
	for _, cert := range certs {
		alias := CertificateAlias{
			Alias: cert.Subject.String(),
			Name:  cert.Subject.CommonName,
		}
		aliases = append(aliases, alias)
	}

	return aliases
}
