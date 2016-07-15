// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/setekhid/jormungand/http/comet"
	"github.com/setekhid/jormungand/jorm/sel"
)

var (
	jman = (*comet.JungleMan)(nil)

	installers = []Installer{
		&VeryInstaller{},
		&ManagerInstaller{},
		&FakerInstaller{},
		&SwaggerInstaller{},
	}
)

func Installers() []Installer {

	if jman == nil {
		jman = &comet.JungleMan{Author: sel.Router()}
	}
	return installers
}

type VeryInstaller struct {
}

// installing exchange and download api, for vpn
func (ins *VeryInstaller) Install(restC *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path("/op").
		Consumes(restful.MIME_OCTET).
		Produces(restful.MIME_OCTET)

	ws.Route(
		ws.POST("/xch/{dumpFile}").To(ins.exchange).
			Doc("exchange a file with server").
			Operation("exchange").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	ws.Route(
		ws.GET("/dl/{dumpFile}").To(ins.download).
			Doc("download a file from server").
			Operation("download").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	restC.Add(ws)
}

// exchange api
func (ins *VeryInstaller) exchange(_req *restful.Request, _resp *restful.Response) {
	token := _req.PathParameter("dumpFile")
	_ = token
	jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

// download api
func (ins *VeryInstaller) download(_req *restful.Request, _resp *restful.Response) {
	token := _req.PathParameter("dumpFile")
	_ = token
	jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

type FakerInstaller struct {
}

// a fake installer
func (ins *FakerInstaller) Install(restC *restful.Container) {

	// TODO
}

type SwaggerInstaller struct {
}

// swagger
func (ins *SwaggerInstaller) Install(restC *restful.Container) {
	swagger.RegisterSwaggerService(swagger.Config{
		WebServices: restC.RegisteredWebServices(),
		ApiPath:     "/apidocs.json",
	}, restC)
}
