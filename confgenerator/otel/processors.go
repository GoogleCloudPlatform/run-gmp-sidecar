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

package otel

import "fmt"

// Helper functions to easily build up processor configs.

// MetricsFilter returns a Component that filters metrics.
// polarity should be "include" or "exclude".
// matchType should be "strict" or "regexp".
func MetricsFilter(polarity, matchType string, metricNames ...string) Component {
	return Component{
		Type: "filter",
		Config: map[string]interface{}{
			"metrics": map[string]interface{}{
				polarity: map[string]interface{}{
					"match_type":   matchType,
					"metric_names": metricNames,
				},
			},
		},
	}
}

// MetricsTransform returns a Component that performs the transformations specified as arguments.
func MetricsTransform(metrics ...map[string]interface{}) Component {
	return Component{
		Type: "metricstransform",
		Config: map[string]interface{}{
			"transforms": metrics,
		},
	}
}

// RenameMetric returns a config snippet that renames old to new, applying zero or more transformations.
func RenameMetric(old, new string, operations ...map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{
		"include":  old,
		"action":   "update",
		"new_name": new,
	}
	if len(operations) > 0 {
		out["operations"] = operations
	}
	return out
}

// CombineMetrics returns a config snippet that renames metrics matching the regex old to new, applying zero or more transformations.
func CombineMetrics(old, new string, operations ...map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{
		"include":       old,
		"match_type":    "regexp",
		"action":        "combine",
		"new_name":      new,
		"submatch_case": "lower",
	}
	if len(operations) > 0 {
		out["operations"] = operations
	}
	return out
}

// ToggleScalarDataType transforms int -> double and double -> int.
var ToggleScalarDataType = map[string]interface{}{"action": "toggle_scalar_data_type"}

// AddLabel adds a label with a fixed value.
func AddLabel(key, value string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "add_label",
		"new_label": key,
		"new_value": value,
	}
}

// RenameLabel renames old to new
func RenameLabel(old, new string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "update_label",
		"label":     old,
		"new_label": new,
	}
}

// AggregateLabels removes all labels except those in the passed list, aggregating values using aggregationType.
func AggregateLabels(aggregationType string, labels ...string) map[string]interface{} {
	return map[string]interface{}{
		"action":           "aggregate_labels",
		"label_set":        labels,
		"aggregation_type": aggregationType,
	}
}

// GroupByGMPAttrs moves the "namespace" and "cluster" metric attributes to
// resource attributes.
//
// Metrics coming from run-gmp-sidecar are written against the
// `prometheus_target` monitored resource in Cloud Monitoring. The labels for
// these monitored resources come from the OTel resource labels. As a result,
// this processor needs to promote certain metric labels to resource labels so
// the translation can happen correctly.
//
// See https://cloud.google.com/monitoring/api/resources#tag_prometheus_target
// for more information about the monitored resource used.
func GroupByGMPAttrs() Component {
	return Component{
		Type: "groupbyattrs",
		Config: map[string]interface{}{
			"keys": []string{"namespace", "cluster"},
		},
	}
}

// GCPResourceDetector returns a resourcedetection processor configured for only GCP.
func GCPResourceDetector() Component {
	config := map[string]interface{}{
		"detectors": []string{"gcp", "env"},
	}

	return Component{
		Type:   "resourcedetection",
		Config: config,
	}
}

// AddResourceAttr adds a resource attribute using a processor.
func AddResourceAttr(key, from string) Component {
	attributeConfig := map[string]interface{}{
		"key":            key,
		"from_attribute": from,
		"action":         "insert",
	}
	config := map[string]interface{}{
		"attributes": attributeConfig,
	}

	return Component{
		Type:   "resource",
		Config: config,
	}
}

// TransformationMetrics returns a transform processor object that contains all the queries passed into it.
func TransformationMetrics(queries ...TransformQuery) Component {
	queryStrings := []string{}
	for _, q := range queries {
		queryStrings = append(queryStrings, string(q))
	}
	return Component{
		Type: "transform",
		Config: map[string]map[string]interface{}{
			"metric_statements": {
				"context":    "datapoint",
				"statements": queryStrings,
			},
		},
	}
}

// TransformQuery is a type wrapper for query expressions supported by the transform
// processor found here: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/transformprocessor
type TransformQuery string

// FlattenResourceAttribute returns an expression that brings down a resource attribute to a
// metric attribute.
func FlattenResourceAttribute(resourceAttribute, metricAttribute string) TransformQuery {
	return TransformQuery(fmt.Sprintf(`set(attributes["%s"], resource.attributes["%s"])`, metricAttribute, resourceAttribute))
}

// PrefixResourceAttribute prefixes the resource attribute with another resource
// attribute.
//
// Note: Mutating the resource attribute results in this update happening for
// each data point.  Since the OTTL statement uses the resource attribute in
// both the target and the source labels, we must make sure after the first
// mutation, the subsequent transformations for the same resource is a no-op.
func PrefixResourceAttribute(destResourceAttribute, srcResourceAttribute, delimiter string) TransformQuery {
	return TransformQuery(fmt.Sprintf(`replace_pattern(resource.attributes["%s"], "^(\\d+)$$", Concat([resource.attributes["%s"], "$$1"], "%s"))`, destResourceAttribute, srcResourceAttribute, delimiter))
}

// AddMetricLabel adds a new metric attribute. If it already exists, then it is overwritten.
func AddMetricLabel(key, val string) TransformQuery {
	return TransformQuery(fmt.Sprintf(`set(attributes["%s"], "%s")`, key, val))
}
