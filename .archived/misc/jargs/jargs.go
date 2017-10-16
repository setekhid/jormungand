// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package jargs

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
)

var (
	args     = Jargs{}
	isParsed = false
)

func parsed() Jargs {

	if !isParsed {

		content, err := ioutil.ReadFile(json_file)
		if err != nil {
			panic(err)
		}

		err = args.Parse(string(content))
		if err != nil {
			panic(err)
		}

		isParsed = true
	}
	return args
}

func unparsed() Jargs {

	if isParsed {

		glog.Warningln("jargs has already been parsed!")
		isParsed = false
	}
	return args
}

func Regist(name string, body interface{}) { unparsed().Regist(name, body) }
func Module(name string) interface{}       { return parsed().Module(name) }

// Jargs type define
type Jargs map[string]interface{}

// regist a module
func (ja Jargs) Regist(name string, body interface{}) {

	if _, ok := ja[name]; ok {
		panic(errors.New(name + " has already been registed"))
	}
	ja[name] = body
}

func (ja Jargs) Parse(jsonStr string) error {

	// parsing module json to RawMessage
	rawModules := map[string]json.RawMessage{}
	err := json.Unmarshal([]byte(jsonStr), &rawModules)
	if err != nil {
		return err
	}

	// parsing RawMessage to module
	for name, raw := range rawModules {

		if module, ok := ja[name]; ok {

			err := json.Unmarshal(raw, module)
			if err != nil {
				return err
			}

			ja[name] = module
		}
	}

	return nil
}

func (ja Jargs) Module(name string) interface{} { return ja[name] }
