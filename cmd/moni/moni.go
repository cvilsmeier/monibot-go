/*
Moni is a command line tool for interacting with the Monibot REST API, see https://monibot.io for details.
It supports a number of commands. To get a list of supported commands, run

	$ moni help
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	verboseEnvKey  = "MONIBOT_VERBOSE"
	verboseFlag    = "v"
	defaultVerbose = false
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
	print("    sample <machineId> [interval]")
	print("        Send resource usage (cpu/mem/disk) samples for machine.")
	print("        This command will stay in background. It monitors resource usage")
	print("        and sends it to monibot periodically, specified in interval.")
	print("        The default interval is 5m. The interval may be lower, but")
	print("        serve-side rate limits may apply.")
	print("        The sample command currently works only on linux.")
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
	// -apiKey 0000000000
	apiKey := os.Getenv(apiKeyEnvKey)
	if apiKey == "" {
		apiKey = defaultApiKey
	}
	flag.StringVar(&apiKey, apiKeyFlag, apiKey, "")
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
	// init Sender and Api
	logger := monibot.NewDiscardLogger()
	if verbose {
		logger = monibot.NewLogger(log.Default())
	}
	sender := monibot.NewSenderWithOptions(apiKey, monibot.SenderOptions{Logger: logger, MonibotUrl: url})
	retrySender := monibot.NewRetrySenderWithOptions(sender, monibot.RetrySenderOptions{Logger: logger})
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
		err := retry(func() error {
			return api.PostWatchdogReset(watchdogId)
		})
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
		// moni sample <machineId> [interval]
		machineId := flag.Arg(1)
		if machineId == "" {
			fatal(2, "empty machineId")
		}
		intervalStr := flag.Arg(2)
		if intervalStr == "" {
			intervalStr = "5m"
		}
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			fatal(2, "cannot parse interval %q: %s", intervalStr, err)
		}
		if interval < 5*time.Second {
			fatal(2, "interval must be >= 5s but was %s", interval)
		}
		if err := sampleMachine(logger, api, machineId, interval); err != nil {
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
		err = retry(func() error {
			return api.PostMetricInc(metricId, value)
		})
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
		err = retry(func() error {
			return api.PostMetricSet(metricId, value)
		})
		if err != nil {
			fatal(1, "%s", err)
		}
	default:
		fatal(2, "unknown command %q, run 'moni help'", command)
	}
}

// print prints a line to stdout.
func print(f string, a ...any) {
	fmt.Printf(f+"\n", a...)
}

// printWatchdogs prints watchdogs.
func printWatchdogs(watchdogs []monibot.Watchdog) {
	print("%-35s | %-25s | %s", "Id", "Name", "IntervalMillis")
	for _, watchdog := range watchdogs {
		print("%-35s | %-25s | %d", watchdog.Id, watchdog.Name, watchdog.IntervalMillis)
	}
}

// printMachines prints machines.
func printMachines(machines []monibot.Machine) {
	print("%-35s | %s", "Id", "Name")
	for _, machine := range machines {
		print("%-35s | %s", machine.Id, machine.Name)
	}
}

// printMetrics prints metrics.
func printMetrics(metrics []monibot.Metric) {
	print("%-35s | %-25s | %s", "Id", "Name", "Type")
	for _, metric := range metrics {
		print("%-35s | %-25s | %d", metric.Id, metric.Name, metric.Type)
	}
}

// fatal prints a message to stdout and exits with exitCode.
func fatal(exitCode int, f string, a ...any) {
	fmt.Printf(f+"\n", a...)
	os.Exit(exitCode)
}

func retry(f func() error) error {
	var err error
	for i := 0; i < 3; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i) * 10 * time.Second)
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}

// sampleMachine samples the local machine (cpu/mem/disk) in an endless loop.
func sampleMachine(logger monibot.Logger, api *monibot.Api, machineId string, interval time.Duration) error {
	_, err := api.GetMachine(machineId)
	if err != nil {
		return err
	}
	lastCpuStat, err := loadCpuStat(logger)
	if err != nil {
		return fmt.Errorf("cannot loadCpuStat: %s", err)
	}
	for {
		logger.Debug("sleep %v", interval)
		time.Sleep(interval)
		// stat cpu
		cpuStat, err := loadCpuStat(logger)
		if err != nil {
			log.Printf("ERROR: cannot loadCpuStat: %s", err)
			continue
		}
		// stat mem
		memStat, err := loadMemStat(logger)
		if err != nil {
			log.Printf("ERROR: cannot loadMemStat: %s", err)
			continue
		}
		// stat disk
		diskStat, err := loadDiskStat(logger)
		if err != nil {
			log.Printf("ERROR: cannot loadDiskStat: %s", err)
			continue
		}
		// POST machine sample
		diffCpuStat := cpuStat.minus(lastCpuStat)
		lastCpuStat = cpuStat
		err = retry(func() error {
			return api.PostMachineSample(
				machineId,
				time.Now().UnixMilli(),
				diffCpuStat.percent(),
				memStat.percent(),
				diskStat.percent(),
			)
		})
		if err != nil {
			log.Printf("ERROR cannot PostMachineSample: %s", err)
		}
	}
}

// loadCpuStat loads cpu usage stat from /proc/stat.
func loadCpuStat(logger monibot.Logger) (sampleStat, error) {
	logger.Debug("loadCpuStat: read /proc/stat")
	f, err := os.Open("/proc/stat")
	if err != nil {
		return sampleStat{}, err
	}
	defer f.Close()
	sca := bufio.NewScanner(f)
	for sca.Scan() {
		line := trimText(sca.Text())
		after, found := strings.CutPrefix(line, "cpu ")
		if found {
			logger.Debug("loadCpuStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 5 {
				return sampleStat{}, fmt.Errorf("want min 5 tokens in %q but was %d", line, len(toks))
			}
			nums := make([]int64, len(toks))
			var total int64
			for i := range toks {
				n, err := strconv.ParseInt(toks[i], 10, 64)
				if err != nil {
					return sampleStat{}, fmt.Errorf("cannot parse toks[%d] %q from line %q: %s", i, toks[i], line, err)
				}
				nums[i] = n
				total += n
			}
			idle := nums[3]
			used := total - idle
			logger.Debug("loadCpuStat: total=%d, used=%d", total, used)
			return sampleStat{total, used}, nil
		}
	}
	return sampleStat{}, fmt.Errorf("prefix 'cpu ' not found in /proc/stat")
}

// loadMemStat uses /usr/bin/free to load mem usage stat.
func loadMemStat(logger monibot.Logger) (sampleStat, error) {
	logger.Debug("loadMemStat: /usr/bin/free -k")
	text, err := execCommand("/usr/bin/free", "-k")
	if err != nil {
		return sampleStat{}, err
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "Mem: ")
		if found {
			logger.Debug("loadMemStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 3 {
				return sampleStat{}, fmt.Errorf("want min 3 tokens int %q but was %d", line, len(toks))
			}
			total, err := strconv.ParseInt(toks[0], 10, 64)
			if err != nil {
				return sampleStat{}, fmt.Errorf("cannot parse toks[0] from %q: %s", line, err)
			}
			used, err := strconv.ParseInt(toks[1], 10, 64)
			if err != nil {
				return sampleStat{}, fmt.Errorf("cannot parse toks[1] from %q: %s", line, err)
			}
			logger.Debug("loadMemStat: total=%d, used=%d", total, used)
			return sampleStat{total, used}, nil
		}
	}
	return sampleStat{}, fmt.Errorf("prefix 'Mem: ' not found in output of /usr/bin/free")
}

// loadMemStat uses /usr/bin/df to load disk usage stat.
func loadDiskStat(logger monibot.Logger) (sampleStat, error) {
	logger.Debug("loadDiskStat: /usr/bin/df --exclude-type=tmpfs --total --output=source,size,used")
	text, _ := execCommand("/usr/bin/df", "--exclude-type=tmpfs", "--total", "--output=source,size,used")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "total ")
		if found {
			logger.Debug("loadDiskStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 2 {
				return sampleStat{}, fmt.Errorf("want 2 toks in %q but has only %d", line, len(toks))
			}
			size, err := strconv.ParseInt(toks[0], 10, 64)
			if err != nil {
				return sampleStat{}, fmt.Errorf("parse toks[0] %q from %q: %w", toks[0], line, err)
			}
			used, err := strconv.ParseInt(toks[1], 10, 64)
			if err != nil {
				return sampleStat{}, fmt.Errorf("parse toks[1] %q from %q: %w", toks[1], line, err)
			}
			logger.Debug("loadDiskStat: size=%d, used=%d", size, used)
			return sampleStat{size, used}, nil
		}
	}
	return sampleStat{}, fmt.Errorf("'total' line not found in df output")
}

// execCommand executes an external binary.
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.WaitDelay = 10 * time.Second
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("cannot run %s: %w", name, err)
	}
	return string(out), err
}

// trimText trims and normalizes a line of text.
func trimText(s string) string {
	s = replace(s, "\t", " ")
	s = replace(s, "\r", "")
	s = replace(s, "\n", "")
	s = replace(s, "  ", " ")
	return strings.TrimSpace(s)
}

func replace(str, old, new string) string {
	for strings.Contains(str, old) {
		str = strings.ReplaceAll(str, old, new)
	}
	return str
}

// A sampleStat holds cpu/mem/disk usage data.
type sampleStat struct {
	total int64
	used  int64
}

// minus calculates s minus o.
func (s sampleStat) minus(o sampleStat) sampleStat {
	return sampleStat{
		s.total - o.total,
		s.used - o.used,
	}
}

// percent calculates usage percent and returns a number between 0 and 100 (inclusive).
func (s sampleStat) percent() int {
	if s.total == 0 {
		return 0
	}
	p := (s.used * 100) / s.total
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}
	return int(p)
}
