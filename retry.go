package monibot

import (
	"context"
	"fmt"
	"time"
)

// RetrySender wraps a Sender and re-sends API requests in case of error.
type RetrySender struct {
	logger    Logger
	sender    Sender
	timeAfter func(time.Duration) <-chan time.Time
	trials    int
	delay     time.Duration
}

var _ Sender = (*RetrySender)(nil)

// TODO
func NewRetrySender(sender Sender) *RetrySender {
	return NewRetrySenderWithOptions(nil, sender, nil, 0, 0)
}

// TODO
func NewRetrySenderWithOptions(logger Logger, sender Sender, timeAfter func(time.Duration) <-chan time.Time, trials int, delay time.Duration) *RetrySender {
	if logger == nil {
		logger = NewDiscardLogger()
	}
	if sender == nil {
		panic("sender == nil")
	}
	if timeAfter == nil {
		timeAfter = time.After
	}
	if trials < 1 {
		trials = 12
	}
	if delay <= 0 {
		delay = 5 * time.Second
	}
	return &RetrySender{logger, sender, timeAfter, trials, delay}
}

func (s *RetrySender) Send(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	// first trial, always
	s.logger.Debug("trial #1 for %s %s", method, path)
	data, err := s.sender.Send(ctx, method, path, body)
	if err == nil {
		return data, err
	}
	// we have to retry again later
	for i := 1; i < s.trials; i++ {
		select {
		case <-s.timeAfter(s.delay):
			// retry
			s.logger.Debug("trial #%d for %s %s", i+1, method, path)
			data, err = s.sender.Send(ctx, method, path, body)
			if err == nil {
				return data, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled")
		}
	}
	return data, err
}
