package main

import (
	"context"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	linkEnabledDesc       = prometheus.NewDesc("switchos_link_enabled", "Whether the port is administratively enabled.", []string{"host", "port"}, nil)
	linkAutonegDesc       = prometheus.NewDesc("switchos_link_autonegotiation_enabled", "Whether auto-negotiation is enabled.", []string{"host", "port"}, nil)
	linkDuplexControlDesc = prometheus.NewDesc("switchos_link_duplex_control_enabled", "Whether full duplex is configured (as opposed to auto/half).", []string{"host", "port"}, nil)
	linkDuplexDesc        = prometheus.NewDesc("switchos_link_duplex", "Whether the port is currently running full duplex (read-only status).", []string{"host", "port"}, nil)
	linkUpDesc            = prometheus.NewDesc("switchos_link_up", "Whether the port currently has link (read-only status; true for both the plain link-up and link-up-with-pause wire states).", []string{"host", "port"}, nil)
	linkSpeedControlDesc  = prometheus.NewDesc("switchos_link_speed_control", "Configured speed, raw enum value (meaning is board/dialect-specific, see switchos-client docs).", []string{"host", "port"}, nil)
	linkSpeedDesc         = prometheus.NewDesc("switchos_link_speed", "Current negotiated speed, raw enum value (read-only, meaning is board/dialect-specific).", []string{"host", "port"}, nil)
	linkInfoDesc          = prometheus.NewDesc("switchos_link_info", "Static per-port info. Always 1.", []string{"host", "port", "name"}, nil)
)

func collectLink(ctx context.Context, ch chan<- prometheus.Metric, client *switchos.Client, host string) {
	l, err := client.GetLink(ctx)
	if err != nil {
		return
	}

	// LinkState (Lite dialect: css610pi/css610g) and LinkStatus (Classic
	// dialect: css106p/css318g/css328p/css354) are mutually exclusive,
	// same-shaped "is there link" bitmasks - merge into one metric.
	up := l.LinkState
	if up == nil {
		up = l.LinkStatus
	}

	n := arrayLenInt32(l.SpeedControl, l.Speed)
	if nn := arrayLenString(l.Name); nn > n {
		n = nn
	}

	for i := 0; i < n; i++ {
		port := portLabel(i)

		if v, ok := bitGet64(l.Enabled, i); ok {
			gaugeBool(ch, linkEnabledDesc, v, host, port)
		}
		if v, ok := bitGet64(l.AutoNegotiation, i); ok {
			gaugeBool(ch, linkAutonegDesc, v, host, port)
		}
		if v, ok := bitGet64(l.DuplexControl, i); ok {
			gaugeBool(ch, linkDuplexControlDesc, v, host, port)
		}
		if v, ok := bitGet32(l.Duplex, i); ok {
			gaugeBool(ch, linkDuplexDesc, v, host, port)
		}
		if v, ok := bitGet32(up, i); ok {
			gaugeBool(ch, linkUpDesc, v, host, port)
		}
		gaugeInt32(ch, linkSpeedControlDesc, atInt32(l.SpeedControl, i), host, port)
		gaugeInt32(ch, linkSpeedDesc, atInt32(l.Speed, i), host, port)
		if name := atString(l.Name, i); name != nil {
			ch <- prometheus.MustNewConstMetric(linkInfoDesc, prometheus.GaugeValue, 1, host, port, *name)
		}
	}
}
