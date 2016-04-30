// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jormungand

import (
	"github.com/setekhid/jormungand/tungo"
	"io"
)

type SysTun struct {
	tungo.IfInfo
	io.ReadWriteCloser
}

func (this *SysTun) Open() error {
	var err error
	this.ReadWriteCloser, err = tungo.NewTunTap(this.IfInfo)
	return err
}
