// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package comet

import (
	"github.com/golang/glog"
	"io"
	"net/http"
)

type Auth2ReadWriteCloser interface {
	Auth(token string) (tunnel io.ReadWriteCloser, contentLen int64)
}

type JungleMan struct {
	A2RWC Auth2ReadWriteCloser
}

// Override http.Handler.ServeHTTP
func (this *JungleMan) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	if _, isFlusher := resp.(http.Flusher); !isFlusher {
		// maybe the response is not in time
		glog.Warningln("the resp object in jungle man isn't a flusher!")
	}

	if req.Method != "POST" && req.Method != "GET" {
		// simplly report unsupported method
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// authorize
	tunnel, contentLen := this.A2RWC.Auth(req.RequestURI)
	if tunnel == nil { // didn't pass the auth
		return
	}
	defer tunnel.Close()

	if req.Method != "GET" { // the post data reading
		_, err := CopyInTime(tunnel, req.Body, req.ContentLength)
		if err != nil && err != io.EOF { // for now, just collecting it in log
			glog.Warningln("failed when jungle man is reading from http body, ", err)
			// keep going, cause this is an alive tunnel, not one request
		}
	}

	// return the data
	_, err := CopyInTime(resp, tunnel, contentLen)
	if err != nil && err != io.EOF { // for now, just collecting it in log
		glog.Warningln("failed when jungle man is writting back, ", err)
	}
}
