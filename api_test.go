package monibot

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	// this test uses a fake HTTP sender
	http := &fakeSender{}
	// create an API connection
	logger := NewLogger(nil)
	sleep := func(time.Duration) {}
	conn := NewApi(logger, http, sleep)
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
		err := conn.PostMachineSample("00000001", tstamp.UnixMilli(), 12, 13, 14)
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST machine/00000001/sample tstamp=1693555200000&cpu=12&mem=13&disk=14", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// POST metric/00000001/inc
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("connect timeout")})
		err := conn.PostMetricInc("00000001", 42)
		assertEq(t, "connect timeout", err.Error())
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
}

// This code is only here to be copied into README.md
func DemoForReadme() {
	// import "github.com/cvilsmeier/monibot-go"
	// init api
	userAgent := "my-app/v1.0.0"
	apiKey := os.Getenv("MONIBOT_API_KEY")
	api := NewDefaultApi(userAgent, apiKey)
	// ping the api
	err := api.GetPing()
	if err != nil {
		log.Fatal(err)
	}
	// reset a watchdog
	err = api.PostWatchdogReset("000000000000001")
	if err != nil {
		log.Fatal(err)
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
