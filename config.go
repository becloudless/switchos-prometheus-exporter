// Config is the exporter's target configuration file format (YAML).
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the exporter's target configuration file format.
type Config struct {
	DefaultUsername string   `yaml:"default_username"`
	DefaultPassword string   `yaml:"default_password"`
	Targets         []Target `yaml:"targets"`
}

// Target is one configured switch.
type Target struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}

	for i := range cfg.Targets {
		if cfg.Targets[i].Host == "" {
			return nil, fmt.Errorf("config %q: target %d has no host", path, i)
		}
		if cfg.Targets[i].Username == "" {
			cfg.Targets[i].Username = cfg.DefaultUsername
		}
		if cfg.Targets[i].Password == "" {
			cfg.Targets[i].Password = cfg.DefaultPassword
		}
	}

	return &cfg, nil
}

// find returns the configured target matching host, if any.
func (c *Config) find(host string) (Target, bool) {
	for _, t := range c.Targets {
		if t.Host == host {
			return t, true
		}
	}
	return Target{}, false
}
