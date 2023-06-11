package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/elasticsearch-cli/cli/cmd/docs"
	"github.com/zmoog/elasticsearch-cli/cli/cmd/version"
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
