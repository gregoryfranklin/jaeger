// Copyright (c) 2017 Uber Technologies, Inc.
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

package integration

import (
	"context"
	"time"

	elastic6 "github.com/olivere/elastic"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/pkg/es/wrapper"
	"github.com/jaegertracing/jaeger/pkg/testutils"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/jaegertracing/jaeger/plugin/storage/es/dependencystore"
	"github.com/jaegertracing/jaeger/plugin/storage/es/spanstore"
)

type ES6StorageIntegration struct {
	StorageIntegration

	client        *elastic6.Client
	bulkProcessor *elastic6.BulkProcessor
	logger        *zap.Logger
}

func (s *ES6StorageIntegration) initializeES(allTagsAsFields, archive bool) error {
	rawClient, err := elastic6.NewClient(
		elastic6.SetURL(queryURL),
		elastic6.SetBasicAuth(username, password),
		elastic6.SetSniff(false))
	if err != nil {
		return err
	}
	s.logger, _ = testutils.NewLogger()

	s.client = rawClient

	s.bulkProcessor, _ = s.client.BulkProcessor().Do(context.Background())
	client := eswrapper.WrapESClient6(s.client, s.bulkProcessor)
	dependencyStore := dependencystore.NewDependencyStore(client, s.logger, indexPrefix)
	s.DependencyReader = dependencyStore
	s.DependencyWriter = dependencyStore
	s.initSpanstore(allTagsAsFields, archive)
	s.CleanUp = func() error {
		return s.esCleanUp(allTagsAsFields, archive)
	}
	s.Refresh = s.esRefresh
	s.esCleanUp(allTagsAsFields, archive)
	return nil
}

func (s *ES6StorageIntegration) esCleanUp(allTagsAsFields, archive bool) error {
	_, err := s.client.DeleteIndex("*").Do(context.Background())
	s.initSpanstore(allTagsAsFields, archive)
	return err
}

func (s *ES6StorageIntegration) initSpanstore(allTagsAsFields, archive bool) {
	bp, _ := s.client.BulkProcessor().BulkActions(1).FlushInterval(time.Nanosecond).Do(context.Background())
	client := eswrapper.WrapESClient6(s.client, bp)
	spanMapping, serviceMapping := es.GetMappings(5, 1, 6)
	s.SpanWriter = spanstore.NewSpanWriter(
		spanstore.SpanWriterParams{
			Client:            client,
			Logger:            s.logger,
			MetricsFactory:    metrics.NullFactory,
			IndexPrefix:       indexPrefix,
			AllTagsAsFields:   allTagsAsFields,
			TagDotReplacement: tagKeyDeDotChar,
			SpanMapping:       spanMapping,
			ServiceMapping:    serviceMapping,
			Archive:           archive,
		})
	s.SpanReader = spanstore.NewSpanReader(spanstore.SpanReaderParams{
		Client:            client,
		Logger:            s.logger,
		MetricsFactory:    metrics.NullFactory,
		IndexPrefix:       indexPrefix,
		MaxSpanAge:        maxSpanAge,
		TagDotReplacement: tagKeyDeDotChar,
		Archive:           archive,
	})
}

func (s *ES6StorageIntegration) esRefresh() error {
	err := s.bulkProcessor.Flush()
	if err != nil {
		return err
	}
	_, err = s.client.Refresh().Do(context.Background())
	return err
}
