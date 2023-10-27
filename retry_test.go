package monibot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestRetrySender(t *testing.T) {
	ass := assert.New(t)
	sender := &fakeSender{}
	timeChan := make(chan time.Time)
	logger := &fakeLogger{t, false}
	retry := NewRetrySenderWithOptions(sender, RetrySenderOptions{
		Logger: logger,
		TimeAfter: func(d time.Duration) <-chan time.Time {
			return timeChan
		},
		Trials: 3,
		Delay:  2 * time.Second,
	})
	// ping
	sender.responses = []fakeResponse{
		{nil, fmt.Errorf("connection refused")},
		{nil, fmt.Errorf("connection refused")},
		{[]byte("ok"), nil},
	}
	go func() {
		timeChan <- time.Now()
		timeChan <- time.Now()
	}()
	data, err := retry.Send(context.Background(), "GET", "/ping", nil)
	ass.Eq(3, len(sender.requests))
	ass.Eq("GET /ping", sender.requests[0])
	ass.Eq("GET /ping", sender.requests[1])
	ass.Eq("GET /ping", sender.requests[2])
	ass.Nil(err)
	ass.Eq("ok", string(data))
}
