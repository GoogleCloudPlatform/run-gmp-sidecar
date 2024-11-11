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

package internal // import "github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/receiver/prometheusreceiver/internal"

import (
	"context"
	"regexp"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

// appendable translates Prometheus scraping diffs into OpenTelemetry format.
type appendable struct {
	sink                 consumer.Metrics
	metricAdjuster       MetricsAdjuster
	useStartTimeMetric   bool
	trimSuffixes         bool
	startTimeMetricRegex *regexp.Regexp
	externalLabels       labels.Labels

	settings receiver.Settings
	obsrecv  *receiverhelper.ObsReport
}

// NewAppendable returns a storage.Appendable instance that emits metrics to the sink.
func NewAppendable(
	sink consumer.Metrics,
	set receiver.Settings,
	gcInterval time.Duration,
	useStartTimeMetric bool,
	startTimeMetricRegex *regexp.Regexp,
	useCreatedMetric bool,
	useCollectorStartTimeFallback bool,
	allowCumulativeResets bool,
	externalLabels labels.Labels,
	trimSuffixes bool) (storage.Appendable, error) {
	var metricAdjuster MetricsAdjuster
	if !useStartTimeMetric {
		metricAdjuster = NewInitialPointAdjuster(set.Logger, gcInterval, useCreatedMetric)
	} else {
		metricAdjuster = NewStartTimeMetricAdjuster(set.Logger, gcInterval, startTimeMetricRegex, useCollectorStartTimeFallback, allowCumulativeResets)
	}

	obsrecv, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{ReceiverID: set.ID, Transport: transport, ReceiverCreateSettings: set})
	if err != nil {
		return nil, err
	}

	return &appendable{
		sink:                 sink,
		settings:             set,
		metricAdjuster:       metricAdjuster,
		useStartTimeMetric:   useStartTimeMetric,
		trimSuffixes:         trimSuffixes,
		startTimeMetricRegex: startTimeMetricRegex,
		externalLabels:       externalLabels,
		obsrecv:              obsrecv,
	}, nil
}

func (o *appendable) Appender(ctx context.Context) storage.Appender {
	return newTransaction(ctx, o.metricAdjuster, o.sink, o.externalLabels, o.settings, o.obsrecv, o.trimSuffixes)
}
