package output

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"io"
	"net"
)

type Tcp struct {
	addr    net.Addr
	backoff backoff.BackOff
	conn    net.Conn
}

func NewTcpWriteCloser(addr net.Addr, bo backoff.BackOff) (io.WriteCloser, error) {
	tcp := &Tcp{
		addr:    addr,
		backoff: bo,
	}

	err := tcp.Connect()
	if err != nil {
		return nil, err
	}

	return tcp, nil
}

func (tcp *Tcp) Connect() (err error) {
	tcp.conn, err = net.Dial(tcp.addr.Network(), tcp.addr.String())
	if err != nil {
		return err
	}
	return nil
}

func (tcp *Tcp) Write(buf []byte) (int, error) {
	bytesWritten := 0
	bytesToWrite := len(buf)

	err := backoff.Retry(func() error {
		if tcp.conn == nil {
			if err := tcp.Connect(); err != nil {
				return err
			}
		}

		for bytesWritten < bytesToWrite {
			n, err := tcp.conn.Write(buf)
			bytesWritten += n
			if err != io.ErrShortWrite {
				_ = tcp.Close()
				return err
			}
		}

		return nil
	}, tcp.backoff)

	return bytesWritten, err
}

func (tcp *Tcp) Close() error {
	if tcp.conn == nil {
		return fmt.Errorf("no connection to close")
	}

	err := tcp.conn.Close()
	tcp.conn = nil
	return err
}
