package commands

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
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
