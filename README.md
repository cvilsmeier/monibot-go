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

// init api
userAgent := "my-app/v1.0.0"
apiKey := os.Getenv("MONIBOT_API_KEY")
api := NewDefaultApi(userAgent, apiKey)
// ping the api
err := api.GetPing()
if err != nil {
    log.Fatal(err)
}
// reset a watchdog
trials := 3
delay := 10 * time.Second
err = api.PostWatchdogReset("000000000000001", trials, delay)
if err != nil {
    log.Fatal(err)
}
```

## Changelog

See CHANGELOG.md

## License

MIT license