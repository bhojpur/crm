package models

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

// Priority
const (
	Min         = -300
	Low         = -200
	BelowNormal = -100
	Normal      = 0
	AboveNormal = 100
	High        = 200
	Max         = 300
)

// NewEvent new an basic event instance
func NewEvent(name string, data M) Event {
	if data == nil {
		data = make(map[string]interface{})
	}

	return Event{
		Name: name,
		data: data,
	}
}

// Abort abort event loop exec
func (event *Event) Abort(abort bool) {
	event.aborted = abort
}

// Fill event data
func (event *Event) Fill(target interface{}, data M) *Event {
	if data != nil {
		event.data = data
	}

	event.target = target
	return event
}

// AttachTo add current event to the event manager.
func (event *Event) AttachTo(em ManagerFace) {
	em.AddEvent(*event)
}

// Get get data by index
func (event *Event) Get(key string) interface{} {
	if v, ok := event.data[key]; ok {
		return v
	}

	return nil
}

// Add value by key
func (event *Event) Add(key string, val interface{}) {
	if _, ok := event.data[key]; !ok {
		event.Set(key, val)
	}
}

// Set value by key
func (event *Event) Set(key string, val interface{}) {
	if event.data == nil {
		event.data = make(map[string]interface{})
	}

	event.data[key] = val
}

// Name get event name
func (event *Event) GetName() string {
	return event.Name
}

// Data get all data
func (event *Event) Data() map[string]interface{} {
	return event.data
}

func (event *Event) RecipientList() []uint {
	if event.recipientList == nil {
		event.recipientList = []uint{}
	}
	return event.recipientList
}

// IsAborted check.
func (event *Event) IsAborted() bool {
	return event.aborted
}

// Target get target
func (event *Event) Target() interface{} {
	return event.target
}

// SetName set event name
func (event *Event) SetName(name string) *Event {
	event.Name = name
	return event
}

func (event *Event) SetData(data M) *Event {
	if data != nil {
		event.data = data
	}
	return event
}

func (event *Event) SetPayload(payload M) *Event {
	event.Set("payload", payload)
	return event
}

func (event *Event) SetRecipientList(recipientList []uint) *Event {
	if recipientList != nil {
		event.recipientList = recipientList
	}
	return event
}

// SetTarget set event target
func (event *Event) SetTarget(target interface{}) *Event {
	event.target = target
	return event
}
