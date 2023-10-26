# monibot-go

Go SDK and CLI for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## CLI Usage

Build from source, needs Go. See also <https://go.dev>.

    $ go install github.com/cvilsmeier/monibot-go/cmd/moni@latest

Show help page with the following command:

    $ moni help

## SDK Usage

    $ go get github.com/cvilsmeier/monibot-go

```go
import "github.com/cvilsmeier/monibot-go"

// init the api
apiKey := os.Getenv("MONIBOT_API_KEY")
api := NewApi(apiKey)
// reset a watchdog
err := api.PostWatchdogReset("2f5f6d47183fdf415a7476837351730c")
if err != nil {
    log.Fatal(err)
}
```

## Changelog

See CHANGELOG.md

## License

MIT license