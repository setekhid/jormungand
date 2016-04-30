// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	jorm "github.com/setekhid/jormungand"
	"net/http"
)

type Entry struct {
	restC *restful.Container
	http.Server
}

func NewEntry(installers []Installer, addr string) *Entry {

	restC := restful.NewContainer()
	for _, installer := range installers {
		installer.Install(restC)
	}

	return &Entry{
		restC: restC,
		Server: http.Server{
			Addr:    addr,
			Handler: restC,
		},
	}
}

type Installer interface {
	Install(restC *restful.Container)
}

type Server struct {
	Entry  *Entry
	Router *jorm.Router
}

type ServerConf struct {
	Router  *jorm.Router
	Address string
}

func NewServer(conf ServerConf) *Server {

	installers := NewInstallers(conf.Router)
	entry := NewEntry(installers, conf.Address)
	return &Server{
		Entry:  entry,
		Router: conf.Router,
	}
}

func NewInstallers(router *jorm.Router) []Installer {

	installers := []Installer{}
	installers = append(installers, NewVeryInstaller(router))
	installers = append(installers, NewFakerInstaller())
	installers = append(installers, NewSwaggerInstaller())
	return installers
}

func (this *Server) Start() error {
	go this.Entry.ListenAndServe()
	return nil
}

func (this *Server) Stop() {
	// TODO
	glog.Warningln("Unimplementated yet!")
}
