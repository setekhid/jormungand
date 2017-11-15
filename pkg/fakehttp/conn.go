package fakehttp

import (
	"io"
	"net"
)

type HDXConn interface {
	StartReading() error
	io.Reader

	StartWriting() error
	io.Writer

	io.Closer

	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

func Dial(network, address string) (HDXConn, error) {
	// TODO
	return nil, nil
}

type HDXListener interface {
	Accept() (HDXConn, error)
	io.Closer

	Addr() net.Addr
}

func Listen(network, address string) HDXListener {
	// TODO
	return nil
}
