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
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer/term"
	"github.com/corymonroe-coinbase/aeron-go/aeron/util"
)

type ControlledPollFragmentHandler func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header)

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
	sourceIdentity     string
	logBuffers         *logbuffer.LogBuffers
	exceptionHandler   func(error)
	termBuffers        [logbuffer.PartitionCount]*atomic.Buffer
	subscriberPosition Position
	header             logbuffer.Header
	isClosed           atomic.Bool
	isEos              bool

	termLengthMask             int32
	positionBitsToShift        uint8
	sessionID                  int32
	joinPosition               int64
	finalPosition              int64
	subscriptionRegistrationID int64
	correlationID              int64
}

// NewImage wraps around provided LogBuffers setting up the structures for polling
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
	image.positionBitsToShift = util.NumberOfTrailingZeroes(uint32(capacity))
	image.header.SetInitialTermID(logBuffers.Meta().InitTermID.Get())
	image.header.SetPositionBitsToShift(int32(image.positionBitsToShift))
	image.isClosed.Set(false)

	return image
}

// IsClosed returns whether this image has been closed. No further operations are valid.
func (image *Image) IsClosed() bool {
	return image.isClosed.Get()
}

func (image *Image) Poll(handler term.FragmentHandler, fragmentLimit int) int {
	if image.IsClosed() {
		return 0
	}

	position := image.subscriberPosition.get()
	termOffset := int32(position) & image.termLengthMask
	index := indexByPosition(position, image.positionBitsToShift)
	termBuffer := image.termBuffers[index]

	offset, result := term.Read(termBuffer, termOffset, handler, fragmentLimit, &image.header)

	newPosition := position + int64(offset-termOffset)
	if newPosition > position {
		image.subscriberPosition.set(newPosition)
	}
	return result
}

func (image *Image) PollWithContext(handler term.FragmentHandler, fragmentLimit int) int {
	if image.IsClosed() {
		return 0
	}

	position := image.subscriberPosition.get()
	termOffset := int32(position) & image.termLengthMask
	index := indexByPosition(position, image.positionBitsToShift)
	termBuffer := image.termBuffers[index]

	offset, result := term.ReadWithContext(termBuffer, termOffset, handler, fragmentLimit, &image.header)

	newPosition := position + int64(offset-termOffset)
	if newPosition > position {
		image.subscriberPosition.set(newPosition)
	}
	return result
}

func (image *Image) BoundedPoll(handler term.FragmentHandler, limitPosition int64, fragmentLimit int) int {
	if image.IsClosed() {
		return 0
	}

	position := image.subscriberPosition.get()
	termOffset := int32(position) & image.termLengthMask
	index := indexByPosition(position, image.positionBitsToShift)
	termBuffer := image.termBuffers[index]
	offsetLimit := (limitPosition - position) + int64(termOffset)
	termCapacity := int64(termBuffer.Capacity())
	if offsetLimit > termCapacity {
		offsetLimit = termCapacity
	}

	offset, result := term.BoundedRead(termBuffer, termOffset, int32(offsetLimit), handler, fragmentLimit, &image.header)

	newPosition := position + int64(offset-termOffset)
	if newPosition > position {
		image.subscriberPosition.set(newPosition)
	}
	return result
}

// Position returns the position this Image has been consumed to by the subscriber.
func (image *Image) Position() int64 {
	if image.IsClosed() {
		return image.finalPosition
	}
	return image.subscriberPosition.get()
}

// IsEndOfStream returns if the current consumed position at the end of the stream?
func (image *Image) IsEndOfStream() bool {
	if image.IsClosed() {
		return image.isEos
	}
	return image.subscriberPosition.get() >= image.logBuffers.Meta().EndOfStreamPosOff.Get()
}

// SessionID returns the sessionId for the steam of messages.
func (image *Image) SessionID() int32 {
	return image.sessionID
}

// CorrelationID returns the correlationId for identification of the image with the media driver.
func (image *Image) CorrelationID() int64 {
	return image.correlationID
}

// SubscriptionRegistrationID returns the registrationId for the Subscription of the Image.
func (image *Image) SubscriptionRegistrationID() int64 {
	return image.subscriptionRegistrationID
}

// Close the image and mappings. The image becomes unusable after closing.
func (image *Image) Close() error {
	var err error
	if image.isClosed.CompareAndSet(false, true) {
		image.finalPosition = image.subscriberPosition.get()
		image.isEos = image.finalPosition >=
			image.logBuffers.Meta().EndOfStreamPosOff.Get()
		logger.Debugf("Closing %v", image)
		err = image.logBuffers.Close()
	}
	return err
}

func indexByPosition(position int64, positionBitsToShift uint8) int32 {
	term := uint64(position) >> positionBitsToShift
	return util.FastMod3(term)
}
