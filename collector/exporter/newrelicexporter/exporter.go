// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package newrelicexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type newrelicExporter struct {
	cfg    *Config
	logger *zap.Logger
	client *http.Client
}

func newMetricsExporter(cfg *Config, set exporter.CreateSettings) (exporter.Metrics, error) {
	ne := &newrelicExporter{
		cfg:    cfg,
		logger: set.Logger,
		client: &http.Client{Timeout: cfg.Timeout},
	}

	return exporterhelper.NewMetricsExporter(
		context.Background(),
		set,
		cfg,
		ne.pushMetrics,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.RetrySettings),
		exporterhelper.WithTimeout(cfg.TimeoutSettings),
	)
}

func newLogsExporter(cfg *Config, set exporter.CreateSettings) (exporter.Logs, error) {
	ne := &newrelicExporter{
		cfg:    cfg,
		logger: set.Logger,
		client: &http.Client{Timeout: cfg.Timeout},
	}

	return exporterhelper.NewLogsExporter(
		context.Background(),
		set,
		cfg,
		ne.pushLogs,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.RetrySettings),
		exporterhelper.WithTimeout(cfg.TimeoutSettings),
	)
}

func newTracesExporter(cfg *Config, set exporter.CreateSettings) (exporter.Traces, error) {
	ne := &newrelicExporter{
		cfg:    cfg,
		logger: set.Logger,
		client: &http.Client{Timeout: cfg.Timeout},
	}

	return exporterhelper.NewTracesExporter(
		context.Background(),
		set,
		cfg,
		ne.pushTraces,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.RetrySettings),
		exporterhelper.WithTimeout(cfg.TimeoutSettings),
	)
}

func (ne *newrelicExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	payload, err := ne.transformMetrics(md)
	if err != nil {
		return fmt.Errorf("failed to transform metrics: %w", err)
	}

	return ne.sendToNewRelic(ctx, ne.cfg.GetMetricsURL(), payload)
}

func (ne *newrelicExporter) pushLogs(ctx context.Context, ld plog.Logs) error {
	payload, err := ne.transformLogs(ld)
	if err != nil {
		return fmt.Errorf("failed to transform logs: %w", err)
	}

	return ne.sendToNewRelic(ctx, ne.cfg.GetLogsURL(), payload)
}

func (ne *newrelicExporter) pushTraces(ctx context.Context, td ptrace.Traces) error {
	payload, err := ne.transformTraces(td)
	if err != nil {
		return fmt.Errorf("failed to transform traces: %w", err)
	}

	return ne.sendToNewRelic(ctx, ne.cfg.GetTracesURL(), payload)
}

func (ne *newrelicExporter) transformMetrics(md pmetric.Metrics) ([]byte, error) {
	var metrics []map[string]interface{}

	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)
		resourceAttrs := attributesToMap(rm.Resource().Attributes())
		
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			sm := rm.ScopeMetrics().At(j)
			
			for k := 0; k < sm.Metrics().Len(); k++ {
				metric := sm.Metrics().At(k)
				
				// Convert metric based on type
				switch metric.Type() {
				case pmetric.MetricTypeGauge:
					nrMetrics := ne.convertGauge(metric, resourceAttrs)
					metrics = append(metrics, nrMetrics...)
				case pmetric.MetricTypeSum:
					nrMetrics := ne.convertSum(metric, resourceAttrs)
					metrics = append(metrics, nrMetrics...)
				case pmetric.MetricTypeHistogram:
					nrMetrics := ne.convertHistogram(metric, resourceAttrs)
					metrics = append(metrics, nrMetrics...)
				case pmetric.MetricTypeSummary:
					nrMetrics := ne.convertSummary(metric, resourceAttrs)
					metrics = append(metrics, nrMetrics...)
				}
			}
		}
	}

	payload := map[string]interface{}{
		"metrics": metrics,
	}

	return json.Marshal(payload)
}

func (ne *newrelicExporter) convertGauge(metric pmetric.Metric, resourceAttrs map[string]interface{}) []map[string]interface{} {
	var metrics []map[string]interface{}
	
	for i := 0; i < metric.Gauge().DataPoints().Len(); i++ {
		dp := metric.Gauge().DataPoints().At(i)
		
		nrMetric := map[string]interface{}{
			"name":      metric.Name(),
			"type":      "gauge",
			"value":     getDataPointValue(dp),
			"timestamp": dp.Timestamp().AsTime().Unix(),
		}
		
		// Merge all attributes
		attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(dp.Attributes()))
		if len(attrs) > 0 {
			nrMetric["attributes"] = attrs
		}
		
		metrics = append(metrics, nrMetric)
	}
	
	return metrics
}

func (ne *newrelicExporter) convertSum(metric pmetric.Metric, resourceAttrs map[string]interface{}) []map[string]interface{} {
	var metrics []map[string]interface{}
	
	for i := 0; i < metric.Sum().DataPoints().Len(); i++ {
		dp := metric.Sum().DataPoints().At(i)
		
		metricType := "count"
		if metric.Sum().AggregationTemporality() == pmetric.AggregationTemporalityCumulative {
			metricType = "count"
		}
		
		nrMetric := map[string]interface{}{
			"name":      metric.Name(),
			"type":      metricType,
			"value":     getDataPointValue(dp),
			"timestamp": dp.Timestamp().AsTime().Unix(),
		}
		
		// Merge all attributes
		attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(dp.Attributes()))
		if len(attrs) > 0 {
			nrMetric["attributes"] = attrs
		}
		
		metrics = append(metrics, nrMetric)
	}
	
	return metrics
}

func (ne *newrelicExporter) convertHistogram(metric pmetric.Metric, resourceAttrs map[string]interface{}) []map[string]interface{} {
	var metrics []map[string]interface{}
	
	for i := 0; i < metric.Histogram().DataPoints().Len(); i++ {
		dp := metric.Histogram().DataPoints().At(i)
		
		// New Relic histogram format
		nrMetric := map[string]interface{}{
			"name":      metric.Name(),
			"type":      "distribution",
			"value": map[string]interface{}{
				"count": dp.Count(),
				"sum":   dp.Sum(),
			},
			"timestamp": dp.Timestamp().AsTime().Unix(),
		}
		
		// Add bucket boundaries and counts if available
		if dp.BucketCounts().Len() > 0 && dp.ExplicitBounds().Len() > 0 {
			buckets := make([]map[string]interface{}, 0)
			bounds := dp.ExplicitBounds()
			counts := dp.BucketCounts()
			
			for j := 0; j < counts.Len(); j++ {
				bucket := map[string]interface{}{
					"count": counts.At(j),
				}
				
				if j < bounds.Len() {
					bucket["upperBound"] = bounds.At(j)
				}
				
				buckets = append(buckets, bucket)
			}
			
			nrMetric["value"].(map[string]interface{})["buckets"] = buckets
		}
		
		// Merge all attributes
		attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(dp.Attributes()))
		if len(attrs) > 0 {
			nrMetric["attributes"] = attrs
		}
		
		metrics = append(metrics, nrMetric)
	}
	
	return metrics
}

func (ne *newrelicExporter) convertSummary(metric pmetric.Metric, resourceAttrs map[string]interface{}) []map[string]interface{} {
	var metrics []map[string]interface{}
	
	for i := 0; i < metric.Summary().DataPoints().Len(); i++ {
		dp := metric.Summary().DataPoints().At(i)
		
		// New Relic summary format
		nrMetric := map[string]interface{}{
			"name":      metric.Name(),
			"type":      "summary",
			"value": map[string]interface{}{
				"count": dp.Count(),
				"sum":   dp.Sum(),
			},
			"timestamp": dp.Timestamp().AsTime().Unix(),
		}
		
		// Add quantile values if available
		if dp.QuantileValues().Len() > 0 {
			quantiles := make([]map[string]interface{}, 0)
			
			for j := 0; j < dp.QuantileValues().Len(); j++ {
				qv := dp.QuantileValues().At(j)
				quantile := map[string]interface{}{
					"quantile": qv.Quantile(),
					"value":    qv.Value(),
				}
				quantiles = append(quantiles, quantile)
			}
			
			nrMetric["value"].(map[string]interface{})["quantiles"] = quantiles
		}
		
		// Merge all attributes
		attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(dp.Attributes()))
		if len(attrs) > 0 {
			nrMetric["attributes"] = attrs
		}
		
		metrics = append(metrics, nrMetric)
	}
	
	return metrics
}

func (ne *newrelicExporter) transformLogs(ld plog.Logs) ([]byte, error) {
	var logs []map[string]interface{}

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		rl := ld.ResourceLogs().At(i)
		resourceAttrs := attributesToMap(rl.Resource().Attributes())
		
		for j := 0; j < rl.ScopeLogs().Len(); j++ {
			sl := rl.ScopeLogs().At(j)
			
			for k := 0; k < sl.LogRecords().Len(); k++ {
				logRecord := sl.LogRecords().At(k)
				
				nrLog := map[string]interface{}{
					"timestamp": logRecord.Timestamp().AsTime().UnixMilli(),
					"message":   logRecord.Body().AsString(),
				}
				
				// Add severity
				if logRecord.SeverityText() != "" {
					nrLog["level"] = logRecord.SeverityText()
				}
				
				// Merge all attributes
				attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(logRecord.Attributes()))
				for k, v := range attrs {
					nrLog[k] = v
				}
				
				logs = append(logs, nrLog)
			}
		}
	}

	payload := map[string]interface{}{
		"logs": logs,
	}

	return json.Marshal(payload)
}

func (ne *newrelicExporter) transformTraces(td ptrace.Traces) ([]byte, error) {
	var spans []map[string]interface{}

	for i := 0; i < td.ResourceSpans().Len(); i++ {
		rs := td.ResourceSpans().At(i)
		resourceAttrs := attributesToMap(rs.Resource().Attributes())
		
		for j := 0; j < rs.ScopeSpans().Len(); j++ {
			ss := rs.ScopeSpans().At(j)
			
			for k := 0; k < ss.Spans().Len(); k++ {
				span := ss.Spans().At(k)
				
				nrSpan := map[string]interface{}{
					"id":        span.SpanID().String(),
					"trace.id":  span.TraceID().String(),
					"name":      span.Name(),
					"timestamp": span.StartTimestamp().AsTime().UnixMilli(),
					"duration":  span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Milliseconds(),
				}
				
				// Add parent span ID if present
				if !span.ParentSpanID().IsEmpty() {
					nrSpan["parent.id"] = span.ParentSpanID().String()
				}
				
				// Add span kind
				nrSpan["span.kind"] = span.Kind().String()
				
				// Merge all attributes
				attrs := ne.mergeAttributes(resourceAttrs, attributesToMap(span.Attributes()))
				for k, v := range attrs {
					nrSpan[k] = v
				}
				
				spans = append(spans, nrSpan)
			}
		}
	}

	payload := map[string]interface{}{
		"spans": spans,
	}

	return json.Marshal(payload)
}

func (ne *newrelicExporter) mergeAttributes(resourceAttrs, dataPointAttrs map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Add common attributes first
	for k, v := range ne.cfg.CommonAttributes {
		result[k] = v
	}
	
	// Add resource attributes
	for k, v := range resourceAttrs {
		result[k] = v
	}
	
	// Add data point attributes (these override resource attributes)
	for k, v := range dataPointAttrs {
		result[k] = v
	}
	
	return result
}

func (ne *newrelicExporter) sendToNewRelic(ctx context.Context, url string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", ne.cfg.APIKey)
	req.Header.Set("User-Agent", "run-gmp-sidecar-newrelic-exporter/1.0")

	resp, err := ne.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("New Relic API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	ne.logger.Debug("Successfully sent data to New Relic", zap.String("url", url), zap.Int("status", resp.StatusCode))
	return nil
}

// Helper functions

func attributesToMap(attrs pcommon.Map) map[string]interface{} {
	result := make(map[string]interface{})
	attrs.Range(func(k string, v pcommon.Value) bool {
		switch v.Type() {
		case pcommon.ValueTypeStr:
			result[k] = v.Str()
		case pcommon.ValueTypeInt:
			result[k] = v.Int()
		case pcommon.ValueTypeDouble:
			result[k] = v.Double()
		case pcommon.ValueTypeBool:
			result[k] = v.Bool()
		default:
			result[k] = v.AsString()
		}
		return true
	})
	return result
}

func getDataPointValue(dp pmetric.NumberDataPoint) interface{} {
	switch dp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return dp.IntValue()
	case pmetric.NumberDataPointValueTypeDouble:
		return dp.DoubleValue()
	default:
		return 0
	}
}
