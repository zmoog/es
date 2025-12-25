# es

Simple CLI tool to interact with [Elasticsearch REST APIs](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) from the terminal.

Currently, it supports only a fraction of the Elasticsearch REST APIs.

## Installation

### Using Homebrew (macOS and Linux)

Install this tool using Homebrew:

```sh
brew install zmoog/homebrew-es/es
```

### Using Go

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

### Search

For searching documents in an Elasticsearch cluster, run:

    $ es search my_index -q '{"query": {"match_all": {}}}' | jq
    {
      "took": 0,
      "timed_out": false,
      "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
      },
      "hits": {
        "total": {
          "value": 2,
          "relation": "eq"
        },
        "max_score": 1.0,
        "hits": [
          {
            "_index": ".ds-logs-test-default-2025.01.08-000001",
            "_id": "638aSJQBashsspu13aIM",
            "_score": 1.0,
            "_source": {
              "name": "Maurizio Branca",
              "@timestamp": "2025-01-08T22:48:27.659837121Z"
            }
          },
          {
            "_index": ".ds-logs-test-default-2025.01.08-000001",
            "_id": "0X8DSJQBashsspu145t7",
            "_score": 1.0,
            "_source": {
              "name": "Maurizio Branca",
              "@timestamp": "2025-01-08T22:23:21.211245339Z"
            }
          }
        ]
      }
    }

### Data Stream

For deleting a data stream and its backing indices, run:

    es datastream delete logs-test-default

You can also use wildcard patterns to delete multiple data streams:

    es datastream delete 'logs-*'

To skip the confirmation prompt, use the `--force` flag:

    es datastream delete logs-test-default --force

### Version

For printing the application version, run:

    es version

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
