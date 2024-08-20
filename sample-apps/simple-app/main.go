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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PubSubMessage struct {
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

// Collector that contains the descriptors for the metrics from the app.
type fooBarCollector struct {
	avgDuration          *prometheus.Desc
	errRate              *prometheus.Desc
	failedScansCount     *prometheus.Desc
	maxMemUsed           *prometheus.Desc
	reportsReceivedCount *prometheus.Desc
	scannersCreatedCount *prometheus.Desc
}

func newFooBarCollector() *fooBarCollector {
	return &fooBarCollector{
		avgDuration: prometheus.NewDesc("avg_duration",
			"Average duration of the lambda.",
			nil, nil,
		),
		errRate: prometheus.NewDesc("err_rate",
			"Err rate for the lambda (count of errors/invocations).",
			nil, nil,
		),
		failedScansCount: prometheus.NewDesc("failed_scans_count",
			"Count of failed scans over the last 24 hours.",
			nil, nil,
		),
		maxMemUsed: prometheus.NewDesc("max_mem_used",
			"Max memory used by the lambda (in MB).",
			nil, nil,
		),
		reportsReceivedCount: prometheus.NewDesc("reports_received_count",
			"Count of reports received over the last 24 hours.",
			nil, nil,
		),
		scannersCreatedCount: prometheus.NewDesc("scanners_created_count",
			"Scanners created over the last 24 hours.",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *fooBarCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.avgDuration
	ch <- collector.errRate
	ch <- collector.failedScansCount
	ch <- collector.maxMemUsed
	ch <- collector.reportsReceivedCount
	ch <- collector.scannersCreatedCount
}

// Collect implements required collect function for all prometheus collectors
func (collector *fooBarCollector) Collect(ch chan<- prometheus.Metric) {
	log.Printf("entered fooBarCollector collector ")
	projectID := "ctd-pi-ads-mgmt-prod"
	// subID := "pull-ads-monitoring-bucket-prod-notification-sub"
	subID := "temp-sub-2"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()
	log.Printf("created Pub/Sub client successfully")
	sub := client.Subscription(subID)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Printf("created Pub/Sub client Subscription successfully")
	messageReceived := 0

	for messageReceived < 5 {
		log.Printf("Pub/Sub calling for receiving message")
		err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
			log.Printf("Pub/Sub msgData: %s", msg.Data)
			var jsonObject map[string]interface{}
			err = json.Unmarshal(msg.Data, &jsonObject)
			if err != nil {
				return
			}
			log.Printf("Pub/Sub jsonObject: %s", jsonObject)
			name, ok := jsonObject["name"].(string)
			if !ok || !strings.HasPrefix(name, "metrics/") {
				return
			}

			bucket, ok := jsonObject["bucket"].(string)
			if !ok {
				return
			}

			// tenantprojectid, ok := msg.Attributes["tenantprojectid"]
			// if !ok {
			// 	return
			// }
			// // log.Printf("Hello, registry is here %s!", registry)
			// groupingKey := map[string]string{
			// 	"tenantProjectId": tenantprojectid,
			// }

			storageClient, err := storage.NewClient(ctx)
			if err != nil {
				log.Fatalf("Failed to create client: %v", err)
			}
			defer storageClient.Close()

			obj := storageClient.Bucket(bucket).Object(name)

			rc, err := obj.NewReader(ctx)
			if err != nil {
				log.Fatalf("Failed to create object reader: %v", err)
			}
			defer rc.Close()

			jsondata, err := io.ReadAll(rc)
			if err != nil {
				log.Fatalf("Failed to create object content: %v", err)
			}

			log.Printf("Pub/Sub file jsondata: %s", jsondata)

			var jsonArray []map[string]interface{}
			err = json.Unmarshal(jsondata, &jsonArray)
			if err != nil {
				return
			}
			var metricFields = [6]*prometheus.Desc{collector.avgDuration, collector.errRate, collector.failedScansCount, collector.maxMemUsed, collector.reportsReceivedCount, collector.scannersCreatedCount}
			metricFieldsIndex := 0
			for _, currObject := range jsonArray {
				currObjectName, ok := currObject["name"].(string)
				if !ok || strings.HasPrefix(currObjectName, "go") || (strings.HasPrefix(currObjectName, "process")) {
					continue
				}
				currObjectType, ok := currObject["type"].(float64)
				if !ok {
					continue
				}

				currObjectMetric, ok := currObject["metric"].([]interface{})
				if !ok {
					continue
				}

				if int(currObjectType) == 1 {

					gaugeObject, ok := currObjectMetric[0].(map[string]interface{})["gauge"].(map[string]interface{})
					if !ok {
						continue
					}

					currObjectValue, ok := gaugeObject["value"].(float64)
					if !ok {
						continue
					}
					log.Printf("registering fooBarCollector collector %v ", currObjectName)
					log.Printf("values are going to channel ")
					ch <- prometheus.MustNewConstMetric(metricFields[metricFieldsIndex], prometheus.GaugeValue, currObjectValue)
					log.Printf("values parsed to channel successfully")
					log.Printf("registered fooBarCollector collector %v ", currObjectName)
					log.Printf("Pub/Sub file currObjectName: %v", currObjectName)
					log.Printf("Pub/Sub file currObjectValue: %v", currObjectValue)
					log.Printf("Pub/Sub file gaugeObject: %v", gaugeObject)
					metricFieldsIndex++
				}
			}
			log.Printf("Message is ack successfully")
			msg.Ack() // Acknowledge the message
		})

		if err != nil {
			log.Printf("error encountered while recieveing message from pub/sub")
			log.Fatal(err)
		}
	}

	log.Printf("values parsed to channel successfully")
}

func entrypointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User request received!")
}

func main() {
	foo := newFooBarCollector()
	log.Printf("registring fooBarCollector")
	prometheus.MustRegister(foo)
	log.Printf("registred fooBarCollector")

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
