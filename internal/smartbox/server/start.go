package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func StartServer(modulePaths []ModulePath) (string, error) {
	smartboxServer := SmartBoxServer{}
	smartboxServer.sessions = make(map[string]SmartboxSession)

	modules := smartboxServer.setModulePaths(modulePaths)
	if modules == 0 {
		return "", fmt.Errorf("no valid pkcs11 path loaded")
	}

	address := findAvailablePort()
	if address == "" {
		return "", fmt.Errorf("no valid port available")
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return "", err
	}

	s := &http.Server{
		Handler:      &smartboxServer,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		s.Serve(l)
	}()

	logger.Info(fmt.Sprintf("Smartbox server listening on ws://%v", l.Addr()))

	return l.Addr().String(), nil
}

func findAvailablePort() string {
	bindHost := "127.0.0.1"
	ports := []string{"17165", "20806", "65097"}

	for _, port := range ports {
		address := net.JoinHostPort(bindHost, port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return address
		}
	}

	return ""
}
