package server

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/beevik/etree"
	"github.com/moov-io/signedxml"
	"github.com/moov-io/signedxml/xmlenc"
)

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`

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

	doc := etree.NewDocument()

	err = doc.ReadFromBytes(xmlToSign)
	if err != nil {
		return fmt.Errorf("reading xml: %w", err)
	}

	err = populateSignature(doc.Root(), cert)
	if err != nil {
		return fmt.Errorf("adding Sign element: %w", err)
	}

	xmlSigner, err := signedxml.NewSignerFromDoc(doc)
	if err != nil {
		return fmt.Errorf("creating xml signer: %w", err)
	}

	signedXML, err := xmlSigner.Sign(session)
	if err != nil {
		return fmt.Errorf("creating xml signer: %w", err)
	}

	rsp := Response[GetSignedXmlPayload]{
		Operation: operationGetSignedXml,
		Payload: GetSignedXmlPayload{
			Xml: base64.StdEncoding.EncodeToString([]byte(signedXML)),
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}

func formatSubject(subject pkix.Name) string {
	var serNo1, serNo2, givenName, surname, country string

	for _, name := range subject.Names {
		strB, ok := name.Value.(string)
		if !ok {
			continue
		}

		if slices.Equal(name.Type, []int{2, 5, 4, 5}) {
			if serNo1 == "" {
				serNo1 = string(strB)
			} else {
				serNo2 = string(strB)
			}
		} else if slices.Equal(name.Type, []int{2, 5, 4, 42}) {
			givenName = string(strB)
		} else if slices.Equal(name.Type, []int{2, 5, 4, 4}) {
			surname = string(strB)
		}
	}

	if strings.Contains(serNo2, "PNORS") {
		serNo1, serNo2 = serNo2, serNo1
	}

	if len(subject.Country) > 0 {
		country = subject.Country[0]
	}

	return fmt.Sprintf("CN=%s, GIVENNAME=%s, SURNAME=%s, SERIALNUMBER=%s, SERIALNUMBER=%s, C=%s",
		subject.CommonName,
		givenName,
		surname,
		serNo1,
		serNo2,
		country)
}

func populateSignature(root *etree.Element, cert *x509.Certificate) error {
	signature := root.FindElement("./Signature")
	if signature == nil {
		signature = etree.NewElement("Signature")
		root.AddChild(signature)
	} else {
		return fmt.Errorf("Signature element already exists")
	}

	signature.Child = nil

	signature.CreateAttr("xmlns", xmlenc.NamespaceXMLDSig)

	signedInfo := signature.CreateElement("SignedInfo")

	canon := signedInfo.CreateElement("CanonicalizationMethod")
	canon.CreateAttr("Algorithm", "http://www.w3.org/TR/2001/REC-xml-c14n-20010315")

	sigMethod := signedInfo.CreateElement("SignatureMethod")
	sigMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256")

	ref := signedInfo.CreateElement("Reference")
	ref.CreateAttr("URI", "")

	transforms := ref.CreateElement("Transforms")
	transform := transforms.CreateElement("Transform")
	transform.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")

	digestMethod := ref.CreateElement("DigestMethod")
	digestMethod.CreateAttr("Algorithm", xmlenc.AlgorithmSHA256)

	ref.CreateElement("DigestValue")
	sigValue := signature.CreateElement("SignatureValue")
	sigValue.SetText("")

	keyInfo := signature.CreateElement("KeyInfo")

	x509Data := keyInfo.CreateElement("X509Data")

	x509Cert := x509Data.CreateElement("X509Certificate")
	x509Cert.SetText(base64.RawStdEncoding.EncodeToString(cert.Raw))

	x509IssuerSerial := x509Data.CreateElement("X509IssuerSerial")

	x509IssuerName := x509IssuerSerial.CreateElement("X509IssuerName")
	x509IssuerName.SetText(cert.Issuer.String())

	x509Serial := x509IssuerSerial.CreateElement("X509SerialNumber")
	x509Serial.SetText(cert.SerialNumber.String())

	x509Subject := x509Data.CreateElement("X509SubjectName")
	x509Subject.SetText(formatSubject(cert.Subject))

	return nil

}
