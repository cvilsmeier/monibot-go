# moni-go

A Go SDK and CLI for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## Install CLI

    go install github.com/cvilsmeier/moni-cli/cmd/moni@latest

## CLI Usage

    $ moni help


## SDK Usage

```go

    import "github.com/cvilsmeier/moni-go/api"

    // init the API
	const verbose = true
	const monibotUrl = "http://monibot.io"
	const userAgent = "acme-app/v1.0.0"
	const apiKey = os.Getenc("MONIBOT_API_KEY")    
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