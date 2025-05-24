package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrDirNameEmpty = errors.New("name directory empty")
	ErrOpenFile     = errors.New("open file")
	ErrReadFile     = errors.New("read file")
)

func readFile(fileName string) (EnvValue, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return EnvValue{}, fmt.Errorf("%w: %w", ErrOpenFile, err)
	}
	defer file.Close()

	inf, err := file.Stat()
	if err != nil {
		return EnvValue{}, fmt.Errorf("%w: %w", ErrOpenFile, err)
	}

	if inf.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		if err.Error() != "EOF" {
			return EnvValue{}, fmt.Errorf("%w: %w", ErrReadFile, err)
		}
	}

	line = strings.TrimRight(line, " \t\r\n")
	line = strings.ReplaceAll(line, "\x00", "\n")

	return EnvValue{
		Value: line,
	}, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if dir == "" {
		return nil, ErrDirNameEmpty
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %v: %w", dir, err)
	}

	env := Environment{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.Contains(entry.Name(), "=") {
			continue
		}

		l, err := readFile(filepath.Join(dir, entry.Name()))

		if err != nil {
			return nil, fmt.Errorf("file error %v: %w", entry.Name(), err)
		}

		env[entry.Name()] = l

	}

	return env, nil
}
