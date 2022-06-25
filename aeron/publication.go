/*
Copyright 2016-2018 Stanislav Liberman

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
	"fmt"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/lirm/aeron-go/aeron/util"
)

const (
	// NotConnected indicates that this Publication is not connected to the driver
	NotConnected int64 = -1
	// BackPressured indicates that sending ring buffer is full
	BackPressured int64 = -2
	// AdminAction indicates that terms needs to be rotated. User should retry the Offer
	AdminAction int64 = -3
	// PublicationClosed indicates that this Publication is closed and no further Offers shall succeed
	PublicationClosed int64 = -4
	// MaxPositionExceeded indicates that ...
	MaxPositionExceeded int64 = -5
)

// Publication is a sender structure
type Publication struct {
	conductor                *ClientConductor
	channel                  string
	regID                    int64
	originalRegID            int64
	maxPossiblePosition      int64
	streamID                 int32
	sessionID                int32
	initialTermID            int32
	maxPayloadLength         int32
	maxMessageLength         int32
	positionBitsToShift      int32
	pubLimit                 Position
	channelStatusIndicatorID int32

	isClosed atomic.Bool
	metaData *logbuffer.LogBufferMetaData

	appenders [logbuffer.PartitionCount]*term.Appender
}

// NewPublication is a factory method create new publications
func NewPublication(logBuffers *logbuffer.LogBuffers) *Publication {
	termBufferCapacity := logBuffers.Buffer(0).Capacity()

	pub := new(Publication)
	pub.metaData = logBuffers.Meta()
	pub.initialTermID = pub.metaData.InitTermID.Get()
	pub.maxPayloadLength = pub.metaData.MTULen.Get() - logbuffer.DataFrameHeader.Length
	pub.maxMessageLength = logbuffer.ComputeMaxMessageLength(termBufferCapacity)
	pub.positionBitsToShift = int32(util.NumberOfTrailingZeroes(uint32(termBufferCapacity)))
	pub.maxPossiblePosition = int64(termBufferCapacity) * (1 << 31)

	pub.isClosed.Set(false)

	for i := 0; i < logbuffer.PartitionCount; i++ {
		appender := term.MakeAppender(logBuffers, i)
		logger.Debugf("TermAppender[%d]: %v", i, appender)
		pub.appenders[i] = appender
	}

	return pub
}

// ChannelStatusID returns the counter used to represent the channel status
// for this publication.
func (pub *Publication) ChannelStatusID() int32 {
	return pub.channelStatusIndicatorID
}

// RegistrationID returns the registration id.
func (pub *Publication) RegistrationID() int64 {
	return pub.regID
}

// OriginalRegistrationID returns the original registration id.
func (pub *Publication) OriginalRegistrationID() int64 {
	return pub.originalRegID
}

// Channel returns the media address for delivery to the channel.
func (pub *Publication) Channel() string {
	return pub.channel
}

// StreamID returns Stream identity for scoping within the channel media address.
func (pub *Publication) StreamID() int32 {
	return pub.streamID
}

// SessionID returns the session id for this publication.
func (pub *Publication) SessionID() int32 {
	return pub.sessionID
}

// InitialTermID returns the initial term id assigned when this publication was
// created. This can be used to determine how many terms have passed since
// creation.
func (pub *Publication) InitialTermID() int32 {
	return pub.initialTermID
}

// IsConnected returns whether this publication is connected to the driver (not whether it has any Subscriptions)
func (pub *Publication) IsConnected() bool {
	return !pub.IsClosed() && pub.metaData.IsConnected.Get() == 1
}

// IsClosed returns whether this Publication has been closed
func (pub *Publication) IsClosed() bool {
	return pub.isClosed.Get()
}

// IsOriginal return true if this instance is the first added otherwise false.
func (pub *Publication) IsOriginal() bool {
	return pub.originalRegID == pub.regID
}

// Close will close this publication with the driver. This is a blocking call.
func (pub *Publication) Close() error {
	// FIXME Why can pub be nil?!
	if pub != nil && pub.isClosed.CompareAndSet(false, true) {
		pub.conductor.releasePublication(pub.regID)
	}

	return nil
}

// Position returns the current position to which the publication has advanced
// for this stream or PublicationClosed if closed.
func (pub *Publication) Position() int64 {
	if pub.IsClosed() {
		return PublicationClosed
	}

	// Spelled out for clarity, the compiler will optimize
	termCount := pub.metaData.ActiveTermCountOff.Get()
	termIndex := termCount % logbuffer.PartitionCount
	termAppender := pub.appenders[termIndex]
	rawTail := termAppender.RawTail()
	termOffset := rawTail & 0xFFFFFFFF
	termId := logbuffer.TermID(rawTail)

	return computeTermBeginPosition(termId, pub.positionBitsToShift, pub.initialTermID) + termOffset
}

// Offer is the primary send mechanism on Publication
func (pub *Publication) Offer(buffer *atomic.Buffer, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {

	newPosition := PublicationClosed

	if reservedValueSupplier == nil {
		reservedValueSupplier = term.DefaultReservedValueSupplier
	}

	if !pub.IsClosed() {

		limit := pub.pubLimit.get()
		termCount := pub.metaData.ActiveTermCountOff.Get()
		termIndex := termCount % logbuffer.PartitionCount
		termAppender := pub.appenders[termIndex]
		rawTail := termAppender.RawTail()
		termOffset := rawTail & 0xFFFFFFFF
		termId := logbuffer.TermID(rawTail)
		position := computeTermBeginPosition(termId, pub.positionBitsToShift, pub.initialTermID) + termOffset

		if termCount != (termId - pub.metaData.InitTermID.Get()) {
			return AdminAction
		}

		if logger.IsEnabledFor(logging.DEBUG) {
			logger.Debugf("Offering at %d of %d (pubLmt: %v)", position, limit, pub.pubLimit)
		}
		if position < limit {
			var termOffsetA int64
			var termId int32
			if length <= pub.maxPayloadLength {
				termOffsetA, termId = termAppender.AppendUnfragmentedMessage(buffer, offset, length, reservedValueSupplier)
			} else {
				pub.checkForMaxMessageLength(length)
				termOffsetA, termId = termAppender.AppendFragmentedMessage(buffer, offset, length, pub.maxPayloadLength, reservedValueSupplier)
			}

			newPosition = pub.newPosition(termCount, termOffset, termId, position, termOffsetA)

		} else {
			newPosition = pub.backPressureStatus(position, length)
		}
	}

	return newPosition
}

// Offer2 attempts to publish a message composed of two parts, e.g. a header and encapsulated payload.
func (pub *Publication) Offer2(
	bufferOne *atomic.Buffer, offsetOne int32, lengthOne int32,
	bufferTwo *atomic.Buffer, offsetTwo int32, lengthTwo int32,
	reservedValueSupplier term.ReservedValueSupplier,
) int64 {
	if lengthOne < 0 {
		logger.Debugf("Offered negative length (lengthOne: %d)", lengthOne)
		return 0
	} else if lengthTwo < 0 {
		logger.Debugf("Offered negative length (lengthTwo: %d)", lengthTwo)
		return 0
	}
	length := lengthOne + lengthTwo
	if length < 0 {
		panic(fmt.Sprintf("Length overflow (lengthOne: %d lengthTwo: %d)", lengthOne, lengthTwo))
	}
	newPosition := PublicationClosed

	if reservedValueSupplier == nil {
		reservedValueSupplier = term.DefaultReservedValueSupplier
	}

	if !pub.IsClosed() {

		limit := pub.pubLimit.get()
		termCount := pub.metaData.ActiveTermCountOff.Get()
		termIndex := termCount % logbuffer.PartitionCount
		termAppender := pub.appenders[termIndex]
		rawTail := termAppender.RawTail()
		termOffset := rawTail & 0xFFFFFFFF
		termId := logbuffer.TermID(rawTail)
		position := computeTermBeginPosition(termId, pub.positionBitsToShift, pub.initialTermID) + termOffset

		if termCount != (termId - pub.metaData.InitTermID.Get()) {
			return AdminAction
		}

		if logger.IsEnabledFor(logging.DEBUG) {
			logger.Debugf("Offering at %d of %d (pubLmt: %v)", position, limit, pub.pubLimit)
		}
		if position < limit {
			var termOffsetA int64
			var termId int32
			if length <= pub.maxPayloadLength {
				termOffsetA, termId = termAppender.AppendUnfragmentedMessage2(
					bufferOne, offsetOne, lengthOne,
					bufferTwo, offsetTwo, lengthTwo,
					reservedValueSupplier)
			} else {
				pub.checkForMaxMessageLength(length)
				termOffsetA, termId = termAppender.AppendFragmentedMessage2(
					bufferOne, offsetOne, lengthOne,
					bufferTwo, offsetTwo, lengthTwo,
					pub.maxPayloadLength, reservedValueSupplier)
			}
			newPosition = pub.newPosition(termCount, termOffset, termId, position, termOffsetA)
		} else {
			newPosition = pub.backPressureStatus(position, length)
		}
	}

	return newPosition
}

func (pub *Publication) newPosition(termCount int32, termOffset int64, termId int32, position int64, resultingOffset int64) int64 {
	if resultingOffset > 0 {
		return (position - termOffset) + resultingOffset
	}

	if (position + termOffset) > pub.maxPossiblePosition {
		return MaxPositionExceeded
	}

	logbuffer.RotateLog(pub.metaData, termCount, termId)
	return AdminAction
}

func (pub *Publication) backPressureStatus(currentPosition int64, messageLength int32) int64 {

	if (currentPosition + int64(messageLength)) >= pub.maxPossiblePosition {
		return MaxPositionExceeded
	}

	if pub.metaData.IsConnected.Get() == 1 {
		return BackPressured
	}

	return NotConnected
}

func (pub *Publication) TryClaim(length int32, bufferClaim *logbuffer.Claim) int64 {
	newPosition := PublicationClosed

	if !pub.IsClosed() {
		pub.checkForMaxPayloadLength(length)

		limit := pub.pubLimit.get()
		termCount := pub.metaData.ActiveTermCountOff.Get()
		termAppender := pub.appenders[termCount]
		rawTail := termAppender.RawTail()
		termOffset := rawTail & 0xFFFFFFFF
		position := computeTermBeginPosition(logbuffer.TermID(rawTail), pub.positionBitsToShift, pub.initialTermID) + termOffset

		if position < limit {
			resultingOffset, termId := termAppender.Claim(length, bufferClaim)

			newPosition = pub.newPosition(termCount, termOffset, termId, position, resultingOffset)
		} else {
			newPosition = pub.backPressureStatus(position, length)
		}

	}
	return newPosition
}

func (pub *Publication) checkForMaxMessageLength(length int32) {
	if length > pub.maxMessageLength {
		panic(fmt.Sprintf("Encoded message exceeds maxMessageLength of %d, length=%d", pub.maxMessageLength, length))
	}
}

func (pub *Publication) checkForMaxPayloadLength(length int32) {
	if length > pub.maxPayloadLength {
		panic(fmt.Sprintf("Encoded message exceeds maxPayloadLength of %d, length=%d", pub.maxPayloadLength, length))
	}
}

func computeTermBeginPosition(activeTermID, positionBitsToShift, initialTermID int32) int64 {
	termCount := int64(activeTermID - initialTermID)

	return termCount << uint32(positionBitsToShift)
}

func nextPartitionIndex(currentIndex int32) int32 {
	return util.FastMod3(uint64(currentIndex) + 1)
}
