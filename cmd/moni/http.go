package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Http struct {
	monibotUrl string
	apiKey     string
	version    string
	verbose    bool
}

func NewHttp(url, apiKey, version string, verbose bool) *Http {
	return &Http{url, apiKey, version, verbose}
}

func (h *Http) Get(path string) ([]byte, error) {
	urlpath := h.monibotUrl + "/api/" + path
	h.debug("GET %s", urlpath)
	req, err := http.NewRequest("GET", urlpath, nil)
	if err != nil {
		h.debug("POST %s: %s", urlpath, err)
		return nil, err
	}
	req.Header.Set("User-Agent", "moni/"+h.version)
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.debug("DEBUG POST %s: %s", urlpath, err)
		return nil, err
	}
	return h.readResponse(resp)
}

func (h *Http) Post(path string, body []byte) ([]byte, error) {
	urlpath := h.monibotUrl + "/api/" + path
	h.debug("POST %s", urlpath)
	if len(body) > 0 {
		h.debug("body=%s", string(body))
	}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest("POST", urlpath, bodyReader)
	if err != nil {
		h.debug("POST %s: %s", urlpath, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "moni/"+h.version)
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.debug("POST %s: %s", urlpath, err)
		return nil, err
	}
	return h.readResponse(resp)
}

func (h *Http) readResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response data: %w", err)
	}
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return nil, fmt.Errorf("response status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (h *Http) debug(f string, a ...any) {
	if h.verbose {
		prt("DEBUG: "+f, a...)
	}
}
