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

package confgenerator

import (
	"fmt"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator/otel"
)

type AgentSelfMetrics struct {
	Version string
	Service string
	Port    int
}

func (r AgentSelfMetrics) OTelReceiverPipeline() otel.ReceiverPipeline {
	return otel.ReceiverPipeline{
		Receiver: otel.Component{
			Type: "prometheus",
			Config: map[string]interface{}{
				"config": map[string]interface{}{
					"scrape_configs": []map[string]interface{}{{
						"job_name":        "run-gmp-sidecar-self-metrics",
						"scrape_interval": "1m",
						// Explicitly set fallback_scrape_protocol to PrometheusText0.0.4.
						// The OTel Collector's config validator recursively walks Go structs at startup.
						// Since ScrapeFallbackProtocol is left as Go's zero-value ("") and its type
						// ScrapeProtocol implements a Validate() error method (with no arguments), OTel
						// automatically calls it. Since "" is not a valid scrape protocol, this triggers
						// a validation crash. Setting a valid default prevents this crash.
						"fallback_scrape_protocol": "PrometheusText0.0.4",
						"static_configs": []map[string]interface{}{{
							"targets": []string{fmt.Sprintf("0.0.0.0:%d", r.Port)},
						}},
						"metric_relabel_configs": []map[string]interface{}{{
							"source_labels": []string{"__address__"},
							"target_label":  "instance",
							"replacement":   fmt.Sprintf("%d", r.Port),
							"action":        "replace",
						}},
					}},
				},
			},
		},
		Processors: []otel.Component{
			otel.MetricsFilter(
				"include",
				"strict",
				"otelcol_process_uptime_seconds_total",
				"otelcol_process_memory_rss_bytes",
				"grpc_client_attempt_duration_seconds",
				"googlecloudmonitoring_point_count_total",
			),
			otel.Transform("metric", "metric",
				// create new count metric from histogram metric
				otel.ExtractCountMetric(true, "grpc_client_attempt_duration_seconds"),
			),
			otel.MetricsTransform(
				otel.RenameMetric("otelcol_process_uptime_seconds_total", "agent/uptime",
					// change data type from double -> int64
					otel.ToggleScalarDataType,
					otel.AddLabel("version", r.Version),
					// remove service.version label
					otel.AggregateLabels("sum", "version"),
				),
				otel.RenameMetric("otelcol_process_memory_rss_bytes", "agent/memory_usage",
					// remove service.version label
					otel.AggregateLabels("sum"),
				),
				otel.RenameMetric("grpc_client_attempt_duration_seconds_count", "agent/api_request_count",
					// TODO: below is proposed new configuration for the metrics transform processor
					// ignore any non "google.monitoring" RPCs (note there won't be any other RPCs for now)
					// - action: select_label_values
					//   label: grpc_client_method
					//   value_regexp: ^google\.monitoring
					otel.RenameLabel("grpc_client_status", "state"),
					// delete grpc_client_method dimension & service.version label, retaining only state
					otel.AggregateLabels("sum", "state"),
				),
				otel.RenameMetric("googlecloudmonitoring_point_count_total", "agent/monitoring/point_count",
					// change data type from double -> int64
					otel.ToggleScalarDataType,
					// Remove service.version label
					otel.AggregateLabels("sum", "status"),
				),
			),
			// Add appropriate resource and metric labels.
			otel.GCPResourceDetector(),
			otel.TransformationMetrics(
				otel.AddMetricLabel("namespace", r.Service),
				otel.AddMetricLabel("cluster", "__run__"),
				otel.PrefixResourceAttribute("service.instance.id", "faas.instance", ":"),
			),
			otel.GroupByGMPAttrs(),
		},
	}
}

// intentionally not registered as a component because this is not created by users
