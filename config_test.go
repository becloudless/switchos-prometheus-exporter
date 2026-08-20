package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "switchos-exporter.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoadConfig(t *testing.T) {
	path := writeTempConfig(t, `
default_username: admin
default_password: defaultpass
targets:
  - host: 192.168.88.1
  - host: 192.168.88.2
    username: other
    password: secret
`)

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if len(cfg.Targets) != 2 {
		t.Fatalf("got %d targets, want 2", len(cfg.Targets))
	}

	t1 := cfg.Targets[0]
	if t1.Host != "192.168.88.1" || t1.Username != "admin" || t1.Password != "defaultpass" {
		t.Errorf("target 0 = %+v, want defaults applied", t1)
	}

	t2 := cfg.Targets[1]
	if t2.Host != "192.168.88.2" || t2.Username != "other" || t2.Password != "secret" {
		t.Errorf("target 1 = %+v, want its own credentials preserved", t2)
	}
}

func TestLoadConfigMissingHost(t *testing.T) {
	path := writeTempConfig(t, `
targets:
  - username: admin
`)

	if _, err := loadConfig(path); err == nil {
		t.Fatal("loadConfig: expected error for target with no host, got nil")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	if _, err := loadConfig(filepath.Join(t.TempDir(), "does-not-exist.yaml")); err == nil {
		t.Fatal("loadConfig: expected error for missing file, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	path := writeTempConfig(t, "targets: [this is not valid: yaml:")

	if _, err := loadConfig(path); err == nil {
		t.Fatal("loadConfig: expected error for invalid YAML, got nil")
	}
}

func TestConfigFind(t *testing.T) {
	cfg := &Config{
		Targets: []Target{
			{Host: "192.168.88.1", Username: "admin"},
			{Host: "192.168.88.2", Username: "other"},
		},
	}

	got, ok := cfg.find("192.168.88.2")
	if !ok {
		t.Fatal("find: expected host to be found")
	}
	if got.Username != "other" {
		t.Errorf("find: got username %q, want %q", got.Username, "other")
	}

	if _, ok := cfg.find("10.0.0.1"); ok {
		t.Error("find: expected unknown host to not be found")
	}
}
