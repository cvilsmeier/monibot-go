package monibot

// Watchdog holds data for a Watchdog.
type Watchdog struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IntervalMillis int64  `json:"intervalMillis"`
}

// Machine holds data for a Machine.
type Machine struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// A MachineSample holds data for a machine resource usage sample.
type MachineSample struct {

	// Unix time millis since 1970-01-01T00:00:00Z, always UTC, never local time.
	Tstamp int64

	// Loadavg 1 minute
	Load1 float64

	// Loadavg 5 minutes
	Load5 float64

	// Loadavg 15 minutes
	Load15 float64

	// CPU usage percent 0..100
	CpuPercent int

	// Memory usage percent 0..100
	MemPercent int

	// Disk usage percent 0..100
	DiskPercent int
}

// Metric holds data for a Metric.
type Metric struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Type MetricType `json:"type"` // TypeCounter or TypeGauge
}

// MetricType is the type of a metric. Currently we have 0 (Counter) and 1 (Gauge).
type MetricType int

const (
	TypeCounter MetricType = 0 // Counter Metric (type 0)
	TypeGauge   MetricType = 1 // Gauge Metric (type 1)
)
