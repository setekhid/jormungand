// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package stor

import (
	"errors"
	"github.com/setekhid/jormungand/misc/atexit"
	"github.com/setekhid/jormungand/misc/jargs"
	"io"
	"time"
)

var (
	moduleName = "stor"
	db         = Stor(nil)
	defaultTTL = int64(365 * 24 * time.Hour / time.Second)
)

type BlowStor interface {
	ReadBfKey(ipId uint32) (BfKeyInfo, bool)
}

type BfKeyInfo struct {
	Key []byte
	TTL int64
}

type Stor interface {
	BlowStor
}

type StorConf struct {
	Type  string     `json:"type"`
	Dummy *DummyConf `json:"dummy, omitempty"`
	Redis *RedisConf `json:"redis, omitempty"`
}

func (conf *StorConf) RegistJargs() { jargs.Regist(moduleName, conf) }

func (conf *StorConf) CreateStor() (Stor, error) {

	var st Stor
	var err error

	switch conf.Type {
	case "dummy":
		st, err = NewDummyStor(conf.Dummy)
	case "redis":
		st, err = NewRedisStor(conf.Redis)
	default:
		st, err = nil, errors.New("uncognized stor type")
	}

	if err != nil {
		return st, err
	}

	if cl, ok := st.(io.Closer); ok {
		atexit.Reg(cl)
	}

	return st, err
}

func DB() Stor {

	if db == nil {

		conf := jargs.Module(moduleName).(*StorConf)
		var err error
		db, err = conf.CreateStor()
		if err != nil {
			panic(err) // nothing I can do
		}
	}
	return db
}
