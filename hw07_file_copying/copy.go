package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrPermissionDenied      = errors.New("access denied to file")
	ErrNotExist              = errors.New("not exist")
	ErrIsDir                 = errors.New("is a directory")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegative        = errors.New("offset negative")
	ErrLimitNegative         = errors.New("limit negative")
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

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	pb := NewProgressBar(limit - offset)
	pb.Place()

	from.Seek(offset, io.SeekStart)

	r := io.LimitReader(from, limit)

	io.Copy(to, pb.Reader(r))
	return nil
}
