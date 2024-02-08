// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlemanagedprometheusexporternew // import "exporter/googlemanagedprometheusexporternew"

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/exporter/googlemanagedprometheusexporternew/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	factories, err := otelcoltest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[component.Type(metadata.Type)] = factory
	cfg, err := otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, len(cfg.Exporters), 2)

	r0 := cfg.Exporters[component.NewID(component.Type(metadata.Type))].(*Config)
	assert.Equal(t, r0, factory.CreateDefaultConfig().(*Config))

	r1 := cfg.Exporters[component.NewIDWithName(component.Type(metadata.Type), "customname")].(*Config)
	assert.Equal(t, r1,
		&Config{
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 20 * time.Second,
			},
			GMPConfig: GMPConfig{
				ProjectID: "my-project",
				UserAgent: "opentelemetry-collector-contrib {{version}}",
				MetricConfig: MetricConfig{
					Config: googlemanagedprometheus.Config{
						AddMetricSuffixes: false,
						ExtraMetricsConfig: googlemanagedprometheus.ExtraMetricsConfig{
							EnableTargetInfo: false,
							EnableScopeInfo:  false,
						},
					},
					Prefix: "my-metric-domain.com",
				},
			},
			QueueSettings: exporterhelper.QueueSettings{
				Enabled:      true,
				NumConsumers: 2,
				QueueSize:    10,
			},
		})
}
