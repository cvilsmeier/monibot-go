# monibot-go

[![GoDoc Reference](https://godoc.org/github.com/cvilsmeier/monibot-go?status.svg)](http://godoc.org/github.com/cvilsmeier/monibot-go)
[![Build Status](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml/badge.svg)](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Golang SDK for https://monibot.io

## Install

    $ go get github.com/cvilsmeier/monibot-go


## Usage

```go
import "github.com/cvilsmeier/monibot-go"

func main() {
	// api access requires an apiKey
	apiKey := os.Getenv("MONIBOT_API_KEY")
	// create new api
	api := monibot.NewApi(apiKey)
	// send a watchdog heartbeat
	err := api.PostWatchdogHeartbeat("5f6d343a471d87687f51771530c3f2f4")
	if err != nil {
		log.Fatal(err)
	}
}
```

## Changelog

### v0.0.7

- added netRecv and netSend to machine sample

### v0.0.6

- added diskReads and diskWrites to machine sample

### v0.0.5

- moved moni command line tool into own repo: https://github.com/cvilsmeier/moni

### v0.0.4

- first version


## License

MIT License, see LICENSE
