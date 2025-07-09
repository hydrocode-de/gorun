package mcp

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var sharedCache *cache.Cache

func TestMain(m *testing.M) {
	sharedCache = &cache.Cache{}
	sharedCache.Reset()

	viper.Set("logger", logrus.New())

	code := m.Run()
	os.Exit(code)
}

func TestStdioTransport_Ping(t *testing.T) {
	// Prepare a valid JSON-RPC ping request with Content-Length
	json := `{"jsonrpc":"2.0","method":"ping","id":1}`
	content := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(json), json)

	// Set up input and output buffers
	in := bytes.NewBufferString(content)
	out := &bytes.Buffer{}

	// Create the transport with our buffers
	transport := &StdioTransport{
		reader: bufio.NewReader(in),
		writer: bufio.NewWriter(out),
		cache:  sharedCache,
	}

	// Run the Start loop once (simulate one message)
	go func() {
		_ = transport.Start()
	}()

	// Wait for output
	// (In a real test, you might want to use sync or context for coordination)
	// For now, just check after a short delay
	time.Sleep(100 * time.Millisecond) // Uncomment if needed

	got := out.String()
	if !strings.Contains(got, `"pong"`) {
		t.Errorf("Expected pong in output, got: %s", got)
	}
}

func TestStdioTransport_Initialize(t *testing.T) {
	json := `{"jsonrpc":"2.0","method":"initialize","id":42,"params":{"protocolVersion":"2025-03-26","clientInfo":{"name":"TestClient","version":"1.0.0"},"capabilities":{}}}`
	content := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(json), json)

	in := bytes.NewBufferString(content)
	out := &bytes.Buffer{}

	transport := &StdioTransport{
		reader: bufio.NewReader(in),
		writer: bufio.NewWriter(out),
		cache:  sharedCache,
	}

	go func() {
		_ = transport.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	got := out.String()
	if !strings.Contains(got, `"jsonrpc":"2.0"`) || !strings.Contains(got, `"id":42`) || !strings.Contains(got, `"protocolVersion"`) {
		t.Errorf("Expected valid initialize response, got: %s", got)
	}
}
