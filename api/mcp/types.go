package mcp

import (
	"encoding/json"
	"net/http"
)

const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

type McpRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Id      int         `json:"id,omitempty"`
}

type ClientInitializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
	ClientInfo      struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"clientInfo"`
	Capabilities struct {
		Root struct {
			ListChanged bool `json:"listChanged,omitempty"`
		} `json:"root"`
		Sampling struct {
			ListChanged bool `json:"listChanged,omitempty"`
		} `json:"sampling"`
		Experimental string `json:"experimental,omitempty"`
	} `json:"capabilities"`
}

type InitializeReqest struct {
	McpRequest
	Params ClientInitializeParams `json:"params"`
}

type McpResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id,omitempty"`
	Result  interface{} `json:"result"`
}

type ServerCapabilities struct {
	Prompts struct {
		ListChanged bool `json:"listChanged,omitempty"`
	} `json:"prompts"`
	Resources struct {
		ListChanged bool `json:"listChanged,omitempty"`
		Subscribe   bool `json:"subscribe,omitempty"`
	} `json:"resources"`
	Tools struct {
		ListChanged bool `json:"listChanged,omitempty"`
	} `json:"tools"`
	Logging      interface{} `json:"logging,omitempty"`
	Experimental string      `json:"experimental,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Instructions    []string           `json:"instructions,omitempty"`
	Capabilities    ServerCapabilities `json:"capabilities"`
}

type InitializeResponse struct {
	McpResponse
	Result ServerInitializeResult `json:"result"`
}

type McpNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}

type McpError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type McpErrorResponse struct {
	Jsonrpc string   `json:"jsonrpc"`
	Id      int      `json:"id,omitempty"`
	Error   McpError `json:"error"`
}

func RespondJsonRpcError(w http.ResponseWriter, status, id, code int, message string, data interface{}) {
	errMessage := McpError{
		Code:    code,
		Message: message,
		Data:    data,
	}
	response := McpErrorResponse{
		Jsonrpc: JSONRPC_VERSION,
		Id:      id,
		Error:   errMessage,
	}

	if status == 0 {
		status = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
