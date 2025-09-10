package command

import "github.com/spf13/cobra"

type Runnable func(cmd *cobra.Command, args []string) error
