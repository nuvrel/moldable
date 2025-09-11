package runnable

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/charmbracelet/log"
	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/command"
	"github.com/spf13/cobra"
)

func NewInit(l *log.Logger) command.Runnable {
	return func(cmd *cobra.Command, args []string) error {
		forced, _ := cmd.Flags().GetBool(app.ForceFlag)

		if _, err := os.Stat(app.ConfigFile); !os.IsNotExist(err) && !forced {
			l.Info("config file already exists, skipping creation", "filepath", app.ConfigFile)

			return nil
		}

		tmpl, err := template.ParseFS(app.Templates, "templates/*.tmpl")
		if err != nil {
			return fmt.Errorf("parsing templates: %w", err)
		}

		var buf bytes.Buffer

		if err := tmpl.ExecuteTemplate(&buf, app.ConfigFileTemplate, nil); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}

		if err := os.WriteFile(app.ConfigFile, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("writing config file to disk: %w", err)
		}

		l.Info("config file written", "filepath", app.ConfigFile)

		return nil
	}
}
