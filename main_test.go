package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleMetricsUnknownTarget(t *testing.T) {
	cfg := &Config{
		Targets: []Target{{Host: "192.168.88.1", Username: "admin"}},
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics?target=10.0.0.1", nil)
	w := httptest.NewRecorder()

	handleMetrics(w, req, cfg, time.Second)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
