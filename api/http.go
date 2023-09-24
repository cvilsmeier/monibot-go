package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Http provides a HTTP send function.
type Http interface {

	// Send sends a HTTP request.
	// It returns the raw response data and/or an error.
	Send(method, path string, body []byte) ([]byte, error)
}

type httpImpl struct {
	logger    Logger
	apiUrl    string
	userAgent string
	apiKey    string
}

var _ Http = (*httpImpl)(nil)

// NewHttp creates a new Http implementation.
func NewHttp(logger Logger, monibotUrl, userAgent, apiKey string) Http {
	apiUrl := monibotUrl + "/api/"
	return &httpImpl{logger, apiUrl, userAgent, apiKey}
}

// Send sends a HTTP request. It returns the response data and/or an error.
func (h *httpImpl) Send(method, path string, body []byte) ([]byte, error) {
	urlpath := h.apiUrl + path
	h.logger.Debugf("%s %s", method, urlpath)
	if len(body) > 0 {
		h.logger.Debugf("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, urlpath, bodyReader)
	if err != nil {
		h.logger.Debugf("cannot create request: %s", err)
		return nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", h.userAgent)
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.logger.Debugf("%s %s: %s", req.Method, urlpath, err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response data: %w", err)
	}
	h.logger.Debugf("%d %s", resp.StatusCode, string(data))
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return nil, fmt.Errorf("response status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}
