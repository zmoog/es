package commands

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type IndexDocCommand struct {
	Doc   string
	Index string
}

func (c IndexDocCommand) ExecuteWith(uow UnitOfWork) error {
	value := c.Doc

	switch {
	case value == "-":
		reader := bufio.NewReader(os.Stdin)
		value, _ = reader.ReadString('\n')
	}

	// Build the request.
	req := esapi.IndexRequest{
		Index: c.Index,
		Body:  strings.NewReader(value),
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), uow.Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		os.Exit(1)
	}

	log.Println(res)

	return nil
}
