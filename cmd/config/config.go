// Package config implements the 'es config' command group.
package config

import (
	"github.com/spf13/cobra"
)

// NewCommand returns the 'es config' command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage es configuration",
		// Override PersistentPreRun so that config subcommands do not require
		// an active Elasticsearch connection to be configured.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(newGetContextsCommand())
	cmd.AddCommand(newSetContextCommand())
	cmd.AddCommand(newUseContextCommand())

	return cmd
}
