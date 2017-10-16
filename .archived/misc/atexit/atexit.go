// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package atexit

import (
	"io"
)

var (
	cleaner = NewCleaner()
)

func Reg(c io.Closer) { cleaner.AtExit(c) }
func Exit()           { cleaner.CleanMain() }

type Cleaner struct {
	closers stack
}

func NewCleaner() *Cleaner { return &Cleaner{} }

func (cl *Cleaner) AtExit(c io.Closer) { cl.closers = append(cl.closers, c) }

func (cl *Cleaner) CleanScope() error {

	err := cl.closers.front().(io.Closer).Close()
	cl.closers.pop()
	return err
}

func (cl *Cleaner) CleanMain() []error {

	errs := make([]error, 0, cl.closers.count())
	for !cl.closers.empty() {

		errs = append(errs, cl.CleanScope())
	}

	return errs
}

type stack []interface{}

func (s stack) push(e interface{}) stack { return append(s, e) }
func (s stack) pop() stack               { return s[:len(s)-1] }
func (s stack) empty() bool              { return s.count() <= 0 }
func (s stack) count() int               { return len(s) }
func (s stack) front() interface{}       { return s[len(s)-1] }
