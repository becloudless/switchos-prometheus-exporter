package main

import (
	"context"

	switchos "github.com/becloudless/switchos-client"

	"github.com/prometheus/client_golang/prometheus"
)

func statsDesc(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(name, help, []string{"host", "port"}, nil)
}

var (
	statsRxBytesDesc               = statsDesc("switchos_stats_rx_bytes_total", "Received bytes.")
	statsTxBytesDesc               = statsDesc("switchos_stats_tx_bytes_total", "Transmitted bytes.")
	statsRxUnicastsDesc            = statsDesc("switchos_stats_rx_unicast_packets_total", "Received unicast packets.")
	statsTxUnicastsDesc            = statsDesc("switchos_stats_tx_unicast_packets_total", "Transmitted unicast packets.")
	statsRxBroadcastsDesc          = statsDesc("switchos_stats_rx_broadcast_packets_total", "Received broadcast packets.")
	statsTxBroadcastsDesc          = statsDesc("switchos_stats_tx_broadcast_packets_total", "Transmitted broadcast packets.")
	statsRxMulticastsDesc          = statsDesc("switchos_stats_rx_multicast_packets_total", "Received multicast packets.")
	statsTxMulticastsDesc          = statsDesc("switchos_stats_tx_multicast_packets_total", "Transmitted multicast packets.")
	statsRxPausesDesc              = statsDesc("switchos_stats_rx_pause_frames_total", "Received pause frames.")
	statsTxPausesDesc              = statsDesc("switchos_stats_tx_pause_frames_total", "Transmitted pause frames.")
	statsRxRuntsDesc               = statsDesc("switchos_stats_rx_runt_frames_total", "Received runt frames.")
	statsRxFragmentsDesc           = statsDesc("switchos_stats_rx_fragment_frames_total", "Received fragment frames.")
	statsRxFcsErrorsDesc           = statsDesc("switchos_stats_rx_fcs_errors_total", "Received frames with an FCS error.")
	statsRxTooLongDesc             = statsDesc("switchos_stats_rx_too_long_frames_total", "Received oversize frames. css610pi/css610g/css106p only.")
	statsRxTotalPacketsDesc        = statsDesc("switchos_stats_rx_packets_total", "Total received packets.")
	statsTxTotalPacketsDesc        = statsDesc("switchos_stats_tx_packets_total", "Total transmitted packets.")
	statsTxCollisionsDesc          = statsDesc("switchos_stats_tx_collisions_total", "Transmit collisions.")
	statsTxDeferredDesc            = statsDesc("switchos_stats_tx_deferred_frames_total", "Transmit deferred frames.")
	statsTxExcessiveCollisionsDesc = statsDesc("switchos_stats_tx_excessive_collisions_total", "Transmit excessive collisions.")
	statsTxLateCollisionsDesc      = statsDesc("switchos_stats_tx_late_collisions_total", "Transmit late collisions.")
	statsTxMultipleCollisionsDesc  = statsDesc("switchos_stats_tx_multiple_collisions_total", "Transmit multiple collisions.")

	statsRxFrames64Desc        = statsDesc("switchos_stats_rx_frames_64_total", "Received frames of size 64 bytes.")
	statsRxFrames65To127Desc   = statsDesc("switchos_stats_rx_frames_65_127_total", "Received frames of size 65-127 bytes.")
	statsRxFrames128To255Desc  = statsDesc("switchos_stats_rx_frames_128_255_total", "Received frames of size 128-255 bytes.")
	statsRxFrames256To511Desc  = statsDesc("switchos_stats_rx_frames_256_511_total", "Received frames of size 256-511 bytes.")
	statsRxFrames512To1023Desc = statsDesc("switchos_stats_rx_frames_512_1023_total", "Received frames of size 512-1023 bytes.")
	statsRxFrames1024PlusDesc  = statsDesc("switchos_stats_rx_frames_1024_plus_total", "Received frames of size 1024+ bytes. Not on css106p, see the two split buckets below instead.")

	statsRxOverrunsDesc = statsDesc("switchos_stats_rx_overruns_total", "Receive overruns. Not on css610pi/css610g.")

	statsTxFcsErrorsDesc        = statsDesc("switchos_stats_tx_fcs_errors_total", "Transmitted frames with an FCS error. css610pi/css610g only.")
	statsRxJabberDesc           = statsDesc("switchos_stats_rx_jabber_frames_total", "Received jabber frames. css610pi/css610g only.")
	statsRxErrorsDesc           = statsDesc("switchos_stats_rx_errors_total", "Receive errors. css610pi/css610g only.")
	statsTxSingleCollisionsDesc = statsDesc("switchos_stats_tx_single_collisions_total", "Transmit single collisions. css610pi/css610g only.")

	statsRxFrames1024To1518Desc  = statsDesc("switchos_stats_rx_frames_1024_1518_total", "Received frames of size 1024-1518 bytes. css106p only.")
	statsRxFrames1519PlusDesc    = statsDesc("switchos_stats_rx_frames_1519_plus_total", "Received frames of size 1519+ bytes. css106p only.")
	statsTxFrames64Desc          = statsDesc("switchos_stats_tx_frames_64_total", "Transmitted frames of size 64 bytes. css106p only.")
	statsTxFrames65To127Desc     = statsDesc("switchos_stats_tx_frames_65_127_total", "Transmitted frames of size 65-127 bytes. css106p only.")
	statsTxFrames128To255Desc    = statsDesc("switchos_stats_tx_frames_128_255_total", "Transmitted frames of size 128-255 bytes. css106p only.")
	statsTxFrames256To511Desc    = statsDesc("switchos_stats_tx_frames_256_511_total", "Transmitted frames of size 256-511 bytes. css106p only.")
	statsTxFrames512To1023Desc   = statsDesc("switchos_stats_tx_frames_512_1023_total", "Transmitted frames of size 512-1023 bytes. css106p only.")
	statsTxFrames1024To1518Desc  = statsDesc("switchos_stats_tx_frames_1024_1518_total", "Transmitted frames of size 1024-1518 bytes. css106p only.")
	statsTxFrames1519PlusDesc    = statsDesc("switchos_stats_tx_frames_1519_plus_total", "Transmitted frames of size 1519+ bytes. css106p only.")
	statsRxAlignErrorsDesc       = statsDesc("switchos_stats_rx_align_errors_total", "Receive alignment errors. css106p only.")
	statsRxTotalErrorsDesc       = statsDesc("switchos_stats_rx_total_errors_total", "Total receive errors. css106p only.")
	statsTxExcessiveDeferredDesc = statsDesc("switchos_stats_tx_excessive_deferred_total", "Transmit excessive deferred frames. css106p only.")
	statsTxTotalErrorsDesc       = statsDesc("switchos_stats_tx_total_errors_total", "Total transmit errors. css106p only.")
	statsTxTooLongDesc           = statsDesc("switchos_stats_tx_too_long_total", "Transmitted oversize frames. css106p only.")

	statsRxMacErrorsDesc = statsDesc("switchos_stats_rx_mac_errors_total", "Receive MAC errors. css318g/css328p/css354 only.")
	statsTxUnderrunsDesc = statsDesc("switchos_stats_tx_underruns_total", "Transmit underruns. css318g/css328p/css354 only.")

	statsTxQueueDesc      = prometheus.NewDesc("switchos_stats_tx_queue_packets", "Current TX queue depth in packets (a gauge, not a cumulative counter). css318g/css328p only.", []string{"host", "port"}, nil)
	statsTxQueueBytesDesc = prometheus.NewDesc("switchos_stats_tx_queue_bytes", "Current TX queue depth in bytes (a gauge, not a cumulative counter). css318g/css328p only.", []string{"host", "port"}, nil)
)

type statsCounterField struct {
	desc  *prometheus.Desc
	value *[]int64
}

func collectStats(ctx context.Context, ch chan<- prometheus.Metric, client *switchos.Client, host string) {
	s, err := client.GetStats(ctx)
	if err != nil {
		return
	}

	fields := []statsCounterField{
		{statsRxBytesDesc, s.RxBytes},
		{statsTxBytesDesc, s.TxBytes},
		{statsRxUnicastsDesc, s.RxUnicasts},
		{statsTxUnicastsDesc, s.TxUnicasts},
		{statsRxBroadcastsDesc, s.RxBroadcasts},
		{statsTxBroadcastsDesc, s.TxBroadcasts},
		{statsRxMulticastsDesc, s.RxMulticasts},
		{statsTxMulticastsDesc, s.TxMulticasts},
		{statsRxPausesDesc, s.RxPauses},
		{statsTxPausesDesc, s.TxPauses},
		{statsRxRuntsDesc, s.RxRunts},
		{statsRxFragmentsDesc, s.RxFragments},
		{statsRxFcsErrorsDesc, s.RxFcsErrors},
		{statsRxTooLongDesc, s.RxTooLong},
		{statsRxTotalPacketsDesc, s.RxTotalPackets},
		{statsTxTotalPacketsDesc, s.TxTotalPackets},
		{statsTxCollisionsDesc, s.TxCollisions},
		{statsTxDeferredDesc, s.TxDeferred},
		{statsTxExcessiveCollisionsDesc, s.TxExcessiveCollisions},
		{statsTxLateCollisionsDesc, s.TxLateCollisions},
		{statsTxMultipleCollisionsDesc, s.TxMultipleCollisions},
		{statsRxFrames64Desc, s.RxFrames64},
		{statsRxFrames65To127Desc, s.RxFrames65To127},
		{statsRxFrames128To255Desc, s.RxFrames128To255},
		{statsRxFrames256To511Desc, s.RxFrames256To511},
		{statsRxFrames512To1023Desc, s.RxFrames512To1023},
		{statsRxFrames1024PlusDesc, s.RxFrames1024Plus},
		{statsRxOverrunsDesc, s.RxOverruns},
		{statsTxFcsErrorsDesc, s.TxFcsErrors},
		{statsRxJabberDesc, s.RxJabber},
		{statsRxErrorsDesc, s.RxErrors},
		{statsTxSingleCollisionsDesc, s.TxSingleCollisions},
		{statsRxFrames1024To1518Desc, s.RxFrames1024To1518},
		{statsRxFrames1519PlusDesc, s.RxFrames1519Plus},
		{statsTxFrames64Desc, s.TxFrames64},
		{statsTxFrames65To127Desc, s.TxFrames65To127},
		{statsTxFrames128To255Desc, s.TxFrames128To255},
		{statsTxFrames256To511Desc, s.TxFrames256To511},
		{statsTxFrames512To1023Desc, s.TxFrames512To1023},
		{statsTxFrames1024To1518Desc, s.TxFrames1024To1518},
		{statsTxFrames1519PlusDesc, s.TxFrames1519Plus},
		{statsRxAlignErrorsDesc, s.RxAlignErrors},
		{statsRxTotalErrorsDesc, s.RxTotalErrors},
		{statsTxExcessiveDeferredDesc, s.TxExcessiveDeferred},
		{statsTxTotalErrorsDesc, s.TxTotalErrors},
		{statsTxTooLongDesc, s.TxTooLong},
		{statsRxMacErrorsDesc, s.RxMacErrors},
		{statsTxUnderrunsDesc, s.TxUnderruns},
	}

	n := 0
	for _, f := range fields {
		if f.value != nil && len(*f.value) > n {
			n = len(*f.value)
		}
	}
	if nn := arrayLenInt64ForStats(s.TxQueue, s.TxQueueBytes); nn > n {
		n = nn
	}

	for i := 0; i < n; i++ {
		lbl := portLabel(i)
		for _, f := range fields {
			counterInt64(ch, f.desc, atInt64(f.value, i), host, lbl)
		}
		gaugeInt64(ch, statsTxQueueDesc, atInt64(s.TxQueue, i), host, lbl)
		gaugeInt64(ch, statsTxQueueBytesDesc, atInt64(s.TxQueueBytes, i), host, lbl)
	}
}

func arrayLenInt64ForStats(ps ...*[]int64) int {
	n := 0
	for _, p := range ps {
		if p != nil && len(*p) > n {
			n = len(*p)
		}
	}
	return n
}
