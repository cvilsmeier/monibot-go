package monibot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestApi(t *testing.T) {
	str := func(a any) string {
		switch x := a.(type) {
		case Watchdog:
			return fmt.Sprintf("Id=%s, Name=%s, IntervalMillis=%d", x.Id, x.Name, x.IntervalMillis)
		}
		return "?"
	}
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
		resp := `[
			{"id":"0001", "name":"Cronjob 1", "intervalMillis": 72000000},
			{"id":"0002", "name":"Cronjob 2", "intervalMillis": 36000000}
		]`
		sender.responses = append(sender.responses, dataAndErr{data: []byte(resp)})
		watchdogs, err := api.GetWatchdogs()
		ass.Nil(err)
		ass.Eq(2, len(watchdogs))
		ass.Eq("Id=0001, Name=Cronjob 1, IntervalMillis=72000000", str(watchdogs[0]))
		ass.Eq("Id=0002, Name=Cronjob 2, IntervalMillis=36000000", str(watchdogs[1]))
		ass.Eq(1, len(sender.requests))
		ass.Eq("GET watchdogs", sender.requests[0])
		ass.Eq(0, len(sender.responses))
		sender.requests, sender.responses = nil, nil
	}
	// GET watchdog/00000001
	{
		resp := `{"id":"0001", "name":"Cronjob 1", "intervalMillis": 72000000}`
		sender.responses = append(sender.responses, dataAndErr{data: []byte(resp)})
		watchdog, err := api.GetWatchdog("00000001")
		ass.Nil(err)
		ass.Eq("Id=0001, Name=Cronjob 1, IntervalMillis=72000000", str(watchdog))
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

func (f *fakeSender) Send(ctx context.Context, method, path string, data []byte) ([]byte, error) {
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
