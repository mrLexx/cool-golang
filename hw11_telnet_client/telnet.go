package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnet struct {
	address string
	timeout time.Duration
	retries int
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{address: address, timeout: timeout, in: in, out: out}
}

func (t *telnet) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return
}

func (t *telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err

}

func (t *telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnet) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
