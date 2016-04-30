// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

// +build darwin

package tungo

import (
	"io"
	"os"
)

func newTun(info IfInfo) (io.ReadWriteCloser, error) {
	return os.OpenFile(info.DevFile, os.O_RDWR, 0)
}

func newTap(info IfInfo) (io.ReadWriteCloser, error) {
	return os.OpenFile(info.DevFile, os.O_RDWR, 0)
}
