package monibot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Api provides access to the Monibot REST API.
type Api struct {
	sender Sender
}

// NewApi creates an Api that sends data to https://monibot.io.
func NewApi(apiKey string) *Api {
	return NewApiWithSender(NewSenderWithOptions(apiKey, SenderOptions{}))
}

// NewApiWithSender creates an Api that uses sender for sending data.
func NewApiWithSender(sender Sender) *Api {
	return &Api{sender}
}

// GetPing calls the /ping endpoint.
// It is used to ensure everything is set up correctly and the API is reachable.
// It returns nil on success or an error if something goes wrong.
func (x *Api) GetPing() error {
	_, err := x.sender.Send(http.MethodGet, "ping", nil)
	return err
}

// GetWatchdogs calls the /watchdogs endpoint.
// It returns a list of watchdogs
func (x *Api) GetWatchdogs() ([]Watchdog, error) {
	data, err := x.sender.Send(http.MethodGet, "watchdogs", nil)
	if err != nil {
		return nil, err
	}
	var watchdogs []Watchdog
	err = json.Unmarshal(data, &watchdogs)
	return watchdogs, err
}

// GetWatchdog calls the /watchdog/:id endpoint.
// It returns a watchdog by id.
func (x *Api) GetWatchdog(watchdogId string) (Watchdog, error) {
	data, err := x.sender.Send(http.MethodGet, "watchdog/"+watchdogId, nil)
	if err != nil {
		return Watchdog{}, err
	}
	var w Watchdog
	err = json.Unmarshal(data, &w)
	return w, err
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
// It resets watchdog with a specific id.
func (x *Api) PostWatchdogReset(watchdogId string) error {
	_, err := x.sender.Send(http.MethodPost, "watchdog/"+watchdogId+"/reset", nil)
	return err
}

// GetMachines calls the /machines endpoint.
// It returns a list of machines.
func (x *Api) GetMachines() ([]Machine, error) {
	data, err := x.sender.Send(http.MethodGet, "machines", nil)
	if err != nil {
		return nil, err
	}
	var machines []Machine
	err = json.Unmarshal(data, &machines)
	return machines, err
}

// GetMachine calls the /machine/:id endpoint.
// It returns a machine by id.
func (x *Api) GetMachine(machineId string) (Machine, error) {
	data, err := x.sender.Send(http.MethodGet, "machine/"+machineId, nil)
	if err != nil {
		return Machine{}, err
	}
	var machine Machine
	err = json.Unmarshal(data, &machine)
	return machine, err
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

// GetMetrics calls the /metrics endpoint.
// It returns a list of metrics as json data.
func (x *Api) GetMetrics() ([]Metric, error) {
	data, err := x.sender.Send(http.MethodGet, "metrics", nil)
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	err = json.Unmarshal(data, &metrics)
	return metrics, err
}

// GetMetric calls the /metric/:id endpoint.
// It returns a metric by id as json data.
func (x *Api) GetMetric(metricId string) (Metric, error) {
	data, err := x.sender.Send(http.MethodGet, "metric/"+metricId, nil)
	if err != nil {
		return Metric{}, err
	}
	var metric Metric
	err = json.Unmarshal(data, &metric)
	return metric, err
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
