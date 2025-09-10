package version_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nuvrel/moldable/pkg/version"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		name := "test"

		root := &cobra.Command{
			Use:     name,
			Version: version.Get().String(),
		}

		root.SetArgs([]string{"version"})

		var buf bytes.Buffer

		c := version.NewCommand(&buf)

		root.AddCommand(c)

		err := c.Execute()

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), fmt.Sprintf("%s version development, commit none, built at unknown.", name))
	})
}
