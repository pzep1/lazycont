package containercli

import (
	"reflect"
	"testing"
)

func TestStatSummaryLinesFormatsAppleStatsShape(t *testing.T) {
	stat := Stat{
		"id":               "web",
		"memoryUsageBytes": float64(47431680),
		"memoryLimitBytes": float64(1073741824),
		"cpuUsageUsec":     float64(1234567),
		"networkRxBytes":   float64(1289011),
		"networkTxBytes":   float64(876544),
		"blockReadBytes":   float64(4718592),
		"blockWriteBytes":  float64(2202009),
		"numProcesses":     float64(3),
	}

	got := stat.SummaryLines()
	want := []string{
		"  CPU time: 1.2s",
		"  Memory:   45.2 MB / 1.0 GB  [#---------------] 4.4%",
		"  Network:  1.2 MB rx / 856.0 KB tx",
		"  Block IO: 4.5 MB read / 2.1 MB write",
		"  PIDs:     3",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("summary mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestStatSummaryLinesShowsCPUPercentWhenAvailable(t *testing.T) {
	stat := Stat{
		"id":             "web",
		"cpuPercent":     float64(125.12),
		"numProcesses":   float64(12),
		"networkRxBytes": float64(0),
	}

	got := stat.SummaryLines()
	want := []string{
		"  CPU:      125.1%  [################]",
		"  Network:  - rx / - tx",
		"  PIDs:     12",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("summary mismatch\nwant: %#v\n got: %#v", want, got)
	}
}
