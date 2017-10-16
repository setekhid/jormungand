// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"crypto/cipher"
	"encoding/binary"
	"github.com/setekhid/jormungand/jorm/stor"
	"golang.org/x/crypto/blowfish"
	"math/rand"
	"time"
)

func BlowfCap(len int) int {

	if len%blowfish.BlockSize == 0 {
		return len
	}
	// return (len / blowfish.BlockSize + 1) * blowfish.BlockSize
	// blowfish.BlockSize == 8
	return (len | (blowfish.BlockSize - 1)) + 1
}

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

func (bf *Blowf) FillCap(pkt []byte) []byte {

	riv := make([]byte, BlowfCap(len(pkt))-len(pkt))
	bf.IV(riv)
	return append(pkt, riv...)
}

func (bf *Blowf) IVSiz(pktLen int) int {

	return blowfish.BlockSize - (pktLen % blowfish.BlockSize)
}

func (bf *Blowf) SplitIV(pkt []byte) (iv []byte, rst []byte) {

	iv = make([]byte, blowfish.BlockSize)
	bf.IV(iv)

	pktLen := len(pkt)
	rstLen := pktLen - (pktLen % blowfish.BlockSize)
	copy(iv, pkt[rstLen:])
	rst = pkt[:rstLen]
	return iv, rst
}

func (bf *Blowf) CombiIV(iv []byte, rst []byte) []byte { return append(rst, iv...) }

func (bf *Blowf) Enc(iv []byte, rst []byte) []byte {

	c := (*blowfish.Cipher)(bf)

	code := make([]byte, blowfish.BlockSize+len(rst))

	// encrypting iv field
	c.Encrypt(code[:blowfish.BlockSize], iv)
	// encrypting rest packet
	cipher.NewCBCEncrypter(c, code[:blowfish.BlockSize]).CryptBlocks(code[blowfish.BlockSize:], rst)

	return code
}

func (bf *Blowf) EncS(ci []byte, iv []byte, pkt []byte) []byte {

	c := (*blowfish.Cipher)(bf)

	cipher.NewCBCEncrypter(c, iv).CryptBlocks(ci, pkt)

	niv := make([]byte, blowfish.BlockSize)
	copy(niv, ci[len(ci)-blowfish.BlockSize:])

	return niv
}

func (bf *Blowf) Dec(pkt []byte) (iv []byte, rst []byte) {

	c := (*blowfish.Cipher)(bf)

	iv = make([]byte, blowfish.BlockSize)
	rst = make([]byte, len(pkt)-blowfish.BlockSize)

	// decrypting rest packet
	cipher.NewCBCDecrypter(c, pkt[:blowfish.BlockSize]).CryptBlocks(rst, pkt[blowfish.BlockSize:])
	// decrypting iv field
	c.Decrypt(iv, pkt[:blowfish.BlockSize])

	return iv, rst
}

func (bf *Blowf) DecS(pkt []byte, iv []byte, ci []byte) []byte {

	c := (*blowfish.Cipher)(bf)

	niv := make([]byte, blowfish.BlockSize)
	copy(niv, ci[len(ci)-blowfish.BlockSize:])

	cipher.NewCBCDecrypter(c, iv).CryptBlocks(pkt, ci)

	return niv
}

// ============================= deprecated above

type FishPool struct {
	fishes  map[uint32]stor.BfKeyInfo
	storage stor.BlowStor
}

func NewFishPool(storage stor.BlowStor) *FishPool {

	return &FishPool{
		fishes:  map[uint32]stor.BfKeyInfo{},
		storage: storage,
	}
}

func (pool *FishPool) Fish(ipId uint32) ([]byte, bool) {

	now := time.Now().Unix()

	if fish, ok := pool.fishes[ipId]; ok {

		if fish.TTL >= now {
			return fish.Key, true
		} else {
			delete(pool.fishes, ipId)
		}
	}

	if fish, ok := pool.storage.ReadBfKey(ipId); ok {

		fish.TTL += now
		pool.fishes[ipId] = fish
		return fish.Key, true
	}

	return nil, false
}

var (
	fishPool = (*FishPool)(nil)
)

func FishPoolInst() *FishPool {

	if fishPool == nil {
		fishPool = NewFishPool(stor.DB())
	}
	return fishPool
}
