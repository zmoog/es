package commands

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

type DeleteDataStreamCommand struct {
	Name  string
	Force bool
}

func (c DeleteDataStreamCommand) ExecuteWith(uow UnitOfWork) error {
	// Check if confirmation is needed
	if !c.Force {
		confirmed, err := confirmDeletion(c.Name)
		if err != nil {
			return fmt.Errorf("error reading confirmation: %w", err)
		}
		if !confirmed {
			log.Println("Deletion cancelled")
			return nil
		}
	}

	// Execute the delete operation
	res, err := uow.Client.Indices.DeleteDataStream(
		[]string{c.Name},
		uow.Client.Indices.DeleteDataStream.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("error deleting data stream: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.Status())
	}

	log.Printf("Successfully deleted data stream(s): %s\n", c.Name)
	return nil
}

func confirmDeletion(name string) (bool, error) {
	var message string
	if strings.Contains(name, "*") {
		message = fmt.Sprintf("Are you sure you want to delete data streams matching '%s'? This will delete all matching data streams and their backing indices. [y/N]: ", name)
	} else {
		message = fmt.Sprintf("Are you sure you want to delete data stream '%s'? This will delete the data stream and all its backing indices. [y/N]: ", name)
	}

	fmt.Print(message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}
