// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"math/rand"
	"testing"
	"time"
)

func TestBlowf(t *testing.T) {

	rand.Seed(time.Now().Unix())

	fish := NewBlowf([]byte("yoyo check it now"))

	msg := []byte("hello world!")
	msglen := len(msg)
	t.Log("original message:")
	t.Log(msg)

	cip := fish.Enc(fish.SplitIV(msg))
	t.Log("cipher message:")
	t.Log(cip)

	cip = fish.CombiIV(fish.Dec(cip))
	t.Log("plain message:")
	t.Log(cip)

	if string(cip[:msglen]) != string(msg) {
		t.FailNow()
	}
	t.Log(string(msg))
}
