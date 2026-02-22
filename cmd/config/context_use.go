package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/es/config"
)

func newUseContextCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "use-context <name>",
		Short: "Set the current context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if _, ok := cfg.Contexts[name]; !ok {
				return fmt.Errorf("context %q not found", name)
			}

			cfg.CurrentContext = name

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Printf("Switched to context %q.\n", name)
			return nil
		},
	}
}
