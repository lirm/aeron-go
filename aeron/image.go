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

//go:norace
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

// BoundedPoll polls for new messages in a stream. If new messages are found
// beyond the last consumed position then they will be delivered to the
// FragmentHandler up to a limited number of fragments as specified or the
// maximum position specified. Use a FragmentAssembler to assemble messages
// which span multiple fragments. Returns the number of fragments that have been
// consumed.
func (image *Image) BoundedPoll(
	handler term.FragmentHandler,
	limitPosition int64,
	fragmentLimit int,
) int {
	if image.IsClosed() {
		return 0
	}

	fragmentsRead := 0
	initialPosition := image.subscriberPosition.get()
	initialOffset := int32(initialPosition) & image.termLengthMask
	offset := initialOffset

	index := indexByPosition(initialPosition, image.positionBitsToShift)
	termBuffer := image.termBuffers[index]

	capacity := termBuffer.Capacity()
	limitOffset := int32(limitPosition-initialPosition) + offset
	if limitOffset > capacity {
		limitOffset = capacity
	}
	header := &image.header
	header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())

	for fragmentsRead < fragmentLimit && offset < limitOffset {
		length := logbuffer.GetFrameLength(termBuffer, offset)
		if length <= 0 {
			break
		}

		frameOffset := offset
		alignedLength := util.AlignInt32(length, logbuffer.FrameAlignment)
		offset += alignedLength

		if logbuffer.IsPaddingFrame(termBuffer, frameOffset) {
			continue
		}
		fragmentsRead++
		header.SetOffset(frameOffset)

		handler(termBuffer, frameOffset+logbuffer.DataFrameHeader.Length,
			length-logbuffer.DataFrameHeader.Length, header)
	}
	resultingPosition := initialPosition + int64(offset-initialOffset)
	if resultingPosition > initialPosition {
		image.subscriberPosition.set(resultingPosition)
	}
	return fragmentsRead
}

// ControlledPoll polls for new messages in a stream. If new messages are found
// beyond the last consumed position then they will be delivered to the
// ControlledFragmentHandler up to a limited number of fragments as
// specified.
//
// To assemble messages that span multiple fragments then use
// ControlledFragmentAssembler. Returns the number of fragments that have been
// consumed.
func (image *Image) ControlledPoll(
	handler term.ControlledFragmentHandler,
	fragmentLimit int,
) int {
	if image.IsClosed() {
		return 0
	}

	fragmentsRead := 0
	initialPosition := image.subscriberPosition.get()
	initialOffset := int32(initialPosition) & image.termLengthMask
	offset := initialOffset

	index := indexByPosition(initialPosition, image.positionBitsToShift)
	termBuffer := image.termBuffers[index]

	capacity := termBuffer.Capacity()
	header := &image.header
	header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())

	for fragmentsRead < fragmentLimit && offset < capacity {
		length := logbuffer.GetFrameLength(termBuffer, offset)
		if length <= 0 {
			break
		}

		frameOffset := offset
		alignedLength := util.AlignInt32(length, logbuffer.FrameAlignment)
		offset += alignedLength

		if logbuffer.IsPaddingFrame(termBuffer, frameOffset) {
			continue
		}
		fragmentsRead++
		header.SetOffset(frameOffset)

		action := handler(termBuffer, frameOffset+logbuffer.DataFrameHeader.Length,
			length-logbuffer.DataFrameHeader.Length, header)
		if action == term.ControlledPollActionAbort {
			fragmentsRead--
			offset -= alignedLength
			break
		}
		if action == term.ControlledPollActionBreak {
			break
		}
		if action == term.ControlledPollActionCommit {
			initialPosition += int64(offset - initialOffset)
			initialOffset = offset
			image.subscriberPosition.set(initialPosition)
		}
	}
	resultingPosition := initialPosition + int64(offset-initialOffset)
	if resultingPosition > initialPosition {
		image.subscriberPosition.set(resultingPosition)
	}
	return fragmentsRead
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

// TermBufferLength returns the length in bytes for each term partition in the log buffer.
func (image *Image) TermBufferLength() int32 {
	return image.termLengthMask + 1
}

// ActiveTransportCount returns the number of observed active
// transports within the image liveness timeout.
//
// Returns 0 if the image is closed, if no datagrams have arrived or the image is IPC
func (image *Image) ActiveTransportCount() int32 {
	return image.logBuffers.Meta().ActiveTransportCount()
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
