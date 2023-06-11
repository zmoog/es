# elasticsearch-cli

Elasticsearch CLI tool to access https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html from the terminal.

## Installation

Install this tool using Go:

    go install github.com/zmoog/elasticsearch-cli@latest

## Usage

### Docs

For sending a single document to an Elasticseach cluster, run:

    elasticsearch-cli docs index --doc "{\"@timestamp\": \"2023-06-11T07:43:48+0200\", \"name\": \"Maurizio Branca\"}"

