# switchos-prometheus-exporter

A Prometheus exporter for MikroTik SwOS switches (CSS106, CSS318, CSS328,
CSS354, CSS610 families), backed by the unified
[`github.com/becloudless/switchos-client`](../switchos-client) Go client.
Each board is auto-detected on connect, so the same exporter binary and
config format work across every supported board without per-board
configuration.

## Building

```sh
go build -o switchos-prometheus-exporter .
```

## Configuration

Targets (device host + credentials) are read from a YAML config file
(`switchos-exporter.yaml` by default):

```yaml
default_username: admin
default_password: ""

targets:
  - host: 192.168.88.1
  - host: 192.168.88.2
    username: admin
    password: "secret"
```

- `default_username` / `default_password` apply to any target that
  doesn't set its own `username` / `password`.
- Each target needs at least a `host`. SwOS has no per-user accounts -
  `admin` with the device's configured password is the only account.

## Running

```sh
./switchos-prometheus-exporter -config switchos-exporter.yaml
```

Flags:

| Flag                   | Default                   | Description                                  |
|------------------------|----------------------------|-----------------------------------------------|
| `-config`              | `switchos-exporter.yaml`  | Path to the target configuration file          |
| `-web.listen-address`  | `:9435`                   | Address to listen on for HTTP requests         |
| `-web.telemetry-path`  | `/metrics`                | Path under which metrics are exposed           |
| `-scrape.timeout`      | `10s`                     | Per-device timeout for a single scrape         |

## Scraping

- `GET /metrics` scrapes **every** configured target on each request (no
  caching/polling - each scrape talks to the devices live).
- `GET /metrics?target=HOST` scrapes only that one target, following the
  standard multi-target exporter pattern used by `blackbox_exporter` /
  `snmp_exporter`. `HOST` must match a configured target's `host`.
- `GET /-/healthy` always returns `200 OK` once the process is up.

### Prometheus scrape config

Single target per instance (simplest, one exporter process per switch or
scrape everything at once):

```yaml
scrape_configs:
  - job_name: switchos
    static_configs:
      - targets: ["switchos-exporter:9435"]
```

Multi-target pattern (one exporter process serving many switches,
matching them up to Prometheus targets via `target` + relabeling, same
idiom as `blackbox_exporter`):

```yaml
scrape_configs:
  - job_name: switchos
    static_configs:
      - targets:
          - 192.168.88.1
          - 192.168.88.2
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: switchos-exporter:9435
```

## Metrics

All metrics carry a `host` label (the configured target's `host`). Most
per-port metrics additionally carry a `port` label (`"1"`, `"2"`, ...).
Field/metric availability differs by board - see each metric's `HELP`
text (also reproduced below) for which boards expose it; metrics for
fields not supported by a given board are simply omitted from that
target's scrape.

### Scrape / identity

| Metric | Description |
|---|---|
| `switchos_up` | Whether the last scrape of the device succeeded (1) or failed (0). |
| `switchos_scrape_duration_seconds` | Time taken to scrape the device. |
| `switchos_board_info{board}` | Detected board type. Always 1. |
| `switchos_sys_info{identity,model,board_name,mac_address,serial_number}` | Static device identity. Always 1. |
| `switchos_sys_software_info{version,build}` | Installed SwOS firmware version/build. Always 1 when known. |
| `switchos_sys_software_build_raw` | Firmware build identifier, raw wire value (small counter on Classic-dialect boards, Unix build timestamp on Lite-dialect boards - not directly comparable across dialects). |
| `switchos_sys_boot_time_seconds` | Unix timestamp of device boot time (`scrape time - uptime`), harmonized to seconds across board dialects. |

### Health / sensors

| Metric | Description |
|---|---|
| `switchos_sys_cpu_temperature_celsius` | CPU temperature. css610pi/css610g/css328p/css354 only. |
| `switchos_sys_temperature_celsius` | Main board temperature sensor. css106p/css318g only. |
| `switchos_sys_board_temperature1_celsius` | Board temperature sensor 1. css328p/css354 only. |
| `switchos_sys_fan1_rpm` ... `switchos_sys_fan4_rpm` | Fan speed. css328p/css354 only (fan3/fan4: css354 only). |
| `switchos_sys_psu1_voltage_raw`, `switchos_sys_psu2_voltage_raw` | PSU voltage, raw wire value (scale 100). css610pi/css328p/css354 only. |
| `switchos_sys_psu1_current_milliamps`, `switchos_sys_psu2_current_milliamps` | PSU current. css610pi/css328p/css354 only. |
| `switchos_sys_psu1_status` | PSU1 status, raw enum: 0=failed, 1=ok. css354 only. |
| `switchos_sys_power_consumption_raw` | Total power consumption, raw wire value. css610pi only. |

### Per-port link

| Metric | Description |
|---|---|
| `switchos_link_enabled` | Whether the port is administratively enabled. |
| `switchos_link_autonegotiation_enabled` | Whether auto-negotiation is enabled. |
| `switchos_link_duplex_control_enabled` | Whether full duplex is configured (vs. auto/half). |
| `switchos_link_duplex` | Whether the port is currently running full duplex (read-only status). |
| `switchos_link_up` | Whether the port currently has link. |
| `switchos_link_speed_control` | Configured speed, raw enum value (board/dialect-specific). |
| `switchos_link_speed` | Current negotiated speed, raw enum value (read-only). |
| `switchos_link_info{name}` | Static per-port info (port name). Always 1. |

### Per-port traffic counters (`switchos_stats_*_total`)

Standard RMON-style counters: `rx_bytes`, `tx_bytes`, `rx_unicast_packets`,
`tx_unicast_packets`, `rx_broadcast_packets`, `tx_broadcast_packets`,
`rx_multicast_packets`, `tx_multicast_packets`, `rx_pause_frames`,
`tx_pause_frames`, `rx_runt_frames`, `rx_fragment_frames`,
`rx_fcs_errors`, `rx_too_long_frames`, `rx_packets`, `tx_packets`,
`tx_collisions`, `tx_deferred_frames`, `tx_excessive_collisions`,
`tx_late_collisions`, `tx_multiple_collisions`, per-size RX/TX frame
histograms (`rx_frames_64`, `rx_frames_65_127`, ... `tx_frames_1519_plus`),
plus several board-specific error counters (`rx_overruns`,
`tx_fcs_errors`, `rx_jabber_frames`, `rx_errors`,
`tx_single_collisions`, `rx_align_errors`, `rx_total_errors`,
`tx_excessive_deferred`, `tx_total_errors`, `tx_too_long`,
`rx_mac_errors`, `tx_underruns`). See `metrics_stats.go` for the full
list and exact per-board availability.

| Metric | Description |
|---|---|
| `switchos_stats_tx_queue_packets` | Current TX queue depth in packets (gauge). css318g/css328p only. |
| `switchos_stats_tx_queue_bytes` | Current TX queue depth in bytes (gauge). css318g/css328p only. |

### Per-port PoE (css610pi/css328p/css354, boards with PoE hardware only)

| Metric | Description |
|---|---|
| `switchos_poe_out` | PoE output enable, raw enum: 0=off, 1=on, 2=auto. |
| `switchos_poe_priority` | PoE priority, raw enum value. |
| `switchos_poe_voltage_level` | PoE voltage level, raw enum: 0=auto, 1=low, 2=high. |
| `switchos_poe_status` | PoE status, raw enum. |
| `switchos_poe_current_milliamps` | PoE output current. |
| `switchos_poe_voltage_volts` | PoE output voltage. |
| `switchos_poe_power_watts` | PoE output power. |
| `switchos_poe_lldp_power_watts` | PoE power allocated via LLDP negotiation. css610pi/css328p only. |
| `switchos_poe_lldp_enabled` | Whether PoE-LLDP negotiation is enabled. css610pi/css328p only. |

### Per-port SFP/SFP+ (only emitted for ports with a transceiver present)

| Metric | Description |
|---|---|
| `switchos_sfp_info{vendor,part_number,revision,serial,date,type}` | Static transceiver info. Always 1. |
| `switchos_sfp_temperature_celsius` | Transceiver temperature. |
| `switchos_sfp_voltage_volts` | Transceiver supply voltage. |
| `switchos_sfp_tx_bias_milliamps` | Transceiver laser bias current. |
| `switchos_sfp_tx_power_dbm` | Transceiver transmit optical power. |
| `switchos_sfp_rx_power_dbm` | Transceiver receive optical power. |

## Grafana dashboard

A ready-to-import dashboard is provided at
[`grafana/switchos-dashboard.json`](grafana/switchos-dashboard.json).
It includes:

- an overview row (up/down, board/firmware info, boot time)
- health/sensor panels (temperatures, fans, PSU)
- per-port link status and negotiated speed
- per-port RX/TX throughput and error-rate panels
- per-port PoE power/current/status (for boards with PoE)
- per-port SFP optical diagnostics

Import via Grafana's **Dashboards → Import**, using a Prometheus
datasource that scrapes this exporter. The dashboard defines a `host`
template variable (populated from `switchos_up`) to filter all panels
down to one switch at a time.
