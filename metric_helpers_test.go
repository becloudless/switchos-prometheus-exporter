package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// collectOne drains exactly one metric from ch (failing the test if none
// arrives) and decodes it into a dto.Metric for inspection.
func collectOne(t *testing.T, ch chan prometheus.Metric) *dto.Metric {
	t.Helper()
	select {
	case m := <-ch:
		var out dto.Metric
		if err := m.Write(&out); err != nil {
			t.Fatalf("writing metric: %v", err)
		}
		return &out
	default:
		t.Fatal("expected a metric to be emitted, got none")
		return nil
	}
}

func assertNoneEmitted(t *testing.T, ch chan prometheus.Metric) {
	t.Helper()
	select {
	case m := <-ch:
		var out dto.Metric
		_ = m.Write(&out)
		t.Fatalf("expected no metric to be emitted, got %v", &out)
	default:
	}
}

func TestGaugeInt64(t *testing.T) {
	desc := prometheus.NewDesc("test_gauge_int64", "help", []string{"host"}, nil)
	ch := make(chan prometheus.Metric, 1)

	gaugeInt64(ch, desc, nil, "h1")
	assertNoneEmitted(t, ch)

	v := int64(42)
	gaugeInt64(ch, desc, &v, "h1")
	m := collectOne(t, ch)
	if got := m.GetGauge().GetValue(); got != 42 {
		t.Errorf("gauge value = %v, want 42", got)
	}
}

func TestGaugeInt32(t *testing.T) {
	desc := prometheus.NewDesc("test_gauge_int32", "help", []string{"host"}, nil)
	ch := make(chan prometheus.Metric, 1)

	gaugeInt32(ch, desc, nil, "h1")
	assertNoneEmitted(t, ch)

	v := int32(7)
	gaugeInt32(ch, desc, &v, "h1")
	m := collectOne(t, ch)
	if got := m.GetGauge().GetValue(); got != 7 {
		t.Errorf("gauge value = %v, want 7", got)
	}
}

func TestCounterInt64(t *testing.T) {
	desc := prometheus.NewDesc("test_counter_int64", "help", []string{"host"}, nil)
	ch := make(chan prometheus.Metric, 1)

	counterInt64(ch, desc, nil, "h1")
	assertNoneEmitted(t, ch)

	v := int64(1000)
	counterInt64(ch, desc, &v, "h1")
	m := collectOne(t, ch)
	if got := m.GetCounter().GetValue(); got != 1000 {
		t.Errorf("counter value = %v, want 1000", got)
	}
}

func TestGaugeScaled(t *testing.T) {
	desc := prometheus.NewDesc("test_gauge_scaled", "help", []string{"host"}, nil)
	ch := make(chan prometheus.Metric, 1)

	gaugeScaled(ch, desc, nil, 10, "h1")
	assertNoneEmitted(t, ch)

	v := int64(255)
	gaugeScaled(ch, desc, &v, 10, "h1")
	m := collectOne(t, ch)
	if got := m.GetGauge().GetValue(); got != 25.5 {
		t.Errorf("scaled gauge value = %v, want 25.5", got)
	}
}

func TestGaugeBool(t *testing.T) {
	desc := prometheus.NewDesc("test_gauge_bool", "help", []string{"host"}, nil)
	ch := make(chan prometheus.Metric, 1)

	gaugeBool(ch, desc, true, "h1")
	m := collectOne(t, ch)
	if got := m.GetGauge().GetValue(); got != 1 {
		t.Errorf("gaugeBool(true) value = %v, want 1", got)
	}

	gaugeBool(ch, desc, false, "h1")
	m = collectOne(t, ch)
	if got := m.GetGauge().GetValue(); got != 0 {
		t.Errorf("gaugeBool(false) value = %v, want 0", got)
	}
}
