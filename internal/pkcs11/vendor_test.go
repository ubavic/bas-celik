package pkcs11

import "testing"

func TestCardVendorString(t *testing.T) {
	tests := []struct {
		vendor CardVendor
		want   string
	}{
		{CardVendorHalcom, "Halcom CA"},
		{CardVendorPosta, "Sertifikaciono telo Pošte"},
		{CardVendorEsmart, "E-Smart Systems d.o.o. Beograd (ESS QCA)"},
		{CardVendorMup, "Ministarstvo unutrašnjih poslova CA"},
		{CardVendorPks, "Privredna komora Srbije CA"},
		{CardVendor(99), ""},
	}

	for _, tt := range tests {
		got := tt.vendor.String()
		if got != tt.want {
			t.Errorf("CardVendor(%d).String() = %q, want %q", tt.vendor, got, tt.want)
		}
	}
}

func TestGetDefaultPath(t *testing.T) {
	tests := []struct {
		vendor CardVendor
		os     string
		want   string
	}{
		{CardVendorPosta, "linux", "/usr/lib/libaetpkss.so"},
		{CardVendorEsmart, "linux", "/usr/lib/libeToken.so"},
		{CardVendorHalcom, "linux", ""},
		{CardVendorMup, "linux", ""},

		{CardVendorPosta, "darwin", "/Applications/tokenadmin.app/Contents/Frameworks/libaetpkss.dylib"},
		{CardVendorEsmart, "darwin", "/Library/Frameworks/eToken.framework/Versions/A/libIDPrimePKCS11.dylib"},
		{CardVendorHalcom, "darwin", "/Applications/Personal.app/Contents/Frameworks/libtokenapi.dylib"},
		{CardVendorMup, "darwin", ""},

		{CardVendorHalcom, "windows", `C:\Program Files (x86)\Personal\bin64\personal64.dll`},
		{CardVendorPosta, "windows", `C:\Windows\System32\aetpkss1.dll`},
		{CardVendorEsmart, "windows", `C:\Program Files\SafeNet\Authentication\SAC\x64\IDPrimePKCS1164.dll`},
		{CardVendorMup, "windows", `C:\Program Files\TrustEdgeID\netsetpkcs11_x64.dll`},
		{CardVendorPks, "windows", `C:\Program Files\TrustEdgeID\netsetpkcs11_x64.dll`},

		{CardVendorHalcom, "freebsd", ""},
		{CardVendor(99), "linux", ""},
	}

	for _, tt := range tests {
		got := GetDefaultPath(tt.vendor, tt.os)
		if got != tt.want {
			t.Errorf("GetDefaultPath(%d, %q) = %q, want %q", tt.vendor, tt.os, got, tt.want)
		}
	}
}
