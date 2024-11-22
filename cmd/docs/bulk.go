package docs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zmoog/es/es"
	"github.com/zmoog/es/es/commands"
)

var (
	filePath          string
	indexOrDataStream string
	action            string
	numWorkers        int
)

func initBulkCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "bulk",
		Short: "Index multiple documents in Elasticsearch using a bulk request",
		RunE:  runBulkCommand,
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "-", "File containing the documents to index")
	cmd.Flags().StringVarP(&indexOrDataStream, "index", "i", "logs-generic-default", "Index or data stream name")
	cmd.Flags().StringVarP(&action, "action", "a", "create", "Action to perform on the documents")
	cmd.Flags().IntVarP(&numWorkers, "workers", "w", runtime.NumCPU(), "Number of workers to use for indexing")

	return &cmd
}

// runBulkCommand runs the bulk command.
func runBulkCommand(cmd *cobra.Command, args []string) error {
	runner, err := es.NewRunner()
	if err != nil {
		return err
	}

	var reader io.Reader

	switch filePath {
	case "-":
		reader = bufio.NewReader(os.Stdin)
	default:
		reader, err = os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file [%s]: %w", filePath, err)
		}
	}

	return runner.Run(commands.BulkCommand{
		Reader:            reader,
		IndexOrDataStream: indexOrDataStream,
		Action:            action,
	})
}
