// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package web

import (
	"github.com/emicklei/go-restful"
)

type ManagerInstaller struct {
}

// installing manage apis
func (ins *ManagerInstaller) Install(restC *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path("/mg").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(
		ws.POST("/route/{stat}").To(ins.reportRoute).
			Doc("report a route status").
			Operation("reportRoute").
			Param(ws.PathParameter("stat", "route status").DataType("int")))

	restC.Add(ws)
}

// report a route status
func (ins *ManagerInstaller) reportRoute(req *restful.Request, resp *restful.Response) {

	// TODO
}
