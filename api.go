package monibot

import (
	"fmt"
	"io"
	"net/http"
)

// Api provides access to the Monibot REST API.
type Api struct {
	logger Logger
	sender Sender
}

// NewDefaultApi creates an Api with default implementations.
// This should be suited for most use cases.
func NewDefaultApi(userAgent, apiKey string) *Api {
	logger := NewLogger(io.Discard)
	sender := NewSender(logger, "https://monibot.io", userAgent, apiKey)
	return NewApi(logger, sender)
}

// NewApi creates an Api with a custom logger and sender.
func NewApi(logger Logger, sender Sender) *Api {
	if logger == nil {
		panic("no logger")
	}
	if sender == nil {
		panic("no sender")
	}
	return &Api{logger, sender}
}

// GetPing calls the /ping endpoint.
// It is used to ensure everything is set up correctly and the API is reachable.
func (x *Api) GetPing() error {
	_, err := x.sender.Send(http.MethodGet, "ping", nil)
	return err
}

// GetWatchdog calls the /watchdog/:id endpoint.
// It returns, as json data, the data of a watchdog with a specific id.
func (x *Api) GetWatchdog(watchdogId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "watchdog/"+watchdogId, nil)
	return data, err
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
// It resets watchdog with a specific id.
func (x *Api) PostWatchdogReset(watchdogId string) error {
	_, err := x.sender.Send(http.MethodPost, "watchdog/"+watchdogId+"/reset", nil)
	return err
}

// GetMachine calls the /machine/:id endpoint.
// It returns, as json data, the data of a machine with a specific id.
func (x *Api) GetMachine(machineId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "machine/"+machineId, nil)
	return data, err
}

// PostMachineSample calls the /machine/:id/sample endpoint.
// It is used to upload a cpu/mem/disk usage sample.
// The tstamp parameter is the number of milliseconds since 1970-01-01T00:00:00Z, always UTC, never local time.
// The cpu, mem and disk parameters are usage precentages between 0 and 100 inclusively.
func (x *Api) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	_, err := x.sender.Send(http.MethodPost, "machine/"+machineId+"/sample", []byte(body))
	return err
}

// GetMetric calls the /metric/:id endpoint.
// It returns, as json data, the data of a metric with a specific id.
func (x *Api) GetMetric(metricId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "metric/"+metricId, nil)
	return data, err
}

// PostMetricInc calls the /metric/:id/inc endpoint.
// It is used to increment a counter metric.
// Value is a non-negative int64 value.
func (x *Api) PostMetricInc(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := x.sender.Send(http.MethodPost, "metric/"+metricId+"/inc", []byte(body))
	return err
}

// PostMetricSet calls the /metric/:id/set endpoint.
// It is used to set a gauge metric.
// Value is a non-negative int64 value.
func (x *Api) PostMetricSet(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := x.sender.Send(http.MethodPost, "metric/"+metricId+"/set", []byte(body))
	return err
}
