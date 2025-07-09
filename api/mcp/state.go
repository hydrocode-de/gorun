package mcp

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ClientConnection struct {
	ID          string
	Params      *ClientInitializeParams
	Initialized bool
}

type McpState struct {
	Connections map[string]*ClientConnection
}

var State *McpState

func (s *McpState) AddConnection(params *ClientInitializeParams) string {
	logger := viper.Get("logger").(*logrus.Logger)
	id := s.getNextConnectionID()

	conn := &ClientConnection{
		ID:          id,
		Params:      params,
		Initialized: false,
	}

	s.Connections[id] = conn

	logger.Debugf("Added connection %s: %v", id, conn)
	return id
}

func (s *McpState) getNextConnectionID() string {
	id := uuid.New().String()

	return id
}

func (s *McpState) GetConnection(sessionId string) (*ClientConnection, bool) {
	conn, ok := s.Connections[sessionId]
	return conn, ok
}

func (s *McpState) SetInitialized(sessionId string) {
	logger := viper.Get("logger").(*logrus.Logger)
	logger.Debugf("Setting initialized for connection %s", sessionId)
	s.Connections[sessionId].Initialized = true
}

func (s *McpState) IsInitialized(sessionId string) bool {
	return s.Connections[sessionId].Initialized
}

func (s *McpState) DeleteConnection(sessionId string) {
	delete(s.Connections, sessionId)
}

func init() {
	State = &McpState{
		Connections: make(map[string]*ClientConnection, 10),
	}
}
