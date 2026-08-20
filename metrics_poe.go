package main

import (
	"context"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	poeOutDesc          = prometheus.NewDesc("switchos_poe_out", "PoE output enable, raw enum: 0=off, 1=on, 2=auto.", []string{"host", "port"}, nil)
	poePriorityDesc     = prometheus.NewDesc("switchos_poe_priority", "PoE priority, raw enum value (meaning is board-specific).", []string{"host", "port"}, nil)
	poeVoltageLevelDesc = prometheus.NewDesc("switchos_poe_voltage_level", "PoE voltage level, raw enum: 0=auto, 1=low, 2=high.", []string{"host", "port"}, nil)
	poeStatusDesc       = prometheus.NewDesc("switchos_poe_status", "PoE status, raw enum (e.g. 3=powered on); see switchos-client docs for the full value table.", []string{"host", "port"}, nil)
	poeCurrentDesc      = prometheus.NewDesc("switchos_poe_current_milliamps", "PoE output current.", []string{"host", "port"}, nil)
	poeVoltageDesc      = prometheus.NewDesc("switchos_poe_voltage_volts", "PoE output voltage.", []string{"host", "port"}, nil)
	poePowerDesc        = prometheus.NewDesc("switchos_poe_power_watts", "PoE output power.", []string{"host", "port"}, nil)
	poeLldpPowerDesc    = prometheus.NewDesc("switchos_poe_lldp_power_watts", "PoE power allocated via LLDP negotiation. css610pi/css328p only.", []string{"host", "port"}, nil)
	poeLldpEnabledDesc  = prometheus.NewDesc("switchos_poe_lldp_enabled", "Whether PoE-LLDP negotiation is enabled for this port. css610pi/css328p only.", []string{"host", "port"}, nil)
)

func collectPoe(ctx context.Context, ch chan<- prometheus.Metric, client *switchos.Client, host string) {
	p, err := client.GetPoe(ctx)
	if err != nil {
		// ErrUnsupported on boards without PoE hardware, or a real
		// error either way - nothing to export.
		return
	}

	for i, port := range p.Ports {
		lbl := portLabel(i)
		gaugeInt32(ch, poeOutDesc, port.Out, host, lbl)
		gaugeInt32(ch, poePriorityDesc, port.Priority, host, lbl)
		gaugeInt32(ch, poeVoltageLevelDesc, port.VoltageLevel, host, lbl)
		gaugeInt32(ch, poeStatusDesc, port.Status, host, lbl)
		gaugeInt64(ch, poeCurrentDesc, port.Current, host, lbl)
		gaugeScaled(ch, poeVoltageDesc, port.Voltage, 10, host, lbl)
		gaugeScaled(ch, poePowerDesc, port.Power, 10, host, lbl)
		gaugeScaled(ch, poeLldpPowerDesc, port.LldpPower, 10, host, lbl)
		if v, ok := bitGet64(p.LldpEnabled, i); ok {
			gaugeBool(ch, poeLldpEnabledDesc, v, host, lbl)
		}
	}
}
