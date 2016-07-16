// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"github.com/setekhid/jormungand/jorm/payload"
	"github.com/setekhid/jormungand/misc/jargs"
)

const (
	moduleName = "jorm"
)

var (
	conf = &Config{}
)

type Config struct {
	Payload     string          `json:"payload"`
	payloadInst payload.Payload `json:"-"`
	Mtu         uint16          `json:"mtu"`
}

func (c *Config) payload() payload.Payload {

	if c.payloadInst == nil {
		c.payloadInst = payload.NewPayload(c.Payload)
	}
	return c.payloadInst
}

func Conf() *Config { return jargs.Module(moduleName).(*Config) }
