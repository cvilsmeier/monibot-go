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

// RetrySenderOptions hold RetrySender opptions.
type RetrySenderOptions struct {

	// Default logs nothing.
	Logger Logger

	// Default is time.After.
	TimeAfter func(time.Duration) <-chan time.Time

	// Default is 12.
	Trials int

	// Default is 5s.
	Delay time.Duration
}

// NewRetrySender creates a new RetrySender that does max. 12 trials with a delay of 5 seconds in between.
func NewRetrySender(sender Sender) *RetrySender {
	return NewRetrySenderWithOptions(sender, RetrySenderOptions{})
}

// NewRetrySenderWithOptions creates a new RetrySender with custom options.
func NewRetrySenderWithOptions(sender Sender, opt RetrySenderOptions) *RetrySender {
	if sender == nil {
		panic("sender == nil")
	}
	if opt.Logger == nil {
		opt.Logger = NewDiscardLogger()
	}
	if opt.TimeAfter == nil {
		opt.TimeAfter = time.After
	}
	if opt.Trials < 1 {
		opt.Trials = 12
	}
	if opt.Delay <= 0 {
		opt.Delay = 5 * time.Second
	}
	return &RetrySender{opt.Logger, sender, opt.TimeAfter, opt.Trials, opt.Delay}
}

func (s *RetrySender) Send(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	// first trial, for sure
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
