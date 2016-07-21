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

package buffers

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/util"
)

const INSUFFICIENT_CAPACITY int32 = -2

type ManyToOneRingBuffer struct {
	buffer                    *Atomic
	capacity                  int32
	maxMsgLength              int32
	headPositionIndex         int32
	headCachePositionIndex    int32
	tailPositionIndex         int32
	correlationIdCounterIndex int32
	consumerHeartbeatIndex    int32
}

func (buf *ManyToOneRingBuffer) Init(buffer *Atomic) *ManyToOneRingBuffer {

	buf.buffer = buffer
	buf.capacity = buffer.Capacity() - RingBufferDescriptor.TRAILER_LENGTH

	util.IsPowerOfTwo(buf.capacity)

	buf.maxMsgLength = buf.capacity / 8
	buf.tailPositionIndex = buf.capacity + RingBufferDescriptor.TAIL_POSITION_OFFSET
	buf.headCachePositionIndex = buf.capacity + RingBufferDescriptor.HEAD_CACHE_POSITION_OFFSET
	buf.headPositionIndex = buf.capacity + RingBufferDescriptor.HEAD_POSITION_OFFSET
	buf.correlationIdCounterIndex = buf.capacity + RingBufferDescriptor.CORRELATION_COUNTER_OFFSET
	buf.consumerHeartbeatIndex = buf.capacity + RingBufferDescriptor.CONSUMER_HEARTBEAT_OFFSET

	return buf
}

func (buf *ManyToOneRingBuffer) NextCorrelationId() int64 {
	return buf.buffer.GetAndAddInt64(buf.correlationIdCounterIndex, 1)
}

func (buf *ManyToOneRingBuffer) SetConsumerHeartbeatTime(time int64) {
	buf.buffer.PutInt64Ordered(buf.consumerHeartbeatIndex, time)
}

func (buf *ManyToOneRingBuffer) ConsumerHeartbeatTime() int64 {
	return buf.buffer.GetInt64Volatile(buf.consumerHeartbeatIndex)
}

func (buf *ManyToOneRingBuffer) ProducerPosition() int64 {
	return buf.buffer.GetInt64Volatile(buf.tailPositionIndex)
}

func (buf *ManyToOneRingBuffer) ConsumerPosition() int64 {
	return buf.buffer.GetInt64Volatile(buf.headPositionIndex)
}

func (buf *ManyToOneRingBuffer) Capacity() int32 {
	return buf.capacity
}

func (buf *ManyToOneRingBuffer) claimCapacity(requiredCapacity int32) int32 {

	mask := buf.capacity - 1
	head := buf.buffer.GetInt64Volatile(buf.headCachePositionIndex)

	var tail int64
	var tailIndex int32
	var padding int32

	for ok := true; ok; ok = !buf.buffer.CompareAndSetInt64(buf.tailPositionIndex, tail, tail+int64(requiredCapacity)+int64(padding)) {
		tail = buf.buffer.GetInt64Volatile(buf.tailPositionIndex)
		availableCapacity := buf.capacity - int32(tail-head)

		if requiredCapacity > availableCapacity {
			head = buf.buffer.GetInt64Volatile(buf.headPositionIndex)

			if requiredCapacity > (buf.capacity - int32(tail-head)) {
				return INSUFFICIENT_CAPACITY
			}

			buf.buffer.PutInt64Ordered(buf.headCachePositionIndex, head)
		}

		padding = 0
		tailIndex = int32(tail & int64(mask))
		toBufferEndLength := buf.capacity - tailIndex

		if requiredCapacity > toBufferEndLength {
			headIndex := int32(head & int64(mask))

			if requiredCapacity > headIndex {
				head = buf.buffer.GetInt64Volatile(buf.headPositionIndex)
				headIndex = int32(head & int64(mask))

				if requiredCapacity > headIndex {
					return INSUFFICIENT_CAPACITY
				}

				buf.buffer.PutInt64Ordered(buf.headCachePositionIndex, head)
			}

			padding = toBufferEndLength
		}
	}

	if 0 != padding {
		buf.buffer.PutInt64Ordered(tailIndex, MakeHeader(int32(padding), RecordDescriptor.PADDING_MSG_TYPE_ID))
		tailIndex = 0
	}

	return tailIndex
}

func (buf *ManyToOneRingBuffer) checkMsgLength(length int32) {
	if length > buf.maxMsgLength {
		panic(fmt.Sprintf("encoded message exceeds maxMsgLength of %d, length=%d", buf.maxMsgLength, length))
	}
}

func (buf *ManyToOneRingBuffer) Write(msgTypeId int32, srcBuffer *Atomic, srcIndex int32, length int32) bool {

	var isSuccessful bool = false

	CheckMsgTypeId(msgTypeId)
	buf.checkMsgLength(length)

	recordLength := length + RecordDescriptor.HEADER_LENGTH
	requiredCapacity := util.AlignInt32(recordLength, RecordDescriptor.RECORD_ALIGNMENT)
	recordIndex := buf.claimCapacity(requiredCapacity)

	if INSUFFICIENT_CAPACITY != recordIndex {
		buf.buffer.PutInt64Ordered(recordIndex, MakeHeader(-recordLength, msgTypeId))
		buf.buffer.PutBytes(EncodedMsgOffset(recordIndex), srcBuffer, srcIndex, length)
		buf.buffer.PutInt32Ordered(LengthOffset(recordIndex), recordLength)

		isSuccessful = true
	}

	return isSuccessful
}

func (buf *ManyToOneRingBuffer) Read(Handler, messageCountLimit int) int32 {
	panic("Not implemented yet")
	return -1
}
