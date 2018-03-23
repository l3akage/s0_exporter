package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "s0_"

var (
	wattDesc *prometheus.Desc
	timeDesc *prometheus.Desc

	lastTime = time.Now()
)

func init() {
	wattDesc = prometheus.NewDesc(prefix+"power", "Power used between scrapes (in watt)", nil, nil)
	timeDesc = prometheus.NewDesc(prefix+"time", "Time since last scrape (in seconds)", nil, nil)
}

type s0Collector struct {
}

func (c s0Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- wattDesc
}

func (c s0Collector) Collect(ch chan<- prometheus.Metric) {
	delta := time.Since(lastTime)
	counter := CounterReset()
	value := (counter / float64(*wattPerPulse)) * 1000 * delta.Seconds()
	ch <- prometheus.MustNewConstMetric(wattDesc, prometheus.GaugeValue, float64(value))
	ch <- prometheus.MustNewConstMetric(timeDesc, prometheus.GaugeValue, float64(delta.Seconds()))
	lastTime = time.Now()
}
