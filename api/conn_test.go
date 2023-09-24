package api

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	// this test uses a fake HTTP implementation
	http := &fakeHttp{}
	// create an API connection
	logger := NewLogger(io.Discard, false)
	conn := NewConn(logger, http, nil)
	// GET ping
	{
		http.responses = append(http.responses, dataAndErr{})
		err := conn.GetPing()
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET ping", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// GET watchdog/00000001
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{data: []byte("{\"id\": \"00000001\"}")})
		data, err := conn.GetWatchdog("00000001")
		assertNil(t, err)
		assertEq(t, "{\"id\": \"00000001\"}", string(data))
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET watchdog/00000001", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// POST machine/00000001/sample
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{})
		tstamp := time.Date(2023, 9, 1, 10, 0, 0, 0, time.Local)
		err := conn.PostMachineSample("00000001", tstamp.UnixMilli(), 12, 13, 14, 1, 0)
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST machine/00000001/sample tstamp=1693555200000&cpu=12&mem=13&disk=14", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// POST metric/00000001/inc (with 3 trials)
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("error 1")})
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("error 2")})
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("error 3")})
		err := conn.PostMetricInc("00000001", 42, 3, 10*time.Second)
		assertEq(t, "error 3", err.Error())
		assertEq(t, 3, len(http.requests))
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[0])
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[1])
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[2])
		assertEq(t, 0, len(http.responses))
	}
}

// fake http

type fakeHttp struct {
	requests  []string
	responses []dataAndErr
}

var _ Http = (*fakeHttp)(nil)

func (f *fakeHttp) Send(method, path string, data []byte) ([]byte, error) {
	req := fmt.Sprintf("%s %s", method, path)
	if len(data) > 0 {
		req += fmt.Sprintf(" %s", string(data))
	}
	f.requests = append(f.requests, req)
	if len(f.responses) == 0 {
		return nil, fmt.Errorf("empty response for %s %s", method, path)
	}
	tmp := f.responses[0]
	f.responses = f.responses[1:]
	return tmp.data, tmp.err
}

type dataAndErr struct {
	data []byte
	err  error
}
