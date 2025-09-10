package logger

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func New(level string) (*log.Logger, error) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("parsing log level: %w", err)
	}

	l := log.NewWithOptions(os.Stderr, log.Options{
		Level: lvl,
	})

	return l, nil
}
