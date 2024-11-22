package es

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/zmoog/es/es/commands"

	"github.com/elastic/go-elasticsearch/v8"
)

// Runner is a command runner.
type Runner struct {
	uow commands.UnitOfWork
}

// Run executes the given command.
func (r Runner) Run(command commands.Command) error {
	err := command.ExecuteWith(r.uow)
	if err != nil {
		return err
	}

	return nil
}

// NewRunner creates a new runner that can execute commands.
func NewRunner() (*Runner, error) {
	endpoints, ok := os.LookupEnv("ELASTICSEARCH_ENDPOINTS")
	if !ok {
		return nil, fmt.Errorf("ELASTICSEARCH_ENDPOINTS is not set")
	}

	apiKey, ok := os.LookupEnv("ELASTICSEARCH_API_KEY")
	if !ok {
		return nil, fmt.Errorf("ELASTICSEARCH_API_KEY is not set")
	}

	retryBackoff := backoff.NewExponentialBackOff()

	//
	// Create the Elasticsearch client.
	//

	cfg := elasticsearch.Config{
		Addresses:     strings.Split(endpoints, ","),
		APIKey:        apiKey,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}

			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the elasticsearch client: %w", err)
	}

	runner := Runner{
		uow: commands.UnitOfWork{
			Client: client,
		},
	}

	return &runner, nil
}
