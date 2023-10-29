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
	err := api.PostWatchdogHeartbeat("5f6d343f517715a471d8768730c3f2f4")
	if err != nil {
		log.Fatal(err)
	}
}
