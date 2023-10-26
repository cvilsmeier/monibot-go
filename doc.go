/*
Package monibot provides a SDK to interact with the Monibot REST API,
see https://monibot.io for details.

	import "github.com/cvilsmeier/monibot-go"

	// init the api
	apiKey := os.Getenv("MONIBOT_API_KEY")
	api := NewApi(apiKey)
	// reset a watchdog
	err := api.PostWatchdogReset(context.Background(), "2f5f6d47183fdf415a7476837351730c")
	if err != nil {
		log.Fatal(err)
	}

Monibot monitors your web app and notifies you if something goes wrong.
It monitors Site Reachability, SSL/TLS Certificates, Watchdog Heartbeats,
Machine Resource Usage, Database Size, Number of newly registered users,
Number of sold articles, Number of failed login attempts, and many more.
*/
package monibot
