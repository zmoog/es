package config

import (
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage contexts",
	}

	cmd.AddCommand(newContextSetCommand())
	cmd.AddCommand(newContextUseCommand())
	cmd.AddCommand(newContextListCommand())

	return cmd
}
