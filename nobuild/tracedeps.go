//usr/bin/env go run $0 -- $@; exit $?

// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

// +build nobuild

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"strings"
)

func main() {

	flag.Parse()
	if flag.NArg() < 1 {
		panic("usage: command pkg_path [pkg_path..]")
	}

	// all dependencies packages
	pkgs := map[string]string{}
	for _, flag_Arg := range flag.Args() {

		for _, pkg := range goDeps(flag_Arg) {

			real_path := resolvPath(pkg)
			if len(real_path) > 0 {
				pkgs[pkg] = real_path
			}
		}
	}

	for hasNew := true; hasNew; {

		// find out new dependencies
		newPkgs := map[string]string{}
		for _, pkgPath := range pkgs {

			for _, depPkg := range goDeps(pkgPath) {

				if _, exists := pkgs[depPkg]; exists { // already recorded
					continue
				}

				depPath := resolvPath(depPkg)
				if len(depPath) <= 0 { // the real_path doesn't exists, ignore
					continue
				}

				newPkgs[depPkg] = depPath
			}
		}

		// add newly dependencies to pkgs
		hasNew = len(newPkgs) > 0
		for pkg, path := range newPkgs {
			pkgs[pkg] = path
		}
	}

	// generate json
	locals := &PkgLocals{}
	for pkg, src_path := range pkgs {
		locals.Pkgs = append(locals.Pkgs, PkgLocal{
			Name:  pkg,
			Local: src_path,
		})
	}
	// print out
	outStr, err := json.Marshal(locals)
	if err != nil {
		panic(err)
	}
	echo(string(outStr))
}

type PkgLocals struct {
	Pkgs []PkgLocal `json:"pkgs"`
}

type PkgLocal struct {
	Name  string `json:"name"`
	Local string `json:"local"`
}

func echo(str string) { fmt.Println(str) }

func goDeps(pkgPath string) []string {

	pkgInfo, err := build.ImportDir(pkgPath, 0)
	if err != nil {
		log.Println(err)
		return nil
	}

	return pkgInfo.Imports
}

var gopath_env = os.Getenv("GOPATH")
var gopaths = strings.Split(gopath_env, ":")

func resolvPath(pkg string) string {

	if len(gopath_env) <= 0 {
		panic("GOPATH doesn't exist")
	}

	for _, gopath := range gopaths {

		real_path := gopath + "/src/" + pkg
		if pathExists(real_path) {
			return real_path
		}
	}

	return ""
}

func pathExists(path string) bool {

	_, err := os.Stat(path)
	return err == nil
}
