package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("22", func(t *testing.T) {
		c := RunCmd([]string{"./testdata/myecho.sh"}, Environment{})
		fmt.Println(c)
		require.Equal(t, 0, c)
	})

	t.Run("ok code", func(t *testing.T) {
		c := RunCmd([]string{"echo", "foo"}, Environment{})
		fmt.Println(c)
		require.Equal(t, 0, c)
	})

	t.Run("error code", func(t *testing.T) {
		c := RunCmd([]string{"grep", "foo"}, Environment{})
		require.Equal(t, 1, c)
	})
}
