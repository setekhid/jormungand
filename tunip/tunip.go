// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package tungo

//#include "tun_fd.h"
import "C"

import (
	"errors"
	"io"
	"os"
)

var IFPKT_OFFSET = int(C.IFPKT_OFFSET)

func NewTunIP(mtu int) (nam string, tun io.ReadWriteCloser, err error) {

	info, err := C.create_tun_fd(C.int(mtu))
	if err != nil {
		return "", nil, err
	}

	nam = C.GoString(&info.nam[0])
	tun = os.NewFile(uintptr(info.fd), nam)
	if tun == nil {
		return "", nil, errors.New("tun fd is not valid")
	}

	return nam, tun, nil
}
