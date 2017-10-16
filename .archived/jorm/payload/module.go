// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package payload

import (
	"github.com/setekhid/jormungand/misc/jargs"
)

const (
	moduleName = "payload"
)

var (
	conf  = &Config{}
	funcs = Payload(nil)
)

type Config struct {
	Type string `json:"type"`
	Mtu  uint16 `json:"mtu"`
}

func (c *Config) fillEmpty() *Config {

	if len(c.Type) <= 0 {
		c.Type = "tunip"
	}
	if c.Mtu <= 0 {
		c.Mtu = 1500
	}
	return c
}

func Funcs() Payload {

	if funcs == nil {
		funcs = NewPayload(Conf().Type)
	}
	return funcs
}

func Conf() *Config { return jargs.Module(moduleName).(*Config).fillEmpty() }
