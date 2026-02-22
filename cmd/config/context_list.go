package config

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/zmoog/es/config"
)

func newGetContextsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-contexts",
		Short: "List configured contexts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if len(cfg.Contexts) == 0 {
				fmt.Println("No contexts configured.")
				return nil
			}

			// Sort names for deterministic output.
			names := make([]string, 0, len(cfg.Contexts))
			for name := range cfg.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)

			for _, name := range names {
				marker := "  "
				if name == cfg.CurrentContext {
					marker = "* "
				}
				fmt.Printf("%s%s\n", marker, name)
			}

			return nil
		},
	}
}
