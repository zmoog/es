package commands

import (
	"github.com/elastic/go-elasticsearch/v8"
)

// Command represents a command that can be executed.
type Command interface {
	ExecuteWith(uow UnitOfWork) error
}

// UnitOfWork is contains the dependencies needed to execute commands.
type UnitOfWork struct {
	Client *elasticsearch.Client
}
