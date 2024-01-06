# monibot-go

[![GoDoc Reference](https://godoc.org/github.com/cvilsmeier/monibot-go?status.svg)](http://godoc.org/github.com/cvilsmeier/monibot-go)
[![Build Status](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml/badge.svg)](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Golang SDK for https://monibot.io - Easy Server and Application Monitoring.

This module provides a SDK to interact with the Monibot REST API.
Monibot is a service that monitors your web apps, servers and
metrics, and notifies you if something goes wrong.


## Usage

    $ go get github.com/cvilsmeier/monibot-go

```go
import "github.com/cvilsmeier/monibot-go"

func main() {
	// api access requires an apiKey
	apiKey := os.Getenv("MONIBOT_API_KEY")
	// create new api
	api := monibot.NewApi(apiKey)
	// send a watchdog heartbeat
	api.PostWatchdogHeartbeat("5f6d343a471d87687f51771530c3f2f4")
	// increment a counter metric
	api.PostMetricInc("c3f2fefae7f6d3e387f1d8761ff17730", 42)
}
```


## Changelog

### v0.1.1

- add histogram values functions

### v0.1.0

- add histogram metric values

### v0.0.9

- add machine text

### v0.0.8

- fix disk/net usage samples

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
