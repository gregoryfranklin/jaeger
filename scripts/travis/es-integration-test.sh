#!/bin/bash

set -e

# Default to elasticsearch 5.6.12 if ES_VERSION is not set
ES_VERSION=${ES_VERSION:-5.6.12}

docker pull docker.elastic.co/elasticsearch/elasticsearch:${ES_VERSION}
CID=$(docker run -d -p 9200:9200 -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" docker.elastic.co/elasticsearch/elasticsearch:${ES_VERSION})
export STORAGE=elasticsearch
make storage-integration-test
docker kill $CID
