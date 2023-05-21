package psi

import "fmt"

type Resource string

const (
	CPU    Resource = "cpu"
	Memory Resource = "memory"
	IO     Resource = "io"
)

type Metric string

const (
	Avg10  Metric = "avg10"
	Avg60  Metric = "avg60"
	Avg300 Metric = "avg300"
)

const lineFormat = "avg10=%f avg60=%f avg300=%f total=%d"

// PSILine is a single line of values as returned by `/proc/pressure/*`.
//
// The Avg entries are averages over n seconds, as a percentage.
// The Total line is in microseconds.
type PSILine struct {
	Avg10  float64
	Avg60  float64
	Avg300 float64
	Total  uint64 // in microseconds, very accurate starvation metric
}

func (l PSILine) GetMetric(metric Metric) float64 {
	switch metric {
	case Avg10:
		return l.Avg10
	case Avg60:
		return l.Avg60
	case Avg300:
		return l.Avg300
	default:
		panic("unexpected metric")
	}
}

func (l PSILine) String() string {
	return fmt.Sprintf(lineFormat, l.Avg10, l.Avg60, l.Avg300, l.Total)
}

// PSIStats represent pressure stall information from /proc/pressure/*
//
// "Some" indicates the share of time in which at least some tasks are stalled.
// "Full" indicates the share of time in which all non-idle tasks are stalled simultaneously.
type PSIStats struct {
	Some *PSILine
	Full *PSILine
}

func (s PSIStats) String() string {
	return fmt.Sprintf("some %s\nfull %s", s.Some, s.Full)
}

func (s PSIStats) HasFull() bool {
	return s.Full != nil
}

type PSIStatsResource struct {
	Memory *PSIStats
	CPU    *PSIStats
	IO     *PSIStats
}

func (s PSIStatsResource) String() string {
	return fmt.Sprintf("memory:\n%s\n\ncpu:\n%s\n\nio:\n%s", s.Memory, s.CPU, s.IO)
}

func ResourceToPath(resource Resource) string {
	return fmt.Sprintf("/proc/pressure/%s", resource)
}
