// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jargs

import (
	"errors"
)

var (
	entries map[string]ModuleEntry
)

type ModuleEntry func()

func RunInMain(modules []string) {
	for _, m := range modules {
		entries[m]()
	}
}

func RegistEntry(module string, entry ModuleEntry) {

	if _, ok := entries[module]; ok {
		panic(errors.New("duplicated module entry " + module))
	}

	entries[module] = entry
}
