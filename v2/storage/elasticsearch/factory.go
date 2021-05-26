// Copyright (c) 2020 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elasticsearch

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jaegertracing/jaeger/plugin/storage/es"
)

const (
	// TypeStr defines type of the Elasticsearch exporter.
	typeStr = "jaeger_elasticsearch"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTracesExporter))
}

// CreateDefaultConfig returns default configuration of Factory.
// This function implements OTEL component.ExporterFactoryBase interface.
func createDefaultConfig() config.Exporter {
	opts := es.NewOptions("es")
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewID(typeStr)),
		TimeoutSettings:  exporterhelper.DefaultTimeoutSettings(),
		RetrySettings:    exporterhelper.DefaultRetrySettings(),
		QueueSettings:    exporterhelper.DefaultQueueSettings(),
		Options:          *opts,
	}
}

// CreateTracesExporter creates Jaeger Elasticsearch trace exporter.
// This function implements OTEL component.ExporterFactory interface.
func createTracesExporter(
	ctx context.Context,
	params component.ExporterCreateParams,
	cfg config.Exporter,
) (component.TracesExporter, error) {
	esCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("could not cast configuration to %s", typeStr)
	}
	return newExporter(ctx, esCfg, params)
}
