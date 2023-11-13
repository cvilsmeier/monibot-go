package sending

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type Transport struct {
	logger  debugLogger
	version string
	apiUrl  string
	apiKey  string
}

func NewTransport(logger debugLogger, version, monibotUrl, apiKey string) *Transport {
	return &Transport{logger, version, monibotUrl + "/api/", apiKey}
}

func (s *Transport) Send(ctx context.Context, method, path string, body []byte) (int, []byte, error) {
	urlpath := s.apiUrl + path
	s.logger.Debug("%s %s", method, urlpath)
	if len(body) > 0 {
		s.logger.Debug("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, method, urlpath, bodyReader)
	if err != nil {
		s.logger.Debug("cannot create request: %s", err)
		return 0, nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", "monibot/"+s.version)
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Debug("%s %s: %s", req.Method, urlpath, err)
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot read response data: %w", err)
	}
	if len(data) > 256 {
		s.logger.Debug("%d (%d bytes) %s", resp.StatusCode, len(data), string(data)[:256]+"...")
	} else {
		s.logger.Debug("%d (%d bytes) %s", resp.StatusCode, len(data), string(data))
	}
	return resp.StatusCode, data, nil
}
