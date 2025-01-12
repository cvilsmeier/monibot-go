# monibot-go

[![GoDoc Reference](https://godoc.org/github.com/cvilsmeier/monibot-go?status.svg)](http://godoc.org/github.com/cvilsmeier/monibot-go)
[![Build Status](https://github.com/cvilsmeier/monibot-go/actions/workflows/go-linux.yml/badge.svg)](https://github.com/cvilsmeier/monibot-go/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Golang SDK for https://monibot.io - Website-, Server- and Application Monitoring.

This module provides a SDK to interact with the Monibot REST API.
Monibot is a service that monitors your web apps, servers and
application metrics, and notifies you if something goes wrong.


## Usage

    go get github.com/cvilsmeier/monibot-go

```go
import "github.com/cvilsmeier/monibot-go"

func main() {
	// api access requires an apiKey
	apiKey := os.Getenv("MONIBOT_API_KEY")
	// create new api
	api := monibot.NewApi(apiKey)
	// send a watchdog heartbeat
	api.PostWatchdogHeartbeat("a749ff35891ecb36")
	// increment a counter metric by 42
	api.PostMetricInc("ffe31498bc7193a4", 42)
}
```


## Changelog

### v0.2.0

- replace MachineSample DiskReads/Writes (number of sectors) with DiskRead/Write (number of bytes)

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
