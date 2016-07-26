// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package payload

import (
	"errors"
)

type Payload interface {
	HeaderSize() uint16
	PacketSize(hdr []byte) uint16
	HeaderInfo(hdr []byte) HeaderInfo
	SetHeader(hdr []byte, info HeaderInfo, siz uint16)
}

type HeaderInfo struct {
	SrcIP uint32
	DstIP uint32
}

func NewPayload(name string) Payload {
	switch name {
	case "tunip":
		return TunIP{}
	case "socks5":
		return Socks5{}
	default:
		panic(errors.New("unknow payload type"))
	}
}
