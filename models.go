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
	Tstamp      int64        // Unix time millis since 1970-01-01T00:00:00Z, always UTC, never local time.
	Load1       float64      // Loadavg 1 minute.
	Load5       float64      // Loadavg 5 minutes.
	Load15      float64      // Loadavg 15 minutes.
	CpuPercent  int          // CPU usage percent 0..100 since last sample.
	MemPercent  int          // Memory usage percent 0..100.
	Disks       []DiskSample // Disk samples, one for each disk.
	DiskPercent int          // Disk usage percent 0..100.
	DiskRead    int64        // Number of disk bytes read since last sample.
	DiskWrite   int64        // Number of disk bytes written since last sample.
	Nets        []NetSample  // Network samples, one for each network device.
	NetRecv     int64        // Number of network bytes received since last sample.
	NetSend     int64        // Number of network bytes sent since last sample.
}

// A DiskSample is part of a MachineSample.
type DiskSample struct {
	Device      string // Disk device, e.g. "/dev/sda1"
	Mountpoint  string // Mount point, e.g. "/mnt/HC_Volume_102242198"
	Total       int64  // Total size in bytes, 0..MAX_I64
	Used        int64  // Used size in bytes, 0..MAX_I64
	UsedPercent int    // Usage in percent, 0..100
	ReadBytes   int64  // Number of bytes read since last sample, 0..MAX_I64
	WriteBytes  int64  // Number of bytes written since last sample, 0..MAX_I64
}

// A NetSample is part of a MachineSample.
type NetSample struct {
	Device    string // Disk device, e.g. "/dev/sda1"
	RecvBytes int64  // Number of bytes received since last sample, 0..MAX_I64
	SendBytes int64  // Number of bytes sent since last sample, 0..MAX_I64
}

// Metric holds data for a Metric.
type Metric struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"` // 0=Counter, 1=Gauge, 2=Histogram
}

const (
	MetricTypeCounter   int = 0
	MetricTypeGauge     int = 1
	MetricTypeHistogram int = 2
)
