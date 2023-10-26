package main

import (
	"log"
	"os"

	"github.com/cvilsmeier/monibot-go"
)

func main() {
	log.Printf("intest")
	apiKey := os.Getenv("MONIBOT_API_KEY")
	if apiKey == "" {
		log.Fatalf("need MONIBOT_API_KEY")
	}
	murl := os.Getenv("MONIBOT_URL")
	if murl == "" {
		log.Fatalf("need MONIBOT_URL")
	}
	// logger := log.New(os.Stdout, "Monibot-Debug: ", log.LstdFlags)
	sender := monibot.NewSenderWithOptions(apiKey, monibot.SenderOptions{
		MonibotUrl: murl,
		// Logger:     monibot.NewLogger(logger),
	})
	api := monibot.NewApiWithSender(sender)
	// ping
	err := api.GetPing()
	if err != nil {
		log.Fatalf("cannot GetPing: %s", err)
	}
	log.Printf("ping ok")
	// watchdogs
	watchdogs, err := api.GetWatchdogs()
	if err != nil {
		log.Fatalf("cannot GetWatchdogs: %s", err)
	}
	log.Printf("watchdogs: [%d]", len(watchdogs))
}