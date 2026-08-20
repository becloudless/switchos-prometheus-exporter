package main

import (
	"context"
	"strconv"
	"time"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	sysInfoDesc = prometheus.NewDesc(
		"switchos_sys_info", "Static device identity info. Always 1.",
		[]string{"host", "identity", "model", "board_name", "mac_address", "serial_number"}, nil,
	)
	sysBootTimeDesc          = prometheus.NewDesc("switchos_sys_boot_time_seconds", "Unix timestamp of device boot time (scrape time minus uptime), harmonized to seconds across board dialects: css610pi/css610g's uptime is native seconds; Classic dialect's (css106p/css318g/css328p/css354) is centiseconds (raw/100), per the SwOS UI's own {scale:100} uptime field metadata.", []string{"host"}, nil)
	sysCpuTemperatureDesc    = prometheus.NewDesc("switchos_sys_cpu_temperature_celsius", "CPU temperature. css610pi/css610g/css328p/css354 only.", []string{"host"}, nil)
	sysTemperatureDesc       = prometheus.NewDesc("switchos_sys_temperature_celsius", "Main board temperature sensor. css106p/css318g only.", []string{"host"}, nil)
	sysBoardTemperature1Desc = prometheus.NewDesc("switchos_sys_board_temperature1_celsius", "Board temperature sensor 1. css328p/css354 only.", []string{"host"}, nil)
	sysFan1Desc              = prometheus.NewDesc("switchos_sys_fan1_rpm", "Fan 1 speed. css328p/css354 only.", []string{"host"}, nil)
	sysFan2Desc              = prometheus.NewDesc("switchos_sys_fan2_rpm", "Fan 2 speed. css328p/css354 only.", []string{"host"}, nil)
	sysFan3Desc              = prometheus.NewDesc("switchos_sys_fan3_rpm", "Fan 3 speed. css354 only.", []string{"host"}, nil)
	sysFan4Desc              = prometheus.NewDesc("switchos_sys_fan4_rpm", "Fan 4 speed. css354 only.", []string{"host"}, nil)
	sysPsu1VoltageDesc       = prometheus.NewDesc("switchos_sys_psu1_voltage_raw", "PSU1 voltage, raw wire value (scale 100). css610pi/css328p/css354 only.", []string{"host"}, nil)
	sysPsu1CurrentDesc       = prometheus.NewDesc("switchos_sys_psu1_current_milliamps", "PSU1 current. css610pi/css328p/css354 only.", []string{"host"}, nil)
	sysPsu2VoltageDesc       = prometheus.NewDesc("switchos_sys_psu2_voltage_raw", "PSU2 voltage, raw wire value (scale 100). css610pi/css328p/css354 only.", []string{"host"}, nil)
	sysPsu2CurrentDesc       = prometheus.NewDesc("switchos_sys_psu2_current_milliamps", "PSU2 current. css610pi/css328p/css354 only.", []string{"host"}, nil)
	sysPsu1StatusDesc        = prometheus.NewDesc("switchos_sys_psu1_status", "PSU1 status, raw enum: 0=failed, 1=ok. css354 only.", []string{"host"}, nil)
	sysPowerConsumptionDesc  = prometheus.NewDesc("switchos_sys_power_consumption_raw", "Total power consumption, raw wire value. css610pi only.", []string{"host"}, nil)
	sysSoftwareInfoDesc      = prometheus.NewDesc("switchos_sys_software_info", "Installed SwOS firmware version/build. Always 1 when known.", []string{"host", "version", "build"}, nil)
	sysSoftwareBuildDesc     = prometheus.NewDesc("switchos_sys_software_build_raw", "Installed SwOS firmware build identifier, raw wire value. On Classic-dialect boards (css106p/css318g/css328p/css354) this is a small sequential build counter (may be absent on older firmware); on Lite-dialect boards (css610pi/css610g) this is a Unix build timestamp instead.", []string{"host"}, nil)
)

func collectSys(ctx context.Context, ch chan<- prometheus.Metric, client *switchos.Client, host string) {
	s, err := client.GetSys(ctx)
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(sysInfoDesc, prometheus.GaugeValue, 1,
		host, derefString(s.Identity), derefString(s.Model), derefString(s.BoardName), derefString(s.MacAddress), derefString(s.SerialNumber))

	uptimeSeconds := int64(0)
	haveUptime := false
	if s.Uptime != nil {
		uptimeSeconds = *s.Uptime
		haveUptime = true
	} else if s.ClassicUptime != nil {
		uptimeSeconds = *s.ClassicUptime / 100
		haveUptime = true
	}
	if haveUptime {
		bootTime := time.Now().Unix() - uptimeSeconds
		ch <- prometheus.MustNewConstMetric(sysBootTimeDesc, prometheus.GaugeValue, float64(bootTime), host)
	}
	gaugeInt64(ch, sysCpuTemperatureDesc, s.CpuTemperature, host)
	gaugeInt64(ch, sysTemperatureDesc, s.Temperature, host)
	gaugeInt64(ch, sysBoardTemperature1Desc, s.BoardTemperature1, host)
	gaugeInt64(ch, sysFan1Desc, s.Fan1, host)
	gaugeInt64(ch, sysFan2Desc, s.Fan2, host)
	gaugeInt64(ch, sysFan3Desc, s.Fan3, host)
	gaugeInt64(ch, sysFan4Desc, s.Fan4, host)
	gaugeInt64(ch, sysPsu1VoltageDesc, s.Psu1Voltage, host)
	gaugeInt64(ch, sysPsu1CurrentDesc, s.Psu1Current, host)
	gaugeInt64(ch, sysPsu2VoltageDesc, s.Psu2Voltage, host)
	gaugeInt64(ch, sysPsu2CurrentDesc, s.Psu2Current, host)
	gaugeInt32(ch, sysPsu1StatusDesc, s.Psu1Status, host)
	gaugeInt64(ch, sysPowerConsumptionDesc, s.PowerConsumption, host)

	if s.SoftwareVersion != nil {
		build := ""
		if s.SoftwareBuild != nil {
			build = strconv.FormatInt(*s.SoftwareBuild, 10)
		}
		ch <- prometheus.MustNewConstMetric(sysSoftwareInfoDesc, prometheus.GaugeValue, 1, host, *s.SoftwareVersion, build)
	}
	gaugeInt64(ch, sysSoftwareBuildDesc, s.SoftwareBuild, host)
}
