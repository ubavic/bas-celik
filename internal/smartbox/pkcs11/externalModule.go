package pkcs11

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sync"

	"github.com/miekg/pkcs11"
)

type NamedCert struct {
	Id          []byte
	Certificate *x509.Certificate
}

// wrapper around pkcs11.SessionHandle
// caches module certificates
type PkcsModuleSession struct {
	context *pkcs11.Ctx
	session pkcs11.SessionHandle
	certs   []NamedCert
}

// wrapper around pkcs11.Ctx
// it is used only for reference counting
type pkcsModuleCtx struct {
	context  *pkcs11.Ctx
	refCount uint
}

var mu sync.Mutex
var gModuleContexts map[string]pkcsModuleCtx

func init() {
	gModuleContexts = make(map[string]pkcsModuleCtx)
}

func NewPkcsExternalModule(modulePath string) (PkcsModuleSession, error) {
	mu.Lock()
	defer mu.Unlock()

	mc, ok := gModuleContexts[modulePath]
	if !ok {
		pkcsCtx := pkcs11.New(modulePath)

		err := pkcsCtx.Initialize()
		if err != nil {
			return PkcsModuleSession{}, fmt.Errorf("failed to initialize PKCS#11: %w", err)
		}

		mc = pkcsModuleCtx{context: pkcsCtx}
	}

	mc.refCount += 1
	gModuleContexts[modulePath] = mc

	return PkcsModuleSession{context: mc.context}, nil
}

func (pm *PkcsModuleSession) ListSlots() ([]uint, []string, error) {
	slots, err := pm.context.GetSlotList(true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get slot list: %w", err)
	}

	slotIds := make([]uint, 0, len(slots))
	slotNames := make([]string, 0, len(slots))

	for _, slot := range slots {
		info, err := pm.context.GetSlotInfo(slot)
		if err != nil {
			continue
		}

		// some modules (looking at you NetSet) don't properly set SlotInfo flags
		// therefore this should be more reliable way to check ig slot has token
		_, err = pm.context.GetTokenInfo(slot)
		if err != nil {
			continue
		}

		slotIds = append(slotIds, slot)
		slotNames = append(slotNames, info.SlotDescription)
	}

	return slotIds, slotNames, nil
}

func (pm *PkcsModuleSession) OpenSessionAndLogin(pin string, slotId int) error {
	if slotId < 0 {
		return fmt.Errorf("invalid slot id: %d", slotId)
	}

	session, err := pm.context.OpenSession(uint(slotId), pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		pm.context.Destroy()
		return fmt.Errorf("failed to open PKCS#11 session: %w", err)
	}

	err = pm.context.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		pm.context.CloseSession(session)
		pm.context.Destroy()
		return fmt.Errorf("failed to login to smart card: %w", err)
	}

	pm.session = session

	return nil
}

func (pm *PkcsModuleSession) getRawCertificates() ([][]byte, [][]byte, error) {
	searchTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_CERTIFICATE),
	}

	getTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, nil),
	}

	err := pm.context.FindObjectsInit(pm.session, searchTemplate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize object search: %w", err)
	}

	objects, _, err := pm.context.FindObjects(pm.session, 10)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find objects: %w", err)
	}

	err = pm.context.FindObjectsFinal(pm.session)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to finalize object search: %w", err)
	}

	ids := make([][]byte, 0, len(objects))
	certificates := make([][]byte, 0, len(objects))
	allErrors := []error{}
	for _, object := range objects {
		attr, err := pm.context.GetAttributeValue(pm.session, object, getTemplate)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		ids = append(ids, attr[0].Value)
		certificates = append(certificates, attr[1].Value)
	}

	return ids, certificates, errors.Join(allErrors...)
}

func (pm *PkcsModuleSession) GetCertificates() ([]NamedCert, error) {
	if len(pm.certs) > 0 {
		return pm.certs, nil
	}

	ids, rawCertificates, err := pm.getRawCertificates()
	if err != nil {
		return nil, err
	}

	pm.certs = []NamedCert{}
	allErrors := []error{}
	for i, rawCertificate := range rawCertificates {
		cert, err := x509.ParseCertificate(rawCertificate)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		pm.certs = append(pm.certs, NamedCert{Certificate: cert, Id: ids[i]})
	}

	return pm.certs, errors.Join(allErrors...)
}

func (pm *PkcsModuleSession) CloseSession() error {
	mu.Lock()
	defer mu.Unlock()

	err1 := pm.context.Logout(pm.session)
	err2 := pm.context.CloseSession(pm.session)

	for modulePath, mc := range gModuleContexts {
		if mc.context == pm.context {
			mc.refCount = mc.refCount - 1

			if mc.refCount > 0 {
				gModuleContexts[modulePath] = mc
			}

			break
		}
	}

	return errors.Join(err1, err2)
}

func (pm *PkcsModuleSession) Sign(certId []byte, message []byte) ([]byte, error) {
	err := pm.context.FindObjectsInit(pm.session, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_ID, certId),
	})
	if err != nil {
		return nil, fmt.Errorf("private key search initialization: %w", err)
	}

	objects, _, err := pm.context.FindObjects(pm.session, 1)
	if err != nil {
		return nil, fmt.Errorf("object search: %w", err)
	}

	if len(objects) == 0 {
		return nil, fmt.Errorf("no private key found")
	}

	err = pm.context.FindObjectsFinal(pm.session)
	if err != nil {
		return nil, fmt.Errorf("finalizing object search: %w", err)
	}

	mech := []*pkcs11.Mechanism{
		pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil),
	}

	err = pm.context.SignInit(pm.session, mech, objects[0])
	if err != nil {
		return nil, fmt.Errorf("sign initialization: %w", err)
	}

	sig, err := pm.context.Sign(pm.session, message)
	if err != nil {
		return nil, fmt.Errorf("signing message: %w", err)
	}

	return sig, nil
}
