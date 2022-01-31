/*
Copyright 2016 Stanislav Liberman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aeron

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
)

// Position is a wrapper for a buffer location of a position counter
type Position struct {
	buffer *atomic.Buffer
	id     int32
	offset int32
}

// NewPosition is a factory method to create new Position wrappers
func NewPosition(buffer *atomic.Buffer, id int32) Position {
	var pos Position

	pos.buffer = buffer
	pos.id = id
	pos.offset = id * 2 * util.CacheLineLength

	logger.Debugf("<+ Counter[%d]: %d", id, pos.buffer.GetInt64(pos.offset))

	return pos
}

//go:norace
func (pos *Position) get() int64 {
	p := pos.buffer.GetInt64(pos.offset)
	return p
}

//go:norace
func (pos *Position) set(value int64) {
	pos.buffer.PutInt64(pos.offset, value)
}
