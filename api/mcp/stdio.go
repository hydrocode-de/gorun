package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/toolImage"
)

type StdioTransport struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	sessionId string
	cache     *cache.Cache
}

func NewStdioTransport(cache *cache.Cache) *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
		cache:  cache,
	}
}

func (t *StdioTransport) readContentLength() (int, error) {
	for {
		line, err := t.reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		var length int
		if _, err = fmt.Sscanf(line, "Content-Length: %d", &length); err != nil {
			return 0, fmt.Errorf("invalid Content-Length: %s", line)
		}
		return length, nil
	}
}

func (t *StdioTransport) writeResponse(response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}
	payload := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	_, err = t.writer.WriteString(payload)
	if err != nil {
		return fmt.Errorf("failed to write response: %v", err)
	}
	return t.writer.Flush()

}
func (t *StdioTransport) writeError(id int, code int, message string, data interface{}) error {
	errorResponse := McpErrorResponse{
		Jsonrpc: JSONRPC_VERSION,
		Id:      id,
		Error: McpError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return t.writeResponse(errorResponse)
}

func (t *StdioTransport) readMessage(contentLength int) ([]byte, error) {
	if contentLength > 0 {
		// Read and discard the blank line after headers
		_, err := t.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		// Optionally check that sep is just "\r\n" or "\n"
		buf := make([]byte, contentLength)
		_, err = t.reader.Read(buf)
		return buf, err
	}

	// MCP fallback: read a single line (no embedded newlines allowed)
	line, err := t.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, "[MCP WARNING] Content-Length missing, using single-line fallback for message boundary.")
	return []byte(strings.TrimRight(line, "\r\n")), nil
}

func (t *StdioTransport) handleMessage(message []byte) error {
	method, id, err := validateMCPMethod(message)
	if err != nil {
		t.writeError(id, MethodNotFound, err.Error(), nil)
		return nil
	}

	switch method {
	case "ping":
		if err := t.writeResponse(Pong(id)); err != nil {
			t.writeError(id, InternalError, err.Error(), nil)
			return nil
		}
		return nil
	case "initialize":
		sessionId, response, err := handleInitialize(message)
		if err != nil {
			t.writeError(id, ParseError, err.Error(), nil)
			return nil
		}
		if err := t.writeResponse(response); err != nil {
			return err
		}
		t.sessionId = sessionId
		return nil
	case "notifications/initialized":
		if err := handleInitialized(t.sessionId); err != nil {
			t.writeError(id, InternalError, err.Error(), nil)
			return nil
		}
		return nil
	case "tools/list":
		result, err := ListTools(context.Background(), id)
		if err != nil {
			t.writeError(id, InternalError, err.Error(), nil)
			return nil
		}
		if err := t.writeResponse(result); err != nil {
			return err
		}

	}

	t.writeError(id, MethodNotFound, "Method not found", nil)
	return nil
}

func (t *StdioTransport) Start() error {
	go func() {
		toolImage.ReadAllTools(context.Background(), t.cache, false)
	}()

	for {
		contentLength, err := t.readContentLength()
		if err != nil {
			if err == io.EOF {
				return nil // Clean shutdown
			}
			return err
		}

		message, err := t.readMessage(contentLength)
		if err != nil {
			return err
		}

		if err := t.handleMessage(message); err != nil {
			return err
		}
	}
}
