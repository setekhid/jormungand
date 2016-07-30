// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
	"github.com/setekhid/jormungand/http/comet"
	"github.com/setekhid/jormungand/jorm/sel"
)

var (
	jman = (*comet.JungleMan)(nil)

	installers = []Installer{
		&veryInstaller{},
		&managerInstaller{},
		&fakerInstaller{},
	}
)

func Installers() []Installer {

	if jman == nil {
		jman = &comet.JungleMan{Author: sel.Router()}
	}
	return installers
}

type veryInstaller struct {
}

// installing exchange and download api, for vpn
func (ins *veryInstaller) Install(restC *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path(URI_UPDOWN_PATH).
		Consumes(restful.MIME_OCTET).
		Produces(restful.MIME_OCTET)

	ws.Route(
		ws.POST("/{dumpFile}").To(ins.exchange).
			Doc("exchange a file with server").
			Operation("exchange").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	ws.Route(
		ws.GET("/{dumpFile}").To(ins.download).
			Doc("download a file from server").
			Operation("download").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	restC.Add(ws)
}

// exchange api
func (ins *veryInstaller) exchange(_req *restful.Request, _resp *restful.Response) {

	token := _req.PathParameter("dumpFile")
	_ = token
	jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

// download api
func (ins *veryInstaller) download(_req *restful.Request, _resp *restful.Response) {

	token := _req.PathParameter("dumpFile")
	_ = token
	jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

type fakerInstaller struct {
}

// a fake installer
func (ins *fakerInstaller) Install(restC *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path(URI_UPDOWN_PATH).
		Consumes(restful.MIME_OCTET).
		Produces(restful.MIME_OCTET)

		// TODO

	restC.Add(ws)
}
