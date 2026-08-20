package main

import (
	"context"
	"sync"
	"time"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	upDesc = prometheus.NewDesc(
		"switchos_up", "Whether the last scrape of the device succeeded (1) or failed (0).",
		[]string{"host"}, nil,
	)
	boardInfoDesc = prometheus.NewDesc(
		"switchos_board_info", "Static info about the detected board. Always 1.",
		[]string{"host", "board"}, nil,
	)
	scrapeDurationDesc = prometheus.NewDesc(
		"switchos_scrape_duration_seconds", "Time taken to scrape the device.",
		[]string{"host"}, nil,
	)
)

// collector scrapes a fixed set of targets on every Collect call. It's
// built fresh per HTTP request (see main.go's handleMetrics) so the
// "target" query parameter can select a subset of the configured
// targets.
type collector struct {
	targets []Target
	timeout time.Duration
}

func newCollector(targets []Target, timeout time.Duration) *collector {
	return &collector{targets: targets, timeout: timeout}
}

// Describe intentionally sends nothing on descs: the metrics this
// collector emits have dynamic label sets (per-port labels, and
// per-board field availability), so this is an "unchecked" collector -
// a standard, supported pattern for exporters like this one (see
// blackbox_exporter, snmp_exporter).
func (c *collector) Describe(_ chan<- *prometheus.Desc) {}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	for _, t := range c.targets {
		wg.Add(1)
		go func(t Target) {
			defer wg.Done()
			c.collectTarget(ch, t)
		}(t)
	}
	wg.Wait()
}

func (c *collector) collectTarget(ch chan<- prometheus.Metric, t Target) {
	start := time.Now()
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), t.Host)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	client, err := switchos.Dial(ctx, "http://"+t.Host, t.Username, t.Password)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0, t.Host)
		return
	}
	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1, t.Host)
	ch <- prometheus.MustNewConstMetric(boardInfoDesc, prometheus.GaugeValue, 1, t.Host, client.Board().String())

	collectSys(ctx, ch, client, t.Host)
	collectLink(ctx, ch, client, t.Host)
	collectPoe(ctx, ch, client, t.Host)
	collectSfp(ctx, ch, client, t.Host)
	collectStats(ctx, ch, client, t.Host)
}
