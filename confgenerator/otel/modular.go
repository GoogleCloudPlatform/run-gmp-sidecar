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

// Package otel provides data structures to represent and generate otel configuration.
package otel

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// ReceiverPipeline represents a single OT receiver and zero or more processors that must be chained after that receiver.
type ReceiverPipeline struct {
	Receiver   Component
	Processors []Component
}

// Component represents a single OT component (receiver, processor, exporter, etc.)
type Component struct {
	// Type is the string type needed to instantiate the OT component.
	Type string
	// Config is an object which can be serialized by mapstructure into the configuration for the component.
	// This can either be a map[string]interface{} or a Config struct from OT.
	Config interface{}
}

func (c Component) name(suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s/%s", c.Type, suffix)
	}
	return c.Type
}

type ModularConfig struct {
	LogLevel          string
	ReceiverPipelines map[string]ReceiverPipeline
	Exporter          Component
	SelfMetricsPort   int
}

func (c ModularConfig) Generate() (string, error) {
	receivers := map[string]interface{}{}
	processors := map[string]interface{}{}
	exporters := map[string]interface{}{}
	pipelines := map[string]interface{}{}
	service := map[string]map[string]interface{}{
		"pipelines": pipelines,
		"telemetry": {
			"metrics": map[string]interface{}{
				"address": fmt.Sprintf("0.0.0.0:%d", c.SelfMetricsPort),
			},
		},
	}

	configMap := map[string]interface{}{
		"receivers":  receivers,
		"processors": processors,
		"exporters":  exporters,
		"service":    service,
	}

	for key, receiverPipeline := range c.ReceiverPipelines {
		receiverName := receiverPipeline.Receiver.name(key)
		var receiverProcessorNames []string
		for i, processor := range receiverPipeline.Processors {
			name := processor.name(fmt.Sprintf("%s_%d", key, i))
			receiverProcessorNames = append(receiverProcessorNames, name)
			processors[name] = processor.Config
		}
		receivers[receiverName] = receiverPipeline.Receiver.Config

		// Keep track of all the processors we're adding to the config.
		var processorNames []string
		processorNames = append(processorNames, receiverProcessorNames...)

		// Add the resource detector
		processors["resourcedetection"] = GCPResourceDetector().Config
		processorNames = append(processorNames, "resourcedetection")

		// Add the serverless instance id as a metric label
		transformProcessor := TransformationMetrics(FlattenResourceAttribute("faas.id", "cloud_run_instance"), PrefixResourceAttribute("service.instance.id", "cloud_run_instance", ":"))

		processors["transform/instance"] = transformProcessor.Config
		processorNames = append(processorNames, "transform/instance")

		exporters["googlemanagedprometheus"] = c.Exporter.Config
		pipelines["metrics/"+key] = map[string]interface{}{
			"receivers":  []string{receiverName},
			"processors": processorNames,
			"exporters":  []string{"googlemanagedprometheus"},
		}
	}

	out, err := configToYaml(configMap)
	// TODO: Return []byte
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// configToYaml converts a tree of structs into a YAML file.
// To match OT's built-in config parsing, we use mapstructure to convert the tree of structs into a tree of maps.
// This allows the direct use of OT's config types at any level of the hierarchy.
func configToYaml(config interface{}) ([]byte, error) {
	outMap := make(map[string]interface{})
	if err := mapstructure.Decode(config, &outMap); err != nil {
		return nil, err
	}
	return yaml.Marshal(outMap)
}
