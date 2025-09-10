package logger_test

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/nuvrel/moldable/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("parsing log level", func(t *testing.T) {
		t.Parallel()

		l, err := logger.New("invalid")
		assert.ErrorContains(t, err, "parsing log level")

		assert.Nil(t, l)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		l, err := logger.New("debug")
		assert.NoError(t, err)

		assert.Equal(t, log.DebugLevel, l.GetLevel())
	})
}
