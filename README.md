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
const userAgent = "my-app/v1.0.0"
const apiKey = os.Getenv("MONIBOT_API_KEY")    
conn := api.NewDefaultConn(userAgent, apiKey)
// use the API
err := conn.GetPing()
if err != nil {
    log.Fatal(err)
}
```


## License

MIT license