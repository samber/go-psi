package psi

import (
	"log"
	"reflect"
	"time"
)

// Inotify does not support pseudo FS like sysfs and procfs.
// This is a workaround to poll for changes in PSI stats.

const PollingInterval = 100 * time.Millisecond

// Notify notify callers when PSI stats change.
func Notify(resource Resource) (<-chan PSIStats, chan<- struct{}, error) {
	ticker := time.NewTicker(PollingInterval)

	var last PSIStats

	stats := make(chan PSIStats)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				close(done)
				close(stats)
				return
			case <-ticker.C:
				current, err := PSIStatsForResource(resource)
				if err != nil {
					// bearer:disable go_lang_logger_leak
					log.Println(err.Error())
					continue
				}

				if reflect.DeepEqual(last, current) {
					continue
				}

				stats <- current
				last = current
			}
		}
	}()

	return stats, done, nil
}

type StarvationAlert struct {
	Resource      Resource
	Metric        Metric
	LowThreshold  int
	HighThreshold int

	Starved bool
	Stats   PSIStats
	Current float64
}

// NotifyStarvation notify callers when a resource is starved. Starvation is enabled when the metric
// is above the high threshold and disabled when the metric is below the low threshold.
func NotifyStarvation(resource Resource, metric Metric, lowThreshold int, highThreshold int) (<-chan StarvationAlert, chan<- struct{}, error) {
	statsIn, doneIn, err := Notify(resource)
	if err != nil {
		return nil, nil, err
	}

	alerts := make(chan StarvationAlert)
	doneOut := make(chan struct{})

	starved := false

	go func() {
		for {
			select {
			case <-doneOut:
				doneIn <- struct{}{}
				close(alerts)
				close(doneOut)
				return
			case stats := <-statsIn:
				alert := StarvationAlert{
					Resource:      resource,
					Metric:        metric,
					LowThreshold:  lowThreshold,
					HighThreshold: highThreshold,

					Starved: !starved,
					Stats:   stats,
					Current: stats.Some.GetMetric(metric),
				}

				if starved {
					if alert.Current < float64(lowThreshold) {
						starved = false
						alerts <- alert
					}
				} else {
					if alert.Current > float64(highThreshold) {
						starved = true
						alerts <- alert
					}
				}
			}
		}
	}()

	return alerts, doneOut, nil
}
