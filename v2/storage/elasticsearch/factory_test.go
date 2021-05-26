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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

func TestCreateTraceExporter(t *testing.T) {
	config := createDefaultConfig().(*Config)
	config.Primary.Servers = []string{"http://foobardoesnotexists.test"}
	exporter, err := createTracesExporter(context.Background(), component.ExporterCreateParams{Logger: zap.NewNop()}, config)
	require.Nil(t, exporter)
	require.Error(t, err)
}

func TestCreateTraceExporter_nilConfig(t *testing.T) {
	exporter, err := createTracesExporter(context.Background(), component.ExporterCreateParams{}, nil)
	require.Nil(t, exporter)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not cast configuration to jaeger_elasticsearch")
}
