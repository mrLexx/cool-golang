package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	ErrIsDir                 = errors.New("is a directory")
	ErrLimitNegative         = errors.New("limit negative")
	ErrNotExist              = errors.New("not exist")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegative        = errors.New("offset negative")
	ErrPathEmpty             = errors.New("file path empty")
	ErrPermissionDenied      = errors.New("access denied to file")
	ErrSomeFile              = errors.New("from and to are the same file")
	ErrUnsupportedFile       = errors.New("unsupported file")
)

func charDevice(file *os.File) (bool, error) {
	inf, err := file.Stat()
	if err != nil {
		return false, err
	}
	return (inf.Mode() & fs.ModeCharDevice) != 0, nil
}

func checkOpen(err error) error {
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotExist
		}
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}
	return nil
}

func checkFrom(file *os.File, offset int64) error {
	st, err := file.Stat()
	if err != nil {
		return fmt.Errorf("%v: %w", file.Name(), err)
	}
	if st.Mode().IsDir() {
		return fmt.Errorf("%v: %w", file.Name(), ErrIsDir)
	}
	if offset > st.Size() {
		return fmt.Errorf("%v, size - %v, offset - %v: %w", file.Name(), st.Size(), offset, ErrOffsetExceedsFileSize)
	}

	isCharDevice, err := charDevice(file)
	if err != nil {
		return fmt.Errorf("%v: %w", file.Name(), err)
	}
	if isCharDevice && limit == 0 {
		return fmt.Errorf("%v: %w", file.Name(), ErrUnsupportedFile)
	}
	return nil
}

func checkTo(file *os.File) error {
	st, err := file.Stat()
	if err != nil {
		return fmt.Errorf("%v: %w", file.Name(), err)
	}
	if st.Mode().IsDir() {
		return fmt.Errorf("%v: %w", file.Name(), ErrIsDir)
	}
	return nil
}

func checkPath(fromPath, toPath string) error {
	if fromPath == "" {
		return fmt.Errorf("file from: %w", ErrPathEmpty)
	}

	if toPath == "" {
		return fmt.Errorf("file to: %w", ErrPathEmpty)
	}

	f, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}
	t, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	if f == t {
		return fmt.Errorf("%w", ErrSomeFile)
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := checkPath(fromPath, toPath); err != nil {
		return err
	}

	if offset < 0 {
		return fmt.Errorf("offset - %v: %w", offset, ErrOffsetNegative)
	}

	if limit < 0 {
		return fmt.Errorf("limit - %v: %w", limit, ErrLimitNegative)
	}

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err := checkOpen(err); err != nil {
		return fmt.Errorf("%v: %w", fromPath, err)
	}
	defer fromFile.Close()
	if err := checkFrom(fromFile, offset); err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err := checkOpen(err); err != nil {
		return fmt.Errorf("%v: %w", toPath, err)
	}
	defer toFile.Close()
	if err := checkTo(toFile); err != nil {
		return err
	}

	err = process(fromFile, toFile, offset, limit)
	if err != nil {
		return err
	}

	return nil
}

func process(from *os.File, to io.Writer, offset, limit int64) error {
	i, err := from.Stat()
	if err != nil {
		return err
	}

	if limit == 0 {
		limit = i.Size()
	}

	if limit > i.Size()-offset {
		limit = i.Size() - offset
	}

	pb := NewProgressBar(limit)
	pb.Place()

	from.Seek(offset, io.SeekStart)

	r := io.LimitReader(from, limit)

	io.Copy(to, pb.Reader(r))
	return nil
}
