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
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jaegertracing/jaeger/plugin/storage/es"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	e0 := cfg.Exporters[config.NewID(typeStr)]
	assert.Equal(t, e0, createDefaultConfig())

	e1 := cfg.Exporters[config.NewIDWithName(typeStr, "2")]
	opts := es.NewOptions("es")
	opts.Primary.Servers = []string{"http://someUrl"}
	opts.Primary.Username = "user"
	opts.Primary.Password = "pass"
	opts.Primary.Sniffer = true
	opts.Primary.Tags.AllAsFields = true
	opts.Primary.Tags.DotReplacement = "O"
	opts.Primary.Tags.File = "/etc/jaeger"
	opts.Primary.UseReadWriteAliases = true
	opts.Primary.CreateIndexTemplates = false

	assert.Equal(t,
		&Config{
			ExporterSettings: config.NewExporterSettings(config.NewIDWithName(typeStr, "2")),
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 5 * time.Second,
			},
			RetrySettings: exporterhelper.RetrySettings{
				Enabled:         true,
				InitialInterval: 5 * time.Second,
				MaxInterval:     30 * time.Second,
				MaxElapsedTime:  300 * time.Second,
			},
			QueueSettings: exporterhelper.QueueSettings{
				Enabled:      true,
				NumConsumers: 10,
				QueueSize:    5000,
			},
			Options: *opts,
		}, e1,
	)
}
