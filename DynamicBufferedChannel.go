// Buffered channel functionality that grows in size as needed.
// Closing Put will close all buffers and eventually close Get.
// You will then have to call Init() to reuse the buffer.
// TODO: remove channels from list as total values become less.
package goBuffers

import (
	"container/list"
)

type DynamicBufferedChannel struct {
	Put chan interface{}   // User sends to this channel.
	Get <-chan interface{} // User receives from this channel.

	buffers *list.List
}

// Create new buffer with initial size.
func New(size int) *DynamicBufferedChannel {
	if size < 1 {
		size = 1
	}
	return (&DynamicBufferedChannel{buffers: new(list.List)}).Init(size)
}

// Reset buffer with initial size.
func (dbc *DynamicBufferedChannel) Init(size int) *DynamicBufferedChannel {
	dbc.Put = make(chan interface{})
	dbc.buffers.Init()
	dbc.Get = dbc.buffers.PushBack(make(chan interface{}, size)).Value.(chan interface{})

	go handle(dbc)
	return dbc
}

func handle(dbc *DynamicBufferedChannel) {
	for val := range dbc.Put {
		lastBuf := dbc.buffers.Back().Value.(chan interface{})
		if len(lastBuf) == cap(lastBuf) {
			newLastBuffer := dbc.buffers.PushBack(make(chan interface{}, cap(lastBuf)*2)).Value.(chan interface{})
			go func(to, from chan interface{}) {
				for val := range from {
					to <- val
				}
				close(to)
			}(lastBuf, newLastBuffer)
			newLastBuffer <- val
		} else {
			lastBuf <- val
		}
	}
	close(dbc.buffers.Back().Value.(chan interface{}))
}
