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
	is := assert.New(t)
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
	sender := NewTransport(logger, server.URL, "api-key-123", "0.2.3")
	// send ok
	status, data, err := sender.Send(context.Background(), "GET", "/ok", nil)
	is.Nil(err)
	is.Eq(200, status)
	is.Eq("ok", string(data))
	// send 500
	status, data, err = sender.Send(context.Background(), "POST", "/500", nil)
	is.Nil(err)
	is.Eq(500, status)
	is.Eq("", string(data))
}

type fakeSenderLogger struct{}

func (f *fakeSenderLogger) Debug(format string, args ...any) {}
