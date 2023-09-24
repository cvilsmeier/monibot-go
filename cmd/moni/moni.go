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

	"github.com/cvilsmeier/moni-cli/api"
)

const (
	version = "v0.0.1"

	urlEnvKey  = "MONIBOT_URL"
	urlFlag    = "url"
	defaultUrl = "https://monibot.io"

	apiKeyEnvKey  = "MONIBOT_API_KEY"
	apiKeyFlag    = "apiKey"
	defaultApiKey = ""

	intervalEnvKey     = "MONIBOT_INTERVAL"
	intervalFlag       = "interval"
	defaultIntervalStr = "5m"

	verboseEnvKey  = "MONIBOT_VERBOSE"
	verboseFlag    = "v"
	defaultVerbose = false
)

func main() {
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
	// -interval 1m
	intervalStr := os.Getenv(intervalEnvKey)
	if intervalStr == "" {
		intervalStr = defaultIntervalStr
	}
	flag.StringVar(&intervalStr, intervalFlag, intervalStr, "")
	// -v
	verboseStr := os.Getenv(verboseEnvKey)
	verbose := verboseStr == "true"
	flag.BoolVar(&verbose, verboseFlag, verbose, "")
	// parse flags
	flag.Usage = usage
	flag.Parse()
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		fatal(2, "invalid interval %q: %s", intervalStr, err)
	}
	// execute command
	command := flag.Arg(0)
	switch command {
	case "", "help":
		usage()
		os.Exit(0)
	case "version":
		print("moni %s", version)
		os.Exit(0)
	}
	// init the API
	logger := api.NewLogger(verbose)
	http := api.NewHttp(logger, url, "moni/"+version, apiKey)
	conn := api.NewConn(http)
	switch command {
	case "ping":
		// ping
		err := conn.GetPing()
		if err != nil {
			fatal(1, "%s", err)
		}
	case "watchdog":
		// watchdog <watchdogId>
		watchdogId := flag.Arg(1)
		if watchdogId == "" {
			fatal(2, "empty watchdogId")
		}
		data, err := conn.GetWatchdog(watchdogId)
		if err != nil {
			fatal(1, "%s", err)
		}
		print("%s", string(data))
	case "reset":
		// reset <watchdogId>
		watchdogId := flag.Arg(1)
		if watchdogId == "" {
			fatal(2, "empty watchdogId")
		}
		err := conn.PostWatchdogReset(watchdogId)
		if err != nil {
			fatal(1, "%s", err)
		}
	case "machine":
		// machine <machineId>
		machineId := flag.Arg(1)
		if machineId == "" {
			fatal(2, "empty machineId")
		}
		data, err := conn.GetMachine(machineId)
		if err != nil {
			fatal(1, "%s", err)
		}
		print("%s", string(data))
	case "sample":
		// sample <machineId>
		machineId := flag.Arg(1)
		if machineId == "" {
			fatal(2, "empty machineId")
		}
		if err := sampleMachine(logger, conn, machineId, interval); err != nil {
			fatal(1, "%s", err)
		}
	case "metric":
		// metric <metricId>
		metricId := flag.Arg(1)
		if metricId == "" {
			fatal(2, "empty metricId")
		}
		data, err := conn.GetMetric(metricId)
		if err != nil {
			fatal(1, "%s", err)
		}
		print("%s", string(data))
	case "inc":
		// inc <metricId> <value>
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
		err = conn.PostMetricInc(metricId, value)
		if err != nil {
			fatal(1, "%s", err)
		}
	case "set":
		// inc <metricId> <value>
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
		err = conn.PostMetricSet(metricId, value)
		if err != nil {
			fatal(1, "%s", err)
		}
	default:
		fatal(2, "unknown command %q, run 'moni help'", command)
	}
}

func usage() {
	print("moni %s", version)
	print("")
	print("Monibot command line tool, see https://monibot.io.")
	print("")
	print("Usage")
	print("")
	print("    moni [flags] command")
	print("")
	print("Flags")
	print("")
	print("    -url")
	print("        Monibot URL, default is %q.", defaultUrl)
	print("")
	print("    -apiKey")
	print("        Monibot API Key, default is %q.", defaultApiKey)
	print("        You can set this also via environment variable %s.", apiKeyEnvKey)
	print("        You can find your API Key in your profile on https://monibot.io.")
	print("")
	print("    -interval")
	print("        Machine sampling interval, default is %q.", defaultIntervalStr)
	print("        This is used for 'sample' command. The minimum allowed value is 5m.")
	print("")
	print("    -v")
	print("        Verbose output, default is %t.", defaultVerbose)
	print("        You can set this also via environment variable %s.", verboseEnvKey)
	print("")
	print("Commands")
	print("")
	print("    ping")
	print("        Ping the Monibot API.")
	print("")
	print("    watchdog <watchdogId>")
	print("        Get and print watchdog info.")
	print("")
	print("    reset <watchdogId>")
	print("        Reset a watchdog.")
	print("")
	print("    machine <machineId>")
	print("        Get and print machine info.")
	print("")
	print("    sample <machineId>")
	print("        Send resource usage (cpu/mem/disk) samples for machine.")
	print("        This command will stay in background. It monitors resource usage")
	print("        and sends it to monibot every 5 minutes. It works only")
	print("        on linux.")
	print("")
	print("    metric <metricId>")
	print("        Get and print metric info.")
	print("")
	print("    inc <metricId> <value>")
	print("        Increment a Counter metric. Value must be a non-negative 64-bit")
	print("        integer value.")
	print("")
	print("    set <metricId> <value>")
	print("        Set a Gauge metric. Value must be a non-negative 64-bit integer")
	print("        value.")
	print("")
	print("    version")
	print("        Show program version.")
	print("")
	print("    help")
	print("        Show this help page.")
	print("")
	print("Exit Codes")
	print("    0   ok")
	print("    1   error")
	print("    2   wrong user input")
	print("")
	print("")
}

// print prints a line to stdout.
func print(f string, a ...any) {
	fmt.Printf(f+"\n", a...)
}

// fatal prints a message to stdout and exits with exitCode.
func fatal(exitCode int, f string, a ...any) {
	fmt.Printf(f+"\n", a...)
	os.Exit(exitCode)
}

// sampleMachine samples the local machine (cpu/mem/disk) endlessly.
func sampleMachine(logger api.Logger, conn *api.Conn, machineId string, interval time.Duration) error {
	_, err := conn.GetMachine(machineId)
	if err != nil {
		return err
	}
	lastCpuStat, err := loadCpuStat(logger)
	if err != nil {
		return fmt.Errorf("cannot loadCpuStat: %s", err)
	}
	for {
		logger.Debugf("sleep %v", interval)
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
		diffCpuStat := cpuStat.Minus(lastCpuStat)
		lastCpuStat = cpuStat
		err = conn.PostMachineSample(
			machineId,
			time.Now().UnixMilli(),
			diffCpuStat.Percent(),
			memStat.Percent(),
			diskStat.Percent(),
		)
		if err != nil {
			log.Printf("ERROR cannot PostMachineSample: %s", err)
		}
	}
}

// loadCpuStat loads cpu stat.
func loadCpuStat(logger api.Logger) (SampleStat, error) {
	logger.Debugf("loadCpuStat: read /proc/stat")
	f, err := os.Open("/proc/stat")
	if err != nil {
		return SampleStat{}, err
	}
	defer f.Close()
	sca := bufio.NewScanner(f)
	for sca.Scan() {
		line := trimText(sca.Text())
		after, found := strings.CutPrefix(line, "cpu ")
		if found {
			logger.Debugf("loadCpuStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 5 {
				return SampleStat{}, fmt.Errorf("want min 5 tokens in %q but was %d", line, len(toks))
			}
			nums := make([]int64, len(toks))
			var total int64
			for i := range toks {
				n, err := strconv.ParseInt(toks[i], 10, 64)
				if err != nil {
					return SampleStat{}, fmt.Errorf("cannot parse toks[%d] %q from line %q: %s", i, toks[i], line, err)
				}
				nums[i] = n
				total += n
			}
			idle := nums[3]
			used := total - idle
			logger.Debugf("loadCpuStat: total=%d, used=%d", total, used)
			return SampleStat{total, used}, nil
		}
	}
	return SampleStat{}, fmt.Errorf("prefix 'cpu ' not found in /proc/stat")
}

func loadMemStat(logger api.Logger) (SampleStat, error) {
	logger.Debugf("loadMemStat: /usr/bin/free -k")
	text, err := execCommand("/usr/bin/free", "-k")
	if err != nil {
		return SampleStat{}, err
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "Mem: ")
		if found {
			logger.Debugf("loadMemStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 3 {
				return SampleStat{}, fmt.Errorf("want min 3 tokens int %q but was %d", line, len(toks))
			}
			total, err := strconv.ParseInt(toks[0], 10, 64)
			if err != nil {
				return SampleStat{}, fmt.Errorf("cannot parse toks[0] from %q: %s", line, err)
			}
			used, err := strconv.ParseInt(toks[1], 10, 64)
			if err != nil {
				return SampleStat{}, fmt.Errorf("cannot parse toks[1] from %q: %s", line, err)
			}
			logger.Debugf("loadMemStat: total=%d, used=%d", total, used)
			return SampleStat{total, used}, nil
		}
	}
	return SampleStat{}, fmt.Errorf("prefix 'Mem: ' not found in output of /usr/bin/free")
}

func loadDiskStat(logger api.Logger) (SampleStat, error) {
	logger.Debugf("loadDiskStat: /usr/bin/df --exclude-type=tmpfs --total --output=source,size,used")
	text, _ := execCommand("/usr/bin/df", "--exclude-type=tmpfs", "--total", "--output=source,size,used")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "total ")
		if found {
			logger.Debugf("loadDiskStat: parse %q", line)
			toks := strings.Split(after, " ")
			if len(toks) < 2 {
				return SampleStat{}, fmt.Errorf("want 2 toks in %q but has only %d", line, len(toks))
			}
			size, err := strconv.ParseInt(toks[0], 10, 64)
			if err != nil {
				return SampleStat{}, fmt.Errorf("parse toks[0] %q from %q: %w", toks[0], line, err)
			}
			used, err := strconv.ParseInt(toks[1], 10, 64)
			if err != nil {
				return SampleStat{}, fmt.Errorf("parse toks[1] %q from %q: %w", toks[1], line, err)
			}
			logger.Debugf("loadDiskStat: size=%d, used=%d", size, used)
			return SampleStat{size, used}, nil
		}
	}
	return SampleStat{}, fmt.Errorf("'total' line not found in df output")
}

func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.WaitDelay = 10 * time.Second
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("cannot run %s: %w", name, err)
	}
	return string(out), err
}

func trimText(s string) string {
	for strings.Contains(s, "\t") {
		s = strings.ReplaceAll(s, "\t", " ")
	}
	for strings.Contains(s, "\r") {
		s = strings.ReplaceAll(s, "\r", "")
	}
	for strings.Contains(s, "\n") {
		s = strings.ReplaceAll(s, "\n", "")
	}
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

type SampleStat struct {
	Total int64
	Used  int64
}

func (s SampleStat) Minus(o SampleStat) SampleStat {
	return SampleStat{
		s.Total - o.Total,
		s.Used - o.Used,
	}
}

func (s SampleStat) Percent() int {
	if s.Total == 0 {
		return 0
	}
	p := (s.Used * 100) / s.Total
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}
	return int(p)
}
