package sending

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestRetrySender(t *testing.T) {
	ass := assert.New(t)
	transport := &fakeTransport{}
	timeChan := make(chan time.Time)
	logger := &fakeLogger{t, false}
	timeAfter := func(d time.Duration) <-chan time.Time {
		return timeChan
	}
	trials := 3
	delay := 2 * time.Second
	sender := NewSender(transport, logger, trials, delay, timeAfter)
	// must retry if network error
	transport.responses = []fakeResponse{
		{0, nil, fmt.Errorf("connection refused")},
		{0, nil, fmt.Errorf("connection refused")},
		{200, []byte("{\"ok\":true}"), nil},
	}
	go func() {
		timeChan <- time.Now()
		timeChan <- time.Now()
	}()
	data, err := sender.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(3, len(transport.calls))
	ass.Eq("GET /ping", transport.calls[0])
	ass.Eq("GET /ping", transport.calls[1])
	ass.Eq("GET /ping", transport.calls[2])
	ass.Nil(err)
	ass.Eq("{\"ok\":true}", string(data))
	transport.calls = nil
	// must retry max trials
	transport.responses = []fakeResponse{
		{0, nil, fmt.Errorf("connection error 1")},
		{0, nil, fmt.Errorf("connection error 2")},
		{0, nil, fmt.Errorf("connection error 3")},
	}
	go func() {
		timeChan <- time.Now()
		timeChan <- time.Now()
	}()
	_, err = sender.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(3, len(transport.calls))
	ass.Eq("GET /ping", transport.calls[0])
	ass.Eq("GET /ping", transport.calls[1])
	ass.Eq("GET /ping", transport.calls[2])
	ass.Eq("connection error 3", err.Error())
	transport.calls = nil
	// must retry if status 502 (bad gateway)
	transport.responses = []fakeResponse{
		{502, nil, nil},
		{502, nil, nil},
		{200, []byte("{\"ok\":true}"), nil},
	}
	go func() {
		timeChan <- time.Now()
		timeChan <- time.Now()
	}()
	data, err = sender.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(3, len(transport.calls))
	ass.Eq("GET /ping", transport.calls[0])
	ass.Eq("GET /ping", transport.calls[1])
	ass.Eq("GET /ping", transport.calls[2])
	ass.Nil(err)
	ass.Eq("{\"ok\":true}", string(data))
	transport.calls = nil
	// must not retry if authorization error
	transport.responses = []fakeResponse{
		{401, []byte("401 - Unauthorized (invalid apiKey)"), nil},
	}
	_, err = sender.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(1, len(transport.calls))
	ass.Eq("GET /ping", transport.calls[0])
	ass.Eq("status 401: 401 - Unauthorized (invalid apiKey)", err.Error())
	transport.calls = nil
	// must not retry if 404 (not found) but give error
	transport.responses = []fakeResponse{
		{404, nil, nil},
	}
	_, err = sender.Send(context.Background(), "GET", "/wrongUrl", nil)
	ass.Eq(1, len(transport.calls))
	ass.Eq("GET /wrongUrl", transport.calls[0])
	ass.Eq("status 404", err.Error())
	transport.calls = nil
}
