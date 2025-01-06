/*
Package monibot provides a SDK to interact with the Monibot REST API,
see https://monibot.io for details.

	import "github.com/cvilsmeier/monibot-go"

	func main() {
		// api access requires an apiKey
		apiKey := os.Getenv("MONIBOT_API_KEY")
		// create new api
		api := monibot.NewApi(apiKey)
		// send a watchdog heartbeat
		api.PostWatchdogHeartbeat("a749ff35891ecb36")
		// increment a counter metric by 42
		api.PostMetricInc("ffe31498bc7193a4", 42)
	}
*/
package monibot
