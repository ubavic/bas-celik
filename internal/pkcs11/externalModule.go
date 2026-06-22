package pkcs11

import (
	"crypto/x509"
	"errors"
	"fmt"
	"slices"
	"strings"
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

var gModuleContextsMu sync.Mutex
var gModuleContexts map[CardVendor]*pkcs11.Ctx

var ErrNilPkcsContext = fmt.Errorf("nil PKCS#11 context")
var ErrInvalidPkcsSession = fmt.Errorf("invalid PKCS#11 session")

func init() {
	gModuleContexts = make(map[CardVendor]*pkcs11.Ctx)
}

func GetLoadedVendors() []CardVendor {
	gModuleContextsMu.Lock()
	defer gModuleContextsMu.Unlock()

	vendors := make([]CardVendor, 0, len(gModuleContexts))
	for v, ctx := range gModuleContexts {
		if ctx != nil {
			vendors = append(vendors, v)
		}
	}
	return vendors
}

func Deinit() {
	gModuleContextsMu.Lock()
	defer gModuleContextsMu.Unlock()

	for _, c := range gModuleContexts {
		if c == nil {
			continue
		}

		slots, _ := c.GetSlotList(true)
		for _, s := range slots {
			c.CloseAllSessions(s)
		}

		c.Destroy()
	}
}

func LoadModules(modulePaths []ModulePath) ([]CardVendor, error) {
	var loadedVendors []CardVendor
	var errs []error

	var foundVendor []CardVendor
	for _, mp := range modulePaths {
		if slices.Contains(foundVendor, mp.Vendor) {
			return nil, fmt.Errorf("duplicate vendor found: %d %s", mp.Vendor, mp.Path)
		}
		foundVendor = append(foundVendor, mp.Vendor)
	}

	gModuleContextsMu.Lock()
	defer gModuleContextsMu.Unlock()

	for _, mp := range modulePaths {
		pkcsCtx := pkcs11.New(mp.Path)
		if pkcsCtx == nil {
			errs = append(errs, fmt.Errorf("loading module `%s`", mp.Path))
			continue
		}

		err := pkcsCtx.Initialize()
		if err != nil {
			errs = append(errs, fmt.Errorf("initializing module `%s`", mp.Path))
			continue
		}

		gModuleContexts[mp.Vendor] = pkcsCtx
		loadedVendors = append(loadedVendors, mp.Vendor)
	}

	return loadedVendors, errors.Join(errs...)
}

func GetPkcsSession(vendor CardVendor) (PkcsModuleSession, error) {
	gModuleContextsMu.Lock()
	defer gModuleContextsMu.Unlock()

	var ctx *pkcs11.Ctx

	ctx, ok := gModuleContexts[vendor]
	if !ok {
		return PkcsModuleSession{}, fmt.Errorf("module context not found")
	}

	return PkcsModuleSession{context: ctx}, nil
}

func (pm *PkcsModuleSession) ListSlots() ([]uint, []string, error) {
	if pm.context == nil {
		return nil, nil, ErrNilPkcsContext
	}

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
		// therefore this should be more reliable way to check if slot has token
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
	if pm.context == nil {
		return ErrNilPkcsContext
	}

	if slotId < 0 {
		return fmt.Errorf("invalid slot id: %d", slotId)
	}

	session, err := pm.context.OpenSession(uint(slotId), pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		pm.context.Destroy()
		pm.context = nil
		return fmt.Errorf("failed to open PKCS#11 session: %w", err)
	}

	err = pm.context.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		if strings.Contains(err.Error(), "CKR_USER_ALREADY_LOGGED_IN") {
			pm.session = session
			return nil
		}

		pm.context.CloseSession(session)
		return fmt.Errorf("failed to login to smart card: %w", err)
	}

	pm.session = session

	return nil
}

func (pm *PkcsModuleSession) getRawCertificates() ([][]byte, [][]byte, error) {
	if pm.context == nil {
		return nil, nil, ErrNilPkcsContext
	}

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
	if pm.context == nil {
		return nil, ErrNilPkcsContext
	}

	if pm.session == 0 {
		return nil, ErrInvalidPkcsSession
	}

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

func (pm *PkcsModuleSession) SignDigest(certId []byte, sha256digest []byte) ([]byte, error) {

	if pm.context == nil {
		return nil, ErrNilPkcsContext
	}

	if pm.session == 0 {
		return nil, ErrInvalidPkcsSession
	}

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
		pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil),
	}

	err = pm.context.SignInit(pm.session, mech, objects[0])
	if err != nil {
		return nil, fmt.Errorf("sign initialization: %w", err)
	}

	pkcsMessage := []byte{
		0x30, 0x31,
		0x30, 0x0d,
		0x06, 0x09,
		0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
		0x05, 0x00,
		0x04, 0x20,
	}
	pkcsMessage = append(pkcsMessage, sha256digest...)

	sig, err := pm.context.Sign(pm.session, pkcsMessage)
	if err != nil {
		return nil, fmt.Errorf("signing message: %w", err)
	}

	return sig, nil
}

func (pm *PkcsModuleSession) CloseSession() error {
	if pm.context == nil {
		return ErrNilPkcsContext
	}

	err1 := pm.context.Logout(pm.session)
	err2 := pm.context.CloseSession(pm.session)

	return errors.Join(err1, err2)
}
