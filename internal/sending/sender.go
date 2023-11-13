package sending

import (
	"context"
	"fmt"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/logging"
)

// TimeAfterFunc is the function type of time.After.
type TimeAfterFunc func(time.Duration) <-chan time.Time

type senderTransport interface {
	Send(ctx context.Context, method, path string, body []byte) (int, []byte, error)
}

type Sender struct {
	transport senderTransport
	logger    logging.Logger
	trials    int
	delay     time.Duration
	timeAfter TimeAfterFunc
}

func NewSender(transport senderTransport, logger logging.Logger, trials int, delay time.Duration, timeAfter TimeAfterFunc) *Sender {
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
			if err == nil && !done {
				err = fmt.Errorf("status %d", status)
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
		// technical error
		// -> retry
		return false
	}
	if status == 200 {
		// success
		// -> done
		return true
	}
	if status == 429 {
		// rate limit
		// -> retry
		return false
	}
	if 400 <= status && status <= 499 {
		// error in request data
		// -> done, because next trial will probably bring the same result
		return true
	}
	// other status code (5xx) (server maintenance, nginx bad gateway, ...)
	// -> retry
	return false
}
