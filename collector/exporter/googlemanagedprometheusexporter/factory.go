// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package googlemanagedprometheusexporter // import "exporter/googlemanagedprometheusexporter"

import (
	"context"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/exporter/googlemanagedprometheusexporter/internal/metadata"
)

const (
	defaultTimeout = 12 * time.Second // Consistent with Cloud Monitoring's timeout
)

// NewFactory creates a factory for the googlemanagedprometheus exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
	)
}

// createDefaultConfig creates the default configuration for exporter.
func createDefaultConfig() component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.TimeoutConfig{Timeout: defaultTimeout},
		QueueSettings:   exporterhelper.NewDefaultQueueConfig(),
		GMPConfig: GMPConfig{
			MetricConfig: MetricConfig{
				Config: googlemanagedprometheus.DefaultConfig(),
			},
		},
	}
}

// createMetricsExporter creates a metrics exporter based on this config.
func createMetricsExporter(
	ctx context.Context,
	params exporter.Settings,
	cfg component.Config) (exporter.Metrics, error) {
	eCfg := cfg.(*Config)

	// We turn off normalization for serverless environments.
	collectorConfig, err := eCfg.GMPConfig.toCollectorConfig()
	if err != nil {
		return nil, err
	}
	collectorConfig.MetricConfig.CumulativeNormalization = false
	mExp, err := collector.NewGoogleCloudMetricsExporter(ctx, collectorConfig, params.TelemetrySettings.Logger, params.BuildInfo.Version, eCfg.TimeoutSettings.Timeout)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewMetricsExporter(
		ctx,
		params,
		cfg,
		mExp.PushMetrics,
		exporterhelper.WithShutdown(mExp.Shutdown),
		// Disable exporterhelper Timeout, since we are using a custom mechanism
		// within exporter itself
		exporterhelper.WithTimeout(exporterhelper.TimeoutConfig{Timeout: 0}),
		exporterhelper.WithQueue(eCfg.QueueSettings))
}
