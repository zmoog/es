# es

Simple CLI tool to interact with [Elasticsearch REST APIs](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) from the terminal.

Currently, it supports only a fraction of the Elasticsearch REST APIs.

## Installation

Install this tool using Go:

    go install github.com/zmoog/es@latest

## Usage

### Docs

For sending a single document to an Elasticseach cluster, run:

    es docs index --doc "{\"@timestamp\": \"2023-06-11T07:43:48+0200\", \"name\": \"Maurizio Branca\"}"

For sending multiple documents to an Elasticsearch cluster, run:

    es docs bulk --file docs.json

You can also send documents from stdin (use `-` as the filename):

    cat docs.json | es docs bulk -f -
