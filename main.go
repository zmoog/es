package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {

	//
	// Get the Elasticsearch client settings from the environment variables.
	//

	endpoints, ok := os.LookupEnv("ELASTICSEARCH_ENDPOINTS")
	if !ok {
		log.Fatal("ELASTICSEARCH_ENDPOINTS is not set")
		os.Exit(1)
	}

	apiKey, ok := os.LookupEnv("ELASTICSEARCH_API_KEY")
	if !ok {
		log.Fatal("ELASTICSEARCH_API_KEY is not set")
		os.Exit(1)
	}

	//
	// Create the Elasticsearch client.
	//

	cfg := elasticsearch.Config{
		Addresses: strings.Split(endpoints, ","),
		APIKey:    apiKey,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		os.Exit(1)
	}

	//
	// Get the Elasticsearch client info and cluster info.
	//

	// log.Println(elasticsearch.Version)

	// info, err := client.Info()
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// 	os.Exit(1)
	// }
	// defer info.Body.Close()
	// log.Println(info)

	//
	// Index a document in Elasticsearch from stdin.
	//

	// Read from stdin until a newline is encountered.
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	log.Println(text)

	// Build the request.
	req := esapi.IndexRequest{
		Index: "logs-generic-default",
		Body:  strings.NewReader(text),
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		os.Exit(1)
	}

	log.Println(res)
}
