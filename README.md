# moni-go

A Go SDK and CLI for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## CLI Usage

Build from source, needs Go. See also <https://go.dev/>.

    $ go install github.com/cvilsmeier/moni-go/cmd/moni@latest

Show help page with the following command:

    $ moni help

## SDK Usage

    $ go get github.com/cvilsmeier/moni-go/api

```go
import "github.com/cvilsmeier/moni-go/api"

// init the API
userAgent := "my-app/v1.0.0"
apiKey := os.Getenv("MONIBOT_API_KEY")
conn := api.NewDefaultConn(userAgent, apiKey)
// ping the API
err := conn.GetPing()
if err != nil {
    log.Fatal(err)
}
// reset a watchdog
err = conn.PostWatchdogReset("9f9f679a44f5f7a0486817ed524c9791")
if err != nil {
    log.Fatal(err)
}
```

## Changelog

See CHANGELOG.md

## License

MIT license