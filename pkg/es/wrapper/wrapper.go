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

package eswrapper

import (
	"context"
	"fmt"

	elastic6 "github.com/olivere/elastic"
	elastic5 "gopkg.in/olivere/elastic.v5"

	"github.com/jaegertracing/jaeger/pkg/es"
)

// This file avoids lint because the Id and Json are required to be capitalized, but must match an outside library.

// ClientWrapper is a wrapper around elastic.Client
type ClientWrapper struct {
	client5      *elastic5.Client
	bulkService5 *elastic5.BulkProcessor
	client6      *elastic6.Client
	bulkService6 *elastic6.BulkProcessor
	client7      *elastic6.Client
	bulkService7 *elastic6.BulkProcessor
	version      int
}

// GetVersion returns the ElasticSearch major version
func (c ClientWrapper) GetVersion() int {
	return c.version
}

// WrapESClient creates a v5 ESClient out of *elastic.Client.
func WrapESClient5(client *elastic5.Client, s *elastic5.BulkProcessor) ClientWrapper {
	return ClientWrapper{client5: client, bulkService5: s, version: 5}
}

// WrapESClient creates a v6 ESClient out of *elastic.Client.
func WrapESClient6(client *elastic6.Client, s *elastic6.BulkProcessor) ClientWrapper {
	return ClientWrapper{client6: client, bulkService6: s, version: 6}
}

// WrapESClient creates a v6 ESClient out of *elastic.Client.
func WrapESClient7(client *elastic6.Client, s *elastic6.BulkProcessor) ClientWrapper {
	return ClientWrapper{client7: client, bulkService7: s, version: 7, client6: client, bulkService6: s}
}

// IndexExists calls this function to internal client.
func (c ClientWrapper) IndexExists(index string) es.IndicesExistsService {
	return WrapESIndicesExistsService(&c, index)
}

// CreateIndex calls this function to internal client.
func (c ClientWrapper) CreateIndex(index string) es.IndicesCreateService {
	return WrapESIndicesCreateService(&c, index)
}

// Index calls this function to internal client.
func (c ClientWrapper) Index() es.IndexService {
	return WrapESIndexService(&c)
}

// Search calls this function to internal client.
func (c ClientWrapper) Search(indices ...string) es.SearchService {
	return WrapESSearchService(c.client6.Search(indices...))
}

// MultiSearch calls this function to internal client.
func (c ClientWrapper) MultiSearch() es.MultiSearchService {
	return WrapESMultiSearchService(c.client6.MultiSearch())
}

// Close closes ESClient and flushes all data to the storage.
func (c ClientWrapper) Close() error {
	var allerr error
	if c.bulkService5 != nil {
		if err := c.bulkService5.Close(); err != nil {
			allerr = err
		}
	}
	if c.bulkService6 != nil {
		if err := c.bulkService6.Close(); err != nil {
			allerr = err
		}
	}
	if c.bulkService7 != nil {
		if err := c.bulkService7.Close(); err != nil {
			allerr = err
		}
	}
	return allerr
}

// ---

// IndicesExistsServiceWrapper is a wrapper around elastic.IndicesExistsService
type IndicesExistsServiceWrapper struct {
	wrapper *ClientWrapper
	index   string
}

// WrapESIndicesExistsService creates an ESIndicesExistsService out of *elastic.IndicesExistsService.
func WrapESIndicesExistsService(c *ClientWrapper, index string) IndicesExistsServiceWrapper {
	return IndicesExistsServiceWrapper{wrapper: c, index: index}
}

// Do calls this function to internal service.
func (e IndicesExistsServiceWrapper) Do(ctx context.Context) (bool, error) {
	if e.wrapper.version == 5 {
		return e.wrapper.client5.IndexExists(e.index).Do(ctx)
	} else if e.wrapper.version == 6 {
		return e.wrapper.client6.IndexExists(e.index).Do(ctx)
	} else if e.wrapper.version == 6 {
		return e.wrapper.client7.IndexExists(e.index).Do(ctx)
	} else {
		return false, fmt.Errorf("Unsupported ElasticSearch version")
	}
}

// ---

// IndicesCreateServiceWrapper is a wrapper around elastic.IndicesCreateService
type IndicesCreateServiceWrapper struct {
	wrapper *ClientWrapper
	index   string
	mapping string
}

// WrapESIndicesCreateService creates an ESIndicesCreateService out of *elastic.IndicesCreateService.
func WrapESIndicesCreateService(c *ClientWrapper, index string) IndicesCreateServiceWrapper {
	return IndicesCreateServiceWrapper{wrapper: c, index: index}
}

// Body calls this function to internal service.
func (c IndicesCreateServiceWrapper) Body(mapping string) es.IndicesCreateService {
	c.mapping = mapping
	return c
}

// Do calls this function to internal service.
func (c IndicesCreateServiceWrapper) Do(ctx context.Context) (*es.IndicesCreateResult, error) {
	if c.wrapper.version == 5 {
		result, err := c.wrapper.client5.CreateIndex(c.index).Body(c.mapping).Do(ctx)
		if err != nil {
			return nil, err
		}
		return &es.IndicesCreateResult{
			Acknowledged:       result.Acknowledged,
			ShardsAcknowledged: result.ShardsAcknowledged,
		}, nil
	} else if c.wrapper.version == 6 {
		result, err := c.wrapper.client6.CreateIndex(c.index).Body(c.mapping).Do(ctx)
		if err != nil {
			return nil, err
		}
		return &es.IndicesCreateResult{
			Acknowledged:       result.Acknowledged,
			ShardsAcknowledged: result.ShardsAcknowledged,
			Index:              result.Index,
		}, nil
	} else if c.wrapper.version == 7 {
		result, err := c.wrapper.client7.CreateIndex(c.index).Body(c.mapping).IncludeTypeName(true).Do(ctx)
		if err != nil {
			return nil, err
		}
		return &es.IndicesCreateResult{
			Acknowledged:       result.Acknowledged,
			ShardsAcknowledged: result.ShardsAcknowledged,
			Index:              result.Index,
		}, nil
	} else {
		return nil, fmt.Errorf("Unsupported ElasticSearch Version")
	}
}

// ---

// IndexServiceWrapper is a wrapper around elastic.ESIndexService.
// See wrapper_nolint.go for more functions.
type IndexServiceWrapper struct {
	wrapper *ClientWrapper
	index   string
	typ     string
	id      string
	body    interface{}
}

// WrapESIndexService creates an ESIndexService out of *elastic.ESIndexService.
func WrapESIndexService(c *ClientWrapper) IndexServiceWrapper {
	return IndexServiceWrapper{wrapper: c}
}

// Index calls this function to internal service.
func (i IndexServiceWrapper) Index(index string) es.IndexService {
	i.index = index
	return i
}

// Type calls this function to internal service.
func (i IndexServiceWrapper) Type(typ string) es.IndexService {
	i.typ = typ
	return i
}

// Add adds the request to bulk service
func (i IndexServiceWrapper) Add() {
	if i.wrapper.version == 5 {
		i.wrapper.bulkService5.Add(elastic5.NewBulkIndexRequest().Index(i.index).Type(i.typ).Id(i.id).Doc(i.body))
	} else if i.wrapper.version == 6 {
		i.wrapper.bulkService6.Add(elastic6.NewBulkIndexRequest().Index(i.index).Type(i.typ).Id(i.id).Doc(i.body))
	} else if i.wrapper.version == 7 {
		i.wrapper.bulkService7.Add(elastic6.NewBulkIndexRequest().Index(i.index).Type(i.typ).Id(i.id).Doc(i.body))
	} else {
		fmt.Println("Unspported ElasticSearch Version")
	}
}

// ---

// SearchServiceWrapper is a wrapper around elastic.ESSearchService
type SearchServiceWrapper struct {
	searchService *elastic6.SearchService
}

// WrapESSearchService creates an ESSearchService out of *elastic.ESSearchService.
func WrapESSearchService(searchService *elastic6.SearchService) SearchServiceWrapper {
	return SearchServiceWrapper{searchService: searchService}
}

// Type calls this function to internal service.
func (s SearchServiceWrapper) Type(typ string) es.SearchService {
	return WrapESSearchService(s.searchService.Type(typ))
}

// Size calls this function to internal service.
func (s SearchServiceWrapper) Size(size int) es.SearchService {
	return WrapESSearchService(s.searchService.Size(size))
}

// Aggregation calls this function to internal service.
func (s SearchServiceWrapper) Aggregation(name string, aggregation elastic6.Aggregation) es.SearchService {
	return WrapESSearchService(s.searchService.Aggregation(name, aggregation))
}

// IgnoreUnavailable calls this function to internal service.
func (s SearchServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) es.SearchService {
	return WrapESSearchService(s.searchService.IgnoreUnavailable(ignoreUnavailable))
}

// Query calls this function to internal service.
func (s SearchServiceWrapper) Query(query elastic6.Query) es.SearchService {
	return WrapESSearchService(s.searchService.Query(query))
}

// Do calls this function to internal service.
func (s SearchServiceWrapper) Do(ctx context.Context) (*elastic6.SearchResult, error) {
	return s.searchService.Do(ctx)
}

// MultiSearchServiceWrapper is a wrapper around elastic.ESMultiSearchService
type MultiSearchServiceWrapper struct {
	multiSearchService *elastic6.MultiSearchService
}

// WrapESMultiSearchService creates an ESSearchService out of *elastic.ESSearchService.
func WrapESMultiSearchService(multiSearchService *elastic6.MultiSearchService) MultiSearchServiceWrapper {
	return MultiSearchServiceWrapper{multiSearchService: multiSearchService}
}

// Add calls this function to internal service.
func (s MultiSearchServiceWrapper) Add(requests ...*elastic6.SearchRequest) es.MultiSearchService {
	return WrapESMultiSearchService(s.multiSearchService.Add(requests...))
}

// Index calls this function to internal service.
func (s MultiSearchServiceWrapper) Index(indices ...string) es.MultiSearchService {
	return WrapESMultiSearchService(s.multiSearchService.Index(indices...))
}

// Do calls this function to internal service.
func (s MultiSearchServiceWrapper) Do(ctx context.Context) (*elastic6.MultiSearchResult, error) {
	return s.multiSearchService.Do(ctx)
}
