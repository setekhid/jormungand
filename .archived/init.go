// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package main

import (
	_ "github.com/setekhid/jormungand/jorm/payload"
	_ "github.com/setekhid/jormungand/jorm/sel"
	_ "github.com/setekhid/jormungand/jorm/sox"
	_ "github.com/setekhid/jormungand/jorm/stor"
	_ "github.com/setekhid/jormungand/jorm/tun"
	_ "github.com/setekhid/jormungand/jorm/web"
)

func init() {

	registJargs()
}
