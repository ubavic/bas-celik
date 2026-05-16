package smartbox

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ubavic/bas-celik/v2/internal/xmlsign"
)

type GetSignedXmlInput struct {
	Certificate CertificateAlias `json:"certificate"`
	Pin         string           `json:"pin"`
	Xml         string           `json:"xml"`
}

type GetSignedXmlPayload struct {
	Xml string `json:"xml"`
}

func (s *SmartBoxServer) handleGetSignedXml(session *SmartboxSession, data []byte, w io.Writer) error {
	msg := Message[GetSignedXmlInput]{}

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshaling sign request: %w", err)
	}

	if session.moduleSession == nil {
		return fmt.Errorf("pkcs11 module not loaded")
	}

	session.certificateId, err = hex.DecodeString(msg.Input.Certificate.Alias)
	if err != nil {
		return fmt.Errorf("decoding certificate alias: %w", err)
	}

	cert := session.Certificate()
	if cert == nil {
		return fmt.Errorf("certificate %s not found", hex.EncodeToString(session.certificateId))
	}

	xmlToSign, err := base64.RawStdEncoding.DecodeString(strings.TrimRight(msg.Input.Xml, "="))
	if err != nil {
		return fmt.Errorf("decoding base64 payload: %w", err)
	}

	signedXML, err := xmlsign.SignXMLBytes(xmlToSign, session, cert)
	if err != nil {
		return fmt.Errorf("signing xml: %w", err)
	}

	rsp := Response[GetSignedXmlPayload]{
		Operation: operationGetSignedXml,
		Payload: GetSignedXmlPayload{
			Xml: base64.StdEncoding.EncodeToString([]byte(signedXML)),
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}
