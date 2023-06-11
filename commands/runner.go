package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// Runner is a command runner.
type Runner struct {
	uow UnitOfWork
}

// Run executes the given command.
func (r Runner) Run(command Command) error {
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

	//
	// Create the Elasticsearch client.
	//

	cfg := elasticsearch.Config{
		Addresses: strings.Split(endpoints, ","),
		APIKey:    apiKey,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the elasticsearch client: %w", err)
	}

	runner := Runner{
		uow: UnitOfWork{
			Client: client,
		},
	}

	return &runner, nil
}
