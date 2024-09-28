# es

`es` is a simple CLI tool to use a subset of the [Elasticsearch REST API](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) from the terminal.

## Installation

Install this tool using Go:

    go install github.com/zmoog/es@latest

## Usage

### Docs

For sending a single document to an Elasticseach cluster, run:

    es docs index --doc "{\"@timestamp\": \"2023-06-11T07:43:48+0200\", \"name\": \"Maurizio Branca\"}"

