module github.com/jaegertracing/jaeger/v2

go 1.16

require (
	github.com/elastic/go-elasticsearch/v6 v6.8.10
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/jaegertracing/jaeger v1.22.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector v0.22.0
	go.uber.org/zap v1.16.0
)

replace github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.4

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.4
