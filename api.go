package monibot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Api provides access to the Monibot REST API.
type Api struct {
	sender Sender
}

// NewApi creates an Api that sends data to https://monibot.io and retries on error.
func NewApi(apiKey string) *Api {
	return NewApiWithSender(NewRetrySender(NewSender(apiKey)))
}

// NewApiWithSender creates an Api that uses sender for sending data.
func NewApiWithSender(sender Sender) *Api {
	return &Api{sender}
}

// GetPing wraps GetPingWithContext using context.Background.
func (a *Api) GetPing() error {
	return a.GetPingWithContext(context.Background())
}

// GetPingWithContext calls the /ping endpoint.
// It is used to ensure everything is set up correctly and the API is reachable.
// It returns nil on success or an error if something goes wrong.
func (a *Api) GetPingWithContext(ctx context.Context) error {
	_, err := a.sender.Send(ctx, http.MethodGet, "ping", nil)
	return err
}

// GetWatchdogs wraps GetWatchdogsWithContext using context.Background.
func (a *Api) GetWatchdogs() ([]Watchdog, error) {
	return a.GetWatchdogsWithContext(context.Background())
}

// GetWatchdogsWithContext calls the /watchdogs endpoint.
// It returns a list of watchdogs
func (a *Api) GetWatchdogsWithContext(ctx context.Context) ([]Watchdog, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "watchdogs", nil)
	if err != nil {
		return nil, err
	}
	var watchdogs []Watchdog
	err = json.Unmarshal(data, &watchdogs)
	return watchdogs, err
}

// GetWatchdog wraps GetWatchdogWithContext using context.Background.
func (a *Api) GetWatchdog(watchdogId string) (Watchdog, error) {
	return a.GetWatchdogWithContext(context.Background(), watchdogId)
}

// GetWatchdogWithContext calls the /watchdog/:id endpoint.
// It returns a watchdog by id.
func (a *Api) GetWatchdogWithContext(ctx context.Context, watchdogId string) (Watchdog, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "watchdog/"+watchdogId, nil)
	if err != nil {
		return Watchdog{}, err
	}
	var w Watchdog
	err = json.Unmarshal(data, &w)
	return w, err
}

// PostWatchdogReset wraps PostWatchdogResetWithContext using context.Background.
func (a *Api) PostWatchdogReset(watchdogId string) error {
	return a.PostWatchdogResetWithContext(context.Background(), watchdogId)
}

// PostWatchdogReset calls the /watchdog/:id/reset endpoint.
// It resets watchdog with a specific id.
func (a *Api) PostWatchdogResetWithContext(ctx context.Context, watchdogId string) error {
	_, err := a.sender.Send(ctx, http.MethodPost, "watchdog/"+watchdogId+"/reset", nil)
	return err
}

// GetMachines wraps GetMachinesWithContext using context.Background.
func (a *Api) GetMachines() ([]Machine, error) {
	return a.GetMachinesWithContext(context.Background())
}

// GetMachines calls the /machines endpoint.
// It returns a list of machines.
func (a *Api) GetMachinesWithContext(ctx context.Context) ([]Machine, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "machines", nil)
	if err != nil {
		return nil, err
	}
	var machines []Machine
	err = json.Unmarshal(data, &machines)
	return machines, err
}

// GetMachine wraps GetMachineWithContext using context.Background.
func (a *Api) GetMachine(machineId string) (Machine, error) {
	return a.GetMachineWithContext(context.Background(), machineId)
}

// GetMachine calls the /machine/:id endpoint.
// It returns a machine by id.
func (a *Api) GetMachineWithContext(ctx context.Context, machineId string) (Machine, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "machine/"+machineId, nil)
	if err != nil {
		return Machine{}, err
	}
	var machine Machine
	err = json.Unmarshal(data, &machine)
	return machine, err
}

// PostMachineSample wraps PostMachineSampleWithContext using context.Background.
func (a *Api) PostMachineSample(machineId string, tstamp int64, cpu, mem, disk int) error {
	return a.PostMachineSampleWithContext(context.Background(), machineId, tstamp, cpu, mem, disk)
}

// PostMachineSample calls the /machine/:id/sample endpoint.
// It is used to upload a cpu/mem/disk usage sample.
// The tstamp parameter is the number of milliseconds since
// 1970-01-01T00:00:00Z, always UTC, never local time.
// The cpu, mem and disk parameters are usage precentages between
// 0 and 100 inclusively.
func (a *Api) PostMachineSampleWithContext(ctx context.Context, machineId string, tstamp int64, cpu, mem, disk int) error {
	body := fmt.Sprintf("tstamp=%d&cpu=%d&mem=%d&disk=%d", tstamp, cpu, mem, disk)
	_, err := a.sender.Send(ctx, http.MethodPost, "machine/"+machineId+"/sample", []byte(body))
	return err
}

// GetMetrics wraps GetMetricsWithContext using context.Background.
func (a *Api) GetMetrics() ([]Metric, error) {
	return a.GetMetricsWithContext(context.Background())
}

// GetMetrics calls the /metrics endpoint.
// It returns a list of metrics.
func (a *Api) GetMetricsWithContext(ctx context.Context) ([]Metric, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "metrics", nil)
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	err = json.Unmarshal(data, &metrics)
	return metrics, err
}

// GetMetric wraps GetMetricWithContext using context.Background.
func (a *Api) GetMetric(metricId string) (Metric, error) {
	return a.GetMetricWithContext(context.Background(), metricId)
}

// GetMetric calls the /metric/:id endpoint.
// It returns a metric by id.
func (a *Api) GetMetricWithContext(ctx context.Context, metricId string) (Metric, error) {
	data, err := a.sender.Send(ctx, http.MethodGet, "metric/"+metricId, nil)
	if err != nil {
		return Metric{}, err
	}
	var metric Metric
	err = json.Unmarshal(data, &metric)
	return metric, err
}

// PostMetricInc wraps PostMetricIncWithContext using context.Background.
func (a *Api) PostMetricInc(metricId string, value int64) error {
	return a.PostMetricIncWithContext(context.Background(), metricId, value)
}

// PostMetricInc calls the /metric/:id/inc endpoint.
// It is used to increment a counter metric.
// Value is a non-negative int64 value.
func (a *Api) PostMetricIncWithContext(ctx context.Context, metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := a.sender.Send(ctx, http.MethodPost, "metric/"+metricId+"/inc", []byte(body))
	return err
}

// PostMetricSet wraps PostMetricSetWithContext using context.Background.
func (a *Api) PostMetricSet(metricId string, value int64) error {
	return a.PostMetricSetWithContext(context.Background(), metricId, value)
}

// PostMetricSet calls the /metric/:id/set endpoint.
// It is used to set a gauge metric.
// Value is a non-negative int64 value.
func (a *Api) PostMetricSetWithContext(ctx context.Context, metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := a.sender.Send(ctx, http.MethodPost, "metric/"+metricId+"/set", []byte(body))
	return err
}
