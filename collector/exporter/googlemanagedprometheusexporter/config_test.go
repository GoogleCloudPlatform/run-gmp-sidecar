// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlemanagedprometheusexporter // import "github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/exporter/googlemanagedprometheusexporter"

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/exporter/googlemanagedprometheusexporter/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	factories, err := otelcoltest.NopFactories()
	assert.Nil(t, err)

	factory := NewFactory()
	factories.Exporters[metadata.Type] = factory
	cfg, err := otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, len(cfg.Exporters), 2)

	r0 := cfg.Exporters[component.NewID(metadata.Type)].(*Config)
	assert.Equal(t, r0, factory.CreateDefaultConfig().(*Config))

	r1 := cfg.Exporters[component.NewIDWithName(metadata.Type, "customname")].(*Config)
	expectedCfg := factory.CreateDefaultConfig().(*Config)
	expectedCfg.TimeoutSettings.Timeout = 20 * time.Second
	expectedCfg.ProjectID = "my-project"
	expectedCfg.UserAgent = "opentelemetry-collector-contrib {{version}}"
	expectedCfg.MetricConfig.Prefix = "my-metric-domain.com"
	expectedCfg.MetricConfig.Config.AddMetricSuffixes = false
	expectedCfg.MetricConfig.Config.ExtraMetricsConfig.EnableTargetInfo = false
	expectedCfg.MetricConfig.Config.ExtraMetricsConfig.EnableScopeInfo = false
	expectedCfg.QueueSettings.NumConsumers = 2
	expectedCfg.QueueSettings.QueueSize = 10

	assert.Equal(t, expectedCfg, r1)
}
