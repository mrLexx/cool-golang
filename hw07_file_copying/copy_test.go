package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorFrom(t *testing.T) {
	t.Run("char device urandom", func(t *testing.T) {
		from := "/dev/urandom"
		toFile, err := os.CreateTemp("", "out.txt")
		require.Nil(t, err)

		err = Copy(from, toFile.Name(), 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)

	})

	t.Run("char device null", func(t *testing.T) {
		from := "/dev/null"
		toFile, err := os.CreateTemp("", "out.txt")
		require.Nil(t, err)

		err = Copy(from, toFile.Name(), 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)

	})

	t.Run("not exist file", func(t *testing.T) {
		from := "randomfile.txt"
		toFile, err := os.CreateTemp("", "out.txt")
		require.Nil(t, err)

		err = Copy(from, toFile.Name(), 0, 0)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("copy from dir", func(t *testing.T) {
		from := t.TempDir()
		toFile, err := os.CreateTemp("", "out.txt")
		require.Nil(t, err)

		err = Copy(from, toFile.Name(), 0, 0)
		require.ErrorIs(t, err, ErrIsDir)
	})

	t.Run("access denied to file", func(t *testing.T) {
		from, err := os.CreateTemp("", "from.txt")
		require.Nil(t, err)
		os.Chmod(from.Name(), 0000)

		toFile, err := os.CreateTemp("", "out.txt")
		require.Nil(t, err)

		err = Copy(from.Name(), toFile.Name(), 0, 0)
		require.ErrorIs(t, err, ErrPermissionDenied)
	})
}
