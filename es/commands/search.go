package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type SearchCommand struct {
	Index string
	Query string
}

func (c SearchCommand) ExecuteWith(uow UnitOfWork) error {
	res, err := uow.Client.Search(
		uow.Client.Search.WithIndex(c.Index),
		uow.Client.Search.WithBody(strings.NewReader(c.Query)),
	)
	if err != nil {
		return fmt.Errorf("error searching: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error searching: %s", res.Status())
	}

	// write the body to stdout
	io.Copy(os.Stdout, res.Body)

	return nil
}
