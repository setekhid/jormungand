// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package stor

import (
	"encoding/binary"
)

type DummyStor struct {
}

type DummyConf bool

func NewDummyStor(conf *DummyConf) (*DummyStor, error) { return &DummyStor{}, nil }

var (
	dummyBfKeys = map[uint32][]byte{
		binary.BigEndian.Uint32([]byte{192, 168, 89, 123}): []byte{192, 168, 89, 123},
		binary.BigEndian.Uint32([]byte{192, 168, 89, 133}): []byte{192, 168, 89, 133},
		binary.BigEndian.Uint32([]byte{192, 168, 89, 143}): []byte{192, 168, 89, 143},
	}
)

func (db *DummyStor) ReadBfKey(ipId uint32) (BfKeyInfo, bool) {

	if key, ok := dummyBfKeys[ipId]; ok {
		return BfKeyInfo{Key: key, TTL: defaultTTL}, true
	}
	return BfKeyInfo{}, false
}
