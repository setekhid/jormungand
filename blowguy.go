// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jormungand

import (
	"github.com/golang/glog"
	_ "golang.org/x/crypto/blowfish" // TODO
	"math"
)

type BlowGuy struct {
	keys map[uint32][]byte
}

func NewBlowGuy() *BlowGuy {

	return &BlowGuy{
		keys: map[uint32][]byte{},
	}
}

func (this *BlowGuy) Encrypt(id uint32, pkt []byte) {

	// TODO
}

func (this *BlowGuy) Decrypt(id uint32, pkt []byte) {

	// TODO
}

func (this *BlowGuy) Checksum(pkt []byte) uint32 {

	// TODO
	return 0
}

func (this *BlowGuy) BlockCeil(len int) int {
	return int(math.Ceil(float64(len)/8.0) * 8.0)
}

func (this *BlowGuy) BlockFloor(len int) int {
	return int(math.Floor(float64(len)/8.0) * 8.0)
}

func (this *BlowGuy) RegKey(id uint32, key []byte) {

	if _, exists := this.keys[id]; exists {
		glog.Warningln("blow guy already knew this id ", id)
	}

	this.keys[id] = key
}
