package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/es/config"
)

func newContextSetCommand() *cobra.Command {
	var cloudID string
	var elasticsearchURL string
	var kibanaURL string
	var apiKey string

	cmd := &cobra.Command{
		Use:   "set <name>",
		Short: "Create or update a named context",
		Long: `Create or update a context entry in the config file.

Examples:
  # Cloud deployment
  es config context set prod --cloud-id 'name:base64...' --api-key 'encoded-key'

  # Local cluster
  es config context set local --elasticsearch-url https://localhost:9200 --api-key 'encoded-key'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Start from the existing entry (if any) so that unspecified
			// flags preserve their stored values.
			ctx := cfg.Contexts[name]

			if cmd.Flags().Changed("cloud-id") {
				ctx.CloudID = cloudID
			}
			if cmd.Flags().Changed("elasticsearch-url") {
				ctx.ElasticsearchURL = elasticsearchURL
			}
			if cmd.Flags().Changed("kibana-url") {
				ctx.KibanaURL = kibanaURL
			}
			if cmd.Flags().Changed("api-key") {
				ctx.APIKey = apiKey
			}

			cfg.Contexts[name] = ctx

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Printf("Context %q saved.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&cloudID, "cloud-id", "", "Elastic Cloud ID")
	cmd.Flags().StringVar(&elasticsearchURL, "elasticsearch-url", "", "Elasticsearch base URL")
	cmd.Flags().StringVar(&kibanaURL, "kibana-url", "", "Kibana base URL (optional)")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key")

	return cmd
}
