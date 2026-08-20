package main

import (
	"context"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	sfpInfoDesc        = prometheus.NewDesc("switchos_sfp_info", "Static SFP/SFP+ transceiver info for a port with a module present. Always 1.", []string{"host", "port", "vendor", "part_number", "revision", "serial", "date", "type"}, nil)
	sfpTemperatureDesc = prometheus.NewDesc("switchos_sfp_temperature_celsius", "SFP transceiver temperature.", []string{"host", "port"}, nil)
	sfpVoltageDesc     = prometheus.NewDesc("switchos_sfp_voltage_volts", "SFP transceiver supply voltage.", []string{"host", "port"}, nil)
	sfpTxBiasDesc      = prometheus.NewDesc("switchos_sfp_tx_bias_milliamps", "SFP transceiver laser bias current.", []string{"host", "port"}, nil)
	sfpTxPowerDesc     = prometheus.NewDesc("switchos_sfp_tx_power_dbm", "SFP transceiver transmit optical power.", []string{"host", "port"}, nil)
	sfpRxPowerDesc     = prometheus.NewDesc("switchos_sfp_rx_power_dbm", "SFP transceiver receive optical power.", []string{"host", "port"}, nil)
)

// sfpTemperatureSentinel is the documented "no module present" sentinel
// value for SfpPort.Temperature.
const sfpTemperatureSentinel = -128

func collectSfp(ctx context.Context, ch chan<- prometheus.Metric, client *switchos.Client, host string) {
	sfps, err := client.GetSfp(ctx)
	if err != nil {
		return
	}

	for i, s := range sfps {
		lbl := portLabel(i)

		if s.Temperature == nil || *s.Temperature == sfpTemperatureSentinel {
			// No transceiver present in this port.
			continue
		}

		ch <- prometheus.MustNewConstMetric(sfpInfoDesc, prometheus.GaugeValue, 1,
			host, lbl, derefString(s.Vendor), derefString(s.PartNumber), derefString(s.Revision), derefString(s.Serial), derefString(s.Date), derefString(s.Type))

		gaugeInt64(ch, sfpTemperatureDesc, s.Temperature, host, lbl)
		gaugeScaled(ch, sfpVoltageDesc, s.Voltage, 1000, host, lbl)
		gaugeInt64(ch, sfpTxBiasDesc, s.TxBias, host, lbl)
		gaugeScaled(ch, sfpTxPowerDesc, s.TxPower, 10000, host, lbl)
		gaugeScaled(ch, sfpRxPowerDesc, s.RxPower, 10000, host, lbl)
	}
}
