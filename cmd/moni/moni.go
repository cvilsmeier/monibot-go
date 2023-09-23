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
		log.Fatalf("invalid interval %q: %s", intervalStr, err)
	}
	// execute command
	command := flag.Arg(0)
	switch command {
	case "", "help":
		usage()
		os.Exit(0)
	case "version":
		prt("moni %s", version)
		os.Exit(0)
	}
	// the following commands need http
	http := NewHttp(url, apiKey, version, verbose)
	switch command {
	case "ping":
		// ping
		if err := runPing(http); err != nil {
			log.Fatal(err)
		}
	case "watchdog":
		// watchdog <watchdogId>
		watchdogId := flag.Arg(1)
		if err := runWatchdog(http, watchdogId); err != nil {
			log.Fatal(err)
		}
	case "reset":
		// reset <watchdogId>
		watchdogId := flag.Arg(1)
		if err := runWatchdogReset(http, watchdogId); err != nil {
			log.Fatal(err)
		}
	case "machine":
		// machine <machineId>
		machineId := flag.Arg(1)
		if err := runMachine(http, machineId); err != nil {
			log.Fatal(err)
		}
	case "sample":
		// sample <machineId>
		machineId := flag.Arg(1)
		if err := runMachineSample(http, machineId, interval, verbose); err != nil {
			log.Fatal(err)
		}
	case "metric":
		// metric <metricId>
		metricId := flag.Arg(1)
		if err := runMetric(http, metricId); err != nil {
			log.Fatal(err)
		}
	case "inc":
		// inc <metricId> <value>
		metricId := flag.Arg(1)
		valueStr := flag.Arg(2)
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse value %q: %s", valueStr, err)
		}
		if err := runMetricInc(http, metricId, value); err != nil {
			log.Fatal(err)
		}
	case "set":
		// set <metricId> <value>
		metricId := flag.Arg(1)
		valueStr := flag.Arg(2)
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse value %q: %s", valueStr, err)
		}
		if err := runMetricSet(http, metricId, value); err != nil {
			log.Fatal(err)
		}
	default:
		prt("unknown command %q", command)
		prt("run 'moni help' to get a list of known commands")
		os.Exit(2)
	}
}

func usage() {
	prt("moni %s", version)
	prt("")
	prt("Monibot command line tool, see https://monibot.io.")
	prt("")
	prt("Usage")
	prt("")
	prt("    moni [flags] command")
	prt("")
	prt("Flags")
	prt("")
	prt("    -url")
	prt("        Monibot URL, default is %q.", defaultUrl)
	prt("")
	prt("    -apiKey")
	prt("        Monibot API Key, default is %q.", defaultApiKey)
	prt("        You can set this also via environment variable %s.", apiKeyEnvKey)
	prt("        You can find your API Key in your profile on https://monibot.io.")
	prt("")
	prt("    -interval")
	prt("        Machine sampling interval, default is %q.", defaultIntervalStr)
	prt("        This is used for 'sample' command. The minimum allowed value is 5m.")
	prt("")
	prt("    -v")
	prt("        Verbose output, default is %t.", defaultVerbose)
	prt("        You can set this also via environment variable %s.", verboseEnvKey)
	prt("")
	prt("Commands")
	prt("")
	prt("    ping")
	prt("        Ping the Monibot API.")
	prt("")
	prt("    watchdog <watchdogId>")
	prt("        Get and print watchdog info.")
	prt("")
	prt("    reset <watchdogId>")
	prt("        Reset a watchdog.")
	prt("")
	prt("    machine <machineId>")
	prt("        Get and print machine info.")
	prt("")
	prt("    sample <machineId>")
	prt("        Send resource usage (cpu/mem/disk) samples for machine.")
	prt("        This command will stay in background. It monitors resource usage")
	prt("        and sends it to monibot every 5 minutes. It works only")
	prt("        on linux.")
	prt("")
	prt("    metric <metricId>")
	prt("        Get and print metric info.")
	prt("")
	prt("    inc <metricId> <value>")
	prt("        Increment a Counter metric. Value must be a non-negative 64-bit")
	prt("        integer value.")
	prt("")
	prt("    set <metricId> <value>")
	prt("        Set a Gauge metric. Value must be a non-negative 64-bit integer")
	prt("        value.")
	prt("")
	prt("    version")
	prt("        Show program version.")
	prt("")
	prt("    help")
	prt("        Show this help page.")
	prt("")
}

// prt prints a line to stdout
func prt(f string, a ...any) {
	fmt.Printf(f+"\n", a...)
}

// debug prints a line to stdout if verbose is true.
func debug(verbose bool, f string, a ...any) {
	if verbose {
		fmt.Printf("DEBUG: "+f+"\n", a...)
	}
}

// commands

func runPing(http *Http) error {
	_, err := http.Get("ping")
	return err
}

func runWatchdog(http *Http, watchdogId string) error {
	data, err := http.Get("watchdog/" + watchdogId)
	if err != nil {
		return err
	}
	prt(string(data))
	return err
}

func runWatchdogReset(http *Http, watchdogId string) error {
	_, err := http.Post("watchdog/"+watchdogId+"/reset", nil)
	return err
}

func runMachine(http *Http, machineId string) error {
	data, err := http.Get("machine/" + machineId)
	if err != nil {
		return err
	}
	prt(string(data))
	return err
}

func runMachineSample(http *Http, machineId string, interval time.Duration, verbose bool) error {
	_, err := http.Get("machine/" + machineId)
	if err != nil {
		return err
	}
	lastCpuStat, err := loadCpuStat(verbose)
	if err != nil {
		return fmt.Errorf("cannot loadCpuStat: %s", err)
	}
	for {
		debug(verbose, "sleep %v", interval)
		time.Sleep(interval)
		// stat cpu
		cpuStat, err := loadCpuStat(verbose)
		if err != nil {
			log.Printf("ERROR: cannot loadCpuStat: %s", err)
			continue
		}
		// stat mem
		memStat, err := loadMemStat(verbose)
		if err != nil {
			log.Printf("ERROR: cannot loadMemStat: %s", err)
			continue
		}
		// stat disk
		diskStat, err := loadDiskStat(verbose)
		if err != nil {
			log.Printf("ERROR: cannot loadDiskStat: %s", err)
			continue
		}
		// POST machine sample
		diffCpuStat := cpuStat.Minus(lastCpuStat)
		lastCpuStat = cpuStat
		body := fmt.Sprintf(
			"tstamp=%d&cpu=%d&mem=%d&disk=%d",
			time.Now().UnixMilli(),
			diffCpuStat.Percent(),
			memStat.Percent(),
			diskStat.Percent(),
		)
		_, err = http.Post("machine/"+machineId+"/sample", []byte(body))
		if err != nil {
			log.Printf("ERROR cannot POST machineSampleApi: %s", err)
		}
	}
}

func runMetric(http *Http, metricId string) error {
	data, err := http.Get("metric/" + metricId)
	if err != nil {
		return err
	}
	prt(string(data))
	return err
}

func runMetricInc(http *Http, metricId string, value int64) error {
	if value == 0 {
		return nil
	}
	if value < 0 {
		return fmt.Errorf("cannot inc negative value %d", value)
	}
	body := fmt.Sprintf("value=%d", value)
	_, err := http.Post("metric/"+metricId+"/inc", []byte(body))
	return err
}

func runMetricSet(http *Http, metricId string, value int64) error {
	if value < 0 {
		return fmt.Errorf("cannot set negative value %d", value)
	}
	body := fmt.Sprintf("value=%d", value)
	_, err := http.Post("metric/"+metricId+"/set", []byte(body))
	return err
}

// machine sampling

func loadCpuStat(verbose bool) (SampleStat, error) {
	debug(verbose, "loadCpuStat: read /proc/stat")
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
			debug(verbose, "loadCpuStat: parse %q", line)
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
			debug(verbose, "loadCpuStat: total=%d, used=%d", total, used)
			return SampleStat{total, used}, nil
		}
	}
	return SampleStat{}, fmt.Errorf("prefix 'cpu ' not found in /proc/stat")
}

func loadMemStat(verbose bool) (SampleStat, error) {
	debug(verbose, "loadMemStat: /usr/bin/free -k")
	text, err := execCommand("/usr/bin/free", "-k")
	if err != nil {
		return SampleStat{}, err
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "Mem: ")
		if found {
			debug(verbose, "loadMemStat: parse %q", line)
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
			debug(verbose, "loadMemStat: total=%d, used=%d", total, used)
			return SampleStat{total, used}, nil
		}
	}
	return SampleStat{}, fmt.Errorf("prefix 'Mem: ' not found in output of /usr/bin/free")
}

func loadDiskStat(verbose bool) (SampleStat, error) {
	debug(verbose, "loadDiskStat: /usr/bin/df --exclude-type=tmpfs --total --output=source,size,used")
	text, _ := execCommand("/usr/bin/df", "--exclude-type=tmpfs", "--total", "--output=source,size,used")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = trimText(line)
		after, found := strings.CutPrefix(line, "total ")
		if found {
			debug(verbose, "loadDiskStat: parse %q", line)
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
			debug(verbose, "loadDiskStat: size=%d, used=%d", size, used)
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
