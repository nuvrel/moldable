package runnable

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/command"
	"github.com/nuvrel/moldable/internal/config"
	"github.com/spf13/cobra"
)

func NewRoot(l *log.Logger) command.Runnable {
	return func(cmd *cobra.Command, args []string) error {
		filepath, _ := cmd.Flags().GetString(app.ConfigFileFlag)

		cfg, err := config.LoadYaml[app.Config](filepath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if err := cfg.Check(); err != nil {
			return fmt.Errorf("checking config: %w", err)
		}

		l.Info("moldable", "config", cfg)

		return nil
	}
}
