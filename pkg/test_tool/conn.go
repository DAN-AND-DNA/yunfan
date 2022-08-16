package test_tool

import (
	"bytes"
	"net"
	"time"
)

type Test_addr string

func (a Test_addr) Network() string {
	return string(a)
}

func (a Test_addr) String() string {
	return string(a)
}

type Test_conn struct {
	r bytes.Buffer
	w bytes.Buffer
}

func (c *Test_conn) Read(b []byte) (int, error)       { return c.r.Read(b) }
func (c *Test_conn) Write(b []byte) (int, error)      { return c.w.Write(b) }
func (c *Test_conn) Close() error                     { return nil }
func (c *Test_conn) LocalAddr() net.Addr              { return Test_addr("local-addr") }
func (c *Test_conn) RemoteAddr() net.Addr             { return Test_addr("remote-addr") }
func (c *Test_conn) SetDeadline(time.Time) error      { return nil }
func (c *Test_conn) SetReadDeadline(time.Time) error  { return nil }
func (c *Test_conn) SetWriteDeadline(time.Time) error { return nil }

func (c *Test_conn) Reset() {
	c.r.Reset()
	c.w.Reset()
}

func (c *Test_conn) Input_req(bytes_req []byte) error {
	if _, err := c.r.Write(bytes_req); err != nil {
		return err
	}

	return nil
}

func (c *Test_conn) Output_resp(buffer []byte) (int, error) {
	n, err := c.w.Read(buffer)
	if err != nil {
		return 0, err
	}

	return n, nil
}
