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

	ports := []uint{17165, 20806, 65097}
	port := findAvailablePort(ports)
	if port == 0 {
		return "", fmt.Errorf("no valid port available")
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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

func findAvailablePort(ports []uint) uint {
	for _, port := range ports {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port
		}
	}

	return 0
}
