package sending

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// logger prints debug messages
type debugLogger interface {
	Debug(format string, args ...any)
}

// TimeAfterFunc is the function type of time.After.
type TimeAfterFunc func(time.Duration) <-chan time.Time

type senderTransport interface {
	Send(ctx context.Context, method, path string, body []byte) (int, []byte, error)
}

type Sender struct {
	transport senderTransport
	logger    debugLogger
	trials    int
	delay     time.Duration
	timeAfter TimeAfterFunc
}

func NewSender(transport senderTransport, logger debugLogger, trials int, delay time.Duration, timeAfter TimeAfterFunc) *Sender {
	if trials < 1 {
		trials = 1
	}
	if delay < 0 {
		delay = 0
	}
	return &Sender{transport, logger, trials, delay, timeAfter}
}

func (s *Sender) Send(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var trial int
	for {
		trial++
		s.logger.Debug("trial #%d/%d for %s %s", trial, s.trials, method, path)
		status, data, err := s.transport.Send(ctx, method, path, body)
		done := isDone(status, err)
		if done || trial >= s.trials {
			if err == nil && status != 200 {
				msg := fmt.Sprintf("status %d", status)
				if len(data) > 0 {
					msg += ": " + string(data)
				}
				err = errors.New(msg)
			}
			return data, err
		}
		select {
		case <-s.timeAfter(s.delay):
			// retry now
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled")
		}
	}
}

func isDone(status int, err error) bool {
	if err != nil {
		// newtwork error, e.g. connect failed
		// -> not done, retry
		return false
	}
	if status == 200 {
		// success
		// -> done
		return true
	}
	if status == 429 {
		// rate limit
		// -> not done, retry
		return false
	}
	if 400 <= status && status <= 499 {
		// not found, wrong apiKey, bad request, etc.
		// -> done, because next trial will probably bring the same result
		return true
	}
	// other status code (5xx) (server maintenance, bad gateway, internal server error...)
	// -> not done, retry
	return false
}
