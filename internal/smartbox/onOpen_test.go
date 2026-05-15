package smartbox

import (
	"bytes"
	"encoding/json"
	"runtime"
	"testing"
)

func TestOnOpen(t *testing.T) {
	var buf bytes.Buffer
	err := onOpen(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rsp Response[OnOpenPayload]
	if err := json.Unmarshal(buf.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rsp.Operation != onOpenEvent {
		t.Errorf("operation = %q, want %q", rsp.Operation, onOpenEvent)
	}

	if rsp.Payload.AppName != "SmartBox" {
		t.Errorf("appName = %q, want %q", rsp.Payload.AppName, "SmartBox")
	}

	if rsp.Payload.AppVersion != "2.0.0" {
		t.Errorf("appVersion = %q, want %q", rsp.Payload.AppVersion, "2.0.0")
	}

	if rsp.Payload.AppBuild != 1 {
		t.Errorf("appBuild = %d, want %d", rsp.Payload.AppBuild, 1)
	}

	if rsp.Payload.OsName != runtime.GOOS {
		t.Errorf("osName = %q, want %q", rsp.Payload.OsName, runtime.GOOS)
	}

	if rsp.Payload.OsArch != runtime.GOARCH {
		t.Errorf("osArch = %q, want %q", rsp.Payload.OsArch, runtime.GOARCH)
	}
}
