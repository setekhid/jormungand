// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jargs

import (
	"testing"
)

func TestJargsParsing(t *testing.T) {

	type Module1 struct {
		A int    `json:"a"`
		B string `json:"b"`
	}

	type Module2 struct {
		C float64 `json:"c"`
		D []int   `json:"d"`
	}

	type Module3 struct {
		K bool
	}

	ja := Jargs{}
	ja.Regist("m2", &Module2{})
	ja.Regist("m1", &Module1{})
	ja.Regist("m3", &Module3{})

	err := ja.Parse(`
{
	"m1": {
		"a": 10,
		"b": "abccc"
	},
	"m2": {
		"c": 0.232,
		"d": [1,2,3,4]
	},
	"yoyo": "feafeafea"
}
	`)
	if err != nil {
		t.Error(err)
	}

	t.Log(ja.Module("m2").(*Module2))
	t.Log(ja.Module("m1").(*Module1))

	if ja.Module("m2").(*Module2).C != 0.232 ||
		ja.Module("m2").(*Module2).D[2] != 3 ||
		ja.Module("m1").(*Module1).A != 10 ||
		ja.Module("m1").(*Module1).B != "abccc" {
		t.Fail()
	}
}
