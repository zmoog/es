package search

import (
	"fmt"
	"io"
	"os"

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

			// If query is "-" or empty, read from stdin
			var queryBody string
			if query == "-" || query == "" {
				queryBytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("error reading query from stdin: %w", err)
				}
				queryBody = string(queryBytes)
			} else {
				queryBody = query
			}

			return runner.Run(commands.SearchCommand{
				Index: args[0],
				Query: queryBody,
			})
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "-", "Query to search for (use '-' or omit to read from stdin)")

	return &cmd
}
