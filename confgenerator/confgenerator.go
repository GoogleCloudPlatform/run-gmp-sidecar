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
	"context"
	"log"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator/otel"
)

// GenerateOtelConfig generates the complete collector config including the agent self metrics.
func (rc *RunMonitoringConfig) GenerateOtelConfig(ctx context.Context, selfMetricsPort int) (string, error) {
	userAgent, _ := UserAgent("Google-Cloud-Run-GMP-Sidecar", "run-gmp", Version)
	metricVersionLabel, _ := VersionLabel("run-gmp-sidecar")
	receiverPipelines := make(map[string]otel.ReceiverPipeline)
	sidecarPipeline, err := rc.OTelReceiverPipeline()
	if err != nil {
		return "", err
	}
	receiverPipelines["application-metrics"] = *sidecarPipeline
	log.Printf("confgenerator: using port %d for self metrics", selfMetricsPort)

	receiverPipelines["run-gmp-self-metrics"] = AgentSelfMetrics{
		Version: metricVersionLabel,
		Port:    selfMetricsPort,
		Service: rc.Env.Service,
	}.OTelReceiverPipeline()

	otelConfig, err := otel.ModularConfig{
		ReceiverPipelines: receiverPipelines,
		Exporter:          newRelicExporter(),
		SelfMetricsPort:   selfMetricsPort,
	}.Generate()
	if err != nil {
		return "", err
	}
	return otelConfig, nil
}

func googleManagedPrometheusExporter(userAgent string) otel.Component {
	return otel.Component{
		Type: "googlemanagedprometheus",
		Config: map[string]interface{}{
			"user_agent": userAgent,
			// The exporter has the config option addMetricSuffixes with default value true. It will add Prometheus
			// style suffixes to metric names, e.g., `_total` for a counter; set to false to collect metrics as is
			"metric": map[string]interface{}{
				"add_metric_suffixes": false,
			},
		},
	}
}

func newRelicExporter() otel.Component {
	return otel.Component{
		Type: "newrelic",
		Config: map[string]interface{}{
			"api_key": "${NEW_RELIC_API_KEY}",
			"common_attributes": map[string]interface{}{
				"service.name": "${SERVICE_NAME:-run-gmp-sidecar}",
				"service.version": "${SERVICE_VERSION:-1.0.0}",
				"deployment.environment": "${ENVIRONMENT:-production}",
			},
		},
	}
}
