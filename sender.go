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

type senderImpl struct {
	logger    Logger
	apiUrl    string
	userAgent string
	apiKey    string
}

var _ Sender = (*senderImpl)(nil)

// NewSender creates a new HTTP Sender.
func NewSender(logger Logger, monibotUrl, userAgent, apiKey string) Sender {
	if logger == nil {
		panic("logger is nil")
	}
	apiUrl := monibotUrl + "/api/"
	return &senderImpl{logger, apiUrl, userAgent, apiKey}
}

// Send sends a HTTP request.
// It returns the raw response data or an error.
func (x *senderImpl) Send(method, path string, body []byte) ([]byte, error) {
	urlpath := x.apiUrl + path
	x.logger.Debug("%s %s", method, urlpath)
	if len(body) > 0 {
		x.logger.Debug("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, urlpath, bodyReader)
	if err != nil {
		x.logger.Debug("cannot create request: %s", err)
		return nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", x.userAgent)
	req.Header.Set("Authorization", "Bearer "+x.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		x.logger.Debug("%s %s: %s", req.Method, urlpath, err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response data: %w", err)
	}
	x.logger.Debug("%d %s", resp.StatusCode, string(data))
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return nil, fmt.Errorf("response status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}
