package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrPermissionDenied      = errors.New("access denied to file")
	ErrNotExist              = errors.New("not exist")
	ErrIsDir                 = errors.New("is a directory")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func checkSize(file *os.File, offset int64) error {
	st, _ := file.Stat()
	switch {
	case st.Mode().IsDir():
		return ErrIsDir
	case offset > st.Size():
		return ErrOffsetExceedsFileSize

	}
	return nil
}

func checkOpen(err error) error {
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotExist
		}
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
	}
	return nil
}

func checkCharDevice(file *os.File) (bool, error) {

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return (fileInfo.Mode() & fs.ModeCharDevice) != 0, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err := checkOpen(err); err != nil {
		return fmt.Errorf("%v: %w", fromPath, err)
	}
	defer close(fromFile)

	switch st, _ := fromFile.Stat(); {
	case st.Mode().IsDir():
		return fmt.Errorf("%v: %w", fromPath, ErrIsDir)
	case offset > st.Size():
		return fmt.Errorf("%v, size - %v, offset - %v: %w", fromPath, st.Size(), offset, ErrOffsetExceedsFileSize)
	}

	isCharDevice, err := checkCharDevice(fromFile)
	if err != nil {
		return fmt.Errorf("%v: %w", fromPath, err)
	}
	if isCharDevice && limit == 0 {
		return fmt.Errorf("%v: %w", fromPath, ErrUnsupportedFile)
	}

	toFile, err := os.Create(toPath)
	if err := checkOpen(err); err != nil {
		return fmt.Errorf("%v: %w", toFile, err)
	}
	defer close(toFile)

	switch st, _ := toFile.Stat(); {
	case st.Mode().IsDir():
		return fmt.Errorf("%v: %w", fromPath, ErrIsDir)
	}

	process(fromFile, toFile, offset, limit)

	return nil
}

func process(from *os.File, to io.Writer, offset, limit int64) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed: ", err, offset, limit)
		}
	}()
	i, _ := from.Stat()

	if limit == 0 {
		limit = i.Size()
	}

	pb := NewProgressBar(limit - offset)
	pb.Place()

	from.Seek(offset, io.SeekStart)

	r := io.LimitReader(from, limit)

	io.Copy(to, pb.Reader(r))
}

func close(f *os.File) {
	if err := f.Close(); err != nil {
		slog.Error("close file", slog.String("name", f.Name()))
	}
}
