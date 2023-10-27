package monibot

// A Watchdog represents a Watchdog.
type Watchdog struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IntervalMillis int64  `json:"intervalMillis"`
}

// A Machine represents a Machine.
type Machine struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// A Metric represents a Metric.
type Metric struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Type MetricType `json:"type"` // TypeCounter or TypeGauge
}

type MetricType = int

const (
	TypeCounter MetricType = 0
	TypeGauge   MetricType = 1
)
