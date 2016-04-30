// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jormungand

import (
	_ "golang.org/x/crypto/blowfish" // TODO
)

type BlowGuy struct {
}

func NewBlowGuy() *BlowGuy {

	return &BlowGuy{}
}

func (this *BlowGuy) Encrypt(pkg []byte) []byte {

	// TODO
	rtn := make([]byte, len(pkg))
	copy(rtn, pkg)
	return rtn
}

func (this *BlowGuy) Decrypt(pkg []byte) []byte {

	// TODO
	rtn := make([]byte, len(pkg))
	copy(rtn, pkg)
	return rtn
}

func (this *BlowGuy) RegKey(id []byte, key []byte) {

	// TODO
}
