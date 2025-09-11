package main

import (
	"log"
	"os"

	"github.com/nuvrel/moldable/cmd/moldable/app/command"
	"github.com/nuvrel/moldable/cmd/moldable/app/runnable"
	"github.com/nuvrel/moldable/pkg/logger"
	"github.com/nuvrel/moldable/pkg/version"
)

func main() {
	l, err := logger.New("info")
	if err != nil {
		log.Fatalf("creating logger: %v", err)
	}

	l.SetReportTimestamp(true)

	root := command.NewRoot(runnable.NewRoot(l))
	version := version.NewCommand(root.OutOrStderr())
	init := command.NewInit(runnable.NewInit(l))

	root.AddCommand(version, init)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
