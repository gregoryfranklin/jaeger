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
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	elastic6 "github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger/model"
)

const (
	host            = "0.0.0.0"
	queryPort       = "9200"
	queryHostPort   = host + ":" + queryPort
	queryURL        = "http://" + queryHostPort
	username        = "elastic"  // the elasticsearch default username
	password        = "changeme" // the elasticsearch default password
	indexPrefix     = "integration-test"
	tagKeyDeDotChar = "@"
	maxSpanAge      = time.Hour * 72
)

func healthCheck() error {
	for i := 0; i < 200; i++ {
		if _, err := http.Get(queryURL); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("elastic search is not ready")
}

func version() (int, error) {
	rawClient, err := elastic6.NewClient(
		elastic6.SetURL(queryURL),
		elastic6.SetBasicAuth(username, password),
		elastic6.SetSniff(false))
	if err != nil {
		return 0, err
	}
	result, _, err := elastic6.NewPingService(rawClient).Do(context.Background())
	if err != nil {
		return 0, err
	}
	version, err := strconv.Atoi(string(result.Version.Number[0]))
	if err != nil {
		return 0, err
	}
	return version, nil
}

func testElasticsearchStorage(t *testing.T, allTagsAsFields, archive bool) {
	if os.Getenv("STORAGE") != "elasticsearch" {
		t.Skip("Integration test against ElasticSearch skipped; set STORAGE env var to elasticsearch to run this")
	}
	if err := healthCheck(); err != nil {
		t.Fatal(err)
	}
	var esVersion int
	if v, err := version(); err != nil {
		t.Fatal(err)
	} else {
		esVersion = v
	}
	if esVersion == 5 {
		s := &ES5StorageIntegration{}
		require.NoError(t, s.initializeES(allTagsAsFields, archive))

		if archive {
			t.Run("ArchiveTrace", s.testArchiveTrace)
		} else {
			s.IntegrationTestAll(t)
		}
	} else if esVersion == 6 {
		s := &ES6StorageIntegration{}
		require.NoError(t, s.initializeES(allTagsAsFields, archive))

		if archive {
			t.Run("ArchiveTrace", s.testArchiveTrace)
		} else {
			s.IntegrationTestAll(t)
		}
	} else if esVersion == 7 {
		s := &ES7StorageIntegration{}
		require.NoError(t, s.initializeES(allTagsAsFields, archive))

		if archive {
			t.Run("ArchiveTrace", s.testArchiveTrace)
		} else {
			s.IntegrationTestAll(t)
		}
	} else {
		t.Fatal("Unsupported ElasticSearch version")
	}
}

func TestElasticsearchStorage(t *testing.T) {
	testElasticsearchStorage(t, false, false)
}

func TestElasticsearchStorage_AllTagsAsObjectFields(t *testing.T) {
	testElasticsearchStorage(t, true, false)
}

func TestElasticsearchStorage_Archive(t *testing.T) {
	testElasticsearchStorage(t, false, true)
}

func (s *StorageIntegration) testArchiveTrace(t *testing.T) {
	defer s.cleanUp(t)
	tId := model.NewTraceID(uint64(11), uint64(22))
	expected := &model.Span{
		OperationName: "archive_span",
		StartTime:     time.Now().Add(-maxSpanAge * 5),
		TraceID:       tId,
		SpanID:        model.NewSpanID(55),
		References:    []model.SpanRef{},
		Process:       model.NewProcess("archived_service", model.KeyValues{}),
	}

	require.NoError(t, s.SpanWriter.WriteSpan(expected))
	s.refresh(t)

	var actual *model.Trace
	found := s.waitForCondition(t, func(t *testing.T) bool {
		var err error
		actual, err = s.SpanReader.GetTrace(context.Background(), tId)
		return err == nil && len(actual.Spans) == 1
	})
	if !assert.True(t, found) {
		CompareTraces(t, &model.Trace{Spans: []*model.Span{expected}}, actual)
	}
}
