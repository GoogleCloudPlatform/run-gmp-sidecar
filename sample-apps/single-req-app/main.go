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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Collector that contains the descriptors for the metrics from the app.
type shortLivedCollector struct {
	gauge   *prometheus.Desc
	counter *prometheus.Desc
}

// Store whether or not we've served our request yet.
var handledSingleRequest bool

// Create channel to listen for signals.
var signalChan chan (os.Signal) = make(chan os.Signal, 1)

func newShortLivedCollector() *shortLivedCollector {
	return &shortLivedCollector{
		gauge: prometheus.NewDesc("gauge",
			"A gauge event has occurred",
			nil, nil,
		),
		counter: prometheus.NewDesc("counter",
			"A counter event has occured",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *shortLivedCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.gauge
	ch <- collector.counter
}

// Collect implements required collect function for all prometheus collectors
func (collector *shortLivedCollector) Collect(ch chan<- prometheus.Metric) {
	m1 := prometheus.MustNewConstMetric(collector.gauge, prometheus.GaugeValue, float64(time.Now().Unix()))
	m2 := prometheus.MustNewConstMetric(collector.counter, prometheus.CounterValue, float64(time.Now().Unix()))
	ch <- m1
	ch <- m2
}

func entrypointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User request received!")
	log.Printf("single request app: handled one and only request")
	handledSingleRequest = true
}

func livenessProbeHandler(w http.ResponseWriter, r *http.Request) {
	if handledSingleRequest {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte("412 - Service has turned down"))
		log.Printf("single request app: handled liveness probe. Returned that it is not live.")
		return
	}
	fmt.Fprintln(w, "Liveness probe received and was successful.")
	log.Printf("single request app: handled liveness probe. Liveness succesful.")
}

func startupProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Startup probe received")
	log.Printf("single request app: handled liveness probe. Liveness succesful.")
}

func main() {
	foo := newShortLivedCollector()
	prometheus.MustRegister(foo)

	entrypointMux := http.NewServeMux()
	entrypointMux.HandleFunc("/", entrypointHandler)
	entrypointMux.HandleFunc("/startup", startupProbeHandler)
	entrypointMux.HandleFunc("/liveness", livenessProbeHandler)

	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())

	mainSrv := http.Server{
		Addr:    ":8000",
		Handler: entrypointMux,
	}
	promSrv := http.Server{
		Addr:    ":8080",
		Handler: promMux,
	}

	go func() {
		mainSrv.ListenAndServe()
	}()

	go func() {
		promSrv.ListenAndServe()
	}()

	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// Receive output from signalChan.
	sig := <-signalChan
	log.Printf("single request app: %s signal caught", sig)

	// Timeout if waiting for connections to return idle.
	mainCtx, mainCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer mainCancel()
	if err := mainSrv.Shutdown(mainCtx); err != nil {
		log.Printf("single request app: server shutdown failed: %+v", err)
	}

	promCtx, promCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer promCancel()
	if err := promSrv.Shutdown(promCtx); err != nil {
		log.Printf("single request app: prom server shutdown failed: %+v", err)
	}

	log.Printf("single request app: shutdown complete")
}
