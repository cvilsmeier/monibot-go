package monibot

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// A Sender sends HTTP requests and receives HTTP responses.
// It is used by Api to send HTTP requests.
type Sender interface {

	// Send sends a HTTP request.
	// It returns the raw response data or an error.
	Send(method, path string, body []byte) ([]byte, error)
}

// SenderOptions hold custom options for a Sender.
type SenderOptions struct {

	// The URL to send data to. Default is "https://monibot.io".
	MonibotUrl string

	// The UserAgent. Default is "monibot/v0.0.0" (whatever the current version is).
	UserAgent string

	// The Logger for verbose debug logging. Default logs nothing.
	Logger Logger
}

type httpSender struct {
	logger    Logger
	apiUrl    string
	userAgent string
	apiKey    string
}

var _ Sender = (*httpSender)(nil)

// NewSender creates a Sender that sends data to https://monibot.io.
func NewSender(apiKey string) Sender {
	return NewSenderWithOptions(apiKey, SenderOptions{})
}

// NewSenderWithOptions creates a new Sender with custom options.
func NewSenderWithOptions(apiKey string, options SenderOptions) Sender {
	if options.MonibotUrl == "" {
		options.MonibotUrl = "https://monibot.io"
	}
	if options.UserAgent == "" {
		options.UserAgent = "monibot/" + Version
	}
	if options.Logger == nil {
		options.Logger = NewDiscardLogger()
	}
	return &httpSender{options.Logger, options.MonibotUrl + "/api/", options.UserAgent, apiKey}
}

// Send sends a HTTP request.
// It returns the raw response data or an error.
func (s *httpSender) Send(method, path string, body []byte) ([]byte, error) {
	urlpath := s.apiUrl + path
	s.logger.Debug("%s %s", method, urlpath)
	if len(body) > 0 {
		s.logger.Debug("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, urlpath, bodyReader)
	if err != nil {
		s.logger.Debug("cannot create request: %s", err)
		return nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("X-Monibot-Version", Version)
	req.Header.Set("X-Monibot-Trial", "1") // TODO weak-code this
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
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
