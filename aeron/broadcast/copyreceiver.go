package broadcast

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/buffers"
)

type CopyReceiver struct {
	receiver      *Receiver
	scratchBuffer *buffers.Atomic
}

func (bcast *CopyReceiver) Init(receiver *Receiver) {
	bcast.receiver = receiver
	bcast.scratchBuffer = buffers.MakeAtomic(make([]byte, 4096))

	for bcast.receiver.receiveNext() {
	}
}

func (bcast *CopyReceiver) Receive(handler buffers.Handler) int {
	messagesReceived := 0
	lastSeenLappedCount := bcast.receiver.GetLappedCount()

	if bcast.receiver.receiveNext() {
		if lastSeenLappedCount != bcast.receiver.GetLappedCount() {
			panic("Unable to keep up with broadcast buffer")
		}

		length := bcast.receiver.length()
		if length > bcast.scratchBuffer.Capacity() {
			panic(fmt.Sprintf("Buffer required size %d but only has %d", length, bcast.scratchBuffer.Capacity()))
		}

		msgTypeId := bcast.receiver.typeId()
		bcast.scratchBuffer.PutBytes(0, bcast.receiver.buffer, bcast.receiver.offset(), length)

		if !bcast.receiver.Validate() {
			panic("Unable to keep up with broadcast buffer")
		}

		// DEBUG
		// fmt.Printf("CopyBroadcastReceiver.Receive: ")
		// bytes := bcast.scratchBuffer.GetBytesArray(0, length)
		// for i := 0; i < int(length); i++ {
		// 	fmt.Printf("%x ", bytes[i])
		// }
		// fmt.Printf("\n")
		// END DEBUG

		handler(msgTypeId, bcast.scratchBuffer, 0, length)

		messagesReceived = 1
	}

	return messagesReceived
}
