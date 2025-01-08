package search

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/es/es"
	"github.com/zmoog/es/es/commands"
)

var query string

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "search [index or data stream]",
		Short: "Search using the low-level API",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runner, err := es.NewRunner()
			if err != nil {
				return err
			}

			return runner.Run(commands.SearchCommand{
				Index: args[0],
				Query: query,
			})
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Query to search for")
	cmd.MarkFlagRequired("query")

	return &cmd
}
