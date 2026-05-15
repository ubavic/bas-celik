package smartbox

import (
	"encoding/json"
	"io"
	"runtime"
)

func onOpen(w io.Writer) error {
	rsp := Response[OnOpenPayload]{
		Operation: onOpenEvent,
		Payload: OnOpenPayload{
			AppName:    "SmartBox",
			AppVersion: "2.0.0",
			AppBuild:   1,
			OsName:     runtime.GOOS,
			OsVersion:  "1",
			OsArch:     runtime.GOARCH,
		},
	}

	return json.NewEncoder(w).Encode(rsp)
}
