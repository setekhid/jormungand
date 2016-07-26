// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"crypto/cipher"
	"github.com/setekhid/jormungand/jorm/payload"
	"github.com/setekhid/jormungand/misc/refc"
	_ "golang.org/x/net/ipv4"
	"io"
)

const (
	PKT_CACHE_COUNT = 128
)

type Network struct {
	Mtu     uint16
	Hub     chan<- []byte
	Payload payload.Payload
}

type Link struct {
	Rcv     chan []byte
	Tx      chan<- []byte
	Mtu     uint16
	Payload payload.Payload
}

func (l *Link) init(net *Network) {

	l.Rcv = make(chan []byte, PKT_CACHE_COUNT)
	l.Tx = net.Hub
	l.Mtu = net.Mtu
	l.Payload = net.Payload
}

func (l *Link) Pusher() chan<- []byte { return l.Rcv }

func (l *Link) ReadPacket() <-chan []byte  { return l.Rcv }
func (l *Link) WritePacket() chan<- []byte { return l.Tx }

type RouterLink struct {
	Link

	R      *Router
	IpId   uint32
	Nets   []IPv4Net
	Cipher cipher.Block
}

func NewRouterLink(r *Router, info *RoutingInfo) *RouterLink {

	rl := &RouterLink{}
	rl.init(r, info)
	return rl
}

func (l *RouterLink) init(r *Router, info *RoutingInfo) {

	l.Link.init(&r.Network)
	l.R = r
	l.IpId = info.IpId
	l.Nets = info.Nets
	l.Cipher = info.Cipher
}

// Override io.Closer.Close
func (l *RouterLink) Close() error {
	l.R.Kick(l.IpId, l.Nets) // notify the router style gc
	return nil
}

type StreamLink struct {
	cc refc.CountCloser

	rl        *RouterLink
	RcvStream cipher.Stream
	TxStream  cipher.Stream
	txCache   []byte
}

func NewStreamLink(cc refc.CountCloser, ivRcv []byte, ivTx []byte) *StreamLink {

	rl := cc.Obj.(*RouterLink)

	return &StreamLink{
		cc:        cc,
		rl:        rl,
		RcvStream: cipher.NewCFBEncrypter(rl.Cipher, ivRcv),
		TxStream:  cipher.NewCFBDecrypter(rl.Cipher, ivTx),
	}
}

// Override io.Reader.Read, block reader
func (l *StreamLink) Read(p []byte) (n int, err error) {

	if uint16(len(p)) < l.rl.Mtu {
		return 0, io.ErrShortBuffer
	}

	plain := <-l.rl.ReadPacket()
	if l.rl.Mtu > uint16(len(plain)) {
		panic(io.ErrShortBuffer) // should not be here
	}

	l.RcvStream.XORKeyStream(p[:len(plain)], plain)
	return len(plain), nil
}

// Override io.Writer.Write, stream writer
func (l *StreamLink) Write(p []byte) (int, error) {

	total := len(p)

	wn := 0
	err := error(nil)
	for p = p[wn:]; len(p) > 0 && err == nil; p = p[wn:] {
		wn, err = l.writePartial(p)
	}

	return total - len(p), err
}

func (l *StreamLink) writePartial(p []byte) (wc int, err error) {

	pl := l.rl.Payload
	mtu := l.rl.Mtu

	if pl.HeaderSize() > uint16(len(l.txCache)+len(p)) {

		// not enough for a header
		l.cache(p)
		wc += len(p)
		return wc, nil
	}

	if pl.HeaderSize() > uint16(len(l.txCache)) {

		// cache the header first
		hdrRstSiz := pl.HeaderSize() - uint16(len(l.txCache))
		l.cache(p[:hdrRstSiz])
		wc += int(hdrRstSiz)
		p = p[hdrRstSiz:]
	}

	// parse the packet size
	pktSiz := pl.PacketSize(l.txCache)
	if pktSiz > mtu {
		return wc, io.ErrShortWrite
	}
	restSiz := pktSiz - uint16(len(l.txCache))

	if restSiz > uint16(len(p)) {

		// not enough for the rest
		l.cache(p)
		wc += len(p)
		return wc, nil
	}

	// cache the whole packet
	l.cache(p[:restSiz])
	wc += int(restSiz)
	// transmit the packet
	l.rl.WritePacket() <- l.txCache
	l.txCache = make([]byte, 0, mtu)
	return wc, nil
}

func (l *StreamLink) cache(p []byte) {

	cacheOff := len(l.txCache)
	l.txCache = l.txCache[:cacheOff+len(p)]
	l.TxStream.XORKeyStream(l.txCache[cacheOff:], p)
}

// Override io.Closer.Close
func (l *StreamLink) Close() error { return l.cc.Close() }
