// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/setekhid/jormungand/http/comet"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type Router struct {
	rtab *RTable
	ulns map[uint32]LinkRef
}

func NewRouter() *Router {

	return &Router{
		rtab: NewRTable(),
		ulns: map[uint32]LinkRef{},
	}
}

// Override comet.TunnelAuthor.Auth
func (r *Router) Auth(tokStr string, inlen int64) (tun io.ReadWriteCloser, outlen int64, err error) {

	// TODO
	return nil, 0, nil
}

func (r *Router) auth_unsafe(tokStr string, inlen int64) (tun io.ReadWriteCloser, outlen int64, err error) {

	// account's token json
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
		uln := r.refLink(NewLink(), r.mkDereference(tok.AccountIP, ipnets))
		r.ulns[tok.AccountIP] = uln
		// regist self route
		r.rtab.LinkIP(tok.AccountIP, uln.Link)
		// regist ipv4 network routes
		for _, ipnet := range ipnets {
			r.rtab.LinkTun(ipnet, tok.AccountIP, uln.Link)
		}
	} else {

		// exists, just clone
		uln = uln.Clone()
	}

	return uln, comet.HTTP_DL_NORMAL_LEN, nil
}

func (r *Router) Kickoff(ip uint32, ipnets []IPv4Net) {

	// TODO
}

func (r *Router) mkDereference(ip uint32, ipnets []IPv4Net) func() {

	return func() {
		r.Kickoff(ip, ipnets)
	}
}

func (r *Router) kickoff_unsafe(ip uint32, ipnets []IPv4Net) {

	// discard ipv4 network routes
	for _, ipnet := range ipnets {
		r.rtab.DiscardTun(ipnet, ip)
	}
	// discard self route
	r.rtab.DiscardIP(ip)
	// discard link reference
	delete(r.ulns, ip)
}

func (r *Router) Routing() {

	// TODO
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

	// TODO checksum

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
