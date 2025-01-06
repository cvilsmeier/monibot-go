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
	transport.responses = []fakeTransportResponse{
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
	transport.responses = []fakeTransportResponse{
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
	transport.responses = []fakeTransportResponse{
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
	transport.responses = []fakeTransportResponse{
		{401, []byte("401 - Unauthorized (invalid apiKey)"), nil},
	}
	_, err = sender.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(1, len(transport.calls))
	ass.Eq("GET /ping", transport.calls[0])
	ass.Eq("status 401: 401 - Unauthorized (invalid apiKey)", err.Error())
	transport.calls = nil
	// must not retry if 404 (not found) but give error
	transport.responses = []fakeTransportResponse{
		{404, nil, nil},
	}
	_, err = sender.Send(context.Background(), "GET", "/wrongUrl", nil)
	ass.Eq(1, len(transport.calls))
	ass.Eq("GET /wrongUrl", transport.calls[0])
	ass.Eq("status 404", err.Error())
	transport.calls = nil
}

// fakeTransport is a Transport for unit tests
type fakeTransport struct {
	calls     []string
	responses []fakeTransportResponse
}

func (f *fakeTransport) Send(ctx context.Context, method, path string, body []byte) (int, []byte, error) {
	call := fmt.Sprintf("%s %s", method, path)
	if len(body) > 0 {
		call += fmt.Sprintf(" %s", string(body))
	}
	f.calls = append(f.calls, call)
	if len(f.responses) == 0 {
		return 0, nil, fmt.Errorf("fakeSender is out of responses for request %s %s", method, path)
	}
	re := f.responses[0]
	f.responses = f.responses[1:]
	return re.status, re.data, re.err
}

type fakeTransportResponse struct {
	status int
	data   []byte
	err    error
}

// fakeLogger is a Logger for unit tests
type fakeLogger struct {
	t       testing.TB
	enabled bool
}

func (f *fakeLogger) Debug(format string, args ...any) {
	if f.enabled {
		f.t.Logf(format, args...)
	}
}
