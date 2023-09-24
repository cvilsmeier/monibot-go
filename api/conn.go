package api

import (
	"fmt"
	"io"
)

// Conn is a Monibot API connection.
// It's a logical connection, not a physical one.
type Conn struct {
	http Http
}

// NewDefaultConn creates a Conn with a default Http implementation.
func NewDefaultConn(userAgent, apiKey string) *Conn {
	logger := NewLogger(io.Discard, false)
	http := NewHttp(logger, "https://monibot.io", userAgent, apiKey)
	return &Conn{http}
}

// NewConn creates a Conn with a custom Http implementation.
func NewConn(http Http) *Conn {
	return &Conn{http}
}

// GetPing calls the /ping endpoint.
func (c *Conn) GetPing() error {
	_, err := c.http.Get("ping")
	return err
}

// GetWatchdog calls the /watchdog/:id endpoint.
func (c *Conn) GetWatchdog(watchdogId string) ([]byte, error) {
	data, err := c.http.Get("watchdog/" + watchdogId)
	return data, err
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
func (c *Conn) PostWatchdogReset(watchdogId string) error {
	_, err := c.http.Post("watchdog/"+watchdogId+"/reset", nil)
	return err
}

// GetMachine calls the /machine/:id endpoint.
func (c *Conn) GetMachine(machineId string) ([]byte, error) {
	data, err := c.http.Get("machine/" + machineId)
	return data, err
}

// PostMachineSample calls the /machine/:id/sample endpoint.
func (c *Conn) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	_, err := c.http.Post("machine/"+machineId+"/sample", []byte(body))
	return err
}

// GetMetric calls the /metric/:id endpoint.
func (c *Conn) GetMetric(metricId string) ([]byte, error) {
	data, err := c.http.Get("metric/" + metricId)
	return data, err
}

// PostMetricInc calls the /metric/:id/inc endpoint.
func (c *Conn) PostMetricInc(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := c.http.Post("metric/"+metricId+"/inc", []byte(body))
	return err
}

// PostMetricSet calls the /metric/:id/set endpoint.
func (c *Conn) PostMetricSet(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := c.http.Post("metric/"+metricId+"/set", []byte(body))
	return err
}
