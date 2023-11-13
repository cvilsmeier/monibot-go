package main

import (
	"log"
	"os"

	"github.com/cvilsmeier/monibot-go"
)

func main() {
	// api access requires an apiKey
	apiKey := os.Getenv("MONIBOT_API_KEY")
	// create new api
	api := monibot.NewApi(apiKey)
	// send a watchdog heartbeat
	err := api.PostWatchdogHeartbeat("5f6d343a471d87687f51771530c3f2f4")
	if err != nil {
		log.Fatal(err)
	}
}
