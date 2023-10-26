package monibot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestSend(t *testing.T) {
	ass := assert.New(t)
	// setup fake api http server
	var pingOk atomic.Bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		if !pingOk.Load() {
			w.WriteHeader(500)
			return
		}
		fmt.Fprintf(w, "ok")
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	// init sender
	sender := NewSenderWithOptions("123456", SenderOptions{MonibotUrl: server.URL})
	// send ping - good
	pingOk.Store(true)
	data, err := sender.Send("GET", "/ping", nil)
	ass.Nil(err)
	ass.Eq("ok", string(data))
	// send ping - error
	pingOk.Store(false)
	data, err = sender.Send("GET", "/ping", nil)
	ass.Eq("response status 500", err.Error())
	ass.Eq("", string(data))
}
