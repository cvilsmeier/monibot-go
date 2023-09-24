package api

import (
	"fmt"
)

type Conn struct {
	http Http
}

func NewConn(http Http) *Conn {
	return &Conn{http}
}

func (c *Conn) GetPing() error {
	_, err := c.http.Get("ping")
	return err
}

func (c *Conn) GetWatchdog(watchdogId string) ([]byte, error) {
	data, err := c.http.Get("watchdog/" + watchdogId)
	return data, err
}

func (c *Conn) PostWatchdogReset(watchdogId string) error {
	_, err := c.http.Post("watchdog/"+watchdogId+"/reset", nil)
	return err
}

func (c *Conn) GetMachine(machineId string) ([]byte, error) {
	data, err := c.http.Get("machine/" + machineId)
	return data, err
}

func (c *Conn) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	_, err := c.http.Post("machine/"+machineId+"/sample", []byte(body))
	return err
}

func (c *Conn) GetMetric(metricId string) ([]byte, error) {
	data, err := c.http.Get("metric/" + metricId)
	return data, err
}

func (c *Conn) PostMetricInc(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := c.http.Post("metric/"+metricId+"/inc", []byte(body))
	return err
}

func (c *Conn) PostMetricSet(metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := c.http.Post("metric/"+metricId+"/set", []byte(body))
	return err
}
