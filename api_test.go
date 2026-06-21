package monibot

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestApi(t *testing.T) {
	str := func(a any) string {
		switch x := a.(type) {
		case Watchdog:
			return fmt.Sprintf("Id=%v, Name=%v, IntervalMillis=%v", x.Id, x.Name, x.IntervalMillis)
		case Machine:
			return fmt.Sprintf("Id=%v, Name=%v", x.Id, x.Name)
		case Metric:
			return fmt.Sprintf("Id=%v, Name=%v, Type=%v", x.Id, x.Name, x.Type)
		}
		return "???"
	}
	is := assert.New(t)
	// this test uses a fake HTTP sender
	sender := &fakeSender{}
	// create Api
	api := &Api{sender}
	// GET ping
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{})
		err := api.GetPing()
		is.Nil(err)
		is.Eq(1, len(sender.calls))
		is.Eq("GET ping", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET watchdogs
	{
		sender.calls = nil
		resp := `[
			{"id":"0001", "name":"Cronjob 1", "intervalMillis": 72000000},
			{"id":"0002", "name":"Cronjob 2", "intervalMillis": 36000000}
		]`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		watchdogs, err := api.GetWatchdogs()
		is.Nil(err)
		is.Eq(2, len(watchdogs))
		is.Eq("Id=0001, Name=Cronjob 1, IntervalMillis=72000000", str(watchdogs[0]))
		is.Eq("Id=0002, Name=Cronjob 2, IntervalMillis=36000000", str(watchdogs[1]))
		is.Eq(1, len(sender.calls))
		is.Eq("GET watchdogs", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET watchdog/00000001
	{
		sender.calls = nil
		resp := `{"id":"0001", "name":"Cronjob 1", "intervalMillis": 72000000}`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		watchdog, err := api.GetWatchdog("00000001")
		is.Nil(err)
		is.Eq("Id=0001, Name=Cronjob 1, IntervalMillis=72000000", str(watchdog))
		is.Eq(1, len(sender.calls))
		is.Eq("GET watchdog/00000001", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST watchdog/00000001/heartbeat
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{})
		err := api.PostWatchdogHeartbeat("00000001")
		is.Nil(err)
		is.Eq(1, len(sender.calls))
		is.Eq("POST watchdog/00000001/heartbeat", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET machines
	{
		sender.calls = nil
		resp := `[
			{"id":"01", "name":"Server 1"},
			{"id":"02", "name":"Server 2"}
		]`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		machines, err := api.GetMachines()
		is.Nil(err)
		is.Eq(2, len(machines))
		is.Eq("Id=01, Name=Server 1", str(machines[0]))
		is.Eq("Id=02, Name=Server 2", str(machines[1]))
		is.Eq(1, len(sender.calls))
		is.Eq("GET machines", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET machine/01
	{
		sender.calls = nil
		resp := `{"id":"01", "name":"Server 1"}`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		machine, err := api.GetMachine("01")
		is.Nil(err)
		is.Eq("Id=01, Name=Server 1", str(machine))
		is.Eq(1, len(sender.calls))
		is.Eq("GET machine/01", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST machine/00000001/sample (old style without disks[] and nets[])
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{})
		tstamp := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
		sample := MachineSample{
			Tstamp:      tstamp.UnixMilli(),
			Load1:       1.01,
			Load5:       0.78,
			Load15:      0.12,
			CpuPercent:  12,
			MemPercent:  34,
			DiskPercent: 12,
			DiskRead:    678,
			DiskWrite:   567,
			NetRecv:     13,
			NetSend:     14,
		}
		err := api.PostMachineSample("00000001", sample)
		is.Nil(err)
		is.Eq(1, len(sender.calls))
		is.Eq("POST machine/00000001/sample tstamp=1698400800000"+
			"&load1=1.010&load5=0.780&load15=0.120"+
			"&cpu=12"+
			"&mem=34"+
			"&disk=12"+
			"&diskRead=678"+
			"&diskWrite=567"+
			"&netRecv=13"+
			"&netSend=14", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST machine/00000001/sample (new style with disks[] and nets[])
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{})
		tstamp := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
		sample := MachineSample{
			Tstamp:     tstamp.UnixMilli(),
			Load1:      1.01,
			Load5:      0.78,
			Load15:     0.12,
			CpuPercent: 12,
			MemPercent: 34,
			Disks: []DiskSample{
				{
					Device:      "/dev/sda",
					Total:       21,
					Used:        22,
					UsedPercent: 23,
					ReadBytes:   24,
					WriteBytes:  25,
				},
				{
					Device:      "/dev/sdb", // string // e.g. "/dev/sda1" // from disk.Partitions()
					Total:       31,         // int64  // 0..MAX_I64       // from disk.Usage()
					Used:        32,         // int64  // 0..MAX_I64       // from disk.Usage()
					UsedPercent: 33,         // int    // 0..100           // from disk.Usage()
					ReadBytes:   34,         // int64  // 0..MAX_I64       // from disk.IOCounters()
					WriteBytes:  35,         // int64  // 0..MAX_I64       // from disk.IOCounters()
				},
			},
			DiskPercent: 12,
			DiskRead:    678,
			DiskWrite:   567,
			Nets: []NetSample{
				{
					Device:    "eth0",
					RecvBytes: 24,
					SendBytes: 25,
				},
				{
					Device:    "eth1",
					RecvBytes: 34,
					SendBytes: 35,
				},
			},
			NetRecv: 24 + 34, // 58
			NetSend: 25 + 35, // 60
		}
		err := api.PostMachineSample("00000001", sample)
		is.Nil(err)
		is.Eq(1, len(sender.calls))
		is.Eq("POST machine/00000001/sample tstamp=1698400800000"+
			"&load1=1.010&load5=0.780&load15=0.120"+
			"&cpu=12"+
			"&mem=34"+
			"&disks=2"+
			"&disks[0].device=/dev/sda"+
			"&disks[0].total=21"+
			"&disks[0].used=22"+
			"&disks[0].usedPercent=23"+
			"&disks[0].readBytes=24"+
			"&disks[0].writeBytes=25"+
			"&disks[1].device=/dev/sdb"+
			"&disks[1].total=31"+
			"&disks[1].used=32"+
			"&disks[1].usedPercent=33"+
			"&disks[1].readBytes=34"+
			"&disks[1].writeBytes=35"+
			"&disk=12"+
			"&diskRead=678"+
			"&diskWrite=567"+
			"&nets=2"+
			"&nets[0].device=eth0"+
			"&nets[0].recvBytes=24"+
			"&nets[0].sendBytes=25"+
			"&nets[1].device=eth1"+
			"&nets[1].recvBytes=34"+
			"&nets[1].sendBytes=35"+
			"&netRecv=58"+
			"&netSend=60", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST machine/00000001/text
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{})
		text := "line1\nline2\n\n"
		err := api.PostMachineText("00000001", text)
		is.Nil(err)
		is.Eq(1, len(sender.calls))
		is.Eq("POST machine/00000001/text text=line1%0Aline2%0A%0A", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET metrics
	{
		sender.calls = nil
		resp := `[
			{"id":"01", "name":"Metric 1", "type": 0},
			{"id":"02", "name":"Metric 2", "type": 1}
		]`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		metrics, err := api.GetMetrics()
		is.Nil(err)
		is.Eq(2, len(metrics))
		is.Eq("Id=01, Name=Metric 1, Type=0", str(metrics[0]))
		is.Eq("Id=02, Name=Metric 2, Type=1", str(metrics[1]))
		is.Eq(1, len(sender.calls))
		is.Eq("GET metrics", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// GET metric/01
	{
		sender.calls = nil
		resp := `{"id":"01", "name":"Metric 1", "type": 0}`
		sender.responses = append(sender.responses, fakeResponse{data: []byte(resp)})
		metric, err := api.GetMetric("01")
		is.Nil(err)
		is.Eq("Id=01, Name=Metric 1, Type=0", str(metric))
		is.Eq(1, len(sender.calls))
		is.Eq("GET metric/01", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST metric/00000001/inc
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricInc("00000001", 42)
		is.Eq("connect timeout", err.Error())
		is.Eq(1, len(sender.calls))
		is.Eq("POST metric/00000001/inc value=42", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST metric/00000001/set
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricSet("00000001", 113)
		is.Eq("connect timeout", err.Error())
		is.Eq(1, len(sender.calls))
		is.Eq("POST metric/00000001/set value=113", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
	// POST metric/00000042/values
	{
		sender.calls = nil
		sender.responses = append(sender.responses, fakeResponse{nil, fmt.Errorf("connect timeout")})
		err := api.PostMetricValues("010101", []int64{3, 5, 2, 5, 0, 3, 4, 3, 1})
		is.Eq("connect timeout", err.Error())
		is.Eq(1, len(sender.calls))
		// "0%2C1%2C2%2C3%3A3%2C4%2C5%3A2" = urlEncode("0,1,2,3:3,4,5:2")
		is.Eq("POST metric/010101/values values=0%2C1%2C2%2C3%3A3%2C4%2C5%3A2", sender.calls[0])
		is.Eq(0, len(sender.responses))
	}
}

type fakeSender struct {
	calls     []string
	responses []fakeResponse
}

func (f *fakeSender) Send(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	call := strings.TrimSpace(method + " " + path + " " + string(body))
	f.calls = append(f.calls, call)
	if len(f.responses) == 0 {
		return nil, fmt.Errorf("no response for %s %s", method, path)
	}
	re := f.responses[0]
	f.responses = f.responses[1:]
	return re.data, re.err
}

type fakeResponse struct {
	data []byte
	err  error
}
