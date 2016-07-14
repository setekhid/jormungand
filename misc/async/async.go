// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package async

import (
	"sync"
)

const (
	ASYNCER_WAKEUP = 1
)

type Asyncer struct {
	ops    []func()
	locker sync.Mutex
	signal chan int
}

func NewAsyncer() *Asyncer {
	return &Asyncer{
		signal: make(chan int),
	}
}

func StartAsyncer() *Asyncer {
	asyncer := NewAsyncer()
	go asyncer.Loop()
	return asyncer
}

func (this *Asyncer) Loop() {

	for _ = range this.signal {

		var ops []func()

		func() {

			this.locker.Lock()
			defer this.locker.Unlock()
			ops = this.ops
			// TODO algorithm to increase
			this.ops = make([]func(), 0, len(ops)*4/3)
		}()

		for _, op := range ops {
			op()
		}
	}
}

func (this *Asyncer) LoopOnce() {

	var ops []func()

	func() {

		this.locker.Lock()
		defer this.locker.Unlock()
		ops = this.ops
		this.ops = make([]func(), 0, cap(ops))
	}()

	for _, op := range ops {
		op()
	}
}

func (this *Asyncer) Stop() {
	close(this.signal)
}

func (this *Asyncer) InvokeAsync(op func()) {

	func() {
		this.locker.Lock()
		defer this.locker.Unlock()

		this.ops = append(this.ops, op)
	}()

	select {
	case this.signal <- ASYNCER_WAKEUP:
	default:
	}
}

func (this *Asyncer) InvokeSync(op func()) {

	locker := sync.Mutex{}
	cond := sync.Cond{L: &locker}

	locker.Lock()
	defer locker.Unlock()

	this.InvokeAsync(func() {

		locker.Lock()
		defer locker.Unlock()
		defer cond.Signal()

		op()
	})

	cond.Wait()
}
