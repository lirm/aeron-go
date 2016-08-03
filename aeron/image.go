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
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/aeron/util"
)

type ControlledPollFragmentHandler func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header)

const (
	ImageClosed int = -1
)

var ControlledPollAction = struct {
	/**
	 * Abort the current polling operation and do not advance the position for this fragment.
	 */
	ABORT int

	/**
	 * Break from the current polling operation and commit the position as of the end of the current fragment
	 * being handled.
	 */
	BREAK int

	/**
	 * Continue processing but commit the position as of the end of the current fragment so that
	 * flow control is applied to this point.
	 */
	COMMIT int

	/**
	 * Continue processing taking the same approach as the in fragment_handler_t.
	 */
	CONTINUE int
}{
	1,
	2,
	3,
	4,
}

type Image struct {
	termBuffers [logbuffer.PartitionCount]*atomic.Buffer
	header      logbuffer.Header

	subscriberPosition Position

	logBuffers *logbuffer.LogBuffers

	sourceIdentity string
	isClosed       atomic.Bool

	exceptionHandler func(error)

	correlationID              int64
	subscriptionRegistrationID int64
	sessionID                  int32
	termLengthMask             int32
	positionBitsToShift        uint8
}

func NewImage(sessionID int32, correlationID int64, logBuffers *logbuffer.LogBuffers) *Image {

	image := new(Image)

	image.correlationID = correlationID
	image.sessionID = sessionID
	image.logBuffers = logBuffers
	for i := 0; i < logbuffer.PartitionCount; i++ {
		image.termBuffers[i] = logBuffers.Buffer(i)
	}
	capacity := logBuffers.Buffer(0).Capacity()
	image.termLengthMask = capacity - 1
	image.positionBitsToShift = util.NumberOfTrailingZeroes(capacity)
	image.header.SetInitialTermID(logbuffer.InitialTermID(logBuffers.Buffer(logbuffer.LogMetaDataSectionIndex)))
	image.header.SetPositionBitsToShift(int32(image.positionBitsToShift))
	image.isClosed.Set(false)

	return image
}

func (image *Image) IsClosed() bool {
	return image.isClosed.Get()
}

func (image *Image) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	result := ImageClosed

	if !image.IsClosed() {
		position := image.subscriberPosition.get()
		//logger.Debugf("Image position: %d, mask:%X", position, image.termLengthMask)
		termOffset := int32(position) & image.termLengthMask
		index := logbuffer.IndexByPosition(position, image.positionBitsToShift)
		termBuffer := image.termBuffers[index]
		//logger.Debugf("Selected Term buffer: %v", termBuffer)

		var readOutcome term.ReadOutcome
		term.Read(&readOutcome, termBuffer, termOffset, handler, fragmentLimit, &image.header)

		newPosition := position + int64(readOutcome.Offset()-termOffset)
		if newPosition > position {
			image.subscriberPosition.set(newPosition)
			//logger.Debugf("New position: %d, term offset: %d", newPosition, termOffset)
		}

		result = readOutcome.FragmentsRead()
	}

	return result
}

func (image Image) Close() error {
	var err error
	if image.isClosed.CompareAndSet(false, true) {
		logger.Debugf("Closing %v", image)
		err = image.logBuffers.Close()
	}
	return err
}
