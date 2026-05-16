package docSigning

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"io"

	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

type pkcs11Signer struct {
	session *pkcs11.PkcsModuleSession
	certId  []byte
	cert    *x509.Certificate
}

func (s *pkcs11Signer) Public() crypto.PublicKey { return s.cert.PublicKey }

func (s *pkcs11Signer) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	if opts.HashFunc() != crypto.SHA256 {
		return nil, fmt.Errorf("only SHA-256 supported")
	}
	return s.session.SignDigest(s.certId, digest)
}
