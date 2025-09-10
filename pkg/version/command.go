package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func NewCommand(w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `Print the version and build information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(w, "%s version %s\n", cmd.Root().Name(), cmd.Root().Version)

			return nil
		},
	}
}
