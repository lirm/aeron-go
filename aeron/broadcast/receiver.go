package broadcast

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/util"
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

	// std::atomic<long> lappedCount
	lappedCount int64
}

func (recv *Receiver) Init(buffer *buffers.Atomic) {
	recv.buffer = buffer
	recv.capacity = buffer.Capacity() - BufferDescriptor.TRAILER_LENGTH
	recv.mask = int64(recv.capacity) - 1
	recv.tailIntentCounterIndex = recv.capacity + BufferDescriptor.TAIL_INTENT_COUNTER_OFFSET
	recv.tailCounterIndex = recv.capacity + BufferDescriptor.TAIL_COUNTER_OFFSET
	recv.latestCounterIndex = recv.capacity + BufferDescriptor.LATEST_COUNTER_OFFSET
}

func (recv *Receiver) Validate() bool {
	// load fence = acquire()
	// atomic::acquire();

	return recv.validate(recv.cursor)
}

func (recv *Receiver) validate(cursor int64) bool {
	return (cursor + int64(recv.capacity)) > recv.buffer.GetInt64Volatile(recv.tailIntentCounterIndex)
}

func (recv *Receiver) GetLappedCount() int64 {
	// FIXME this needs to be atomic
	return recv.lappedCount
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
			recv.lappedCount += 1
			cursor = recv.buffer.GetInt64(recv.latestCounterIndex)
			recordOffset = int32(cursor & recv.mask)
		}

		recv.cursor = cursor
		recv.nextRecord = cursor + int64(util.AlignInt32(recv.buffer.GetInt32(buffers.LengthOffset(recordOffset)), int32(buffers.RecordDescriptor.RECORD_ALIGNMENT)))

		if buffers.RecordDescriptor.PADDING_MSG_TYPE_ID == recv.buffer.GetInt32(buffers.TypeOffset(recordOffset)) {
			recordOffset = 0
			recv.cursor = recv.nextRecord
			recv.nextRecord += int64(util.AlignInt32(recv.buffer.GetInt32(buffers.LengthOffset(recordOffset)), int32(buffers.RecordDescriptor.RECORD_ALIGNMENT)))
		}

		recv.recordOffset = recordOffset
		isAvailable = true
	}

	return isAvailable
}
