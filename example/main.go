package main

import (
	"os"

	"github.com/cvilsmeier/monibot-go"
)

func main() {
	// api access requires an apiKey
	apiKey := os.Getenv("MONIBOT_API_KEY")
	// create new api
	api := monibot.NewApi(apiKey)
	// send a watchdog heartbeat
	api.PostWatchdogHeartbeat("5f6d343a471d87687f51771530c3f2f4")
	// increment a counter metric
	api.PostMetricInc("c3f2fefae7f6d3e387f1d8761ff17730", 42)
}
