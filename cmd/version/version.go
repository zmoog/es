package version

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/es/es"
	"github.com/zmoog/es/es/commands"
)

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		RunE:  runVersionCommand,
	}
	return &cmd
}

func runVersionCommand(cmd *cobra.Command, args []string) error {

	runner, err := es.NewRunner()
	if err != nil {
		return err
	}

	return runner.Run(commands.VersionCommand{})
}
