package xmlsign

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/beevik/etree"
	"github.com/moov-io/signedxml"
	"github.com/moov-io/signedxml/xmlenc"
)

func FormatSubject(subject pkix.Name) string {
	var serNo1, serNo2, givenName, surname string

	for _, name := range subject.Names {
		strB, ok := name.Value.(string)
		if !ok {
			continue
		}

		if slices.Equal(name.Type, []int{2, 5, 4, 5}) {
			if serNo1 == "" {
				serNo1 = strB
			} else {
				serNo2 = strB
			}
		} else if slices.Equal(name.Type, []int{2, 5, 4, 42}) {
			givenName = strB
		} else if slices.Equal(name.Type, []int{2, 5, 4, 4}) {
			surname = strB
		}
	}

	if strings.Contains(serNo2, "PNORS") {
		serNo1, serNo2 = serNo2, serNo1
	}

	country := ""
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

func PopulateSignature(root *etree.Element, cert *x509.Certificate) error {
	signature := root.FindElement("./Signature")
	if signature != nil {
		return fmt.Errorf("Signature element already exists")
	}

	signature = etree.NewElement("Signature")
	root.AddChild(signature)
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
	x509Subject.SetText(FormatSubject(cert.Subject))

	return nil
}

func SignXMLBytes(xmlBytes []byte, signer crypto.Signer, cert *x509.Certificate) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlBytes); err != nil {
		return "", fmt.Errorf("parsing XML: %w", err)
	}

	if err := PopulateSignature(doc.Root(), cert); err != nil {
		return "", fmt.Errorf("populating signature skeleton: %w", err)
	}

	xmlSigner, err := signedxml.NewSignerFromDoc(doc)
	if err != nil {
		return "", fmt.Errorf("creating XML signer: %w", err)
	}

	signed, err := xmlSigner.Sign(signer)
	if err != nil {
		return "", fmt.Errorf("signing XML: %w", err)
	}

	return signed, nil
}
