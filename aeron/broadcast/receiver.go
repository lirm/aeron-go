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

package broadcast

import (
	"log"

	"github.com/lirm/aeron-go/aeron/atomic"
	rb "github.com/lirm/aeron-go/aeron/ringbuffer"
	"github.com/lirm/aeron-go/aeron/util"
)

var BufferDescriptor = struct {
	tailIntentCounterOffset int32
	tailCounterOffset       int32
	latestCounterOffset     int32
	trailerLength           int32
}{
	0,
	util.SizeOfInt64,
	util.SizeOfInt64 * 2,
	util.CacheLineLength * 2,
}

func checkCapacity(capacity int32) {
	if !util.IsPowerOfTwo(int64(capacity)) {
		log.Fatalf("Capacity must be a positive power of 2 + TRAILER_LENGTH: capacity=%d", capacity)
	}
}

type Receiver struct {
	buffer                 *atomic.Buffer
	capacity               int32
	mask                   int64
	tailIntentCounterIndex int32
	tailCounterIndex       int32
	latestCounterIndex     int32

	recordOffset int32
	cursor       int64
	nextRecord   int64

	lappedCount atomic.Long
}

func NewReceiver(buffer *atomic.Buffer) *Receiver {
	recv := new(Receiver)
	recv.buffer = buffer
	recv.capacity = buffer.Capacity() - BufferDescriptor.trailerLength
	recv.mask = int64(recv.capacity) - 1
	recv.tailIntentCounterIndex = recv.capacity + BufferDescriptor.tailIntentCounterOffset
	recv.tailCounterIndex = recv.capacity + BufferDescriptor.tailCounterOffset
	recv.latestCounterIndex = recv.capacity + BufferDescriptor.latestCounterOffset
	recv.lappedCount.Set(0)

	checkCapacity(recv.capacity)

	return recv
}

func (recv *Receiver) Validate() bool {
	return recv.validate(recv.cursor)
}

func (recv *Receiver) validate(cursor int64) bool {
	return (cursor + int64(recv.capacity)) > recv.buffer.GetInt64Volatile(recv.tailIntentCounterIndex)
}

func (recv *Receiver) GetLappedCount() int64 {
	return recv.lappedCount.Get()
}

func (recv *Receiver) typeID() int32 {
	return recv.buffer.GetInt32(rb.TypeOffset(recv.recordOffset))
}

func (recv *Receiver) offset() int32 {
	return rb.EncodedMsgOffset(recv.recordOffset)
}

func (recv *Receiver) length() int32 {
	return int32(recv.buffer.GetInt32(rb.LengthOffset(recv.recordOffset))) - rb.RecordDescriptor.HeaderLength
}

func (recv *Receiver) receiveNext() bool {
	isAvailable := false

	tail := recv.buffer.GetInt64Volatile(recv.tailCounterIndex)
	cursor := recv.nextRecord

	if tail > cursor {
		recordOffset := int32(cursor & recv.mask)

		if !recv.validate(cursor) {
			recv.lappedCount.Inc()
			cursor = recv.buffer.GetInt64(recv.latestCounterIndex)
			recordOffset = int32(cursor & recv.mask)
		}

		recv.cursor = cursor
		length := recv.buffer.GetInt32(rb.LengthOffset(recordOffset))
		alignedLength := int64(util.AlignInt32(length, rb.RecordDescriptor.RecordAlignment))
		recv.nextRecord = cursor + alignedLength

		if rb.RecordDescriptor.PaddingMsgTypeID == recv.buffer.GetInt32(rb.TypeOffset(recordOffset)) {
			recordOffset = 0
			recv.cursor = recv.nextRecord
			length := recv.buffer.GetInt32(rb.LengthOffset(recordOffset))
			alignedLength := int64(util.AlignInt32(length, rb.RecordDescriptor.RecordAlignment))
			recv.nextRecord += alignedLength
		}

		recv.recordOffset = recordOffset
		isAvailable = true
	}

	return isAvailable
}
