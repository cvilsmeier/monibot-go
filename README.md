# monibot-go

Golang SDK and CLI tool for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## SDK Usage

    $ go get github.com/cvilsmeier/monibot-go

```go
import "github.com/cvilsmeier/monibot-go"

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
```

## CLI Usage

Build from source, needs Go. See also <https://go.dev>.

    $ go install github.com/cvilsmeier/monibot-go/cmd/moni@latest

Show help page with the following command:

    $ moni help


## Changelog

See CHANGELOG.md


## License

MIT License, see LICENSE
