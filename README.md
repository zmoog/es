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

### Configuration

You have the flexibility to set the app setting using the following methods (increasing priority):

- config file
- env vars
- CLI option flags

#### Config file

For example, I want to have a default config in a file (low priority):

    api:
      endpoints: https://localhost:9200
      key: [redacted]
      client:
      max-retries: 1
      retry-on-status:
        - 502
        - 503
        - 504
        - 429
      ca-cert-path: /Users/zmoog/.elastic-package/profiles/default/certs/ca-cert.pem

#### Env vars

In addition, I want to override any setting using an env var (mid priority)

    export ES_API_ENDPOINTS="another endpoint"
    export ES_API_KEY="another key"

    cat docs.json | es docs bulk -f -

#### CLI flags

In addition, I want to override any setting using a flag (high priority):

    cat docs.json | es docs bulk -f - --api-key "[yet another key]"
