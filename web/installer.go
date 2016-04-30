// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/setekhid/jormungand/http/comet"
)

type VeryInstaller struct {
	jman comet.JungleMan
}

func NewVeryInstaller(a2rwc comet.Auth2ReadWriteCloser) Installer {
	return &VeryInstaller{
		jman: comet.JungleMan{
			A2RWC: a2rwc,
		},
	}
}

func (this *VeryInstaller) Install(restC *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path("/op").
		Consumes(restful.MIME_OCTET).
		Produces(restful.MIME_OCTET)

	ws.Route(
		ws.POST("/xch/{dumpFile}").To(this.exchange).
			Doc("exchange a file with server").
			Operation("exchange").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	ws.Route(
		ws.GET("/dl/{dumpFile}").To(this.download).
			Doc("download a file from server").
			Operation("download").
			Param(ws.PathParameter("dumpFile", "file name").DataType("string")))

	restC.Add(ws)
}

func (this *VeryInstaller) exchange(_req *restful.Request, _resp *restful.Response) {
	token := _req.PathParameter("dumpFile")
	_ = token
	this.jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

func (this *VeryInstaller) download(_req *restful.Request, _resp *restful.Response) {
	token := _req.PathParameter("dumpFile")
	_ = token
	this.jman.ServeHTTP(_resp.ResponseWriter, _req.Request)
}

type FakerInstaller struct {
}

func NewFakerInstaller() Installer {
	return &FakerInstaller{}
}

func (this *FakerInstaller) Install(restC *restful.Container) {

	// TODO
}

type SwaggerInstaller struct {
}

func NewSwaggerInstaller() Installer {
	return &SwaggerInstaller{}
}

func (this *SwaggerInstaller) Install(restC *restful.Container) {
	swagger.RegisterSwaggerService(swagger.Config{
		WebServices: restC.RegisteredWebServices(),
		ApiPath:     "/apidocs.json",
	}, restC)
}
