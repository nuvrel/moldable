package runnable

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/command"
	"github.com/nuvrel/moldable/internal/config"
	"github.com/nuvrel/moldable/internal/generator"
	"github.com/nuvrel/moldable/internal/reporter"
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

		gen := generator.New(cfg, reporter.NewLog(l))

		if err := gen.Generate(); err != nil {
			return fmt.Errorf("generating interfaces: %w", err)
		}

		l.Info("interfaces generated successfully")

		return nil
	}
}
