// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package refc

import (
	"io"
	"sync/atomic"
)

type CountCloser struct {
	refc *int32
	Obj  io.Closer
}

func NewCountCloser(obj io.Closer) CountCloser {

	refc := int32(1)
	return CountCloser{
		refc: &refc,
		Obj:  obj,
	}
}

func (cc CountCloser) Clone() CountCloser {

	atomic.AddInt32(cc.refc, 1)
	return cc
}

func (cc CountCloser) RefCount() int { return int(atomic.LoadInt32(cc.refc)) }
func (cc CountCloser) NoRef() bool   { return cc.RefCount() <= 0 }

func (cc CountCloser) Close() error {

	if atomic.AddInt32(cc.refc, -1) <= 0 {
		return cc.Obj.Close()
	}
	return nil
}
