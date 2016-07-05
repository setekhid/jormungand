// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"errors"
	"github.com/ryszard/goskiplist/skiplist"
)

type IPv4 uint32
type IPv4Net struct {
	IP   IPv4
	Mask IPv4
}

type RTable skiplist.SkipList

func NewRTable() *RTable {

	return (*RTable)(skiplist.NewCustomMap(func(l, r interface{}) bool {
		lnet := l.(*IPv4Net)
		rnet := r.(*IPv4Net)
		return lnet.IP < rnet.IP || (lnet.IP == rnet.IP && lnet.Mask < rnet.Mask)
	}))
}

func (rtab *RTable) LinkTun(ipnet IPv4Net, uid uint32, tun interface{}) {

	if ipnet.IP&ipnet.Mask != 0 {
		panic(errors.New("tun's ip net is not match the mask"))
	}

	var tuns map[uint32]interface{}
	if v, ok := (*skiplist.SkipList)(rtab).Get(&ipnet); ok {

		tuns = v.(map[uint32]interface{})
	} else {

		tuns = map[uint32]interface{}{}
	}

	tuns[uid] = tun
	(*skiplist.SkipList)(rtab).Set(&ipnet, tuns)
}

func (rtab *RTable) DiscardTun(ipnet IPv4Net, uid uint32) {

	if v, ok := (*skiplist.SkipList)(rtab).Get(&ipnet); ok {

		tuns := v.(map[uint32]interface{})
		delete(tuns, uid)
	}
}

func (rtab *RTable) RouteIP(ip IPv4) map[uint32]interface{} {

	iter := (*skiplist.SkipList)(rtab).Seek(&IPv4Net{
		IP:   ip,
		Mask: 0,
	})
	if iter == nil {
		return nil
	}
	defer iter.Close()

	for iter.Previous() {

		tuns := iter.Value().(map[uint32]interface{})
		ipnet := iter.Key().(*IPv4Net)
		if ipnet.Mask|ipnet.IP == ipnet.Mask|ip {
			return tuns
		}
	}

	return nil
}
