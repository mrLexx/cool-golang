package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("ok code", func(t *testing.T) {
		c := RunCmd([]string{"echo"}, Environment{})
		require.Equal(t, 0, c)
	})

	t.Run("not found cmd", func(t *testing.T) {
		c := RunCmd([]string{"./testdata/anyscript.sh"}, Environment{})
		require.Equal(t, 127, c)
	})

	t.Run("error code", func(t *testing.T) {
		c := RunCmd([]string{"ls", "/not_exists_dir/"}, Environment{})
		require.Equal(t, 2, c)
	})
}
