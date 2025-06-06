// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package newrelicexporter

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for New Relic exporter.
type Config struct {
	confighttp.ClientConfig       `mapstructure:",squash"`
	exporterhelper.QueueSettings  `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings  `mapstructure:"retry_on_failure"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`

	// APIKey is the New Relic Insights Insert API key
	APIKey string `mapstructure:"api_key"`
	
	// MetricsURLOverride overrides the default New Relic metrics endpoint URL
	MetricsURLOverride string `mapstructure:"metrics_url_override"`
	
	// LogsURLOverride overrides the default New Relic logs endpoint URL  
	LogsURLOverride string `mapstructure:"logs_url_override"`
	
	// TracesURLOverride overrides the default New Relic traces endpoint URL
	TracesURLOverride string `mapstructure:"traces_url_override"`
	
	// CommonAttributes are attributes to be included with every metric
	CommonAttributes map[string]interface{} `mapstructure:"common_attributes"`
}

const (
	defaultMetricsURL = "https://metric-api.newrelic.com/metric/v1"
	defaultLogsURL    = "https://log-api.newrelic.com/log/v1" 
	defaultTracesURL  = "https://trace-api.newrelic.com/trace/v1"
)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.APIKey == "" {
		return errors.New("api_key is required")
	}

	// Validate custom URLs if provided
	if cfg.MetricsURLOverride != "" {
		if _, err := url.Parse(cfg.MetricsURLOverride); err != nil {
			return fmt.Errorf("invalid metrics_url_override: %w", err)
		}
	}
	
	if cfg.LogsURLOverride != "" {
		if _, err := url.Parse(cfg.LogsURLOverride); err != nil {
			return fmt.Errorf("invalid logs_url_override: %w", err)
		}
	}
	
	if cfg.TracesURLOverride != "" {
		if _, err := url.Parse(cfg.TracesURLOverride); err != nil {
			return fmt.Errorf("invalid traces_url_override: %w", err)
		}
	}

	return nil
}

// GetMetricsURL returns the metrics endpoint URL
func (cfg *Config) GetMetricsURL() string {
	if cfg.MetricsURLOverride != "" {
		return cfg.MetricsURLOverride
	}
	return defaultMetricsURL
}

// GetLogsURL returns the logs endpoint URL  
func (cfg *Config) GetLogsURL() string {
	if cfg.LogsURLOverride != "" {
		return cfg.LogsURLOverride
	}
	return defaultLogsURL
}

// GetTracesURL returns the traces endpoint URL
func (cfg *Config) GetTracesURL() string {
	if cfg.TracesURLOverride != "" {
		return cfg.TracesURLOverride
	}
	return defaultTracesURL
}
