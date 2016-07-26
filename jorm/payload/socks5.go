// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package payload

import (
	"encoding/binary"
)

type Socks5 struct{}

func (pl Socks5) HeaderSize() uint16           { return 4 + 4 + 2 }
func (pl Socks5) PacketSize(hdr []byte) uint16 { return binary.BigEndian.Uint16(hdr[0:2]) }

func (pl Socks5) HeaderInfo(hdr []byte) HeaderInfo {
	return HeaderInfo{
		SrcIP: binary.BigEndian.Uint32(hdr[2:6]),
		DstIP: binary.BigEndian.Uint32(hdr[6:10]),
	}
}

func (pl Socks5) SetHeader(hdr []byte, info HeaderInfo, siz uint16) {
	binary.BigEndian.PutUint16(hdr[0:2], uint16(siz))
	binary.BigEndian.PutUint32(hdr[2:6], info.SrcIP)
	binary.BigEndian.PutUint32(hdr[6:10], info.DstIP)
}
