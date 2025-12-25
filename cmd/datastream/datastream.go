package datastream

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "datastream",
		Short: "Data stream operations",
	}

	cmd.AddCommand(initDeleteCommand())

	return &cmd
}
