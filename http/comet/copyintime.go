// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package comet

import (
	"io"
	"net/http"
)

func CopyInTime(dst io.Writer, src io.Reader, n int64) (int64, error) {

	// dst flusher
	flush := func() {}
	if f, ok := dst.(http.Flusher); ok {
		flush = func() {
			f.Flush()
		}
	}

	// copied from golang
	written := int64(0)
	err := error(nil)

	buf := make([]byte, 1500)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			flush()
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil { // er is possibly io.EOF
			err = er
			break
		}
	}
	return written, err
}
