// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"testing"
)

func TestBlowf(t *testing.T) {

	fish := NewBlowf([]byte("yoyo check it now"))
	msg := []byte("hello world!")
	msglen := len(msg)
	ivlen := fish.IVSiz(msglen)
	iv := make([]byte, ivlen)
	cip := append(msg, iv...)
	t.Log(cip)
	fish.IV(cip[msglen:])
	cip = fish.Enc(cip)
	t.Log(cip)
	cip = fish.Dec(cip)
	t.Log(cip)
	if string(cip[:msglen]) != string(msg) {
		t.FailNow()
	}
	t.Log(string(msg))
}
