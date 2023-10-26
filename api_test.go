package monibot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestApi(t *testing.T) {
	ass := assert.New(t)
	// this test uses a fake HTTP sender
	sender := &fakeSender{}
	// create Api
	api := NewApiWithSender(sender)
	// GET ping
	{
		sender.responses = append(sender.responses, dataAndErr{})
		err := api.GetPing()
		ass.Nil(err)
		ass.Eq(1, len(sender.requests))
		ass.Eq("GET ping", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// GET watchdogs
	{
		sender.responses = append(sender.responses, dataAndErr{data: []byte("[{\"id\": \"00000001\"}]")})
		data, err := api.GetWatchdogs()
		ass.Nil(err)
		ass.Eq("[{\"id\": \"00000001\"}]", string(data))
		ass.Eq(1, len(sender.requests))
		ass.Eq("GET watchdogs", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// GET watchdog/00000001
	{
		sender.responses = append(sender.responses, dataAndErr{data: []byte("{\"id\": \"00000001\"}")})
		data, err := api.GetWatchdog("00000001")
		ass.Nil(err)
		ass.Eq("{\"id\": \"00000001\"}", string(data))
		ass.Eq(1, len(sender.requests))
		ass.Eq("GET watchdog/00000001", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// POST watchdog/00000001/reset
	{
		sender.responses = append(sender.responses, dataAndErr{})
		err := api.PostWatchdogReset("00000001")
		ass.Nil(err)
		ass.Eq(1, len(sender.requests))
		ass.Eq("POST watchdog/00000001/reset", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// POST machine/00000001/sample
	{
		sender.responses = append(sender.responses, dataAndErr{})
		tstamp := time.Date(2023, 9, 1, 10, 0, 0, 0, time.Local)
		err := api.PostMachineSample("00000001", tstamp.UnixMilli(), 12, 13, 14)
		ass.Nil(err)
		ass.Eq(1, len(sender.requests))
		ass.Eq("POST machine/00000001/sample tstamp=1693555200000&cpu=12&mem=13&disk=14", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// POST metric/00000001/inc
	{
		sender.responses = append(sender.responses, dataAndErr{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricInc("00000001", 42)
		ass.Eq("connect timeout", err.Error())
		ass.Eq(1, len(sender.requests))
		ass.Eq("POST metric/00000001/inc value=42", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
}

// fake http

type fakeSender struct {
	requests  []string
	responses []dataAndErr
}

var _ Sender = (*fakeSender)(nil)

func (f *fakeSender) Send(_ context.Context, method, path string, data []byte) ([]byte, error) {
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
