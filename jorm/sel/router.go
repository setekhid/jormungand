// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package sel

import (
	"github.com/setekhid/jormungand/jorm"
)

const (
	moduleName = "sel"
)

var (
	router = (*jorm.Router)(nil)
	rterm  = make(chan struct{}, 8)
)

func Router() *jorm.Router {

	if router == nil {

		router = jorm.NewRouter()
	}
	return router
}

func StopSel() { close(rterm) }
