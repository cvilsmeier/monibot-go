package monibot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// A Sender sends HTTP requests and receives HTTP responses.
// It is used by Api to send HTTP requests.
type Sender interface {

	// Send sends a HTTP request.
	// It returns the raw response data or an error.
	Send(ctx context.Context, method, path string, body []byte) ([]byte, error)
}

// httpSender is a Sender that uses HTTP for sending API requests.
type httpSender struct {
	logger    Logger
	apiUrl    string
	userAgent string
	apiKey    string
}

var _ Sender = (*httpSender)(nil)

// NewSender creates a new Senderthat sends api requests to https://monibot.io.
func NewSender(apiKey string) Sender {
	return NewSenderWithOptions(nil, "", "", apiKey)
}

// NewSenderWithOptions creates a new Sender.
// If logger is nil, it logs nothing.
// If monibotUrl is empty, it sends api requests to https://monibot.io.
// If userAgent is empty, it uses "monibot/v0.0.0" (whatever the current version is).
func NewSenderWithOptions(logger Logger, monibotUrl, userAgent, apiKey string) Sender {
	if logger == nil {
		logger = NewDiscardLogger()
	}
	if monibotUrl == "" {
		monibotUrl = "https://monibot.io"
	}
	if userAgent == "" {
		userAgent = "monibot/" + Version
	}
	return &httpSender{logger, monibotUrl + "/api/", userAgent, apiKey}
}

// Send sends a HTTP request.
// It returns the raw response data or an error.
func (s *httpSender) Send(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	urlpath := s.apiUrl + path
	s.logger.Debug("%s %s", method, urlpath)
	if len(body) > 0 {
		s.logger.Debug("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, method, urlpath, bodyReader)
	if err != nil {
		s.logger.Debug("cannot create request: %s", err)
		return nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("X-Monibot-Client", "monibot-go")
	req.Header.Set("X-Monibot-Version", Version)
	req.Header.Set("X-Monibot-TstampMillis", strconv.FormatInt(time.Now().UnixMilli(), 10))
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Debug("%s %s: %s", req.Method, urlpath, err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response data: %w", err)
	}
	if len(data) > 256 {
		s.logger.Debug("%d (%d bytes) %s", resp.StatusCode, len(data), string(data)[:256]+"...")
	} else {
		s.logger.Debug("%d (%d bytes) %s", resp.StatusCode, len(data), string(data))
	}
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		text := string(data)
		if text == "" {
			return nil, fmt.Errorf("response status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status %d: %s", resp.StatusCode, text)
	}
	return data, nil
}
