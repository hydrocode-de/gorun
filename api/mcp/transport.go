package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func RespondInitialize(w http.ResponseWriter, body []byte) {
	sessionId, response, err := handleInitialize(body)
	if err != nil {
		RespondJsonRpcError(w, http.StatusBadRequest, 0, InternalError, err.Error(), nil)
		return
	}

	w.Header().Set("Mcp-Session-Id", sessionId)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		RespondJsonRpcError(w, http.StatusInternalServerError, 0, InternalError, err.Error(), nil)
		return
	}
}

func RespondInitialized(w http.ResponseWriter, r *http.Request) {
	sessionId, err := validateMcpSession(r.Header)
	if err != nil {
		RespondJsonRpcError(w, http.StatusBadRequest, 0, InvalidRequest, err.Error(), nil)
		return
	}

	err = handleInitialized(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(202)
}

func RespondShutdown(w http.ResponseWriter, r *http.Request) {
	sessionId, err := validateMcpSession(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	State.DeleteConnection(sessionId)
	w.Header().Set("Mcp-Session-Id", sessionId)
	w.WriteHeader(405)
}

func RespondToolsList(w http.ResponseWriter, r *http.Request, id int) {
	sessionId, err := validateMcpSession(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := ListTools(r.Context(), id)
	if err != nil {
		RespondJsonRpcError(w, http.StatusInternalServerError, 0, InternalError, err.Error(), nil)
		return
	}

	w.Header().Set("Mcp-Session-Id", sessionId)
	w.WriteHeader(202)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		RespondJsonRpcError(w, http.StatusInternalServerError, 0, InternalError, err.Error(), nil)
		return
	}
}

func validateMCPHeader(header http.Header) error {
	field := header.Get("Accept")
	accepts := strings.Split(field, ",")

	hasJson := false
	hasEventStream := false

	for _, accept := range accepts {
		if accept == "application/json" {
			hasJson = true
		}
		if accept == "text/event-stream" {
			hasEventStream = true
		}
	}

	if !hasJson {
		return fmt.Errorf("Accept header must contain application/json")
	}

	if hasEventStream {
		return fmt.Errorf("Accept header must contain text/event-stream")
	}

	return nil
}

func validateMcpSession(header http.Header) (string, error) {
	sessionId := header.Get("Mcp-Session-Id")
	if sessionId == "" {
		return "", fmt.Errorf("Mcp-Session-Id not present in header")
	}

	return sessionId, nil
}

func validateMCPMethod(body []byte) (string, int, error) {
	var request McpRequest
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&request)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode request: %w", err)
	}

	if request.Method == "" {
		return "", 0, fmt.Errorf("method is required")
	}

	return request.Method, request.Id, nil
}
