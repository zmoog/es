package es

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/viper"
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
	retryBackoff := backoff.NewExponentialBackOff()

	//
	// Create the Elasticsearch client.
	//
	cfg := elasticsearch.Config{
		APIKey:        viper.GetString("api.key"),
		RetryOnStatus: viper.GetIntSlice("client.retry-on-status"),

		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}

			return retryBackoff.NextBackOff()
		},
		MaxRetries: viper.GetInt("client.max-retries"),
	}

	// Prefer CloudID (set from a context's cloud_id field); fall back to
	// explicit endpoint addresses.
	if cloudID := viper.GetString("api.cloud-id"); cloudID != "" {
		cfg.CloudID = cloudID
	} else if endpoints := viper.GetString("api.endpoints"); endpoints != "" {
		cfg.Addresses = strings.Split(endpoints, ",")
	}

	// Load and set the CA certificate, if a path is provided.
	if viper.IsSet("client.ca-cert-path") {
		caCert, err := os.ReadFile(viper.GetString("client.ca-cert-path"))
		if err != nil {
			return nil, fmt.Errorf("error reading CA certificate: %w", err)
		}

		// Set the CA certificate.
		cfg.CACert = caCert
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
