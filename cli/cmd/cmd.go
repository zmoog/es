package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/es/cli/cmd/docs"
	"github.com/zmoog/es/cli/cmd/version"
)

func Execute() {
	rootCmd := cobra.Command{
		Use: "elasticsearch-cli",
	}

	rootCmd.AddCommand(docs.NewCommand())

	rootCmd.AddCommand(version.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		feedback.Error(err)
	}
}
