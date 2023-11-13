package sending

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestSender(t *testing.T) {
	ass := assert.New(t)
	// setup fake api http server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	mux.HandleFunc("/api/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	// init
	logger := &fakeSenderLogger{}
	sender := NewTransport(logger, "v1.2.3", server.URL, "api-key-123")
	// send ok
	status, data, err := sender.Send(context.Background(), "GET", "/ok", nil)
	ass.Nil(err)
	ass.Eq(200, status)
	ass.Eq("ok", string(data))
	// send 500
	status, data, err = sender.Send(context.Background(), "POST", "/500", nil)
	ass.Nil(err)
	ass.Eq(500, status)
	ass.Eq("", string(data))
}

type fakeSenderLogger struct{}

func (f *fakeSenderLogger) Debug(format string, args ...any) {}
