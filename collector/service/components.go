// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/service/telemetry/otelconftelemetry"
	"go.uber.org/multierr"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/exporter/googlemanagedprometheusexporter"
	"github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/receiver/prometheusreceiver"
)

func components() (otelcol.Factories, error) {
	errs := []error{}
	factories, err := Components()
	if err != nil {
		return otelcol.Factories{}, err
	}

	extensions := []extension.Factory{}
	for _, ext := range factories.Extensions {
		extensions = append(extensions, ext)
	}
	seenExtensions := make(map[string]bool)
	var uniqueExtensions []extension.Factory
	for _, ext := range extensions {
		if !seenExtensions[ext.Type().String()] {
			seenExtensions[ext.Type().String()] = true
			uniqueExtensions = append(uniqueExtensions, ext)
		}
	}
	factories.Extensions, err = otelcol.MakeFactoryMap[extension.Factory](uniqueExtensions...)
	if err != nil {
		errs = append(errs, err)
	}

	receivers := []receiver.Factory{
		prometheusreceiver.NewFactory(),
	}
	for _, rcv := range factories.Receivers {
		receivers = append(receivers, rcv)
	}
	seenReceivers := make(map[string]bool)
	var uniqueReceivers []receiver.Factory
	for _, rcv := range receivers {
		if !seenReceivers[rcv.Type().String()] {
			seenReceivers[rcv.Type().String()] = true
			uniqueReceivers = append(uniqueReceivers, rcv)
		}
	}
	factories.Receivers, err = otelcol.MakeFactoryMap[receiver.Factory](uniqueReceivers...)
	if err != nil {
		errs = append(errs, err)
	}

	exporters := []exporter.Factory{
		fileexporter.NewFactory(),
		googlecloudexporter.NewFactory(),
		googlemanagedprometheusexporter.NewFactory(),
	}
	for _, exp := range factories.Exporters {
		exporters = append(exporters, exp)
	}
	seenExporters := make(map[string]bool)
	var uniqueExporters []exporter.Factory
	for _, exp := range exporters {
		if !seenExporters[exp.Type().String()] {
			seenExporters[exp.Type().String()] = true
			uniqueExporters = append(uniqueExporters, exp)
		}
	}
	factories.Exporters, err = otelcol.MakeFactoryMap[exporter.Factory](uniqueExporters...)
	if err != nil {
		errs = append(errs, err)
	}

	processors := []processor.Factory{
		filterprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		metricstransformprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		transformprocessor.NewFactory(),
		groupbyattrsprocessor.NewFactory(),
	}
	for _, pr := range factories.Processors {
		processors = append(processors, pr)
	}
	seenProcessors := make(map[string]bool)
	var uniqueProcessors []processor.Factory
	for _, pr := range processors {
		if !seenProcessors[pr.Type().String()] {
			seenProcessors[pr.Type().String()] = true
			uniqueProcessors = append(uniqueProcessors, pr)
		}
	}
	factories.Processors, err = otelcol.MakeFactoryMap[processor.Factory](uniqueProcessors...)
	if err != nil {
		errs = append(errs, err)
	}

	return factories, multierr.Combine(errs...)
}

func Components() (
	otelcol.Factories,
	error,
) {
	var errs error

	extensions, err := otelcol.MakeFactoryMap[extension.Factory](
		zpagesextension.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	receivers, err := otelcol.MakeFactoryMap[receiver.Factory](
		otlpreceiver.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	exporters, err := otelcol.MakeFactoryMap[exporter.Factory](
		debugexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	processors, err := otelcol.MakeFactoryMap[processor.Factory](
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Telemetry:  otelconftelemetry.NewFactory(),
	}

	return factories, errs
}
