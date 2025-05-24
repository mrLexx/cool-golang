package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorDir(t *testing.T) {
	t.Run("name directory empty", func(t *testing.T) {
		_, err := ReadDir("")
		require.ErrorIs(t, err, ErrDirNameEmpty)
	})

	t.Run("no such directory", func(t *testing.T) {
		dir := t.TempDir()
		os.Remove(dir)

		_, err := ReadDir(dir)
		require.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("not a directory", func(t *testing.T) {
		file, err := os.CreateTemp("", "from.txt")
		require.NoError(t, err)

		_, err = ReadDir(file.Name())
		require.ErrorContains(t, err, "not a directory")
	})

	t.Run("permission denied", func(t *testing.T) {
		dir, err := os.Open(t.TempDir())
		require.NoError(t, err)

		err = dir.Chmod(0000)
		require.NoError(t, err)

		_, err = ReadDir(dir.Name())
		require.ErrorContains(t, err, "permission denied")
	})

}

func TestErrorFile(t *testing.T) {
	t.Run("permission denied", func(t *testing.T) {
		dir, err := os.Open(t.TempDir())
		require.NoError(t, err)

		fromFile, err := os.CreateTemp(dir.Name(), "from.txt")
		require.NoError(t, err)

		err = fromFile.Chmod(0000)
		require.NoError(t, err)

		_, err = ReadDir(dir.Name())
		require.ErrorContains(t, err, ErrOpenFile.Error())
		require.ErrorIs(t, err, ErrOpenFile)
		require.ErrorContains(t, err, "permission denied")
	})
}

func TestRead(t *testing.T) {
	t.Run("load env", func(t *testing.T) {
		test := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, test)
	})

	t.Run("skip dir", func(t *testing.T) {
		err := os.Mkdir("./testdata/env/skipdir", 0755)
		require.NoError(t, err)

		test := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, test)
		err = os.Remove("./testdata/env/skipdir")
		require.NoError(t, err)
	})

	t.Run("skip =", func(t *testing.T) {
		file, err := os.Create("./testdata/env/FILE=NAME")
		require.NoError(t, err)
		err = file.Close()
		require.NoError(t, err)

		test := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, test)

		err = os.Remove("./testdata/env/FILE=NAME")
		require.NoError(t, err)
	})

}

func TestReadDir(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		for k, v := range env {
			fmt.Println(k, ": ", v)
		}

		require.NoError(t, err)
	})

}
