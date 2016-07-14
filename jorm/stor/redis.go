// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package stor

import (
	"github.com/setekhid/jormungand/misc/jargs"
)

type RedisStor struct {
}

type RedisConf struct {
}

func (conf *RedisConf) RegistJargs() { jargs.Regist(moduleName, conf) }

func NewRedisStor(conf *RedisConf) *RedisStor {

	return &RedisStor{}
}

func (db *RedisStor) ReadBfKey(ipId uint32) []byte {

	// TODO
	return nil
}

func (db *RedisStor) WriteBfKey(ipId uint32, key []byte) {

	// TODO
}
