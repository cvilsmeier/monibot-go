package monibot

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	// this test uses a fake HTTP sender
	http := &fakeSender{}
	// create Api
	logger := NewLogger(io.Discard)
	api := NewApi(logger, http)
	// GET ping
	{
		http.responses = append(http.responses, dataAndErr{})
		err := api.GetPing()
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET ping", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
	// GET watchdogs
	{
		http.responses = append(http.responses, dataAndErr{data: []byte("[{\"id\": \"00000001\"}]")})
		data, err := api.GetWatchdogs()
		assertNil(t, err)
		assertEq(t, "[{\"id\": \"00000001\"}]", string(data))
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET watchdogs", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
	// GET watchdog/00000001
	{
		http.responses = append(http.responses, dataAndErr{data: []byte("{\"id\": \"00000001\"}")})
		data, err := api.GetWatchdog("00000001")
		assertNil(t, err)
		assertEq(t, "{\"id\": \"00000001\"}", string(data))
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET watchdog/00000001", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
	// POST watchdog/00000001/reset
	{
		http.responses = append(http.responses, dataAndErr{})
		err := api.PostWatchdogReset("00000001")
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST watchdog/00000001/reset", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
	// POST machine/00000001/sample
	{
		http.responses = append(http.responses, dataAndErr{})
		tstamp := time.Date(2023, 9, 1, 10, 0, 0, 0, time.Local)
		err := api.PostMachineSample("00000001", tstamp.UnixMilli(), 12, 13, 14)
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST machine/00000001/sample tstamp=1693555200000&cpu=12&mem=13&disk=14", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
	// POST metric/00000001/inc
	{
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricInc("00000001", 42)
		assertEq(t, "connect timeout", err.Error())
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[0])
		assertEq(t, 0, len(http.responses))
		http.requests, http.responses = nil, nil
	}
}

// fake http

type fakeSender struct {
	requests  []string
	responses []dataAndErr
}

var _ Sender = (*fakeSender)(nil)

func (f *fakeSender) Send(method, path string, data []byte) ([]byte, error) {
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
