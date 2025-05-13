package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorFrom(t *testing.T) {
	t.Run("from and to are some file", func(t *testing.T) {
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy(to, to, 0, 0)
		require.ErrorIs(t, err, ErrSomeFile)
	})

	t.Run("empty path", func(t *testing.T) {
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy("", to, 0, 0)
		require.ErrorIs(t, err, ErrPathEmpty)
	})

	t.Run("char device urandom", func(t *testing.T) {
		from := "/dev/urandom"
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("char device null", func(t *testing.T) {
		from := "/dev/null"
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("not exist file", func(t *testing.T) {
		from := "randomfile.txt"
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("copy from dir", func(t *testing.T) {
		from := t.TempDir()
		to := filepath.Join(t.TempDir(), "out.txt")

		err := Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrIsDir)
	})

	t.Run("access denied to file", func(t *testing.T) {
		fromFile, err := os.CreateTemp("", "from.txt")
		require.NoError(t, err)

		//nolint:gofumpt
		fromFile.Chmod(0000)
		from := fromFile.Name()

		to := filepath.Join(t.TempDir(), "out.txt")

		err = Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrPermissionDenied)
	})
}

func TestErrorTo(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		fromFile, err := os.CreateTemp("", "from.txt")
		require.NoError(t, err)
		from := fromFile.Name()

		err = Copy(from, "", 0, 0)
		require.ErrorIs(t, err, ErrPathEmpty)
	})

	t.Run("copy to dir", func(t *testing.T) {
		fromFile, err := os.CreateTemp("", "from.txt")
		require.NoError(t, err)
		from := fromFile.Name()

		to := t.TempDir()

		err = Copy(from, to, 0, 0)
		require.ErrorContains(t, err, ErrIsDir.Error())
	})

	t.Run("access denied to file", func(t *testing.T) {
		fromFile, err := os.CreateTemp("", "from.txt")
		require.NoError(t, err)
		from := fromFile.Name()

		toFile, err := os.CreateTemp("", "out.txt")
		require.NoError(t, err)

		//nolint:gofumpt
		toFile.Chmod(0000)
		to := toFile.Name()

		err = Copy(from, to, 0, 0)
		require.ErrorIs(t, err, ErrPermissionDenied)
	})
}

func TestErrorCopy(t *testing.T) {
	from := "testdata/input.txt"

	tests := []struct {
		offset int64
		limit  int64
		err    error
	}{
		{116000, 1000, ErrOffsetExceedsFileSize},
		{-1, 1000, ErrOffsetNegative},
		{0, -1, ErrLimitNegative},
	}

	for _, tc := range tests {
		offset, limit, dstErr := tc.offset, tc.limit, tc.err

		t.Run(dstErr.Error(), func(t *testing.T) {
			to := filepath.Join(t.TempDir(), "out.txt")
			err := Copy(from, to, offset, limit)
			require.ErrorIs(t, err, dstErr)
		})
	}
}

func TestCopy(t *testing.T) {
	from := "testdata/input.txt"

	tests := []struct {
		offset int64
		limit  int64
	}{
		{0, 0},
		{0, 10},
		{0, 1000},
		{0, 10000},
		{100, 1000},
		{6000, 1000},
	}

	for _, tc := range tests {
		offset, limit := tc.offset, tc.limit

		t.Run(fmt.Sprintf("offset: %d, limit: %d", offset, limit), func(t *testing.T) {
			goldenFileName := filepath.Join("testdata", fmt.Sprintf("out_offset%d_limit%d.txt", offset, limit))

			to := filepath.Join(t.TempDir(), "out.txt")

			err := Copy(from, to, offset, limit)
			require.ErrorIs(t, err, nil)

			s, err := os.ReadFile(to)
			require.NoError(t, err)

			g, err := os.ReadFile(goldenFileName)
			require.NoError(t, err)

			require.Equal(t, string(s), string(g))
		})
	}
}
