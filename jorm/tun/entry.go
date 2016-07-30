// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package tun

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/setekhid/jormungand/misc/jargs"
	"net"
)

const (
	moduleName = "tun"
)

type Config struct {
	IpId_str   string `json:"ip_id"`
	IpId       uint32 `json:"-"`
	Key_base64 string `json:"key"`
	Key        []byte `json:"-"`
}

func (conf *Config) Parse() {

	var err error

	conf.IpId = binary.BigEndian.Uint32(net.ParseIP(conf.IpId_str))
	conf.Key, err = base64.StdEncoding.DecodeString(conf.Key_base64)
	if err != nil {
		panic(err)
	}
}

type Entry struct {
}

func (en *Entry) Initialize(conf *Config) {
	// TODO
}

func (en *Entry) Operating() {
	// TODO
}

func (conf *Config) registJargs() { jargs.Regist(moduleName, conf) }

func (en *Entry) registEntry() {

	jargs.RegistEntry(moduleName, func() {

		conf := jargs.Module(moduleName).(*Config)
		conf.Parse()

		en.Initialize(conf)
		en.Operating()
	})
}
