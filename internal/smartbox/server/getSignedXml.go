package server

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"
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

type envelopeReq struct {
	XMLName   xml.Name `xml:"envelopaEPrijave"`
	XMLNS     string   `xml:"xmlns:ns2,attr"`
	Timestamp string   `xml:"timestamp"`
}

type envelopeResp struct {
	XMLName   xml.Name   `xml:"ns2:envelopaEPrijave"`
	XMLNS     string     `xml:"xmlns:ns2,attr"`
	Timestamp string     `xml:"timestamp"`
	Signature *signature `xml:"signature,omitempty"`
}

type signature struct {
	XMLName        xml.Name
	Xmlns          string `xml:"xmlns,attr"`
	SignedInfo     signedInfo
	SignatureValue string
	KeyInfo        struct {
		X509Data struct {
			X509Certificate  string
			X509IssuerSerial struct {
				X509IssuerName   string
				X509SerialNumber string
			}
			X509SubjectName string
		}
	}
}

type signedInfo struct {
	XMLName                xml.Name
	Xmlns                  string `xml:"xmlns,attr,omitempty"`
	Ns2                    string `xml:"xmlns:ns2,attr,omitempty"`
	CanonicalizationMethod struct {
		Algorithm string `xml:",attr"`
	}
	SignatureMethod struct {
		Algorithm string `xml:",attr"`
	}
	Reference struct {
		Uri        string `xml:"URI,attr"`
		Transforms struct {
			Transform struct {
				Algorithm string `xml:",attr"`
			}
		}
		DigestMethod struct {
			Algorithm string `xml:",attr"`
		}
		DigestValue string
	}
}

func (s *SmartBoxServer) handleGetSignedXml(session *SmartboxSession, data []byte, w io.Writer) error {
	msg := Message[GetSignedXmlInput]{}

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshaling sign request: %w", err)
	}

	if session.module == nil {
		return fmt.Errorf("pkcs11 module not loaded")
	}

	certId, err := hex.DecodeString(msg.Input.Certificate.Alias)
	if err != nil {
		return fmt.Errorf("decoding certificate alias: %w", err)
	}

	signXML, err := signRequest(session.module, certId, msg.Input.Xml)
	if err != nil {
		return err
	}

	rsp := Response[GetSignedXmlPayload]{
		Operation: operationGetCertificates,
		Payload: GetSignedXmlPayload{
			Xml: base64.StdEncoding.EncodeToString(signXML),
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}

func signRequest(module PkcsModuleSession, id []byte, base64XmlRequest string) ([]byte, error) {
	defer module.CloseSession()

	namedCerts, err := module.GetCertificates()
	if err != nil {
		return nil, fmt.Errorf("getting certificates: %w", err)
	}

	var cert *x509.Certificate
	for _, namedCert := range namedCerts {
		if slices.Equal(namedCert.Id, id) {
			cert = namedCert.Certificate
		}
	}

	if cert == nil {
		return nil, fmt.Errorf("certificate %s not found", hex.EncodeToString(id))
	}

	xmlString, err := base64.StdEncoding.DecodeString(base64XmlRequest)
	if err != nil {
		return nil, fmt.Errorf("decoding request: %w", err)
	}

	a, _ := bytes.CutPrefix(xmlString, []byte(xmlHeader))

	hash := sha256.Sum256(a)
	hashBase64 := base64.StdEncoding.EncodeToString(hash[:])

	signedInfo := constructSignedInfo(hashBase64)
	timestamp, err := extractTimestamp(xmlString)
	if err != nil {
		return nil, fmt.Errorf("extracting timestamp: %w", err)
	}

	marshaled := signedInfo.marshal()

	signed, err := module.Sign(id, marshaled)
	if err != nil {
		return nil, fmt.Errorf("signing request: %w", err)
	}

	envelope := constructResponse(cert, timestamp, signedInfo, signed)

	buf := bytes.Buffer{}
	_, err = buf.Write([]byte(xmlHeader))
	if err != nil {
		return nil, fmt.Errorf("writing xml header to the buffer: %w", err)
	}

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "")
	err = enc.Encode(envelope)
	if err != nil {
		return nil, fmt.Errorf("encoding envelope to the buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func extractTimestamp(input []byte) (string, error) {
	var env envelopeReq
	err := xml.Unmarshal(input, &env)
	if err != nil {
		return "", err
	}

	return env.Timestamp, nil
}

func constructResponse(cert *x509.Certificate, timestamp string, signedInfo signedInfo, signatureValue []byte) envelopeResp {
	signatureValueBase64 := base64.StdEncoding.EncodeToString(signatureValue)

	signedInfo.Xmlns = ""
	signedInfo.Ns2 = ""

	signature := signatureXML(cert, signedInfo, signatureValueBase64)

	envelope := envelopeResp{
		XMLNS:     "urn:poreskauprava.gov.rs/zim",
		Timestamp: timestamp,
		Signature: &signature,
	}

	return envelope
}

func constructSignedInfo(digestValue string) signedInfo {
	signedInfo := signedInfo{}

	signedInfo.XMLName = xml.Name{Local: "SignedInfo"}
	signedInfo.Xmlns = "http://www.w3.org/2000/09/xmldsig#"
	signedInfo.Ns2 = "urn:poreskauprava.gov.rs/zim"
	signedInfo.CanonicalizationMethod.Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	signedInfo.SignatureMethod.Algorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	signedInfo.Reference.DigestMethod.Algorithm = "http://www.w3.org/2001/04/xmlenc#sha256"
	signedInfo.Reference.DigestValue = digestValue
	signedInfo.Reference.Transforms.Transform.Algorithm = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"

	return signedInfo
}

func signatureXML(cert *x509.Certificate, signedInfo signedInfo, signatureValue string) signature {
	sig := signature{}
	sig.XMLName = xml.Name{Local: "Signature"}
	sig.Xmlns = "http://www.w3.org/2000/09/xmldsig#"

	sig.SignedInfo = signedInfo
	sig.SignatureValue = signatureValue
	sig.KeyInfo.X509Data.X509Certificate = base64.RawStdEncoding.EncodeToString(cert.Raw)
	sig.KeyInfo.X509Data.X509IssuerSerial.X509IssuerName = cert.Issuer.String()
	sig.KeyInfo.X509Data.X509IssuerSerial.X509SerialNumber = cert.SerialNumber.String()
	sig.KeyInfo.X509Data.X509SubjectName = formatSubject(cert.Subject)

	return sig
}

func (s *signedInfo) marshal() []byte {
	buf := bytes.Buffer{}
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "")
	enc.Encode(s)

	return buf.Bytes()
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
