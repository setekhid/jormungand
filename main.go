// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/setekhid/jormungand/misc/jargs"
	"strings"
)

func main() {

	flag.Parse()

	modules := strings.Split(mainModule().Modules, ",")
	for i := 0; i < len(modules); i++ {
		modules[i] = strings.TrimSpace(modules[i])
	}

	jargs.RunInMain(modules)
}

const (
	moduleName = "main"
)

type MainModule struct {
	Modules string `json:"modules"`
}

func registJargs()            { jargs.Regist(moduleName, &MainModule{}) }
func mainModule() *MainModule { return jargs.Module(moduleName).(*MainModule) }
