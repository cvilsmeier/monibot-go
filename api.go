package monibot

import (
	"fmt"
	"net/http"
	"time"
)

type SleepFunc func(time.Duration)

// Api provides access to the Monibot REST API.
type Api struct {
	logger Logger
	sender Sender
	sleep  SleepFunc
}

// NewDefaultApi creates an Api with default implementations.
// This should be suited for most use cases.
func NewDefaultApi(userAgent, apiKey string) *Api {
	logger := NewLogger(nil)
	sender := NewSender(logger, "https://monibot.io", userAgent, apiKey)
	return NewApi(logger, sender, time.Sleep)
}

// NewApi creates an Api with a custom sender and sleep function.
func NewApi(logger Logger, sender Sender, sleep SleepFunc) *Api {
	if logger == nil {
		panic("no logger")
	}
	if sender == nil {
		panic("no sender")
	}
	if sleep == nil {
		panic("no sleep")
	}
	return &Api{logger, sender, sleep}
}

// GetPing calls the /ping endpoint.
func (x *Api) GetPing() error {
	_, err := x.sender.Send(http.MethodGet, "ping", nil)
	return err
}

// GetWatchdog calls the /watchdog/:id endpoint.
func (x *Api) GetWatchdog(watchdogId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "watchdog/"+watchdogId, nil)
	return data, err
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
func (x *Api) PostWatchdogReset(watchdogId string) error {
	_, err := x.sender.Send(http.MethodPost, "watchdog/"+watchdogId+"/reset", nil)
	return err
}

// GetMachine calls the /machine/:id endpoint.
func (x *Api) GetMachine(machineId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "machine/"+machineId, nil)
	return data, err
}

// PostMachineSample calls the /machine/:id/sample endpoint.
func (x *Api) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	_, err := x.sender.Send(http.MethodPost, "machine/"+machineId+"/sample", []byte(body))
	return err
}

// GetMetric calls the /metric/:id endpoint.
func (x *Api) GetMetric(metricId string) ([]byte, error) {
	data, err := x.sender.Send(http.MethodGet, "metric/"+metricId, nil)
	return data, err
}

// PostMetricInc calls the /metric/:id/inc endpoint.
func (x *Api) PostMetricInc(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := x.sender.Send(http.MethodPost, "metric/"+metricId+"/inc", []byte(body))
	return err
}

// PostMetricSet calls the /metric/:id/set endpoint.
func (x *Api) PostMetricSet(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := x.sender.Send(http.MethodPost, "metric/"+metricId+"/set", []byte(body))
	return err
}
