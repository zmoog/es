package commands

import (
	"fmt"

	"github.com/zmoog/classeviva/adapters/feedback"
)

var (
	version string = "latest"
	commit  string = "unknown"
	date    string = "2025-01-07"
	builtBy string = "zmoog"
)

type VersionCommand struct{}

func (c VersionCommand) ExecuteWith(uow UnitOfWork) error {
	return feedback.PrintResult(VersionResult{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	})
}

type VersionResult struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"built_by"`
}

func (r VersionResult) String() string {
	return fmt.Sprintf("%v (%v) %v by %v", version, commit, date, builtBy)
}

func (r VersionResult) Data() interface{} {
	return r
}
