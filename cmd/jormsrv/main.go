// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package main

import (
	goflag "flag"
	"github.com/golang/glog"
	jorm "github.com/setekhid/jormungand"
	"github.com/setekhid/jormungand/tungo"
	"github.com/setekhid/jormungand/web"
	flag "github.com/spf13/pflag"
)

var (
	webAddr = flag.String("web-addr", "localhost:8888", "web address to serve on")
	ifName  = flag.String("if-name", "tun0", "tuntap interface name")
	devFile = flag.String("dev-file", "/dev/net/tun", "tuntap device file")
	isTap   = flag.Bool("is-tap", false, "if it's tap interface")
)

func main() {

	// fix glog and other library, while using native flag library
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	var err error

	// constructing jorm.Router
	router := jorm.NewRouter(jorm.RouterConf{
		IfInfo: tungo.IfInfo{
			IfName:  *ifName,
			DevFile: *devFile,
			IsTap:   *isTap,
		},
	})
	err = router.Start()
	if err != nil {
		glog.Fatalln("Failed starting router: ", err)
	}
	defer router.Stop()

	// constructing web.Server
	server := web.NewServer(web.ServerConf{
		Router:  router,
		Address: *webAddr,
	})
	err = server.Start()
	if err != nil {
		glog.Fatalln("Failed starting server: ", err)
	}
	defer server.Stop()
}
