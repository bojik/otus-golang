package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	ConnectContext(context.Context) error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	conn       net.Conn
	connected  bool
	context    context.Context
	cancelFunc context.CancelFunc
}

func (c *client) Connect() error {
	return c.ConnectContext(context.Background())
}

func (c *client) ConnectContext(ctx context.Context) error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return c.newError(err)
	}
	c.conn = conn
	c.connected = true
	c.context, c.cancelFunc = context.WithCancel(ctx)
	return nil
}

func (c *client) Close() error {
	if !c.connected {
		return nil
	}
	c.cancelFunc()
	c.connected = false
	if err := c.conn.Close(); err != nil {
		return c.newError(err)
	}
	return nil
}

func (c *client) Send() error {
	if !c.connected {
		return c.newError(ErrIsNotConnected)
	}
	go c.copy(c.conn, c.in)
	return nil
}

func (c *client) Receive() error {
	if !c.connected {
		return c.newError(ErrIsNotConnected)
	}
	go c.copy(c.out, c.conn)
	return nil
}

func (c *client) copy(writer io.Writer, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
OUT:
	for {
		select {
		case <-c.context.Done():
			break OUT
		default:
			if !scanner.Scan() {
				break OUT
			}
			bytes := append(scanner.Bytes(), '\n')
			_, err := writer.Write(bytes)
			if err != nil {
				break OUT
			}
		}
	}
	c.cancelFunc()
}

func (c *client) newError(err error) error {
	return newTelnetError(c.address, err)
}

var _ TelnetClient = (*client)(nil)
