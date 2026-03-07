package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

type Message[I any] struct {
	Operation string `json:"operation"`
	Input     I
}

type Response[P any] struct {
	Operation string `json:"operation"`
	Status    int    `json:"status"`
	Payload   P      `json:"payload"`
}

type OnOpenPayload struct {
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
	AppBuild   int    `json:"appBuild"`
	OsName     string `json:"osName"`
	OsVersion  string `json:"osVersion"`
	OsArch     string `json:"osArch"`
}

var onOpenEvent = "ON_OPEN_EVENT"
var operationGetInfo = "GET_INFO"
var operationGetProviders = "GET_PROVIDERS"
var operationGetTerminals = "GET_TERMINALS"
var operationGetCertificates = "GET_CERTIFICATES"
var operationGetSignedXml = "GET_SIGNED_XML"

type SmartBoxServer struct {
	sessions    map[string]SmartboxSession
	modulePaths []ModulePath
}

func (s *SmartBoxServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"eporezi.purs.gov.rs", "test.purs.gov.rs"},
	})

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.CloseNow()

	err = s.smartBoxHandler(conn)
	if err != nil {
		logger.Error(err)
	}
}

func (s *SmartBoxServer) smartBoxHandler(conn *websocket.Conn) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelCtx()

	w, err := conn.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}

	err = onOpen(w)
	if err != nil {
		return err
	}

	w.Close()

	sessionId := ""

	for {
		err = s.onMessage(&sessionId, ctx, conn)
		if err != nil {
			logger.Error(err)
			break
		}
	}

	return conn.CloseNow()
}

func (s *SmartBoxServer) onMessage(sessionId *string, ctx context.Context, conn *websocket.Conn) error {
	if sessionId == nil {
		return fmt.Errorf("nil session")
	}

	_, data, err := conn.Read(ctx)
	if err != nil {
		return fmt.Errorf("reading message in session %s: %w", *sessionId, err)
	}

	msg := Message[any]{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshaling data: %w", err)
	}

	w, err := conn.Writer(ctx, websocket.MessageText)
	if err != nil {
		return fmt.Errorf("getting writer: %w", err)
	}
	defer w.Close()

	session, ok := s.sessions[*sessionId]
	if !ok && msg.Operation != operationGetInfo {
		return fmt.Errorf("session %s not found, operation %s", *sessionId, msg.Operation)
	}

	switch msg.Operation {
	case operationGetInfo:
		err = s.handleGetInfo(sessionId, data, w)
	case operationGetProviders:
		err = s.handleGetProviders(data, w)
	case operationGetTerminals:
		err = s.handleGetTerminals(&session, data, w)
	case operationGetCertificates:
		err = s.handleGetCertificates(&session, data, w)
	case operationGetSignedXml:
		err = s.handleGetSignedXml(&session, data, w)
	default:
		err = fmt.Errorf("unknown operation %s in session %s", msg.Operation, *sessionId)
	}

	s.sessions[*sessionId] = session

	return err
}
