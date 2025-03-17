package server

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

type GetSignedXmlInput struct {
	Certificate CertificateAlias `json:"certificate"`
	Pin         string           `json:"pin"`
	Xml         string           `json:"xml"`
}

type GetSignedXmlPayload struct {
	Xml string `json:"xml"`
}

type envelope struct {
	XMLName   xml.Name   `xml:"envelopaEPrijave"`
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

func (s *SmartBoxServer) handleGetSignedXml(session *Session, data []byte, w io.Writer) error {
	msg := Message[GetSignedXmlInput]{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if session.module == nil {
		return fmt.Errorf("pkcs11 module not loaded")
	}

	signXML, err := signRequest(session.module, msg.Input.Certificate.Alias, session.terminalId, msg.Input.Pin, msg.Input.Xml)
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

func signRequest(module PkcsModule, certificateAlias string, terminalId int, pin string, base64XmlRequest string) ([]byte, error) {
	ids, _, err := module.GetCertificates(pin, terminalId)
	if err != nil {
		return nil, err
	}

	xmlString, err := base64.StdEncoding.DecodeString(base64XmlRequest)
	if err != nil {
		return nil, err
	}

	a := bytes.Replace(xmlString, []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`), []byte{}, 1)

	hash := sha256.Sum256(a)
	hashBase64 := base64.StdEncoding.EncodeToString(hash[:])

	signedInfo := constructSignedInfo(hashBase64)
	timestamp, err := extractTimestamp(xmlString)
	if err != nil {
		return nil, err
	}

	marshaled := signedInfo.marshal()

	//todo
	signed, err := module.Sign(ids[0], marshaled)
	if err != nil {
		return nil, err
	}

	module.CloseSession(terminalId)

	envelope := constructResponse(timestamp, signedInfo, signed)

	buf := bytes.Buffer{}
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "")
	enc.Encode(envelope)

	return buf.Bytes(), nil
}

func extractTimestamp(input []byte) (string, error) {
	var env envelope
	err := xml.Unmarshal(input, &env)
	if err != nil {
		return "", err
	}

	return env.Timestamp, nil
}

func constructResponse(timestamp string, signedInfo signedInfo, signatureValue []byte) envelope {
	signatureValueBase64 := base64.StdEncoding.EncodeToString(signatureValue)

	signedInfo.Xmlns = ""
	signedInfo.Ns2 = ""

	signature := signatureXML(nil, signedInfo, signatureValueBase64)

	envelope := envelope{
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
	sig.KeyInfo.X509Data.X509SubjectName = cert.Subject.String()

	return sig
}

func (s *signedInfo) marshal() []byte {
	buf := bytes.Buffer{}
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "")
	enc.Encode(s)

	return buf.Bytes()
}
