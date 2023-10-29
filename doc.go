/*
Package monibot provides a SDK to interact with the Monibot REST API,
see https://monibot.io for details.

	import "github.com/cvilsmeier/monibot-go"

	func main() {
		// api access requires an apiKey
		apiKey := os.Getenv("MONIBOT_API_KEY")
		// create new api
		api := monibot.NewApi(apiKey)
		// reset a watchdog by id
		err := api.PostWatchdogReset("5f6d343f517715a471d8768730c3f2f4")
		if err != nil {
			log.Fatal(err)
		}
	}
*/
package monibot
