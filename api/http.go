package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Http interface {
	Get(path string) ([]byte, error)
	Post(path string, data []byte) ([]byte, error)
}

type httpImpl struct {
	logger    Logger
	apiUrl    string
	userAgent string
	apiKey    string
}

func NewHttp(logger Logger, monibotUrl, userAgent, apiKey string) *httpImpl {
	apiUrl := monibotUrl + "/api/"
	return &httpImpl{logger, apiUrl, userAgent, apiKey}
}

func (h *httpImpl) Get(path string) ([]byte, error) {
	urlpath := h.apiUrl + path
	h.logger.Debugf("GET %s", urlpath)
	req, err := http.NewRequest("GET", urlpath, nil)
	if err != nil {
		h.logger.Debugf("GET %s: %s", urlpath, err)
		return nil, err
	}
	req.Header.Set("User-Agent", h.userAgent)
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.logger.Debugf("GET %s: %s", urlpath, err)
		return nil, err
	}
	return h.readResponse(resp)
}

func (h *httpImpl) Post(path string, data []byte) ([]byte, error) {
	urlpath := h.apiUrl + path
	h.logger.Debugf("POST %s", urlpath)
	if len(data) > 0 {
		h.logger.Debugf("data=%s", string(data))
	}
	dataReader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", urlpath, dataReader)
	if err != nil {
		h.logger.Debugf("POST %s: %s", urlpath, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", h.userAgent)
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.logger.Debugf("POST %s: %s", urlpath, err)
		return nil, err
	}
	return h.readResponse(resp)
}

func (h *httpImpl) readResponse(resp *http.Response) ([]byte, error) {
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
