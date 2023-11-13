package sending

import (
	"context"
	"fmt"
	"testing"
)

// fakeTransport is a Transport for unit tests
type fakeTransport struct {
	calls     []string
	responses []fakeResponse
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

type fakeResponse struct {
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
