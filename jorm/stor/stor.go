// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package stor

import (
	"github.com/setekhid/jormungand/misc/jargs"
)

var (
	moduleName = "stor"
	db         = Stor(nil)
)

type BlowStor interface {
	ReadBfKey(ipId uint32) []byte // TODO with ttl
}

type Stor interface {
	BlowStor
}

func DB() Stor {

	if db == nil {

		conf := jargs.Module(moduleName).(*RedisConf)
		db = NewRedisStor(conf)
	}
	return db
}
