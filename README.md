# monibot-go

[![GoDoc Reference](https://godoc.org/github.com/cvilsmeier/monibot-go?status.svg)](http://godoc.org/github.com/cvilsmeier/monibot-go)
[![Build Status](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml/badge.svg)](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Golang SDK and CLI tool for <https://monibot.io>.

PLEASE NOTE: Monibot is still under development, visit <https://monibot.io> for details.

## SDK Usage

    $ go get github.com/cvilsmeier/monibot-go

```go
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
