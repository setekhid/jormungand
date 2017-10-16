// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package payload

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
)

type TunIP struct{}

func (pl TunIP) HeaderSize() uint16           { return ipv4.HeaderLen }
func (pl TunIP) PacketSize(hdr []byte) uint16 { return binary.BigEndian.Uint16(hdr[2:4]) }

func (pl TunIP) HeaderInfo(hdr []byte) HeaderInfo {
	return HeaderInfo{
		SrcIP: binary.BigEndian.Uint32(hdr[12:16]),
		DstIP: binary.BigEndian.Uint32(hdr[16:20]),
	}
}

func (pl TunIP) SetHeader(hdr []byte, info HeaderInfo, siz uint16) {
	// TODO check and panic
}
