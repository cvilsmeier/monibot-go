/*
Moni is a command line tool for interacting with the Monibot REST API, see https://monibot.io for details.
It supports a number of commands. To get a list of supported commands, run

	$ moni help
*/
package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cvilsmeier/monibot-go"
)

const (
	urlEnvKey  = "MONIBOT_URL"
	urlFlag    = "url"
	defaultUrl = "https://monibot.io"

	apiKeyEnvKey  = "MONIBOT_API_KEY"
	apiKeyFlag    = "apiKey"
	defaultApiKey = ""

	verboseEnvKey     = "MONIBOT_VERBOSE"
	verboseFlag       = "v"
	defaultVerbose    = false
	defaultVerboseStr = "false"

	trialsEnvKey     = "MONIBOT_TRIALS"
	trialsFlag       = "trials"
	defaultTrials    = 12
	defaultTrialsStr = "12"

	delayEnvKey     = "MONIBOT_DELAY"
	delayFlag       = "delay"
	defaultDelay    = 5 * time.Second
	defaultDelayStr = "5s"
)

func usage() {
	print("moni %s", monibot.Version)
	print("")
	print("Monibot command line tool, see https://monibot.io.")
	print("")
	print("Usage")
	print("")
	print("    moni [flags] command")
	print("")
	print("Flags")
	print("")
	print("    -%s", urlFlag)
	print("        Monibot URL, default is %q.", defaultUrl)
	print("        You can set this also via environment variable %s.", urlEnvKey)
	print("")
	print("    -%s", apiKeyFlag)
	print("        Monibot API Key, default is %q.", defaultApiKey)
	print("        You can set this also via environment variable %s.", apiKeyEnvKey)
	print("        You can find your API Key in your profile on https://monibot.io.")
	print("")
	print("    -%s", trialsFlag)
	print("        Max. Send trials, default is %d.", defaultTrials)
	print("        You can set this also via environment variable %s.", trialsEnvKey)
	print("")
	print("    -%s", delayFlag)
	print("        Delay between trials, default is %d.", defaultDelay)
	print("        You can set this also via environment variable %s.", delayEnvKey)
	print("")
	print("    -%s", verboseFlag)
	print("        Verbose output, default is %t.", defaultVerbose)
	print("        You can set this also via environment variable %s ('true' or 'false').", verboseEnvKey)
	print("")
	print("Commands")
	print("")
	print("    ping")
	print("        Ping the Monibot API.")
	print("")
	print("    watchdogs")
	print("        List watchdogs.")
	print("")
	print("    watchdog <watchdogId>")
	print("        Get watchdog by id.")
	print("")
	print("    reset <watchdogId>")
	print("        Reset a watchdog.")
	print("")
	print("    machines")
	print("        List machines.")
	print("")
	print("    machine <machineId>")
	print("        Get machine by id.")
	print("")
	print("    sample <machineId>")
	print("        Send resource usage (cpu/mem/disk) samples for machine.")
	print("")
	print("    metrics")
	print("        List metrics.")
	print("")
	print("    metric <metricId>")
	print("        Get and print metric info.")
	print("")
	print("    inc <metricId> <value>")
	print("        Increment a Counter metric.")
	print("        Value must be a non-negative 64-bit integer value.")
	print("")
	print("    set <metricId> <value>")
	print("        Set a Gauge metric.")
	print("        Value must be a non-negative 64-bit integer value.")
	print("")
	print("    config")
	print("        Show config.")
	print("")
	print("    version")
	print("        Show program version.")
	print("")
	print("    help")
	print("        Show this help page.")
	print("")
	print("Exit Codes")
	print("    0 ok")
	print("    1 error")
	print("    2 wrong user input")
	print("")
}

func main() {
	log.SetOutput(os.Stdout)
	// flags
	// -url https://monibot.io
	url := os.Getenv(urlEnvKey)
	if url == "" {
		url = defaultUrl
	}
	flag.StringVar(&url, urlFlag, url, "")
	// -apiKey 007
	apiKey := os.Getenv(apiKeyEnvKey)
	if apiKey == "" {
		apiKey = defaultApiKey
	}
	flag.StringVar(&apiKey, apiKeyFlag, apiKey, "")
	// -trials 12
	trialsStr := os.Getenv(trialsEnvKey)
	if trialsStr == "" {
		trialsStr = defaultTrialsStr
	}
	trials, err := strconv.Atoi(trialsStr)
	if err != nil {
		fatal(2, "cannot parse trials %q: %s", trialsStr, err)
	}
	flag.IntVar(&trials, trialsFlag, trials, "")
	// -delay 5s
	delayStr := os.Getenv(delayEnvKey)
	if delayStr == "" {
		delayStr = defaultDelayStr
	}
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		fatal(2, "cannot parse delay %q: %s", delayStr, err)
	}
	flag.DurationVar(&delay, delayFlag, delay, "")
	// -v
	verboseStr := os.Getenv(verboseEnvKey)
	if verboseStr == "" {
		verboseStr = strconv.FormatBool(defaultVerbose)
	}
	verbose := verboseStr == "true"
	flag.BoolVar(&verbose, verboseFlag, verbose, "")
	// parse flags
	flag.Usage = usage
	flag.Parse()
	// execute non-API commands
	command := flag.Arg(0)
	switch command {
	case "", "help":
		usage()
		os.Exit(0)
	case "config":
		print("url      %v", url)
		print("apiKey   %v", apiKey)
		print("trials   %v", trials)
		print("delay    %v", delay)
		print("verbose  %v", verbose)
		os.Exit(0)
	case "version":
		print("moni %s", monibot.Version)
		os.Exit(0)
	}
	// validate flags
	if url == "" {
		fatal(2, "empty url")
	}
	if apiKey == "" {
		fatal(2, "empty apiKey")
	}
	if trials < 0 {
		fatal(2, "invalid trials: %d", trials)
	}
	if delay < 0 {
		fatal(2, "invalid delay: %s", delay)
	}
	// init Sender and Api
	logger := monibot.NewDiscardLogger()
	if verbose {
		logger = monibot.NewLogger(log.Default())
	}
	sender := monibot.NewSenderWithOptions(apiKey, monibot.SenderOptions{
		Logger:     logger,
		MonibotUrl: url,
	})
	retrySender := monibot.NewRetrySenderWithOptions(sender, monibot.RetrySenderOptions{
		Logger: logger,
		Trials: trials,
		Delay:  delay,
	})
	api := monibot.NewApiWithSender(retrySender)
	// execute API commands
	switch command {
	case "ping":
		// moni ping
		err := api.GetPing()
		if err != nil {
			fatal(1, "%s", err)
		}
	case "watchdogs":
		// moni watchdogs
		watchdogs, err := api.GetWatchdogs()
		if err != nil {
			fatal(1, "%s", err)
		}
		printWatchdogs(watchdogs)
	case "watchdog":
		// moni watchdog <watchdogId>
		watchdogId := flag.Arg(1)
		if watchdogId == "" {
			fatal(2, "empty watchdogId")
		}
		watchdog, err := api.GetWatchdog(watchdogId)
		if err != nil {
			fatal(1, "%s", err)
		}
		printWatchdogs([]monibot.Watchdog{watchdog})
	case "reset":
		// moni reset <watchdogId>
		watchdogId := flag.Arg(1)
		if watchdogId == "" {
			fatal(2, "empty watchdogId")
		}
		err := api.PostWatchdogReset(watchdogId)
		if err != nil {
			fatal(1, "%s", err)
		}
	case "machines":
		// moni machines
		machines, err := api.GetMachines()
		if err != nil {
			fatal(1, "%s", err)
		}
		printMachines(machines)
	case "machine":
		// moni machine <machineId>
		machineId := flag.Arg(1)
		if machineId == "" {
			fatal(2, "empty machineId")
		}
		machine, err := api.GetMachine(machineId)
		if err != nil {
			fatal(1, "%s", err)
		}
		printMachines([]monibot.Machine{machine})
	case "sample":
		// moni sample <machineId>
		machineId := flag.Arg(1)
		if machineId == "" {
			fatal(2, "empty machineId")
		}
		sample, err := loadMachineSample()
		if err != nil {
			fatal(1, "cannot loadMachineSample: %s", err)
		}
		err = api.PostMachineSample(machineId, sample)
		if err != nil {
			fatal(1, "%s", err)
		}
	case "metrics":
		// moni metrics
		metrics, err := api.GetMetrics()
		if err != nil {
			fatal(1, "%s", err)
		}
		printMetrics(metrics)
	case "metric":
		// moni metric <metricId>
		metricId := flag.Arg(1)
		if metricId == "" {
			fatal(2, "empty metricId")
		}
		metric, err := api.GetMetric(metricId)
		if err != nil {
			fatal(1, "%s", err)
		}
		printMetrics([]monibot.Metric{metric})
	case "inc":
		// moni inc <metricId> <value>
		metricId := flag.Arg(1)
		if metricId == "" {
			fatal(2, "empty metricId")
		}
		valueStr := flag.Arg(2)
		if valueStr == "" {
			fatal(2, "empty value")
		}
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			fatal(2, "cannot parse value %q: %s", valueStr, err)
		}
		err = api.PostMetricInc(metricId, value)
		if err != nil {
			fatal(1, "%s", err)
		}
	case "set":
		// moni set <metricId> <value>
		metricId := flag.Arg(1)
		if metricId == "" {
			fatal(2, "empty metricId")
		}
		valueStr := flag.Arg(2)
		if valueStr == "" {
			fatal(2, "empty value")
		}
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			fatal(2, "cannot parse value %q: %s", valueStr, err)
		}
		err = api.PostMetricSet(metricId, value)
		if err != nil {
			fatal(1, "%s", err)
		}
	default:
		fatal(2, "unknown command %q, run 'moni help'", command)
	}
}
