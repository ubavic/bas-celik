package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ubavic/bas-celik/v2/internal/logger"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

const bindHost = "127.0.0.1"

func defaultSessionProvider(vendor pkcs11.CardVendor) (PkcsModuleSession, error) {
	ms, err := pkcs11.GetPkcsSession(vendor)
	if err != nil {
		return nil, err
	}
	return &ms, nil
}

func StartServer(modulePaths []pkcs11.ModulePath) (string, []string, error) {
	smartboxServer := SmartBoxServer{
		sessions:        make(map[string]SmartboxSession),
		sessionProvider: defaultSessionProvider,
	}

	loadedVendors, err := pkcs11.LoadModules(modulePaths)
	if len(loadedVendors) == 0 {
		if err != nil {
			return "", nil, fmt.Errorf("no valid pkcs11 path loaded: %w", err)
		}

		return "", nil, fmt.Errorf("no valid pkcs11 path loaded")
	}

	smartboxServer.loadedVendors = loadedVendors

	port := findAvailablePort()
	if port == "" {
		return "", nil, fmt.Errorf("no valid port available")
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(bindHost, port))
	if err != nil {
		return "", nil, err
	}

	s := &http.Server{
		Handler:      &smartboxServer,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		s.Serve(listener)
	}()

	logger.Info(fmt.Sprintf("Smartbox server listening on ws://%v", listener.Addr()))

	loadedVendorNames := make([]string, 0, len(loadedVendors))
	for _, m := range loadedVendors {
		loadedVendorNames = append(loadedVendorNames, m.String())
	}

	return port, loadedVendorNames, nil
}

func findAvailablePort() string {
	ports := []string{"17165", "20806", "65097"}

	for _, port := range ports {
		address := net.JoinHostPort(bindHost, port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port
		}
	}

	return ""
}
