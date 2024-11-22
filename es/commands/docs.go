package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

// IndexDocCommand represents a command to index a document in Elasticsearch.
type IndexDocCommand struct {
	Doc   string // The document to index.
	Index string // The name of the index or data stream to index the document in.
}

// ExecuteWith executes the command with the given unit of work.
func (c IndexDocCommand) ExecuteWith(uow UnitOfWork) error {
	value := c.Doc

	switch {
	case value == "-":
		reader := bufio.NewReader(os.Stdin)
		value, _ = reader.ReadString('\n')
	}

	// Create the request body.
	reqBody := strings.NewReader(value)

	// Create the index request.
	req := esapi.IndexRequest{
		Index: c.Index,
		Body:  reqBody,
	}

	// Execute the request.
	res, err := req.Do(context.Background(), uow.Client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	log.Println(res)

	return nil
}

// BulkCommand represents a command to index multiple documents in Elasticsearch.
type BulkCommand struct {
	io.Reader                // The reader containing the documents to index.
	IndexOrDataStream string // The name of the index or data stream to index the documents in.
	Action            string // The action to perform on the documents.
	NumWorkers        int    // The number of workers to use for indexing.
}

func (c BulkCommand) ExecuteWith(uow UnitOfWork) error {

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      c.IndexOrDataStream,
		NumWorkers: c.NumWorkers,
		Client:     uow.Client,
	})
	if err != nil {
		return fmt.Errorf("error creating bulk indexer: %w", err)
	}

	scanner := bufio.NewScanner(c.Reader)

	ctx := context.Background()

	for scanner.Scan() {
		line := scanner.Text()

		bulkIndexer.Add(ctx,
			esutil.BulkIndexerItem{
				Action: c.Action,
				Body:   strings.NewReader(line),
				OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
					log.Printf("Successfully indexed document %s\n", bii.DocumentID)
				},
				OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Failed to index document %s: %s\n", bii.DocumentID, err)
					}

					log.Printf("Response: status=%d, error=%s\n", biri.Status, biri.Error)
				},
			},
		)
	}

	if err := bulkIndexer.Close(ctx); err != nil {
		return fmt.Errorf("error closing bulk indexer: %w", err)
	}

	bulkIndexerStats := bulkIndexer.Stats()
	log.Printf("Indexed %d documents\n", bulkIndexerStats.NumAdded)
	log.Printf("Errors: %d\n", bulkIndexerStats.NumFailed)

	return nil
}
