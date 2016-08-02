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

package buffer

import (
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("buffer")

type Position struct {
	buffer *Atomic
	id     int32
	offset int32
}

func NewPosition(buffer *Atomic, id int32) Position {
	pos := new(Position)

	pos.buffer = buffer
	pos.id = id
	pos.offset = id * 2 * util.CACHE_LINE_LENGTH

	logger.Debugf("<+ Counter[%d]: %d", id, pos.buffer.GetInt64(pos.offset))

	return *pos
}

func (pos *Position) Get() int64 {
	p := pos.buffer.GetInt64(pos.offset)
	return p
}

func (pos *Position) Set(value int64) {
	pos.buffer.PutInt64(pos.offset, value)
}
