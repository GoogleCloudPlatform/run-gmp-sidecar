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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator/otel"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/model/relabel"

	yaml "github.com/goccy/go-yaml"
	prommodel "github.com/prometheus/common/model"
	promconfig "github.com/prometheus/prometheus/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RunMonitoringConfig struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`
	Spec              RunMonitoringSpec `yaml:"spec"`

	Env             *CloudRunEnvironment
	SelfMetricsPort int
}

// RunMonitoringSpec contains specification parameters for RunMonitoring.
type RunMonitoringSpec struct {
	// The endpoints to scrape on the selected pods.
	Endpoints []ScrapeEndpoint `yaml:"endpoints"`
	// Labels to add to the Prometheus target for discovered endpoints.
	TargetLabels RunTargetLabels `yaml:"targetLabels,omitempty"`
	// Limits to apply at scrape time.
	Limits *ScrapeLimits `yaml:"limits,omitempty"`
}

// RunTargetLabels specifies the additional metadata about the target
// users can add to their metric. Allowed options are {service, revision
// , configuration}. If not specified, the sidecar defaults to adding all
// of them to every metric.
type RunTargetLabels struct {
	Metadata *[]string `yaml:"metadata,omitempty"`
}

// ScrapeEndpoint specifies a Prometheus metrics endpoint to scrape.
type ScrapeEndpoint struct {
	// Name or number of the port to scrape.
	Port string `yaml:"port"`
	// Protocol scheme to use to scrape.
	Scheme string `yaml:"scheme,omitempty"`
	// HTTP path to scrape metrics from. Defaults to "/metrics".
	Path string `yaml:"path,omitempty"`
	// HTTP GET params to use when scraping.
	Params map[string][]string `yaml:"params,omitempty"`
	// Proxy URL to scrape through. Encoded passwords are not supported.
	ProxyURL string `yaml:"proxyUrl,omitempty"`
	// Interval at which to scrape metrics. Must be a valid Prometheus duration.
	Interval string `yaml:"interval,omitempty"`
	// Timeout for metrics scrapes. Must be a valid Prometheus duration.
	// Must not be larger then the scrape interval.
	Timeout string `yaml:"timeout,omitempty"`
	// Relabeling rules for metrics scraped from this endpoint. Relabeling rules
	// that override protected target labels (project_id, location, cluster,
	// namespace, job, cloud_run_instance, or __address__) are not permitted.
	MetricRelabeling []RelabelingRule `yaml:"metricRelabeling,omitempty"`
}

type RelabelingRule struct {
	// The source labels select values from existing labels. Their content is concatenated
	// using the configured separator and matched against the configured regular expression
	// for the replace, keep, and drop actions.
	SourceLabels []string `yaml:"sourceLabels,omitempty"`
	// Separator placed between concatenated source label values. Defaults to ';'.
	Separator string `yaml:"separator,omitempty"`
	// Label to which the resulting value is written in a replace action.
	// It is mandatory for replace actions. Regex capture groups are available.
	TargetLabel string `yaml:"targetLabel,omitempty"`
	// Regular expression against which the extracted value is matched. Defaults to '(.*)'.
	Regex string `yaml:"regex,omitempty"`
	// Modulus to take of the hash of the source label values.
	Modulus uint64 `yaml:"modulus,omitempty"`
	// Replacement value against which a regex replace is performed if the
	// regular expression matches. Regex capture groups are available. Defaults to '$1'.
	Replacement string `yaml:"replacement,omitempty"`
	// Action to perform based on regex matching. Defaults to 'replace'.
	Action string `yaml:"action,omitempty"`
}

type ScrapeLimits struct {
	// Maximum number of samples accepted within a single scrape.
	// Uses Prometheus default if left unspecified.
	Samples uint64 `yaml:"samples,omitempty"`
	// Maximum number of labels accepted for a single sample.
	// Uses Prometheus default if left unspecified.
	Labels uint64 `yaml:"labels,omitempty"`
	// Maximum label name length.
	// Uses Prometheus default if left unspecified.
	LabelNameLength uint64 `yaml:"labelNameLength,omitempty"`
	// Maximum label value length.
	// Uses Prometheus default if left unspecified.
	LabelValueLength uint64 `yaml:"labelValueLength,omitempty"`
}

var allowedTargetMetadata = []string{"revision", "service", "configuration"}

const kind = "RunMonitoring"
const apiVersion = "monitoring.googleapis.com/v1beta"

// DefaultRunMonitoringConfig creates a config that will be used by default if
// no user config (or an empty one) is found. It scrapes the default location of
// 0.0.0.0:8080/metrics for prometheus metrics.
func DefaultRunMonitoringConfig() *RunMonitoringConfig {
	return &RunMonitoringConfig{
		metav1.TypeMeta{
			Kind:       kind,
			APIVersion: apiVersion,
		},
		metav1.ObjectMeta{
			Name: "run-gmp-sidecar",
		},
		RunMonitoringSpec{
			Endpoints: []ScrapeEndpoint{
				{
					Port:     "8080",
					Path:     "/metrics",
					Interval: "60s",
				},
			},
			TargetLabels: RunTargetLabels{Metadata: &allowedTargetMetadata},
		},
		nil,
		0, /* dynamic port selection for self metrics */
	}
}

// ReadConfigFromFile reads the user config file and returns a RunMonitoringConfig.
// If the user config file does not exist, or is empty - it returns the default
// RunMonitoringConfig.
func ReadConfigFromFile(ctx context.Context, path string) (*RunMonitoringConfig, error) {
	config := DefaultRunMonitoringConfig()

	// Fetch metadata from the available environment variables.
	config.Env = fetchMetadata()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, fmt.Errorf("failed to retrieve the user config file %q: %w", path, err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the user config over the default config. If some options are unspecified
	// the collector uses the default settings for those options. For example, if not specified
	// targetLabels is set to {"revision", "service", "configuration"}
	if err := yaml.UnmarshalContext(ctx, data, config, yaml.Strict()); err != nil {
		return nil, err
	}

	// Validate the RunMonitoring config
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

// OTelReceiverPipeline creates the appropriate OTel pipeline translated from the
// RunMonitoringConfig.
func (rc *RunMonitoringConfig) OTelReceiverPipeline() (*otel.ReceiverPipeline, error) {
	scrapeConfig, err := rc.scrapeConfigs()
	if err != nil {
		return nil, err
	}
	return &otel.ReceiverPipeline{
		Receiver: otel.Component{
			Type: "prometheus",
			Config: map[string]interface{}{
				"preserve_untyped":                  true,
				"use_start_time_metric":             true,
				"use_collector_start_time_fallback": true,
				"allow_cumulative_resets":           true,
				"config": map[string]interface{}{
					"scrape_configs": scrapeConfig,
				},
			},
		},
		Processors: []otel.Component{otel.GroupByGMPAttrs()},
	}, nil
}

// Validate validates the RunMonitoring config.
func (rc *RunMonitoringConfig) Validate() error {
	if rc.APIVersion != apiVersion {
		return fmt.Errorf("apiVersion must be %s", apiVersion)
	}
	if rc.Kind != kind {
		return fmt.Errorf("kind must be %s", kind)
	}

	return nil
}

// scrapeConfigs converts the given RunMonitoringConfig to an equivalent set of Prometheus ScrapeConfigs.
func (rc *RunMonitoringConfig) scrapeConfigs() (res []*promconfig.ScrapeConfig, err error) {
	for i := range rc.Spec.Endpoints {
		c, err := rc.endpointScrapeConfig(i)
		if err != nil {
			return nil, fmt.Errorf("invalid definition for endpoint with index %d: %w", i, err)
		}
		res = append(res, c)
	}
	return res, nil
}

// endpointScrapeConfig creates a scrape config for the endpoint specified.
func (rc *RunMonitoringConfig) endpointScrapeConfig(index int) (*promconfig.ScrapeConfig, error) {
	metadataLabels := map[string]struct{}{}
	if rc.Spec.TargetLabels.Metadata != nil {
		for _, l := range *rc.Spec.TargetLabels.Metadata {
			if !contains(allowedTargetMetadata, l) {
				return nil, fmt.Errorf("metadata label %q not allowed, must be one of %v", l, allowedTargetMetadata)
			}
			metadataLabels[l] = struct{}{}
		}
	}
	relabelCfgs := relabelingsForMetadata(metadataLabels, rc.Env)
	return endpointScrapeConfig(
		fmt.Sprintf("RunMonitoring/%s", rc.Name),
		rc.Spec.Endpoints[index],
		relabelCfgs,
		rc.Spec.Limits,
		rc.Env,
	)
}

func relabelingsForMetadata(keys map[string]struct{}, env *CloudRunEnvironment) (res []*relabel.Config) {
	if env == nil {
		return
	}

	if _, ok := keys["service"]; ok {
		res = append(res, &relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			Replacement:  env.Service,
			TargetLabel:  "cloud_run_service",
		})
	}
	if _, ok := keys["revision"]; ok {
		res = append(res, &relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			Replacement:  env.Revision,
			TargetLabel:  "cloud_run_revision",
		})
	}
	if _, ok := keys["configuration"]; ok {
		res = append(res, &relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			Replacement:  env.Configuration,
			TargetLabel:  "cloud_run_configuration",
		})
	}
	return res
}

func endpointScrapeConfig(id string, ep ScrapeEndpoint, relabelCfgs []*relabel.Config, limits *ScrapeLimits, env *CloudRunEnvironment) (*promconfig.ScrapeConfig, error) {
	if env == nil {
		return nil, fmt.Errorf("metadata from Cloud Run was not found")
	}
	labelSet := make(map[prommodel.LabelName]prommodel.LabelValue)
	labelSet[prommodel.AddressLabel] = prommodel.LabelValue("0.0.0.0:" + ep.Port)
	discoveryCfgs := discovery.Configs{
		discovery.StaticConfig{
			&targetgroup.Group{Targets: []prommodel.LabelSet{labelSet}},
		},
	}
	relabelCfgs = append(relabelCfgs,
		&relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			TargetLabel:  "cluster",
			Replacement:  "__run__",
		},
		&relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			TargetLabel:  "namespace",
			Replacement:  env.Service,
		},
		// The `instance` label will be <faas.id>:<port> in the final metric.
		// But since <faas.id> is unavailable until the gcp resource detector
		// runs later in the pipeline we just populate the port for now.
		//
		// See the usage of PrefixResourceAttribute for when the rest of the
		// instance label is filled in.
		&relabel.Config{
			Action:       relabel.Replace,
			SourceLabels: prommodel.LabelNames{"__address__"},
			TargetLabel:  "instance",
			Replacement:  ep.Port,
		},
	)

	interval, err := prommodel.ParseDuration(ep.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid scrape interval: %w", err)
	}
	timeout := interval
	if ep.Timeout != "" {
		timeout, err = prommodel.ParseDuration(ep.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid scrape timeout: %w", err)
		}
		if timeout > interval {
			return nil, fmt.Errorf("scrape timeout %v must not be greater than scrape interval %v", timeout, interval)
		}
	}

	metricsPath := "/metrics"
	if ep.Path != "" {
		metricsPath = ep.Path
	}

	var metricRelabelCfgs []*relabel.Config
	for _, r := range ep.MetricRelabeling {
		rcfg, err := convertRelabelingRule(r)
		if err != nil {
			return nil, err
		}
		metricRelabelCfgs = append(metricRelabelCfgs, rcfg)
	}

	scrapeCfg := &promconfig.ScrapeConfig{
		// Generate a job name to make it easy to track what generated the scrape configuration.
		// The actual job label attached to its metrics is overwritten via relabeling.
		JobName:                 fmt.Sprintf("%s/%s", id, ep.Port),
		ServiceDiscoveryConfigs: discoveryCfgs,
		MetricsPath:             metricsPath,
		Scheme:                  ep.Scheme,
		Params:                  ep.Params,
		ScrapeInterval:          interval,
		ScrapeTimeout:           timeout,
		RelabelConfigs:          relabelCfgs,
		MetricRelabelConfigs:    metricRelabelCfgs,
	}
	if limits != nil {
		scrapeCfg.SampleLimit = uint(limits.Samples)
		scrapeCfg.LabelLimit = uint(limits.Labels)
		scrapeCfg.LabelNameLengthLimit = uint(limits.LabelNameLength)
		scrapeCfg.LabelValueLengthLimit = uint(limits.LabelValueLength)
	}
	// The Prometheus configuration structs do not generally have validation methods and embed their
	// validation logic in the UnmarshalYAML methods. To keep things reasonable we don't re-validate
	// everything and simply do a final marshal-unmarshal cycle at the end to run all validation
	// upstream provides at the end of this method.
	b, err := yaml.Marshal(scrapeCfg)
	if err != nil {
		return nil, fmt.Errorf("scrape config cannot be marshalled: %w", err)
	}
	var scrapeCfgCopy promconfig.ScrapeConfig
	if err := yaml.Unmarshal(b, &scrapeCfgCopy); err != nil {
		return nil, fmt.Errorf("invalid scrape configuration: %w", err)
	}
	return scrapeCfg, nil
}

// convertRelabelingRule converts the rule to a relabel configuration. An error is returned
// if the rule would modify one of the protected labels.
func convertRelabelingRule(r RelabelingRule) (*relabel.Config, error) {
	rcfg := &relabel.Config{
		// Upstream applies ToLower when digesting the config, so we allow the same.
		Action:      relabel.Action(strings.ToLower(r.Action)),
		TargetLabel: r.TargetLabel,
		Separator:   r.Separator,
		Replacement: r.Replacement,
		Modulus:     r.Modulus,
	}
	for _, n := range r.SourceLabels {
		rcfg.SourceLabels = append(rcfg.SourceLabels, prommodel.LabelName(n))
	}
	// Instantiate the default regex Prometheus uses so that the checks below can be run
	// if no explicit value is provided.
	re := relabel.MustNewRegexp(`(.*)`)

	// We must only set the regex if its not empty. Like in other cases, the Prometheus code does
	// not setup the structs correctly and this would default to the string "null" when marshalled,
	// which is then interpreted as a regex again when read by Prometheus.
	if r.Regex != "" {
		var err error
		re, err = relabel.NewRegexp(r.Regex)
		if err != nil {
			return nil, fmt.Errorf("invalid regex %q: %w", r.Regex, err)
		}
		rcfg.Regex = re
	}

	// Validate that the protected target labels are not mutated by the provided relabeling rules.
	switch rcfg.Action {
	// Default action is "replace" per https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config.
	case relabel.Replace, relabel.HashMod, "":
		// These actions write into the target label and it must not be a protected one.
		if isProtectedLabel(r.TargetLabel) {
			return nil, fmt.Errorf("cannot relabel with action %q onto protected label %q", r.Action, r.TargetLabel)
		}
	case relabel.LabelDrop:
		if matchesAnyProtectedLabel(re) {
			return nil, fmt.Errorf("regex %s would drop at least one of the protected labels %s", r.Regex, strings.Join(protectedLabels, ", "))
		}
	case relabel.LabelKeep:
		// Keep drops all labels that don't match the regex. So all protected labels must
		// match keep.
		if !matchesAllProtectedLabels(re) {
			return nil, fmt.Errorf("regex %s would drop at least one of the protected labels %s", r.Regex, strings.Join(protectedLabels, ", "))
		}
	case relabel.LabelMap:
		// It is difficult to prove for certain that labelmap does not override a protected label.
		// Thus we just prohibit its use for now.
		// The most feasible way to support this would probably be store all protected labels
		// in __tmp_protected_<name> via a replace rule, then apply labelmap, then replace the
		// __tmp label back onto the protected label.
		return nil, fmt.Errorf("relabeling with action %q not allowed", r.Action)
	case relabel.Keep, relabel.Drop:
		// These actions don't modify a series and are OK.
	default:
		return nil, fmt.Errorf("unknown relabeling action %q", r.Action)
	}
	return rcfg, nil
}
