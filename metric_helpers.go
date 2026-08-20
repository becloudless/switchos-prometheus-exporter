package main

import "github.com/prometheus/client_golang/prometheus"

func gaugeInt64(ch chan<- prometheus.Metric, desc *prometheus.Desc, p *int64, labels ...string) {
	if p == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(*p), labels...)
}

func gaugeInt32(ch chan<- prometheus.Metric, desc *prometheus.Desc, p *int32, labels ...string) {
	if p == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(*p), labels...)
}

func counterInt64(ch chan<- prometheus.Metric, desc *prometheus.Desc, p *int64, labels ...string) {
	if p == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(*p), labels...)
}

// gaugeScaled emits *p / scale as a gauge. Used for fields whose
// switchos-client doc comment states a fixed scale factor consistent
// across every board that provides the field (e.g. Sfp voltage/power).
func gaugeScaled(ch chan<- prometheus.Metric, desc *prometheus.Desc, p *int64, scale float64, labels ...string) {
	if p == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(*p)/scale, labels...)
}

func gaugeBool(ch chan<- prometheus.Metric, desc *prometheus.Desc, val bool, labels ...string) {
	v := 0.0
	if val {
		v = 1
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, labels...)
}
