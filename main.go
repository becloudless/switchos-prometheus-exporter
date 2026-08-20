// Command switchos-prometheus-exporter is a Prometheus exporter for
// MikroTik SwOS switches (CSS106, CSS318, CSS328, CSS354, CSS610
// families), backed by the unified github.com/becloudless/switchos-client
// Go client.
//
// Targets (device host + credentials) are read from a YAML config file.
// GET /metrics scrapes every configured target. GET /metrics?target=HOST
// scrapes only that one target (HOST must match a configured target's
// host, so its credentials can be looked up), following the standard
// multi-target exporter pattern (see blackbox_exporter/snmp_exporter).
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		configPath    = flag.String("config", "switchos-exporter.yaml", "path to the target configuration file")
		listenAddress = flag.String("web.listen-address", ":9435", "address to listen on for HTTP requests")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "path under which to expose metrics")
		timeout       = flag.Duration("scrape.timeout", 10*time.Second, "per-device timeout for a single scrape")
	)
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("switchos_exporter: %v", err)
	}

	http.HandleFunc(*metricsPath, func(w http.ResponseWriter, r *http.Request) {
		handleMetrics(w, r, cfg, *timeout)
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `<html><head><title>switchos Prometheus Exporter</title></head><body>
<h1>switchos Prometheus Exporter</h1>
<p><a href="%[1]s">Metrics</a> (scrapes every configured target)</p>
<p><code>%[1]s?target=HOST</code> scrapes only that one target</p>
</body></html>`, *metricsPath)
	})

	log.Printf("switchos_exporter: listening on %s, config=%s, %d target(s) configured", *listenAddress, *configPath, len(cfg.Targets))
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetrics(w http.ResponseWriter, r *http.Request, cfg *Config, timeout time.Duration) {
	var targets []Target

	if host := r.URL.Query().Get("target"); host != "" {
		t, ok := cfg.find(host)
		if !ok {
			http.Error(w, fmt.Sprintf("unknown target %q (not present in config)", host), http.StatusBadRequest)
			return
		}
		targets = []Target{t}
	} else {
		targets = cfg.Targets
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(newCollector(targets, timeout))
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
