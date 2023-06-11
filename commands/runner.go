package commands

import (
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type Runner struct {
	uow UnitOfWork
}

func (r Runner) Run(command Command) error {
	err := command.ExecuteWith(r.uow)
	if err != nil {
		return err
	}

	return nil
}

func NewRunner() (*Runner, error) {
	endpoints, ok := os.LookupEnv("ELASTICSEARCH_ENDPOINTS")
	if !ok {
		log.Fatal("ELASTICSEARCH_ENDPOINTS is not set")
		os.Exit(1)
	}

	apiKey, ok := os.LookupEnv("ELASTICSEARCH_API_KEY")
	if !ok {
		log.Fatal("ELASTICSEARCH_API_KEY is not set")
		os.Exit(1)
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
		log.Fatalf("Error creating the client: %s", err)
		os.Exit(1)
	}

	runner := Runner{
		uow: UnitOfWork{
			Client: client,
		},
	}

	return &runner, nil
}
