package fakehttp

import (
	"context"
	"io"
)

type HTTPFaker interface {
	Request() *RequestHeader
	Response(*RequestHeader) *ResponseHeader
}

type DefaultHTTPFaker struct{}

func (f DefaultHTTPFaker) Request() *RequestHeader {
	// TODO
	return nil
}

func (f DefaultHTTPFaker) Response(req *RequestHeader) *ResponseHeader {
	// TODO
	return nil
}

type RequestHeader struct {
	Method, URI, Version string
	Fields               map[string]string
}

func ReadRequestHeader(r io.Reader) (*RequestHeader, error) {
	// TODO
	return nil, nil
}

func (r *RequestHeader) SetCookie(cookie string) {
	// TODO
}

func (r *RequestHeader) GetCookie() string {
	// TODO
	return ""
}

func (r *RequestHeader) String() string {
	// TODO
	return ""
}

type ResponseHeader struct {
	Version    string
	StatusCode int
	Fields     map[string]string
	HasContent bool
}

func ReadResponseHeader(r io.Reader) (*ResponseHeader, error) {
	// TODO
	return nil, nil
}

func (r *ResponseHeader) String() string {
	// TODO
	return ""
}

type fakeHTTPServant struct {
	faker   HTTPFaker
	payload io.ReadWriteCloser
}

func (f fakeHTTPServant) servant(ctx context.Context, conn HDXConn) {
	// TODO
}

func NewHTTPServant(
	faker HTTPFaker, payload io.ReadWriteCloser,
) func(context.Context, HDXConn) {
	return fakeHTTPServant{faker, payload}.servant
}

func NewDefaultHTTPServant(
	payload io.ReadWriteCloser,
) func(context.Context, HDXConn) {
	return NewHTTPServant(DefaultHTTPFaker{}, payload)
}
