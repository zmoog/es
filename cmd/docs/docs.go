package docs

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "docs",
		Short: "Docs",
	}

	cmd.AddCommand(initIndexCommand())
	cmd.AddCommand(initBulkCommand())

	return &cmd
}
