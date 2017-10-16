// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jargs

import (
	"flag"
)

var (
	json_file string
)

func init() {
	flag.StringVar(&json_file, "f", "./args.json", "config file in json format")
}
