package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Collector that contains the descriptors for the metrics from the app.
// Foo is a gauge with no labels. Bar is a counter with no labels.
type fooBarCollector struct {
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
}

func newFooBarCollector() *fooBarCollector {
	return &fooBarCollector{
		fooMetric: prometheus.NewDesc("foo_metric",
			"A foo event has occurred",
			nil, nil,
		),
		barMetric: prometheus.NewDesc("bar_metric",
			"A bar event has occured",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *fooBarCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.fooMetric
	ch <- collector.barMetric
}

// Collect implements required collect function for all promehteus collectors
func (collector *fooBarCollector) Collect(ch chan<- prometheus.Metric) {
	m1 := prometheus.MustNewConstMetric(collector.fooMetric, prometheus.GaugeValue, float64(rand.Int31()))
	m2 := prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, float64(rand.Int31()))
	ch <- m1
	ch <- m2
}

func entrypointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User request received!")
}

func main() {
	foo := newFooBarCollector()
	prometheus.MustRegister(foo)

	entrypointMux := http.NewServeMux()
	entrypointMux.HandleFunc("/", entrypointHandler)

	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe("localhost:8080", entrypointMux)
	}()

	http.ListenAndServe("localhost:8000", promMux)
}
