package event

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

type Event interface {
	Name() string
	// Target() interface{}
	Get(key string) interface{}
	Add(key string, val interface{})
	Set(key string, val interface{})
	Data() map[string]interface{}
	SetData(M) Event
	Abort(bool)
	IsAborted() bool
}

// BasicEvent a basic event struct define.
type BasicEvent struct {
	// event name
	name string
	// user data.
	data map[string]interface{}
	// target
	target interface{}
	// mark is aborted
	aborted bool
}

// NewBasic new an basic event instance
func NewBasic(name string, data M) *BasicEvent {
	if data == nil {
		data = make(map[string]interface{})
	}

	return &BasicEvent{
		name: name,
		data: data,
	}
}

// Abort abort event loop exec
func (e *BasicEvent) Abort(abort bool) {
	e.aborted = abort
}

// Fill event data
func (e *BasicEvent) Fill(target interface{}, data M) *BasicEvent {
	if data != nil {
		e.data = data
	}

	e.target = target
	return e
}

// AttachTo add current event to the event manager.
func (e *BasicEvent) AttachTo(em ManagerFace) {
	em.AddEvent(e)
}

// Get get data by index
func (e *BasicEvent) Get(key string) interface{} {
	if v, ok := e.data[key]; ok {
		return v
	}

	return nil
}

// Add value by key
func (e *BasicEvent) Add(key string, val interface{}) {
	if _, ok := e.data[key]; !ok {
		e.Set(key, val)
	}
}

// Set value by key
func (e *BasicEvent) Set(key string, val interface{}) {
	if e.data == nil {
		e.data = make(map[string]interface{})
	}

	e.data[key] = val
}

// Name get event name
func (e *BasicEvent) Name() string {
	return e.name
}

// Data get all data
func (e *BasicEvent) Data() map[string]interface{} {
	return e.data
}

// IsAborted check.
func (e *BasicEvent) IsAborted() bool {
	return e.aborted
}

// Target get target
func (e *BasicEvent) Target() interface{} {
	return e.target
}

// SetName set event name
func (e *BasicEvent) SetName(name string) *BasicEvent {
	e.name = name
	return e
}

// SetData set data to the event
func (e *BasicEvent) SetData(data M) Event {
	if data != nil {
		e.data = data
	}
	return e
}

// SetTarget set event target
func (e *BasicEvent) SetTarget(target interface{}) *BasicEvent {
	e.target = target
	return e
}
