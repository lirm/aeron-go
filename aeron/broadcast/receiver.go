package broadcast

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/util"
	"sync/atomic"
)

type Receiver struct {
	buffer                 *buffers.Atomic
	capacity               int32
	mask                   int64
	tailIntentCounterIndex int32
	tailCounterIndex       int32
	latestCounterIndex     int32

	recordOffset int32
	cursor       int64
	nextRecord   int64

	lappedCount int64
}

func NewReceiver(buffer *buffers.Atomic) *Receiver {
	recv := new(Receiver)
	recv.buffer = buffer
	recv.capacity = buffer.Capacity() - BufferDescriptor.TRAILER_LENGTH
	recv.mask = int64(recv.capacity) - 1
	recv.tailIntentCounterIndex = recv.capacity + BufferDescriptor.TAIL_INTENT_COUNTER_OFFSET
	recv.tailCounterIndex = recv.capacity + BufferDescriptor.TAIL_COUNTER_OFFSET
	recv.latestCounterIndex = recv.capacity + BufferDescriptor.LATEST_COUNTER_OFFSET
	recv.lappedCount = 0

	CheckCapacity(recv.capacity)

	return recv
}

func (recv *Receiver) Validate() bool {
	return recv.validate(recv.cursor)
}

func (recv *Receiver) validate(cursor int64) bool {
	return (cursor + int64(recv.capacity)) > recv.buffer.GetInt64Volatile(recv.tailIntentCounterIndex)
}

func (recv *Receiver) GetLappedCount() int64 {
	return atomic.LoadInt64(&recv.lappedCount)
}

func (recv *Receiver) typeId() int32 {
	return recv.buffer.GetInt32(buffers.TypeOffset(recv.recordOffset))
}

func (recv *Receiver) offset() int32 {
	return buffers.EncodedMsgOffset(recv.recordOffset)
}

func (recv *Receiver) length() int32 {
	return int32(recv.buffer.GetInt32(buffers.LengthOffset(recv.recordOffset))) - buffers.RecordDescriptor.HEADER_LENGTH
}

func (recv *Receiver) receiveNext() bool {
	isAvailable := false

	tail := recv.buffer.GetInt64Volatile(recv.tailCounterIndex)
	cursor := recv.nextRecord

	if tail > cursor {
		recordOffset := int32(cursor & recv.mask)

		if !recv.validate(cursor) {
			atomic.AddInt64(&recv.lappedCount, 1)
			cursor = recv.buffer.GetInt64(recv.latestCounterIndex)
			recordOffset = int32(cursor & recv.mask)
		}

		recv.cursor = cursor
		length := recv.buffer.GetInt32(buffers.LengthOffset(recordOffset))
		alignedLength := int64(util.AlignInt32(length, buffers.RecordDescriptor.RECORD_ALIGNMENT))
		recv.nextRecord = cursor + alignedLength

		if buffers.RecordDescriptor.PADDING_MSG_TYPE_ID == recv.buffer.GetInt32(buffers.TypeOffset(recordOffset)) {
			recordOffset = 0
			recv.cursor = recv.nextRecord
			recv.nextRecord += alignedLength
		}

		recv.recordOffset = recordOffset
		isAvailable = true
	}

	return isAvailable
}
