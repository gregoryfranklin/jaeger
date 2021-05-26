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
	"testing"

	"github.com/stretchr/testify/assert"

	jConfig "github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
)

func TestDefaultConfig(t *testing.T) {
	v, _ := jConfig.Viperize(es.NewOptions("es").AddFlags)
	opts := es.NewOptions("es")
	opts.InitFromViper(v)
	defaultCfg := createDefaultConfig().(*Config)
	assert.Equal(t, []string{"http://127.0.0.1:9200"}, defaultCfg.GetPrimary().Servers)
	assert.Equal(t, int64(5), defaultCfg.GetPrimary().NumShards)
	assert.Equal(t, int64(1), defaultCfg.GetPrimary().NumReplicas)
	assert.Equal(t, "@", defaultCfg.GetPrimary().Tags.DotReplacement)
	assert.Equal(t, false, defaultCfg.GetPrimary().TLS.Enabled)
}
