package datastream

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/es/es"
	"github.com/zmoog/es/es/commands"
)

var force bool

func initDeleteCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete [data stream name or pattern]",
		Short: "Delete one or more data streams",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteCommand,
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return &cmd
}

func runDeleteCommand(cmd *cobra.Command, args []string) error {
	runner, err := es.NewRunner()
	if err != nil {
		return err
	}

	return runner.Run(commands.DeleteDataStreamCommand{
		Name:  args[0],
		Force: force,
	})
}
