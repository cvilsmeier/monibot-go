package monibot

import (
	"context"
	"fmt"
	"testing"
)

// fakeSender is a Sender for unit tests
type fakeSender struct {
	requests  []string
	responses []fakeResponse
}

var _ Sender = (*fakeSender)(nil)

func (f *fakeSender) Send(ctx context.Context, method, path string, data []byte) ([]byte, error) {
	req := fmt.Sprintf("%s %s", method, path)
	if len(data) > 0 {
		req += fmt.Sprintf(" %s", string(data))
	}
	f.requests = append(f.requests, req)
	if len(f.responses) == 0 {
		return nil, fmt.Errorf("fakeSender is out of responses, request was %s %s", method, path)
	}
	tmp := f.responses[0]
	f.responses = f.responses[1:]
	return tmp.data, tmp.err
}

type fakeResponse struct {
	data []byte
	err  error
}

// fakeLogger is a Logger for unit tests
type fakeLogger struct {
	t       testing.TB
	enabled bool
}

var _ Logger = (*fakeLogger)(nil)

func (f *fakeLogger) Debug(format string, args ...any) {
	if f.enabled {
		f.t.Logf(format, args...)
	}
}
