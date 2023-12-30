package gui

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	host "github.com/rickypc/native-messaging-host"
	"github.com/ubavic/bas-celik/document"
)

type NativeStatus int

const (
	nsNoError NativeStatus = iota
	nsNoReader
	nsNoCard
	nsCardError
)

var previousSatus NativeStatus
var messagingHost *host.Host

func StartNativeMessaging() {
	exec, _ := os.Executable()
	evaled, _ := filepath.EvalSymlinks(exec)
	executable, _ := filepath.Abs(evaled)

	//TODO: get hash from parameters
	messagingHost = (&host.Host{
		AppName:     "com.github.ubavic.bas_celik",
		AppDesc:     "Tool for reading smart-card documents issued by the government of Serbia",
		AllowedExts: []string{"chrome-extension://{hash}/"},
		ExecName:    executable + " -verbose",
	}).Init()

	if err := messagingHost.Install(); err != nil {
		log.Printf("install error: %v", err)
	}
}

func sendNativeDoc(doc document.Document) {
	data, err := doc.BuildJson()
	if err != nil {
		log.Fatalf("generating json: %w", err)
		return
	}

	var message map[string]interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		log.Fatalf("generating json: %w", err)
		return
	}

	if err := messagingHost.PostMessage(os.Stdout, message); err != nil {
		log.Fatalf("messaging.PostMessage error: %v", err)
	}
}

func setNativeStatus(status NativeStatus) {
	var message *host.H
	if previousSatus != status {
		previousSatus = status
		switch status {
		case nsNoReader:
			message = &host.H{"error": 1, "message": "No Reader"}
		case nsNoCard:
			message = &host.H{"error": 2, "message": "No Card"}
		case nsCardError:
			message = &host.H{"error": 2, "message": "Card Error"}
		}

	}

	if message != nil {
		if err := messagingHost.PostMessage(os.Stdout, message); err != nil {
			log.Fatalf("messaging.PostMessage error: %v", err)
		}
	}

}
