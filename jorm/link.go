// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"sync/atomic"
)

const (
	PKT_CACHE_COUNT = 128
)

type Link struct {
	Rcv chan []byte
	Tx  chan<- []byte
}

func NewLink(tx chan<- []byte) *Link {

	return &Link{
		Rcv: make(chan []byte, PKT_CACHE_COUNT),
		Tx:  tx,
	}
}

func (l *Link) Pusher() chan<- []byte { return l.Rcv }

// Override io.Reader.Read
func (l *Link) Read(p []byte) (n int, err error) {
	p = <-l.Rcv
	return len(p), nil
}

// Override io.Writer.Write
func (l *Link) Write(p []byte) (n int, err error) {
	l.Tx <- p
	return len(p), nil
}

// Override io.Closer.Close
func (l *Link) Close() error {

	l.Rcv = nil
	l.Tx = nil
	return nil
}

// struct Link count refrence
type LinkRef struct {
	refC *int32

	r    *Router
	ipId uint32
	nets []IPv4Net

	*Link
}

func (r *Router) registLink(l *Link, ipId uint32, nets []IPv4Net) LinkRef {

	refC := int32(1)
	return LinkRef{

		refC: &refC,

		r:    r,
		ipId: ipId,
		nets: nets,

		Link: l,
	}
}

func (l LinkRef) Clone() LinkRef {

	atomic.AddInt32(l.refC, 1)
	return l
}

func (l LinkRef) NoRef() bool { return atomic.LoadInt32(l.refC) <= 0 }

func (l LinkRef) Close() error {

	if atomic.AddInt32(l.refC, -1) <= 0 {
		l.r.Kick(l.ipId, l.nets) // notify the router style gc
	}
	return nil
}
