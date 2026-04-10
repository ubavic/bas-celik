package gui

import (
	"testing"

	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

type ModulePathStub func(string) string

func (f ModulePathStub) String(key string) string {
	return f(key)
}

func TestGetValidModulePaths(t *testing.T) {
	modulePath := "mup.dll"

	mps := ModulePathStub(func(key string) string {
		switch key {
		case mupPkcsPathKey:
			return modulePath
		default:
			return ""
		}
	})

	mp := getValidModulePaths(mps)

	if len(mp) != 1 {
		t.Errorf("Expected 1 module, but got %d", len(mp))
	}

	if mp[0].Vendor != pkcs11.CardVendorMup {
		t.Errorf("Expected vendor '%s', but got '%s'", pkcs11.CardVendorMup, mp[0].Vendor)
	}

	if mp[0].Path != modulePath {
		t.Errorf("Expected path '%s', but got '%s'", modulePath, mp[0].Path)
	}
}

func TestGetValidModulePathsNoModules(t *testing.T) {
	mps := ModulePathStub(func(key string) string {
		switch key {
		case mupPkcsPathKey:
			return "  "
		default:
			return ""
		}
	})

	mp := getValidModulePaths(mps)

	if len(mp) != 0 {
		t.Errorf("Expected 0 modules, but got %d", len(mp))
	}
}
