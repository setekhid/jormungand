// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package edgo

type Event interface {
	Type() int
}

type EventFunc func(interface{}, Event)

type EdGo struct {
	self interface{}

	eventCh <-chan Event
	eventFn map[int]EventFunc
}

func NewEdGo(self interface{}, eventCh <-chan Event) *EdGo {

	return &EdGo{
		self: self,

		eventCh: eventCh,
		eventFn: map[int]EventFunc{},
	}
}

func (d *EdGo) Regist(etype int, efunc EventFunc) { d.eventFn[etype] = efunc }
func (d *EdGo) Process(event Event)               { d.eventFn[event.Type()](d.self, event) }

func (d *EdGo) TryFetch() {

	for true {
		select {
		case event := <-d.eventCh:
			d.Process(event)
		default:
			break
		}
	}
}

func (d *EdGo) Dispatching() {

	for event := range d.eventCh {
		d.Process(event)
	}
}
