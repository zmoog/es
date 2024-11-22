package docs

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/es/es"
	"github.com/zmoog/es/es/commands"
)

var (
	doc   string
	index string
)

func initIndexCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "index",
		Short: "Index a document in Elasticsearch",
		RunE:  runIndexCommand,
	}

	cmd.Flags().StringVar(&doc, "doc", "-", "Document to index")
	cmd.Flags().StringVar(&index, "index", "logs-generic-default", "Index or data stream name")

	return &cmd
}

func runIndexCommand(cmd *cobra.Command, args []string) error {
	runner, err := es.NewRunner()
	if err != nil {
		return err
	}

	return runner.Run(commands.IndexDocCommand{
		Doc:   doc,
		Index: index,
	})
}
