package command

import (
	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewInit(r command.Runnable) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a starter moldable.yaml in the current directory",
		Long:  `Creates a starter moldable.yaml in the current directory unless one already exists.`,
		Args:  cobra.NoArgs,
		RunE:  r,
	}

	{
		fs := new(pflag.FlagSet)

		fs.BoolP(app.ForceFlag, "f", false, "force creation of config file")

		cmd.Flags().AddFlagSet(fs)
	}

	return cmd
}
