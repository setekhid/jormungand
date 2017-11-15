package fakehttp

import (
	"io"
)

type HDXClient struct {
	Conn HDXConn
}

func (c *HDXClient) DoReceive(out io.WriteCloser) {
	// TODO
}

func (c *HDXClient) DoSend(in io.Reader) {
	// TODO
}
