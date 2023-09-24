package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type SleepFunc func(time.Duration)

// Conn is a Monibot API connection.
// It's a logical connection, not a physical one.
type Conn struct {
	logger Logger
	http   Http
	sleep  SleepFunc // nil means "do not sleep"
}

// NewDefaultConn creates a Conn with a default Http implementation.
func NewDefaultConn(userAgent, apiKey string) *Conn {
	logger := NewLogger(io.Discard, false)
	http := NewHttp(logger, "https://monibot.io", userAgent, apiKey)
	return NewConn(logger, http, time.Sleep)
}

// NewConn creates a Conn with a custom Http implementation and sleep function.
func NewConn(logger Logger, http Http, sleep SleepFunc) *Conn {
	if sleep == nil {
		sleep = func(time.Duration) {}
	}
	return &Conn{logger, http, sleep}
}

// GetPing calls the /ping endpoint.
func (c *Conn) GetPing() error {
	_, err := c.http.Send(http.MethodGet, "ping", nil)
	return err
}

// GetWatchdog calls the /watchdog/:id endpoint.
func (c *Conn) GetWatchdog(watchdogId string) ([]byte, error) {
	data, err := c.http.Send(http.MethodGet, "watchdog/"+watchdogId, nil)
	return data, err
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
// It tries max trials with a delay between trials.
func (c *Conn) PostWatchdogReset(watchdogId string, trials int, delay time.Duration) error {
	return c.try(func() error {
		_, err := c.http.Send(http.MethodPost, "watchdog/"+watchdogId+"/reset", nil)
		return err
	}, trials, delay)
}

// GetMachine calls the /machine/:id endpoint.
func (c *Conn) GetMachine(machineId string) ([]byte, error) {
	data, err := c.http.Send(http.MethodGet, "machine/"+machineId, nil)
	return data, err
}

// PostMachineSample calls the /machine/:id/sample endpoint.
// It tries max trials with a delay between trials.
func (c *Conn) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int, trials int, delay time.Duration) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	return c.try(func() error {
		_, err := c.http.Send(http.MethodPost, "machine/"+machineId+"/sample", []byte(body))
		return err
	}, trials, delay)
}

// GetMetric calls the /metric/:id endpoint.
func (c *Conn) GetMetric(metricId string) ([]byte, error) {
	data, err := c.http.Send(http.MethodGet, "metric/"+metricId, nil)
	return data, err
}

// PostMetricInc calls the /metric/:id/inc endpoint.
func (c *Conn) PostMetricInc(metricId string, value int64, trials int, delay time.Duration) error {
	body := fmt.Sprintf("value=%d", value)
	return c.try(func() error {
		_, err := c.http.Send(http.MethodPost, "metric/"+metricId+"/inc", []byte(body))
		return err
	}, trials, delay)
}

// PostMetricSet calls the /metric/:id/set endpoint.
func (c *Conn) PostMetricSet(metricId string, value int64, trials int, delay time.Duration) error {
	body := fmt.Sprintf("value=%d", value)
	return c.try(func() error {
		_, err := c.http.Send(http.MethodPost, "metric/"+metricId+"/set", []byte(body))
		return err
	}, trials, delay)
}

func (c *Conn) try(f func() error, trials int, delay time.Duration) error {
	var err error
	for i := 0; i < trials; i++ {
		if i > 0 {
			c.logger.Debugf("will sleep %s, then start trial %d/%d", delay, i+1, trials)
			c.sleep(delay)
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
