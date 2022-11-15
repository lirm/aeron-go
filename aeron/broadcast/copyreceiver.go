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
	"errors"
	"fmt"

	"github.com/lirm/aeron-go/aeron/atomic"
)

type Handler func(int32, *atomic.Buffer, int32, int32) error

type CopyReceiver struct {
	receiver      *Receiver
	scratchBuffer *atomic.Buffer
}

func NewCopyReceiver(receiver *Receiver) *CopyReceiver {
	bcast := new(CopyReceiver)
	bcast.receiver = receiver
	bcast.scratchBuffer = atomic.MakeBuffer(make([]byte, 4096))

	// Scroll to the latest unprocessed
	for bcast.receiver.receiveNext() {
	}
	return bcast
}

func (bcast *CopyReceiver) Receive(handler Handler) (int, error) {
	messagesReceived := 0
	lastSeenLappedCount := bcast.receiver.GetLappedCount()

	if bcast.receiver.receiveNext() {
		if lastSeenLappedCount != bcast.receiver.GetLappedCount() {
			return 0, errors.New("unable to keep up with broadcast buffer")
		}

		length := bcast.receiver.length()
		if length > bcast.scratchBuffer.Capacity() {
			return 0, fmt.Errorf("buffer required size %d but only has %d",
				length, bcast.scratchBuffer.Capacity())
		}

		msgTypeID := bcast.receiver.typeID()
		bcast.scratchBuffer.PutBytes(0, bcast.receiver.buffer, bcast.receiver.offset(), length)

		if !bcast.receiver.Validate() {
			return 0, errors.New("unable to keep up with broadcast buffer")
		}

		if err := handler(msgTypeID, bcast.scratchBuffer, 0, length); err != nil {
			return 0, err
		}

		messagesReceived = 1
	}

	return messagesReceived, nil
}
