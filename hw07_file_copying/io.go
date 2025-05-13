package main

import (
	"io"
)

type Reader struct {
	io.Reader
	pb *ProgressBar
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		r.pb.Finish()
	}
	r.pb.Add(n)
	return
}

func (r *Reader) Close() (err error) {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}
