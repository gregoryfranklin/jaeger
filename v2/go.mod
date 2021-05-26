module github.com/jaegertracing/jaeger/v2

go 1.16

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/elastic/go-elasticsearch/v6 v6.8.10
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/gogo/googleapis v1.3.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/jaegertracing/jaeger v1.22.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector v0.27.0
	go.uber.org/zap v1.16.0
	gopkg.in/ini.v1 v1.57.0 // indirect
)

replace github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.4

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.4
