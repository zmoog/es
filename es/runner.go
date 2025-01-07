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
	apiEndpoints := viper.GetString("api.endpoints")
	apiKey := viper.GetString("api.key")

	retryBackoff := backoff.NewExponentialBackOff()

	// caCert, _ := os.ReadFile("/Users/zmoog/.elastic-package/profiles/default/certs/ca-cert.pem")

	//
	// Create the Elasticsearch client.
	//

	cfg := elasticsearch.Config{
		Addresses: strings.Split(apiEndpoints, ","),
		APIKey:    apiKey,
		// RetryOnStatus: []int{502, 503, 504, 429},
		RetryOnStatus: viper.GetIntSlice("client.retry-on-status"),
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}

			return retryBackoff.NextBackOff()
		},
		MaxRetries: viper.GetInt("client.max-retries"),
		// CACert:     viper.get,
	}

	if viper.IsSet("client.ca-cert-path") {
		cert, err := os.ReadFile(viper.GetString("client.ca-cert-path"))
		if err != nil {
			return nil, fmt.Errorf("error reading CA certificate: %w", err)
		}

		cfg.CACert = cert
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
