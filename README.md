# moni-go

A Go SDK and CLI for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## CLI Usage

    $ go install github.com/cvilsmeier/moni-go/cmd/moni@latest

    $ moni help

## SDK Usage

```go
import "github.com/cvilsmeier/moni-go/api"

// init the API
const verbose = true
const monibotUrl = "http://monibot.io"
const userAgent = "acme-app/v1.0.0"
const apiKey = os.Getenv("MONIBOT_API_KEY")    
logger := api.NewLogger(verbose)
http := api.NewHttp(logger, monibotUrl, userAgent, apiKey)
conn := api.NewConn(http)
// use the API
err := conn.GetPing()
if err != nil {
    log.Fatal(err)
}
```


## License

MIT license