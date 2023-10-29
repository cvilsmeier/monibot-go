package main

import (
	"log"
	"os"

	"github.com/cvilsmeier/monibot-go"
)

func main() {
	// init the api, take apiKey from environment
	apiKey := os.Getenv("MONIBOT_API_KEY")
	api := monibot.NewApi(apiKey)
	// reset a watchdog by id
	err := api.PostWatchdogReset("5f6d343f517715a471d8768730c3f2f4")
	if err != nil {
		log.Fatal(err)
	}
}
