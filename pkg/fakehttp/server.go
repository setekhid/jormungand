package fakehttp

import (
	"context"
)

type HDXServer struct {
	Servant func(context.Context, HDXConn)
}

func (s *HDXServer) ListenAndServe(network, address string) error {
	// TODO
	return s.Serve(nil)
}

func (s *HDXServer) Serve(l HDXListener) error {

	defer l.Close()

	for {

		conn, err := l.Accept()
		if err != nil {
			// TODO
			break
		}

		go s.Servant(nil, conn)
	}

	return nil
}

func (s *HDXServer) Shutdown() error {
	// TODO
	return nil
}

func (s *HDXServer) Close() error {
	// TODO
	return nil
}
