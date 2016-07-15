// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
	"github.com/setekhid/jormungand/misc/jargs"
	"net"
	"net/http"
	"time"
)

const (
	moduleName = "web"
)

var (
	entry = (*EntryServ)(nil)
)

// singleton web server entry
func Entry() *EntryServ {

	if entry == nil {

		module := jargs.Module(moduleName).(*EntryConfig)

		var err error
		entry, err = NewEntryServ(module, Installers())
		if err != nil {
			panic(err)
		}
	}
	return entry
}

type EntryConfig struct {
	Addr string `json:"addr"`
}

type EntryServ struct {
	restC *restful.Container
	webLn net.Listener
	http.Server
}

func NewEntryServ(conf *EntryConfig, installers []Installer) (*EntryServ, error) {

	restC := restful.NewContainer()
	for _, installer := range installers {
		installer.Install(restC)
	}

	webLn, err := webListener(conf.Addr)
	if err != nil {
		return nil, err
	}

	return &EntryServ{
		restC: restC,
		webLn: webLn,
		Server: http.Server{
			Addr:    conf.Addr,
			Handler: restC,
		},
	}, nil
}

func (this *EntryServ) Listener() net.Listener { return this.webLn }
func (this *EntryServ) ListenAndServe() error  { return this.Serve(this.webLn) }
func (this *EntryServ) Close() error           { return this.webLn.Close() }

type Installer interface {
	Install(restC *restful.Container)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func webListener(addr string) (net.Listener, error) {

	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return tcpKeepAliveListener{ln.(*net.TCPListener)}, nil
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
