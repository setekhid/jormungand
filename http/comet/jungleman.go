// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package comet

import (
	"github.com/golang/glog"
	"io"
	"net/http"
	"strconv"
)

const (
	HTTP_DL_NORMAL_LEN = int64(60000)
)

type TunnelAuthor interface {
	// if writable less than 0, the result will not limit output content length
	// if readable less than 0, JungleMan will keep read tunnel untill EOF
	// err marking any error occured, the tunnel should be nil
	Auth(token string, writable int64) (tunnel io.ReadWriteCloser, readable int64, err error)
}

type JungleMan struct {
	Author TunnelAuthor
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
	tunnel, contentLen, _ := this.Author.Auth(req.RequestURI, req.ContentLength)
	if tunnel == nil { // didn't pass the auth, omit the error
		return
	}
	defer tunnel.Close()

	if req.Method == "POST" { // the post data reading

		input := req.Body.(io.Reader)
		if req.ContentLength >= 0 {
			input = io.LimitReader(req.Body, req.ContentLength)
		}
		_, err := CopyInTime(tunnel, input)
		if err != nil && err != io.EOF { // for now, just collecting it in log
			glog.Warningln("failed when jungle man is reading from http body, ", err)
			// keep going, cause this is an alive tunnel, not one request
		}
	}

	input := tunnel.(io.Reader)
	// specified response content length
	if contentLen >= 0 {
		resp.Header().Set("Content-Length", strconv.FormatInt(contentLen, 10))
		input = io.LimitReader(input, contentLen)
	}

	// return the data
	_, err := CopyInTime(resp, input)
	if err != nil && err != io.EOF { // for now, just collecting it in log
		glog.Warningln("failed when jungle man is writting back, ", err)
	}
}
