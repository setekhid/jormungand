// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/setekhid/jormungand/http/comet"
	"github.com/setekhid/jormungand/misc/edgo"
	"github.com/setekhid/jormungand/misc/refc"
	"golang.org/x/crypto/blowfish"
	"hash/crc32"
	"io"
	"net/url"
	"sync"
)

const (
	HUB_CACHE_COUNT = 10240
	ACPT_BACK_COUNT = 32
)

type Router struct {
	Network

	txch chan []byte

	evChan chan edgo.Event
	evProc *edgo.EdGo

	ulns  map[uint32]refc.CountCloser
	rtab  *RTable
	rcach RoutingCacher

	pool *FishPool
}

func NewRouter(pool *FishPool) *Router {

	router := &Router{

		txch: make(chan []byte, HUB_CACHE_COUNT),

		ulns:  map[uint32]refc.CountCloser{},
		rtab:  NewRTable(),
		rcach: NewRoutingCacher(),

		pool: pool,
	}

	evChan := make(chan edgo.Event, ACPT_BACK_COUNT)
	evProc := edgo.NewEdGo(router, evChan)
	evProc.Regist(authEventType, authEventFunc)
	evProc.Regist(closEventType, closEventFunc)

	router.evChan = evChan
	router.evProc = evProc

	router.Network = Network{
		Mtu:     Conf().Mtu,
		Hub:     router.txch,
		Payload: Conf().payload(),
	}

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

func getAuthEvent() *authEvent { return authEventPool.Get().(*authEvent) }
func (e *authEvent) free()     { authEventPool.Put(e) }

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

	event := getAuthEvent()
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

	hdr := r.Payload.HeaderInfo(msg)
	if tun, ok := (RoutingConns{r.rtab, r.rcach}).Routing(hdr); ok {
		tun.(chan<- []byte) <- msg
		return true
	}
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
	uln, ok := r.ulns[tok.IpId]
	if !ok { // user first login

		// reference a new link
		ln := NewRouterLink(r, tok)
		uln = refc.NewCountCloser(ln)
		r.ulns[tok.IpId] = uln
		// regist self route
		r.rtab.LinkIP(tok.IpId, ln.Pusher())
		// regist ipv4 network routes
		for _, ipnet := range tok.Nets {
			r.rtab.LinkTun(ipnet, tok.IpId, ln.Pusher())
		}
	} else {

		// exists, just clone
		uln = uln.Clone()
	}

	return NewStreamLink(uln, tok.IvRcv, tok.IvTx), comet.HTTP_DL_NORMAL_LEN, nil
}

func (r *Router) kick_unsafe(ipId uint32, nets []IPv4Net) {

	uln := r.ulns[ipId]

	if uln.NoRef() {

		for _, ipnet := range nets {
			r.rtab.DiscardTun(ipnet, ipId)
		}
		r.rtab.DiscardIP(ipId)
		delete(r.ulns, ipId)
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

type RoutingInfo struct {
	IpId   uint32
	Cipher cipher.Block

	IvRcv, IvTx []byte

	Nets []IPv4Net
}

func (r *Router) decodeToken(tokStr string) (*RoutingInfo, error) {

	// base64 decode tokj
	tokj, err := base64.StdEncoding.DecodeString(tokStr)
	if err != nil {
		return nil, err
	}

	// ipid and iv
	ipId := binary.BigEndian.Uint32(tokj[:4])
	tokj = tokj[4:]
	iv := tokj[:8]
	tokj = tokj[8:]

	// decrypt tokj
	fishKey, fishKeyExists := r.pool.Fish(ipId)
	if !fishKeyExists {
		return nil, errors.New("IpId doesn't registered")
	}
	fish, err := blowfish.NewCipher(fishKey)
	if err != nil {
		panic(err) // should not be here
	}
	dec := cipher.NewCFBDecrypter(fish, iv)
	dec.XORKeyStream(tokj, tokj)

	// second 4 bytes as checksum
	checksum := binary.BigEndian.Uint32(tokj[:4])
	tokj = tokj[4:]

	// crc32 checksum
	if _checksum := crc32.ChecksumIEEE(tokj); _checksum != checksum {
		return nil, errors.New("crc32 checksum doesn't match")
	}

	// account identify and checksum
	tok := RoutingInfo{
		IpId:   ipId,
		Cipher: fish,
		IvRcv:  tokj[:8],
		IvTx:   tokj[8:16],
	}
	tokj = tokj[16:]

	// rand
	tokj = tokj[8:]

	// masks
	masks := tokj[:4]
	tokj = tokj[4:]

	// routing ip net
	nets := make([]uint32, 0, 4)
	for len(tokj) >= 4 {
		nets = append(nets, binary.BigEndian.Uint32(tokj[:4]))
		tokj = tokj[4:]
	}
	tok.Nets = r.parseIPNet(masks, nets)

	return &tok, nil
}

func (r *Router) parseIPNet(masks []byte, nets []uint32) []IPv4Net {

	ipnets := make([]IPv4Net, 0, 4)

	// parse account's ipnet
	for i := 0; i < len(nets); i++ {

		mask := (uint32(0x1) << masks[i]) - 1
		ipnets = append(ipnets, IPv4Net{
			IP:   IPv4(nets[i]),
			Mask: IPv4(mask),
		})
	}

	return ipnets
}
