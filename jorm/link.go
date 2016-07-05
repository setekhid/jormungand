// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import ()

const (
	PKT_CACHE_COUNT = 128
)

type Link struct {
	Rcv chan []byte
	Tx  chan []byte
}

func NewLink() *Link {

	return &Link{
		Rcv: make(chan []byte, PKT_CACHE_COUNT),
		Tx:  make(chan []byte, PKT_CACHE_COUNT),
	}
}

func (l *Link) Pusher() chan<- []byte { return l.Rcv }

// Override io.Reader.Read
func (l *Link) Read(p []byte) (n int, err error) {
	p = <-l.Rcv
	return len(p), nil
}

func (l *Link) Puller() <-chan []byte { return l.Tx }

// Override io.Writer.Write
func (l *Link) Write(p []byte) (n int, err error) {
	l.Tx <- p
	return len(p), nil
}

// Override io.Closer.Close
func (l *Link) Close() error {
	close(l.Tx)
	close(l.Rcv)
	return nil
}

// struct Link count refrence
type LinkRef struct {
	*Link
	refC *int
}

func NewLinkRef(l *Link) LinkRef {

	refC := 1
	return LinkRef{
		Link: l,
		refC: &refC,
	}
}

func (l LinkRef) Clone() LinkRef {

	*(l.refC)++
	return l
}

func (l LinkRef) Close() error {

	*(l.refC)--
	if *(l.refC) <= 0 {
		return l.Link.Close()
	}
	return nil
}
