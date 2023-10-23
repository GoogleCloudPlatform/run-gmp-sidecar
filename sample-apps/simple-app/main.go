// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"time"

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

// Collect implements required collect function for all prometheus collectors
func (collector *fooBarCollector) Collect(ch chan<- prometheus.Metric) {
	m1 := prometheus.MustNewConstMetric(collector.fooMetric, prometheus.GaugeValue, float64(time.Now().Unix()))
	m2 := prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, float64(time.Now().Unix()))
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
	entrypointMux.HandleFunc("/startup", entrypointHandler)
	entrypointMux.HandleFunc("/liveness", entrypointHandler)

	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":8000", entrypointMux)
	}()

	http.ListenAndServe(":8080", promMux)
}
