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

package collector_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/testbed/testbed"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	performanceResultsSummary testbed.TestResultsSummary = &testbed.PerformanceResults{}
)

func TestPrometheusMetrics(t *testing.T) {
	// Set up prometheus metrics sender, collector using prometheus receiver and
	// OTLP exporter pushing metrics to the testing harness (mock backend).
	col := newTestCollector(100, 100)
	tc, sender, cleanup := setupPrometheusTestCase(t, col)
	defer prometheus.Unregister(col)
	defer cleanup()
	defer sender.Stop()

	// Wait for metrics to start showing up.
	tc.WaitForN(func() bool { return tc.MockBackend.DataItemsReceived() > 0 }, 60*time.Second,
		"datapoints received")

	// Shutdown the harness so no more telemetry gets through.
	tc.Stop()

	// Assert metrics.
	metrics := tc.MockBackend.ReceivedMetrics
	require.Greater(t, len(metrics), 0)
	require.True(t, findMetricByName(metrics[0], "test_gauge"))
	require.True(t, findMetricByName(metrics[0], "test_counter"))
}

func TestShortLivedCollector(t *testing.T) {
	// Set up prometheus metrics sender, collector using prometheus receiver and
	// OTLP exporter pushing metrics to the testing harness (mock backend).
	col := newTestCollector(100, 100)
	tc, sender, cleanup := setupPrometheusTestCase(t, col)
	defer prometheus.Unregister(col)
	defer cleanup()
	defer sender.Stop()

	// Stop the agent after 1s of running.
	time.Sleep(1 * time.Second)
	tc.StopAgent()

	// Give 1s for graceful shutdown.
	time.Sleep(1 * time.Second)

	// Wait for metrics to start showing up.
	tc.WaitFor(func() bool { return tc.MockBackend.DataItemsReceived() > 0 },
		"datapoints received")

	// Shutdown the harness so no more telemetry gets through.
	tc.Stop()

	// Assert metrics.
	metrics := tc.MockBackend.ReceivedMetrics
	require.Greater(t, len(metrics), 0)
	require.True(t, findMetricByName(metrics[0], "test_gauge"))
	require.True(t, findMetricByName(metrics[0], "test_counter"))
}

func setupPrometheusTestCase(t *testing.T, collector prometheus.Collector) (tc *testbed.TestCase, sender *PrometheusDataSender, configCleanup func()) {
	var prometheusPort, otlpPort int
	prometheusPort = getAvailablePort(t, nil)
	otlpPort = getAvailablePort(t, []int{prometheusPort})

	// Set up prometheus exporter and the OTLP receiver that the test harness will use.
	sender = NewPrometheusDataSender("/metrics", prometheusPort, collector)
	receiver := testbed.NewOTLPDataReceiver(otlpPort)

	// Set up the collector.
	agentProc := testbed.NewChildProcessCollector(testbed.WithAgentExePath("../bin/rungmpcol"))
	configStr := createConfigYaml(t, sender, receiver)
	configCleanup, err := agentProc.PrepareConfig(configStr)
	require.NoError(t, err)

	// Set up the test case and data provider.
	dataProvider := NewNoopDataProvider()
	tc = testbed.NewTestCase(
		t,
		dataProvider,
		sender,
		receiver,
		agentProc,
		&testbed.PerfTestValidator{},
		performanceResultsSummary,
	)

	// Control flow for telemetry.
	tc.StartBackend()
	tc.StartAgent()
	tc.EnableRecording()
	require.NoError(t, sender.Start())

	// We bypass the load generator in this test, but make sure to increment the
	// counter since it is used in final reports.
	tc.LoadGenerator.IncDataItemsSent()

	return
}

func findMetricByName(ms pmetric.Metrics, name string) bool {
	rms := ms.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			ms := ilm.Metrics()
			for k := 0; k < ms.Len(); k++ {
				m := ms.At(k)
				if m.Name() == name {
					return true
				}
			}
		}
	}
	return false
}

// testCollector contains the descriptors for the metrics from the app.
type testCollector struct {
	gauge    *prometheus.Desc
	gaugeVal float64

	counter    *prometheus.Desc
	counterVal float64
}

func newTestCollector(gauge, counter float64) *testCollector {
	return &testCollector{
		gauge: prometheus.NewDesc("test_gauge",
			"A gauge event has occurred",
			nil, nil,
		),
		counter: prometheus.NewDesc("test_counter",
			"A counter event has occurred",
			nil, nil,
		),
		gaugeVal:   gauge,
		counterVal: counter,
	}
}

// Describe writes the descriptors to the prometheus desc channel.
func (collector *testCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.gauge
	ch <- collector.counter
}

// Collect implements required collect function for all prometheus collectors
func (collector *testCollector) Collect(ch chan<- prometheus.Metric) {
	m1 := prometheus.MustNewConstMetric(collector.gauge, prometheus.GaugeValue, collector.gaugeVal)
	m2 := prometheus.MustNewConstMetric(collector.counter, prometheus.CounterValue, collector.counterVal)
	ch <- m1
	ch <- m2
}

type PrometheusDataSender struct {
	testbed.DataSender
	consumer.Metrics

	PromCollector prometheus.Collector
	Path          string
	Port          int

	server *http.Server
}

func (p *PrometheusDataSender) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (p *PrometheusDataSender) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	return nil
}

func NewPrometheusDataSender(path string, port int, promCollector prometheus.Collector) *PrometheusDataSender {
	return &PrometheusDataSender{
		PromCollector: promCollector,
		Path:          path,
		Port:          port,

		// server is created on Start()
		server: nil,
	}
}

func (p *PrometheusDataSender) exportMetrics() *http.Server {
	prometheus.MustRegister(p.PromCollector)

	promMux := http.NewServeMux()
	promMux.Handle(p.Path, promhttp.Handler())

	exporter := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Port),
		Handler: promMux,
	}

	go func() {
		err := exporter.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("prometheus exporter closed\n")
		} else if err != nil {
			fmt.Printf("error listening on prometheus server: %s\n", err)
		}
	}()

	return exporter
}

func (p *PrometheusDataSender) Start() error {
	// Start running the prometheus exporter.
	p.server = p.exportMetrics()
	log.Printf("Listening on %d%s", p.Port, p.Path)
	return nil
}

func (p *PrometheusDataSender) Flush() {
	// no-op
}

func (p *PrometheusDataSender) GenConfigYAMLStr() string {
	// Note that this generates a receiver config for agent.
	return fmt.Sprintf(`
  prometheus:
    preserve_untyped: true
    allow_cumulative_resets: true
    use_collector_start_time_fallback: true
    use_start_time_metric: true
    config:
      scrape_configs:
      - job_name: RunMonitoring/run-gmp-sidecar
        metrics_path: "%s"
        static_configs:
        - targets:
          - 0.0.0.0:%d
`, p.Path, p.Port)
}

func (p *PrometheusDataSender) ProtocolName() string {
	return "prometheus"
}

func (p *PrometheusDataSender) GetEndpoint() net.Addr {
	return nil
}

func (p *PrometheusDataSender) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}

	return nil
}

type NoopDataProvider struct {
	dataItemsGenerated *atomic.Uint64
}

// NewFileDataProvider creates an instance of FileDataProvider which generates test data
// loaded from a file.
func NewNoopDataProvider() *NoopDataProvider {
	dp := &NoopDataProvider{
		dataItemsGenerated: &atomic.Uint64{},
	}
	return dp
}

func (dp *NoopDataProvider) SetLoadGeneratorCounters(dataItemsGenerated *atomic.Uint64) {
	dp.dataItemsGenerated = dataItemsGenerated
}

func (dp *NoopDataProvider) GenerateTraces() (ptrace.Traces, bool) {
	dp.dataItemsGenerated.Add(1)
	return ptrace.NewTraces(), false
}

func (dp *NoopDataProvider) GenerateMetrics() (pmetric.Metrics, bool) {
	dp.dataItemsGenerated.Add(1)
	return pmetric.NewMetrics(), false
}

func (dp *NoopDataProvider) GenerateLogs() (plog.Logs, bool) {
	dp.dataItemsGenerated.Add(1)
	return plog.NewLogs(), false
}

// createConfigYaml creates a collector config file that corresponds to the
// sender and receiver used in the test and returns the config file name.
func createConfigYaml(
	t *testing.T,
	sender testbed.DataSender,
	receiver testbed.DataReceiver,
) string {

	// Create a config. Note that our DataSender is used to generate a config for Collector's
	// receiver and our DataReceiver is used to generate a config for Collector's exporter.
	// This is because our DataSender sends to Collector's receiver and our DataReceiver
	// receives from Collector's exporter.

	format := `
receivers:%v
exporters:%v

service:
  pipelines:
    metrics:
      receivers: [%v]
      exporters: [%v]
`

	// Put corresponding elements into the config template to generate the final config.
	return fmt.Sprintf(
		format,
		sender.GenConfigYAMLStr(),
		receiver.GenConfigYAMLStr(),
		sender.ProtocolName(),
		receiver.ProtocolName(),
	)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getAvailablePort(t *testing.T, excludePorts []int) int {
	for {
		port, err := confgenerator.GetFreePort()
		require.NoError(t, err)
		if !contains(excludePorts, port) {
			return port
		}
	}
}
