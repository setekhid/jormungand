// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/setekhid/jormungand/http/comet"
	"github.com/setekhid/jormungand/misc/edgo"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"io"
	"net/url"
	"sync"
)

type Router struct {
	txch chan []byte

	evChan chan edgo.Event
	evProc *edgo.EdGo

	ulns map[uint32]LinkRef
	rtab *RTable
}

func NewRouter() *Router {

	router := &Router{

		txch: make(chan []byte, 10240),

		ulns: map[uint32]LinkRef{},
		rtab: NewRTable(),
	}

	evChan := make(chan edgo.Event, 32)
	evProc := edgo.NewEdGo(router, evChan)
	evProc.Regist(authEventType, authEventFunc)
	evProc.Regist(closEventType, closEventFunc)

	router.evChan = evChan
	router.evProc = evProc

	return router
}

const (
	authEventType = 1
	closEventType = 2
)

type authEvent struct {
	Uri   string
	InLen int64

	Done chan struct{}

	Tun    io.ReadWriteCloser
	OutLen int64
	Err    error
}

var authEventPool = &sync.Pool{
	New: func() interface{} { return &authEvent{Done: make(chan struct{})} },
}

func fishAuthEvent() *authEvent { return authEventPool.Get().(*authEvent) }
func (e *authEvent) free()      { authEventPool.Put(e) }

// Override edgo.Event.Type()
func (e *authEvent) Type() int { return authEventType }

func authEventFunc(self interface{}, event edgo.Event) {

	router := self.(*Router)
	authE := event.(*authEvent)
	authE.Tun, authE.OutLen, authE.Err = router.auth_unsafe(authE.Uri, authE.InLen)
	authE.Done <- struct{}{}
}

type closEvent struct {
	IpId uint32
	Nets []IPv4Net
}

// Override edgo.Event.Type()
func (e *closEvent) Type() int { return closEventType }

func closEventFunc(self interface{}, event edgo.Event) {

	router := self.(*Router)
	closE := event.(*closEvent)
	router.kick_unsafe(closE.IpId, closE.Nets)
}

// Override comet.TunnelAuthor.Auth
func (r *Router) Auth(uri string, inlen int64) (tun io.ReadWriteCloser, outlen int64, err error) {

	event := fishAuthEvent()
	defer event.free()

	event.Uri = uri
	event.InLen = inlen
	r.evChan <- event
	<-event.Done

	return event.Tun, event.OutLen, event.Err
}

func (r *Router) Kick(ipId uint32, nets []IPv4Net) {

	r.evChan <- &closEvent{
		IpId: ipId,
		Nets: nets,
	}
}

func (r *Router) route_unsafe(msg []byte) bool {

	// TODO
	return false
}

func (r *Router) auth_unsafe(uri string, inlen int64) (tun io.ReadWriteCloser, outlen int64, err error) {

	// account's token json
	tokStr, err := r.cutTokenStr(uri)
	if err != nil {
		return nil, 0, err
	}
	tok, err := r.decodeToken(tokStr)
	if err != nil {
		return nil, 0, err
	}

	// user link
	uln, ok := r.ulns[tok.AccountIP]
	if !ok { // user first login

		// routeable v4 network
		ipnets := r.parseIPNet(tok)
		// reference a new link
		uln := r.registLink(NewLink(r.txch), tok.AccountIP, ipnets)
		r.ulns[tok.AccountIP] = uln
		lnPusher := uln.Link.Pusher()
		// regist self route
		r.rtab.LinkIP(tok.AccountIP, lnPusher)
		// regist ipv4 network routes
		for _, ipnet := range ipnets {
			r.rtab.LinkTun(ipnet, tok.AccountIP, lnPusher)
		}
	} else {

		// exists, just clone
		uln = uln.Clone()
	}

	return uln, comet.HTTP_DL_NORMAL_LEN, nil
}

func (r *Router) kick_unsafe(ipId uint32, nets []IPv4Net) {

	uln := r.ulns[ipId]

	if uln.NoRef() {

		for _, ipnet := range nets {
			r.rtab.DiscardTun(ipnet, ipId)
		}
		r.rtab.DiscardIP(ipId)
		delete(r.ulns, ipId)
		uln.Link.Close()
	}
}

func (r *Router) cutTokenStr(uri string) (string, error) {

	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	tok := url.Path[1:]
	return tok, nil
}

type TokJson struct {
	AccountIP uint32 `json:"-" bson:"-"`
	Checksum  uint32 `json:"-" bson:"-"`

	Route1 *uint32 `json:"r1,omitempty" bson:"r1,omitempty"`
	Route2 *uint32 `json:"r2,omitempty" bson:"r2,omitempty"`
	Route3 *uint32 `json:"r3,omitempty" bson:"r3,omitempty"`
	Route4 *uint32 `json:"r4,omitempty" bson:"r4,omitempty"`
	Masks  uint32  `json:"ms,omitempty" bson:"ms,omitempty"`
	Simple string  `json:"sm" bson:"sm"`
}

func (r *Router) decodeToken(tokStr string) (*TokJson, error) {

	// base64 decode tokj
	tokj, err := base64.StdEncoding.DecodeString(tokStr)
	if err != nil {
		return nil, err
	}

	// first 4 bytes as account ip
	accountIP := binary.BigEndian.Uint32(tokj[:4])
	tokj = tokj[4:]

	// TODO decrypt tokj

	// second 4 bytes as checksum
	checksum := binary.BigEndian.Uint32(tokj[:4])
	tokj = tokj[4:]

	// crc32 checksum
	if _checksum := crc32.ChecksumIEEE(tokj); _checksum != checksum {
		return nil, errors.New("crc32 checksum doesn't match")
	}

	// account identify and checksum
	tok := TokJson{
		AccountIP: accountIP,
		Checksum:  checksum,
	}

	// complete tok struct
	err = bson.Unmarshal([]byte(tokj), &tok)
	if err != nil {
		return nil, err
	}

	return &tok, nil
}

func (r *Router) parseIPNet(tok *TokJson) []IPv4Net {

	// parse account's ipnet
	ipnets := make([]IPv4Net, 0, 4)
	if tok.Route1 != nil {
		r1Mask := (uint32(0x1) << ((tok.Masks >> (3 * 8)) & 0xff)) - 1
		ipnets = append(ipnets, IPv4Net{
			IP:   IPv4(*tok.Route1),
			Mask: IPv4(r1Mask),
		})
	}
	if tok.Route2 != nil {
		r2Mask := (uint32(0x1) << ((tok.Masks >> (2 * 8)) & 0xff)) - 1
		ipnets = append(ipnets, IPv4Net{
			IP:   IPv4(*tok.Route2),
			Mask: IPv4(r2Mask),
		})
	}
	if tok.Route3 != nil {
		r3Mask := (uint32(0x1) << ((tok.Masks >> (1 * 8)) & 0xff)) - 1
		ipnets = append(ipnets, IPv4Net{
			IP:   IPv4(*tok.Route3),
			Mask: IPv4(r3Mask),
		})
	}
	if tok.Route4 != nil {
		r4Mask := (uint32(0x1) << ((tok.Masks >> (0 * 8)) & 0xff)) - 1
		ipnets = append(ipnets, IPv4Net{
			IP:   IPv4(*tok.Route4),
			Mask: IPv4(r4Mask),
		})
	}

	return ipnets
}
