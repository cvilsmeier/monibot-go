package monibot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	// this test uses a fake HTTP sender
	http := &fakeSender{}
	// create Api
	logger := NewLogger(nil)
	sleep := func(time.Duration) {}
	api := NewApi(logger, http, sleep)
	// GET ping
	{
		http.responses = append(http.responses, dataAndErr{})
		err := api.GetPing()
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "GET ping", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// GET watchdog/00000001
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{data: []byte("{\"id\": \"00000001\"}")})
		data, err := api.GetWatchdog("00000001")
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
		err := api.PostMachineSample("00000001", tstamp.UnixMilli(), 12, 13, 14)
		assertNil(t, err)
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST machine/00000001/sample tstamp=1693555200000&cpu=12&mem=13&disk=14", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
	// POST metric/00000001/inc
	{
		http.requests = nil
		http.responses = append(http.responses, dataAndErr{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricInc("00000001", 42)
		assertEq(t, "connect timeout", err.Error())
		assertEq(t, 1, len(http.requests))
		assertEq(t, "POST metric/00000001/inc value=42", http.requests[0])
		assertEq(t, 0, len(http.responses))
	}
}

func TestDemoForReadme(t *testing.T) {
	// parse api_test.go
	data, err := os.ReadFile("api_test.go")
	assertNil(t, err)
	want := string(data)
	_, want, found := strings.Cut(want, "// "+"@test-start")
	assertTrue(t, found)
	want, _, found = strings.Cut(want, "// "+"@test-end")
	assertTrue(t, found)
	want = strings.ReplaceAll(want, "\t", "")
	want = strings.TrimSpace(want)
	// parse README.md
	data, err = os.ReadFile("README.md")
	assertNil(t, err)
	have := string(data)
	_, have, found = strings.Cut(have, "import "+"\"github.com/cvilsmeier/monibot-go\"")
	assertTrue(t, found)
	have, _, found = strings.Cut(have, "`"+"`"+"`")
	assertTrue(t, found)
	have = strings.ReplaceAll(have, "\t", "")
	have = strings.ReplaceAll(have, "    ", "")
	have = strings.TrimSpace(have)
	if want != have {
		t.Logf("want %q", want)
		t.Logf("have %q", have)
	}
	assertEq(t, want, have)
}

// This code is only here to be copied into README.md - do not execute.
func DemoForReadme() {
	// ensure it's never executed
	if 2+2 > 1 {
		panic("do not execute")
	}
	// import "github.com/cvilsmeier/monibot-go"
	// @test-start
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
	// @test-end
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
