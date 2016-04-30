// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

// +build linux

package tungo

// #include <linux/if.h>
// #include <linux/if_tun.h>
import "C"

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	LINUX_TUN_DEV = "/dev/net/tun"
)

func newTunTap(info IfInfo, typeFlag int) {

	if len(info.DevFile) <= 0 {
		info.DevFile = LINUX_TUN_DEV
	}

	dev, err := os.OpenFile(info.DevFile, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	err = createIface(dev.Fd(), info.IfName, typeFlag|C.IFF_NO_PI)
	if err != nil {
		dev.Close()
		return nil, err
	}

	return dev, nil
}

func createIface(fd uintptr, ifName string, flags uint16) error {

	var req C.ifreq
	copy(req.ifr_name[:], ifName)
	req.ifr_flags = flags
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		return errno
	}
	return nil
}

func newTun(info IfInfo) (io.ReadWriteCloser, error) {
	return newTunTap(info, C.IFF_TUN)
}

func newTap(info IfInfo) (io.ReadWriteCloser, error) {
	return newTunTap(info, C.IFF_TAP)
}
