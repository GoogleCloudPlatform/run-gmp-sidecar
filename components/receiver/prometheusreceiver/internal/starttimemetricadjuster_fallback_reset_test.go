// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestStartTimeMetricFallback(t *testing.T) {
	const startTime = pcommon.Timestamp(123 * 1e9)
	const currentTime = pcommon.Timestamp(126 * 1e9)
	mockStartTime := time.Now().Add(-10 * time.Hour)
	mockStartTimeSeconds := float64(mockStartTime.Unix())

	tests := []struct {
		name                 string
		inputs               pmetric.Metrics
		startTimeMetricRegex *regexp.Regexp
		expectedStartTime    pcommon.Timestamp
		expectedErr          error
	}{
		{
			name: "regexp_match_sum_metric_fallback",
			inputs: metrics(
				sumMetric("test_sum_metric", doublePoint(nil, startTime, currentTime, 16)),
				histogramMetric("test_histogram_metric", histogramPoint(nil, startTime, currentTime, []float64{1, 2}, []uint64{2, 3, 4})),
				summaryMetric("test_summary_metric", summaryPoint(nil, startTime, currentTime, 10, 100, []float64{10, 50, 90}, []float64{9, 15, 48})),
			),
			startTimeMetricRegex: regexp.MustCompile("^.*_process_start_time_seconds$"),
			expectedStartTime:    timestampFromFloat64(mockStartTimeSeconds),
		},
		{
			name: "match_default_sum_start_time_metric_fallback",
			inputs: metrics(
				sumMetric("test_sum_metric", doublePoint(nil, startTime, currentTime, 16)),
				histogramMetric("test_histogram_metric", histogramPoint(nil, startTime, currentTime, []float64{1, 2}, []uint64{2, 3, 4})),
				summaryMetric("test_summary_metric", summaryPoint(nil, startTime, currentTime, 10, 100, []float64{10, 50, 90}, []float64{9, 15, 48})),
			),
			expectedStartTime: timestampFromFloat64(mockStartTimeSeconds),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcInterval := 10 * time.Millisecond
			stma := NewStartTimeMetricAdjuster(zap.NewNop(), gcInterval, tt.startTimeMetricRegex, true, false)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, stma.AdjustMetrics(tt.inputs), tt.expectedErr)
				return
			}

			// Make sure the right adjuster is used and one that has the fallback time set.
			metricAdjuster, ok := stma.(*startTimeMetricAdjuster)
			assert.Equal(t, ok, true)
			assert.NotNil(t, metricAdjuster.fallbackStartTime)

			// To test that the adjuster is using the fallback correctly, override the fallback time to use
			// directly.
			metricAdjuster.fallbackStartTime = &mockStartTime

			assert.NoError(t, stma.AdjustMetrics(tt.inputs))
			for i := 0; i < tt.inputs.ResourceMetrics().Len(); i++ {
				rm := tt.inputs.ResourceMetrics().At(i)
				for j := 0; j < rm.ScopeMetrics().Len(); j++ {
					ilm := rm.ScopeMetrics().At(j)
					for k := 0; k < ilm.Metrics().Len(); k++ {
						metric := ilm.Metrics().At(k)
						switch metric.Type() {
						case pmetric.MetricTypeSum:
							dps := metric.Sum().DataPoints()
							for l := 0; l < dps.Len(); l++ {
								assert.Equal(t, tt.expectedStartTime, dps.At(l).StartTimestamp())
							}
						case pmetric.MetricTypeSummary:
							dps := metric.Summary().DataPoints()
							for l := 0; l < dps.Len(); l++ {
								assert.Equal(t, tt.expectedStartTime, dps.At(l).StartTimestamp())
							}
						case pmetric.MetricTypeHistogram:
							dps := metric.Histogram().DataPoints()
							for l := 0; l < dps.Len(); l++ {
								assert.Equal(t, tt.expectedStartTime, dps.At(l).StartTimestamp())
							}
						}
					}
				}
			}
		})
	}
}

func TestSumWithFallBackAndReset(t *testing.T) {
	mockStartTime := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	mockTimestamp := pcommon.NewTimestampFromTime(mockStartTime)
	t1 := pcommon.Timestamp(126 * 1e9)
	t2 := pcommon.NewTimestampFromTime(t1.AsTime().Add(1 * time.Hour))
	t3 := pcommon.NewTimestampFromTime(t2.AsTime().Add(1 * time.Hour))
	t4 := pcommon.NewTimestampFromTime(t3.AsTime().Add(1 * time.Hour))
	t5 := pcommon.NewTimestampFromTime(t4.AsTime().Add(1 * time.Hour))
	script := []*metricsAdjusterTest{
		{
			description: "Sum: round 1 - initial instance, start time is established",
			metrics:     metrics(sumMetric("test_sum", doublePoint(nil, t1, t1, 44))),
			adjusted:    metrics(sumMetric("test_sum", doublePoint(nil, mockTimestamp, t1, 44))),
		},
		{
			description: "Sum: round 2 - instance adjusted based on round 1",
			metrics:     metrics(sumMetric("test_sum", doublePoint(nil, t2, t2, 66))),
			adjusted:    metrics(sumMetric("test_sum", doublePoint(nil, mockTimestamp, t2, 66))),
		},
		{
			description: "Sum: round 3 - instance reset (value less than previous value), start time is reset",
			metrics:     metrics(sumMetric("test_sum", doublePoint(nil, t3, t3, 55))),
			adjusted:    metrics(sumMetric("test_sum", doublePoint(nil, t3, t3, 55))),
		},
		{
			description: "Sum: round 4 - instance adjusted based on round 3",
			metrics:     metrics(sumMetric("test_sum", doublePoint(nil, t4, t4, 72))),
			adjusted:    metrics(sumMetric("test_sum", doublePoint(nil, t3, t4, 72))),
		},
		{
			description: "Sum: round 5 - instance adjusted based on round 4",
			metrics:     metrics(sumMetric("test_sum", doublePoint(nil, t5, t5, 72))),
			adjusted:    metrics(sumMetric("test_sum", doublePoint(nil, t3, t5, 72))),
		},
	}
	gcInterval := 10 * time.Millisecond
	stma := NewStartTimeMetricAdjuster(zap.NewNop(), gcInterval, nil, true, true)

	// Make sure the right adjuster is used and one that has the fallback time set.
	metricAdjuster, ok := stma.(*startTimeMetricAdjuster)
	assert.Equal(t, ok, true)
	assert.NotNil(t, metricAdjuster.fallbackStartTime)

	// To test that the adjuster is using the fallback correctly, override the fallback time to use
	// directly.
	metricAdjuster.fallbackStartTime = &mockStartTime
	runScript(t, stma, "job", "0", script)
}
func TestGaugeWithFallbackAndReset(t *testing.T) {
	mockStartTime := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	t1 := pcommon.Timestamp(126 * 1e9)
	t2 := pcommon.NewTimestampFromTime(t1.AsTime().Add(1 * time.Hour))
	t3 := pcommon.NewTimestampFromTime(t2.AsTime().Add(1 * time.Hour))
	script := []*metricsAdjusterTest{
		{
			description: "Gauge: round 1 - gauge not adjusted",
			metrics:     metrics(gaugeMetric("test_gauge", doublePoint(nil, t1, t1, 44))),
			adjusted:    metrics(gaugeMetric("test_gauge", doublePoint(nil, t1, t1, 44))),
		},
		{
			description: "Gauge: round 2 - gauge not adjusted",
			metrics:     metrics(gaugeMetric("test_gauge", doublePoint(nil, t2, t2, 66))),
			adjusted:    metrics(gaugeMetric("test_gauge", doublePoint(nil, t2, t2, 66))),
		},
		{
			description: "Gauge: round 3 - value less than previous value - gauge is not adjusted",
			metrics:     metrics(gaugeMetric("test_gauge", doublePoint(nil, t3, t3, 55))),
			adjusted:    metrics(gaugeMetric("test_gauge", doublePoint(nil, t3, t3, 55))),
		},
	}
	gcInterval := 10 * time.Millisecond
	stma := NewStartTimeMetricAdjuster(zap.NewNop(), gcInterval, nil, true, true)

	// Make sure the right adjuster is used and one that has the fallback time set.
	metricAdjuster, ok := stma.(*startTimeMetricAdjuster)
	assert.Equal(t, ok, true)
	assert.NotNil(t, metricAdjuster.fallbackStartTime)

	// To test that the adjuster is using the fallback correctly, override the fallback time to use
	// directly.
	metricAdjuster.fallbackStartTime = &mockStartTime
	runScript(t, stma, "job", "0", script)
}

func TestSummaryFallBackAndReset(t *testing.T) {
	mockStartTime := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	mockTimestamp := pcommon.NewTimestampFromTime(mockStartTime)
	t1 := pcommon.Timestamp(126 * 1e9)
	t2 := pcommon.NewTimestampFromTime(t1.AsTime().Add(1 * time.Hour))
	t3 := pcommon.NewTimestampFromTime(t2.AsTime().Add(1 * time.Hour))
	t4 := pcommon.NewTimestampFromTime(t3.AsTime().Add(1 * time.Hour))
	script := []*metricsAdjusterTest{
		{
			description: "Summary: round 1 - initial instance, start time is established",
			metrics: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t1, t1, 10, 40, percent0, []float64{1, 5, 8})),
			),
			adjusted: metrics(
				summaryMetric("test_summary", summaryPoint(nil, mockTimestamp, t1, 10, 40, percent0, []float64{1, 5, 8})),
			),
		},
		{
			description: "Summary: round 2 - instance adjusted based on round 1",
			metrics: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t2, t2, 15, 70, percent0, []float64{7, 44, 9})),
			),
			adjusted: metrics(
				summaryMetric("test_summary", summaryPoint(nil, mockTimestamp, t2, 15, 70, percent0, []float64{7, 44, 9})),
			),
		},
		{
			description: "Summary: round 3 - instance reset (count less than previous), start time is reset",
			metrics: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t3, t3, 12, 66, percent0, []float64{3, 22, 5})),
			),
			adjusted: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t3, t3, 12, 66, percent0, []float64{3, 22, 5})),
			),
		},
		{
			description: "Summary: round 4 - instance adjusted based on round 3",
			metrics: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t4, t4, 14, 96, percent0, []float64{9, 47, 8})),
			),
			adjusted: metrics(
				summaryMetric("test_summary", summaryPoint(nil, t3, t4, 14, 96, percent0, []float64{9, 47, 8})),
			),
		},
	}

	gcInterval := 10 * time.Millisecond
	stma := NewStartTimeMetricAdjuster(zap.NewNop(), gcInterval, nil, true, true)

	// Make sure the right adjuster is used and one that has the fallback time set.
	metricAdjuster, ok := stma.(*startTimeMetricAdjuster)
	assert.Equal(t, ok, true)
	assert.NotNil(t, metricAdjuster.fallbackStartTime)

	// To test that the adjuster is using the fallback correctly, override the fallback time to use
	// directly.
	metricAdjuster.fallbackStartTime = &mockStartTime
	runScript(t, stma, "job", "0", script)
}

func TestHistogramFallBackAndReset(t *testing.T) {
	mockStartTime := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	mockTimestamp := pcommon.NewTimestampFromTime(mockStartTime)
	t1 := pcommon.Timestamp(126 * 1e9)
	t2 := pcommon.NewTimestampFromTime(t1.AsTime().Add(1 * time.Hour))
	t3 := pcommon.NewTimestampFromTime(t2.AsTime().Add(1 * time.Hour))
	t4 := pcommon.NewTimestampFromTime(t3.AsTime().Add(1 * time.Hour))
	script := []*metricsAdjusterTest{
		{
			description: "Histogram: round 1 - initial instance, start time is established",
			metrics:     metrics(histogramMetric("test_histogram", histogramPoint(nil, t1, t1, bounds0, []uint64{4, 2, 3, 7}))),
			adjusted:    metrics(histogramMetric("test_histogram", histogramPoint(nil, mockTimestamp, t1, bounds0, []uint64{4, 2, 3, 7}))),
		}, {
			description: "Histogram: round 2 - instance adjusted based on round 1",
			metrics:     metrics(histogramMetric("test_histogram", histogramPoint(nil, t2, t2, bounds0, []uint64{6, 3, 4, 8}))),
			adjusted:    metrics(histogramMetric("test_histogram", histogramPoint(nil, mockTimestamp, t2, bounds0, []uint64{6, 3, 4, 8}))),
		}, {
			description: "Histogram: round 3 - instance reset (value less than previous value), start time is reset",
			metrics:     metrics(histogramMetric("test_histogram", histogramPoint(nil, t3, t3, bounds0, []uint64{5, 3, 2, 7}))),
			adjusted:    metrics(histogramMetric("test_histogram", histogramPoint(nil, t3, t3, bounds0, []uint64{5, 3, 2, 7}))),
		}, {
			description: "Histogram: round 4 - instance adjusted based on round 3",
			metrics:     metrics(histogramMetric("test_histogram", histogramPoint(nil, t4, t4, bounds0, []uint64{7, 4, 2, 12}))),
			adjusted:    metrics(histogramMetric("test_histogram", histogramPoint(nil, t3, t4, bounds0, []uint64{7, 4, 2, 12}))),
		},
	}
	gcInterval := 10 * time.Millisecond
	stma := NewStartTimeMetricAdjuster(zap.NewNop(), gcInterval, nil, true, true)

	// Make sure the right adjuster is used and one that has the fallback time set.
	metricAdjuster, ok := stma.(*startTimeMetricAdjuster)
	assert.Equal(t, ok, true)
	assert.NotNil(t, metricAdjuster.fallbackStartTime)

	// To test that the adjuster is using the fallback correctly, override the fallback time to use
	// directly.
	metricAdjuster.fallbackStartTime = &mockStartTime
	runScript(t, stma, "job", "0", script)
}
