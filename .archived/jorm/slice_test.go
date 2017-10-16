// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"testing"
)

func TestReslice(t *testing.T) {

	a := make([]byte, 10, 200)
	t.Log(len(a))
	a = a[:20]
	t.Log(len(a))
	a = a[:200]
	t.Log(len(a))
	//a = a[:201]
	//t.Log(len(a))
}
