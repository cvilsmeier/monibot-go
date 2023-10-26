package monibot

type Watchdog struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IntervalMillis int64  `json:"intervalMillis"`
}

type Machine struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Metric struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

const (
	TypeCounter int = 0
	TypeGauge   int = 1
)
