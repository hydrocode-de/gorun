package mcp

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var LATEST_PROTOCOL_VERSION string = "2025-03-26"
var MIN_PROTOCOL_VERSION time.Time = time.Date(2025, 3, 26, 0, 0, 0, 0, time.UTC)
var JSONRPC_VERSION string = "2.0"

func HandleMCP(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	body, err := io.ReadAll(r.Body)
	if err != nil {
		RespondJsonRpcError(w, http.StatusBadRequest, 0, InvalidRequest, err.Error(), nil)
		return
	}
	err = validateMCPHeader(header)
	if err != nil {
		RespondJsonRpcError(w, http.StatusBadRequest, 0, InvalidRequest, err.Error(), nil)
		return
	}

	method, id, err := validateMCPMethod(body)
	if err != nil {
		RespondJsonRpcError(w, http.StatusBadRequest, id, InvalidRequest, err.Error(), nil)
		return
	}

	switch method {
	case "initialize":
		RespondInitialize(w, body)
	case "notifications/initialized":
		RespondInitialized(w, r)
	case "tools/list":
		RespondToolsList(w, r, id)
	default:
		RespondJsonRpcError(w, http.StatusBadRequest, id, MethodNotFound, "Invalid method", nil)
	}
}

func HandleMCPStream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
}

func InitMCP(mux *http.ServeMux, prefix string) (*http.ServeMux, error) {

	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:len(prefix)-1]
	}
	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}
	mux.HandleFunc(fmt.Sprintf("POST /%s/mcp", prefix), HandleMCP)
	mux.HandleFunc(fmt.Sprintf("GET /%s/mcp", prefix), HandleMCPStream)
	mux.HandleFunc(fmt.Sprintf("DELETE /%s/mcp", prefix), RespondShutdown)
	return mux, nil
}
