// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package comet

import (
	"github.com/golang/glog"
	"io"
	"net/http"
)

type PreyGuy struct {
	Req *http.Request
	Out io.Writer
}

// if input is not null, prey guy will use post method, otherwise get
func NewPreyGuy(url string, out io.Writer, input *io.LimitedReader) (*PreyGuy, error) {

	if _, ok := out.(http.Flusher); !ok {
		glog.Warningln("prey guy got a non-flusher output.")
	}

	method := "GET"
	contentLen := int64(0)
	if input != nil {
		method = "POST"
		contentLen = input.N
	}

	req, err := http.NewRequest(method, url, input)
	if err != nil {
		// the only error will be possiblly caused is url paring error
		return nil, err
	}
	req.ContentLength = contentLen

	return &PreyGuy{
		Req: req,
		Out: out,
	}, nil
}

func (this *PreyGuy) Loop() error {

	cli := new(http.Client)

	resp, err := cli.Do(this.Req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentLen := resp.ContentLength
	content := io.LimitReader(resp.Body, contentLen)
	copied, err := CopyInTime(this.Out, content)
	if err != nil {
		return err
	}
	if copied != contentLen {
		glog.Warningln("supposed copied ", contentLen, "bytes, actually copied ", copied)
	}
	return nil
}
