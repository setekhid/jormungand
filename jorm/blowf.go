// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"crypto/cipher"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/setekhid/jormungand/jorm/stor"
	"golang.org/x/crypto/blowfish"
	"math"
	"math/rand"
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

//==============================================================>

type Blowf blowfish.Cipher

func NewBlowf(key []byte) *Blowf {

	// guaratee blowfish will not complain
	if k := len(key); k < 1 || k > 56 {
		key = append(key, 0x0)
		key = key[:56]
	}

	fish, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err) // there must be succeed
	}

	return (*Blowf)(fish)
}

func (bf *Blowf) IV(pktTail []byte) {

	iv := make([]byte, blowfish.BlockSize)
	binary.BigEndian.PutUint32(iv[:blowfish.BlockSize/2], rand.Uint32())
	binary.BigEndian.PutUint32(iv[blowfish.BlockSize/2:], rand.Uint32())
	copy(pktTail, iv)
}

func (bf *Blowf) IVSiz(pktLen int) int {

	return blowfish.BlockSize - (pktLen % blowfish.BlockSize)
}

func (bf *Blowf) Enc(pkt []byte) []byte {

	c := (*blowfish.Cipher)(bf)

	pktLen := len(pkt)
	rstLen := pktLen - blowfish.BlockSize
	code := make([]byte, pktLen)

	// encrypting rest packet
	cipher.NewCBCEncrypter(c, pkt[rstLen:]).CryptBlocks(code[:rstLen], pkt[:rstLen])
	// encrypting iv field
	c.Encrypt(code[rstLen:], pkt[rstLen:])

	return code
}

func (bf *Blowf) Dec(pkt []byte) []byte {

	c := (*blowfish.Cipher)(bf)

	pktLen := len(pkt)
	rstLen := pktLen - blowfish.BlockSize
	plain := make([]byte, pktLen)

	// decrypting iv field
	c.Decrypt(plain[rstLen:], pkt[rstLen:])
	// decrypting rest packet
	cipher.NewCBCDecrypter(c, plain[rstLen:]).CryptBlocks(plain[:rstLen], pkt[:rstLen])

	return plain
}

type BlowPool struct {
	fishes  map[uint32]*Blowf
	storage stor.BlowStor
}

func NewBlowPool(storage stor.BlowStor) *BlowPool {

	return &BlowPool{
		fishes:  map[uint32]*Blowf{},
		storage: storage,
	}
}

func (pool *BlowPool) Fish(ipId uint32) *Blowf {

	if fish, ok := pool.fishes[ipId]; ok {
		return fish
	}

	fish := NewBlowf(pool.storage.ReadBfKey(ipId))
	pool.fishes[ipId] = fish
	return fish
}

var (
	fishPool = (*BlowPool)(nil)
)

func FishPool() *BlowPool {

	if fishPool == nil {
		fishPool = NewBlowPool(stor.DB())
	}
	return fishPool
}
