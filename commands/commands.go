package commands

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Command interface {
	ExecuteWith(uow UnitOfWork) error
}

type UnitOfWork struct {
	Client *elasticsearch.Client
}
