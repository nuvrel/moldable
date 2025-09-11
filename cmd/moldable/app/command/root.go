package command

import (
	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/command"
	"github.com/nuvrel/moldable/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRoot(r command.Runnable) *cobra.Command {
	cmd := &cobra.Command{
		Use:          app.Name,
		Short:        "Builds precise interfaces from any package",
		Long:         `Builds precise interfaces from any package so you can plug in mockery, gomock, moq, or any other mock tool you like.`,
		Args:         cobra.NoArgs,
		Version:      version.Get().String(),
		RunE:         r,
		SilenceUsage: true,
	}

	{
		fs := new(pflag.FlagSet)

		fs.StringP(app.ConfigFileFlag, "c", app.ConfigFile, "path to config file")

		cmd.Flags().AddFlagSet(fs)
	}

	return cmd
}
