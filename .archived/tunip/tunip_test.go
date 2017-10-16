// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package tungo

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"testing"
	"time"
)

func TestTunip(t *testing.T) {

	t.Parallel()

	mtu := 1234
	t.Log("mtu:", mtu)

	nam, tun, err := NewTunIP(mtu)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	go func() {

		defer tun.Close()
		<-time.After(3 * time.Minute)
	}()

	t.Log("tun name:", nam)
	t.Log("if packet offset:", IFPKT_OFFSET)

	for true {

		p := make([]byte, mtu)

		n, e := tun.Read(p)
		if e != nil {
			t.Log("failed reading tun packet -", e)
			return
		}

		go func() {

			pkt := p[IFPKT_OFFSET:n]
			t.Logf("tun header: [% x]", p[:IFPKT_OFFSET])
			t.Logf("tun packet: [% x]", p[IFPKT_OFFSET:n])

			iphdr, err := ipv4.ParseHeader(pkt)
			if err != nil {
				t.Log("parsing ipv4 header", err)
				return
			}

			iphdr.Dst, iphdr.Src = iphdr.Src, iphdr.Dst
			iphdr_msg, err := iphdr.Marshal()
			if err != nil {
				t.Log("marshal ipv4 header", err)
				return
			}
			copy(pkt, iphdr_msg)

			icmphdr, err := icmp.ParseMessage(1, pkt[iphdr.Len:])
			if err != nil {
				t.Log("icmp parse error", err)
				return
			}

			icmphdr.Type = ipv4.ICMPTypeEchoReply
			icmphdr_msg, err := icmphdr.Marshal(nil)
			if err != nil {
				t.Log("icmp marshal error", err)
				return
			}
			copy(pkt[iphdr.Len:], icmphdr_msg)

			wn, err := tun.Write(pkt)
			t.Log("write", wn, "got err", err)
		}()
	}
}
