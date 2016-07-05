//usr/bin/env go run $0 -- $@; exit $?

// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	MKDIR_CMD = "mkdir -p"
	CP_CMD    = "cp"
)

func main() {

	flag.Parse()
	if flag.NArg() != 1 {
		panic("usage: command dest_dir")
	}

	// get destination working directory
	pwd := flag.Arg(0)

	// parse package locations
	locals := &PkgLocals{}
	err := json.NewDecoder(os.Stdin).Decode(locals)
	if err != nil {
		panic(err)
	}

	// generate shell script
	sh := &ShellCmds{}
	for _, local := range locals.Pkgs {

		cmds, err := local.DoCopy(pwd)
		if err != nil {
			panic(err)
		}
		sh.AppendCmds(cmds)
	}
	echo(sh.PrintShell())
}

func echo(str string) { fmt.Println(str) }

type ShellCmds struct {
	Cmds []string
}

func (sh *ShellCmds) AppendCmd(cmd string)      { sh.Cmds = append(sh.Cmds, cmd) }
func (sh *ShellCmds) AppendCmds(oth *ShellCmds) { sh.Cmds = append(sh.Cmds, oth.Cmds...) }

func (sh *ShellCmds) PrintShell() string {

	script := ""
	for _, cmd := range sh.Cmds {
		script += fmt.Sprintln(cmd)
	}
	return script
}

type PkgLocals struct {
	Pkgs []PkgLocal `json:"pkgs"`
}

type PkgLocal struct {
	Name  string `json:"name"`
	Local string `json:"local"`
}

var LICENSE_FILES = []string{
	"LICENSE",
	"COPYING",
}

func (l *PkgLocal) DoCopy(dst string) (*ShellCmds, error) {

	pkg := l.Name
	src_path := l.Local
	dst_path := dst + "/" + pkg

	sh := &ShellCmds{}
	sh.AppendCmd(MKDIR_CMD + " " + dst_path) // mkdir -p

	files, err := ioutil.ReadDir(src_path)
	if err != nil {
		return nil, err
	}

	for _, file := range files { // copy all source files

		if file.IsDir() { // not copy package recursively
			continue
		}

		src_file := src_path + "/" + file.Name()
		dst_file := dst_path + "/" + file.Name()
		sh.AppendCmd(CP_CMD + " " + src_file + " " + dst_file)
	}

	// process license
	for dst_path != dst && path.Base(src_path) == path.Base(dst_path) {

		// check license exist
		for _, license_file := range LICENSE_FILES {

			license_src := src_path + "/" + license_file
			license_dst := dst_path + "/" + license_file
			if pathExists(license_src) {
				sh.AppendCmd(CP_CMD + " " + license_src + " " + license_dst)
				break
			}
		}

		// iterate
		src_path = path.Dir(src_path)
		dst_path = path.Dir(dst_path)
	}

	return sh, nil
}

func pathExists(path string) bool {

	_, err := os.Stat(path)
	return err == nil
}
