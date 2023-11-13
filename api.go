package monibot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/logging"
	"github.com/cvilsmeier/monibot-go/internal/sending"
	"github.com/cvilsmeier/monibot-go/internal/version"
)

// A Logger prints debug messages.
type Logger = logging.Logger

// TimeAfterFunc is the function type of time.After.
type TimeAfterFunc = sending.TimeAfterFunc

// Version is monibot-go sdk version.
const Version = version.Version

// ApiOptions holds optional parameters for a Api.
type ApiOptions struct {

	// Default is no logging.
	Logger Logger

	// Default is "https://monibot.io".
	MonibotUrl string

	// Default is 12 trials.
	Trials int

	// Default is 5s delay.
	Delay time.Duration

	// Default time.After
	TimeAfter TimeAfterFunc
}

// An apiSender provides a Send method and can be overridden in unit tests.
type apiSender interface {
	Send(ctx context.Context, method, path string, body []byte) ([]byte, error)
}

// Api provides access to the Monibot REST API.
type Api struct {
	sender apiSender
}

// NewApi creates an Api that sends data to https://monibot.io
// and retries 12 times every 5s if an error occurs,
// and logs nothing.
func NewApi(apiKey string) *Api {
	return NewApiWithOptions(apiKey, ApiOptions{})
}

// NewApiWithOptions creates an Api with custom options.
func NewApiWithOptions(apiKey string, opt ApiOptions) *Api {
	logger := opt.Logger
	if logger == nil {
		logger = logging.NewDiscardLogger()
	}
	monibotUrl := opt.MonibotUrl
	if monibotUrl == "" {
		monibotUrl = "http://monibot.io"
	}
	trials := opt.Trials
	if trials == 0 {
		trials = 12
	}
	delay := opt.Delay
	if delay == 0 {
		delay = 5 * time.Second
	}
	timeAfter := opt.TimeAfter
	if timeAfter == nil {
		timeAfter = time.After
	}
	transport := sending.NewTransport(logger, Version, monibotUrl, apiKey)
	sender := sending.NewSender(transport, logger, trials, delay, timeAfter)
	return &Api{sender}
}

// GetPing is like GetPingWithContext using context.Background.
func (a *Api) GetPing() error {
	return a.GetPingWithContext(context.Background())
}

// GetPingWithContext pings the API.
// It is used to ensure everything is set up correctly and the API is
// reachable. It returns nil on success or a non-nil error if
// something goes wrong.
func (a *Api) GetPingWithContext(ctx context.Context) error {
	_, err := a.sender.Send(ctx, "GET", "ping", nil)
	return err
}

// GetWatchdogs is like GetWatchdogsWithContext using context.Background.
func (a *Api) GetWatchdogs() ([]Watchdog, error) {
	return a.GetWatchdogsWithContext(context.Background())
}

// GetWatchdogsWithContext fetches the list of watchdogs.
func (a *Api) GetWatchdogsWithContext(ctx context.Context) ([]Watchdog, error) {
	data, err := a.sender.Send(ctx, "GET", "watchdogs", nil)
	if err != nil {
		return nil, err
	}
	var watchdogs []Watchdog
	err = json.Unmarshal(data, &watchdogs)
	return watchdogs, err
}

// GetWatchdog is like GetWatchdogWithContext using context.Background.
func (a *Api) GetWatchdog(watchdogId string) (Watchdog, error) {
	return a.GetWatchdogWithContext(context.Background(), watchdogId)
}

// GetWatchdogWithContext fetches a watchdog by id.
func (a *Api) GetWatchdogWithContext(ctx context.Context, watchdogId string) (Watchdog, error) {
	data, err := a.sender.Send(ctx, "GET", "watchdog/"+watchdogId, nil)
	if err != nil {
		return Watchdog{}, err
	}
	var w Watchdog
	err = json.Unmarshal(data, &w)
	return w, err
}

// PostWatchdogHeartbeat is like PostWatchdogHeartbeatWithContext using context.Background.
func (a *Api) PostWatchdogHeartbeat(watchdogId string) error {
	return a.PostWatchdogHeartbeatWithContext(context.Background(), watchdogId)
}

// PostWatchdogHeartbeatWithContext sends a watchdog heartbeat.
func (a *Api) PostWatchdogHeartbeatWithContext(ctx context.Context, watchdogId string) error {
	_, err := a.sender.Send(ctx, "POST", "watchdog/"+watchdogId+"/heartbeat", nil)
	return err
}

// GetMachines is like GetMachinesWithContext using context.Background.
func (a *Api) GetMachines() ([]Machine, error) {
	return a.GetMachinesWithContext(context.Background())
}

// GetMachinesWithContext fetches the list of machines.
func (a *Api) GetMachinesWithContext(ctx context.Context) ([]Machine, error) {
	data, err := a.sender.Send(ctx, "GET", "machines", nil)
	if err != nil {
		return nil, err
	}
	var machines []Machine
	err = json.Unmarshal(data, &machines)
	return machines, err
}

// GetMachine is like GetMachineWithContext using context.Background.
func (a *Api) GetMachine(machineId string) (Machine, error) {
	return a.GetMachineWithContext(context.Background(), machineId)
}

// GetMachineWithContext fetches a machine by id.
func (a *Api) GetMachineWithContext(ctx context.Context, machineId string) (Machine, error) {
	data, err := a.sender.Send(ctx, "GET", "machine/"+machineId, nil)
	if err != nil {
		return Machine{}, err
	}
	var machine Machine
	err = json.Unmarshal(data, &machine)
	return machine, err
}

// PostMachineSample is like PostMachineSampleWithContext using context.Background.
func (a *Api) PostMachineSample(machineId string, sample MachineSample) error {
	return a.PostMachineSampleWithContext(context.Background(), machineId, sample)
}

// PostMachineSampleWithContext uploads a machine sample to the API.
func (a *Api) PostMachineSampleWithContext(ctx context.Context, machineId string, sample MachineSample) error {
	toks := []string{
		fmt.Sprintf("tstamp=%d", sample.Tstamp),
		fmt.Sprintf("load1=%.3f", sample.Load1),
		fmt.Sprintf("load5=%.3f", sample.Load5),
		fmt.Sprintf("load15=%.3f", sample.Load15),
		fmt.Sprintf("cpu=%d", sample.CpuPercent),
		fmt.Sprintf("mem=%d", sample.MemPercent),
		fmt.Sprintf("disk=%d", sample.DiskPercent),
		fmt.Sprintf("diskReads=%d", sample.DiskReads),
		fmt.Sprintf("diskWrites=%d", sample.DiskWrites),
		fmt.Sprintf("netRecv=%d", sample.NetRecv),
		fmt.Sprintf("netSend=%d", sample.NetSend),
	}
	body := strings.Join(toks, "&")
	_, err := a.sender.Send(ctx, "POST", "machine/"+machineId+"/sample", []byte(body))
	return err
}

// GetMetrics is like GetMetricsWithContext using context.Background.
func (a *Api) GetMetrics() ([]Metric, error) {
	return a.GetMetricsWithContext(context.Background())
}

// GetMetricsWithContext fetches the list of metrics.
func (a *Api) GetMetricsWithContext(ctx context.Context) ([]Metric, error) {
	data, err := a.sender.Send(ctx, "GET", "metrics", nil)
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	err = json.Unmarshal(data, &metrics)
	return metrics, err
}

// GetMetric is like GetMetricWithContext using context.Background.
func (a *Api) GetMetric(metricId string) (Metric, error) {
	return a.GetMetricWithContext(context.Background(), metricId)
}

// GetMetricWithContext fetches a metric by id.
func (a *Api) GetMetricWithContext(ctx context.Context, metricId string) (Metric, error) {
	data, err := a.sender.Send(ctx, "GET", "metric/"+metricId, nil)
	if err != nil {
		return Metric{}, err
	}
	var metric Metric
	err = json.Unmarshal(data, &metric)
	return metric, err
}

// PostMetricInc is like PostMetricIncWithContext using context.Background.
func (a *Api) PostMetricInc(metricId string, value int64) error {
	return a.PostMetricIncWithContext(context.Background(), metricId, value)
}

// PostMetricIncWithContext uploads a counter metric increment value to the API.
// It is used to increment a counter metric.
// It is an error to try to increment a non-counter metric.
// The value is a non-negative int64 number.
func (a *Api) PostMetricIncWithContext(ctx context.Context, metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := a.sender.Send(ctx, "POST", "metric/"+metricId+"/inc", []byte(body))
	return err
}

// PostMetricSet is like PostMetricSetWithContext using context.Background.
func (a *Api) PostMetricSet(metricId string, value int64) error {
	return a.PostMetricSetWithContext(context.Background(), metricId, value)
}

// PostMetricSetWithContext uploads a gauge metric value to the API.
// It is used to set a gauge metric.
// It is an error to try to set a non-gauge metric.
// The value is a non-negative int64 number.
func (a *Api) PostMetricSetWithContext(ctx context.Context, metricId string, value int64) error {
	body := fmt.Sprintf("value=%d", value)
	_, err := a.sender.Send(ctx, "POST", "metric/"+metricId+"/set", []byte(body))
	return err
}
