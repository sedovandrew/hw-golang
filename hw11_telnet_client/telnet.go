package main

import (
	"io"
	"net"
	"time"
)

const protocol = "tcp"

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tc := TClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return &tc
}

type TClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *TClient) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout(protocol, tc.address, tc.timeout)
	return err
}

func (tc TClient) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}

func (tc TClient) Close() error {
	err := tc.conn.Close()
	if err != nil {
		return err
	}

	err = tc.in.Close()
	return err
}
