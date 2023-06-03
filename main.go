package main

import (
	"log"
	"os"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func main() {

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

	cfg := elasticsearch.Config{
		Addresses: strings.Split(endpoints, ","),
		APIKey:    apiKey,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		os.Exit(1)
	}

	log.Println(elasticsearch.Version)
	log.Println(client.Info())
}
