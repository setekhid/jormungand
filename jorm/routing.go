// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"bytes"
	"golang.org/x/net/ipv4"
	"io"
	"sync"
)

func Routing(r *Router, term <-chan struct{}) {

	for true {

		select {
		case msg := <-r.txch: // routing packet or message
			r.route_unsafe(msg)
		case ev := <-r.evChan: // processing router event
			r.evProc.Process(ev)

		case <-term:
			break
		}
	}
}

type routingHelper struct{}

var rHelper = routingHelper{}

func (_ routingHelper) send2RcvChan(pusher chan<- []byte, msg []byte) bool {

	var ok bool
	select {
	case pusher <- msg:
		ok = true
	default:
		ok = false
	}
	return ok
}

// ==========================================================>

// implement io.ReadWriteCloser
type Tunnel struct {

	// cache the post data of http comet
	Pollable      int64
	Writer        bytes.Buffer
	WriterGuarder sync.Mutex

	// cache the response body of http
	Pushable      int64
	Reader        bytes.Buffer
	ReaderGuarder sync.Mutex

	// whether holding by a http go routine or not
	IsHolding bool
}

func NewTunnel(pollable, pushable int64) *Tunnel {

	this := &Tunnel{}
	this.initialize(pollable, pushable)
	return this
}

func (this *Tunnel) initialize(pollable, pushable int64) {

	this.Pollable = pollable
	this.Pushable = pushable
	this.IsHolding = false
}

func (this *Tunnel) LeftPollable() int64 { return this.Pollable }

// poll some bytes from buffer, return the bytes count polled
func (this *Tunnel) PollBytes(p []byte) int {

	if this.Pollable <= 0 { // no more bytes exists
		return 0
	}

	// round the len(p) to limit Pollable
	if int64(len(p)) > this.Pollable {
		p = p[:this.Pollable]
	}

	this.WriterGuarder.Lock()
	defer this.WriterGuarder.Unlock()
	data := this.Writer.Next(len(p))
	count := copy(p, data)
	this.Pollable -= int64(count)
	return count
}

func (this *Tunnel) Write(p []byte) (int, error) {

	this.WriterGuarder.Lock()
	defer this.WriterGuarder.Unlock()
	return this.Writer.Write(p)
}

func (this *Tunnel) LeftPushable() int64 { return this.Pushable }

// push the whole packet into buffer, if the buffer can't hold it,
// return io.ErrorShortWrite, otherwise nil
func (this *Tunnel) PushPacket(p []byte) error {

	if int64(len(p)) > this.Pushable { // packet can't be pushed entirely
		return io.ErrShortWrite
	}

	this.Pushable -= int64(len(p))

	this.ReaderGuarder.Lock()
	defer this.ReaderGuarder.Unlock()
	this.Reader.Write(p)
	return nil
}

func (this *Tunnel) Read(p []byte) (int, error) {

	this.ReaderGuarder.Lock()
	defer this.ReaderGuarder.Unlock()
	data := this.Reader.Next(len(p))
	count := copy(p, data)
	return count, nil
}

func (this *Tunnel) Open() {
	this.IsHolding = true
}

func (this *Tunnel) Close() error {
	this.IsHolding = false
	return nil
}

/*
 *  every tunnelled ipv4 packet should be the times of the block size 8, if not,
 *  attach some random data to make it be.
 */

type IPv4Tunnel struct {
	Tunnel

	//Blower *BlowGuy
	BlowId uint32

	FragPacket []byte
	IPv4Header *ipv4.Header
	FragLength int
}

func NewIPv4Tunnel(pollable, pushable int64, mtu int, blower *int, blowId uint32) *IPv4Tunnel {

	this := &IPv4Tunnel{}
	this.initialize(pollable, pushable, mtu, blower, blowId)
	return this
}

func (this *IPv4Tunnel) initialize(pollable, pushable int64, mtu int, blower *int, blowId uint32) {

	this.Tunnel.initialize(pollable, pushable)

	//this.Blower = blower
	this.BlowId = blowId
	this.FragPacket = make([]byte, mtu)
	this.IPv4Header = nil
	this.FragLength = 0
}

func (this *IPv4Tunnel) cacheFragment() {

	// not long enough, try to poll
	copied := this.PollBytes(this.FragPacket[this.FragLength:])
	// decode
	//decodeBegin := this.Blower.BlockFloor(this.FragLength)
	this.FragLength += copied
	//decodeEnd := this.Blower.BlockFloor(this.FragLength)
	//this.Blower.Decrypt(this.BlowId, this.FragPacket[decodeBegin:decodeEnd]) // TODO combine xor the previous block
}

func (this *IPv4Tunnel) PollPacket(p []byte) (int, *ipv4.Header) {

	if this.IPv4Header == nil { // parse ipv4 header

		// completing ipv4 header
		//if this.FragLength < this.Blower.BlockCeil(ipv4.HeaderLen) {
		//	this.cacheFragment()
		//}
		//if this.FragLength < this.Blower.BlockCeil(ipv4.HeaderLen) {
		// this time, we have nothing to do
		//	return 0, nil
		//}

		// parse ipv4 header
		header, err := ipv4.ParseHeader(this.FragPacket)
		if err != nil {
			// only two error will cause, both mean fragment is not long enough
			return 0, nil
		}
		this.IPv4Header = header
	}

	// check packet length is in mtu
	if len(this.FragPacket) < this.IPv4Header.TotalLen {
		// TODO report an error
		return 0, nil
	}

	// we got a header, now cache packet to its length
	//if this.FragLength < this.Blower.BlockCeil(this.IPv4Header.TotalLen) {
	//	this.cacheFragment()
	//}
	//if this.FragLength < this.Blower.BlockCeil(this.IPv4Header.TotalLen) {
	// this time, we have nothing to do
	//	return 0, nil
	//}

	// now we got a whole packet
	copied := copy(p, this.FragPacket[:this.IPv4Header.TotalLen])
	header := this.IPv4Header
	// clean cache
	//this.FragLength = copy(this.FragPacket, this.FragPacket[this.Blower.BlockCeil(this.IPv4Header.TotalLen):0])
	this.IPv4Header = nil

	return copied, header
}
