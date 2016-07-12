// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jorm

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type Person struct {
	Name  string `json:"jnn,omitempty" bson:"bnn,omitempty"`
	Phone string `json:"jfff,omitempty" bson:"bfff,omitempty"`
	Lolo  int
}

func TestBson(t *testing.T) {

	data, err := bson.Marshal(&Person{
		Name: "Bob",
		Lolo: 13213,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q", data)
}

func TestMap(t *testing.T) {

	am := map[string]int{}
	bm := am
	bm["bbb"] = 10
	t.Log(am["bbb"])
}
