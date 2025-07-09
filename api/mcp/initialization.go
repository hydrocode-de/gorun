package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"github.com/hydrocode-de/gorun/version"
)

func handleInitialize(body []byte) (string, *InitializeResponse, error) {
	var params InitializeReqest
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&params)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode initialize request: %w", err)
	}

	if params.Method != "initialize" {
		return "", nil, fmt.Errorf("invalid method: %s", params.Method)
	}

	// the protocol version is a date
	clientVersion, err := time.Parse("2006-01-02", params.Params.ProtocolVersion)
	if err != nil {
		return "", nil, fmt.Errorf("invalid protocol version: %w", err)
	}
	if clientVersion.Before(MIN_PROTOCOL_VERSION) {
		return "", nil, fmt.Errorf("protocol version too old: %s. This server requires at least %s and supports up to %s", params.Params.ProtocolVersion, MIN_PROTOCOL_VERSION.Format("2006-01-02"), LATEST_PROTOCOL_VERSION)
	}

	// if we are here, we can add the connection to the state, as a version could be negotiated
	sessionId := State.AddConnection(&params.Params)

	var negotiatedVersion string
	if params.Params.ProtocolVersion == LATEST_PROTOCOL_VERSION {
		negotiatedVersion = params.Params.ProtocolVersion
	} else {
		negotiatedVersion = LATEST_PROTOCOL_VERSION
	}

	response := InitializeResponse{
		McpResponse: McpResponse{
			Jsonrpc: JSONRPC_VERSION,
			Id:      params.Id,
		},
		Result: ServerInitializeResult{
			ProtocolVersion: negotiatedVersion,
			ServerInfo: ServerInfo{
				Name:    "GoRun",
				Version: version.Version,
			},
			Capabilities: ServerCapabilities{
				Tools: struct {
					ListChanged bool "json:\"listChanged,omitempty\""
				}{
					ListChanged: false,
				},
			},
		},
	}

	return sessionId, &response, nil
}

func handleInitialized(sessionId string) error {
	_, ok := State.GetConnection(sessionId)
	if !ok {
		return fmt.Errorf("session not found")
	}

	State.SetInitialized(sessionId)
	return nil
}
