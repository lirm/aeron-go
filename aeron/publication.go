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
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/op/go-logging"
)

const (
	// NotConnected indicates that this Publication is not connected to the driver
	NotConnected int64 = -1
	// BackPressured indicates that sending ring buffer is full
	BackPressured int64 = -2
	// AdminAction indicates that terms needs to be rotated. User should retry the Offer
	AdminAction int64 = -3
	// PublicationClosed indicates that this Publication is closed an no further Offers shall succeed
	PublicationClosed int64 = -4
)

// Publication is a sender structure
type Publication struct {
	conductor           *ClientConductor
	channel             string
	registrationID      int64
	streamID            int32
	sessionID           int32
	initialTermID       int32
	maxPayloadLength    int32
	maxMessageLength    int32
	positionBitsToShift int32
	publicationLimit    Position

	isClosed atomic.Bool
	metaData *logbuffer.LogBufferMetaData

	appenders [logbuffer.PartitionCount]*term.Appender
}

// NewPublication is a factory method create new publications
func NewPublication(logBuffers *logbuffer.LogBuffers) *Publication {
	pub := new(Publication)
	pub.metaData = logBuffers.Meta()
	pub.initialTermID = pub.metaData.InitTermID.Get()
	pub.maxPayloadLength = pub.metaData.MTULen.Get() - logbuffer.DataFrameHeader.Length
	pub.maxMessageLength = logbuffer.ComputeMaxMessageLength(logBuffers.Buffer(0).Capacity())
	pub.positionBitsToShift = int32(util.NumberOfTrailingZeroes(logBuffers.Buffer(0).Capacity()))

	pub.isClosed.Set(false)

	for i := 0; i < logbuffer.PartitionCount; i++ {
		appender := term.MakeAppender(logBuffers, i)
		logger.Debugf("TermAppender[%d]: %v", i, appender)
		pub.appenders[i] = appender
	}

	return pub
}

// IsConnected returns whether this publication is connected to the driver (not whether it has any Subscriptions)
func (pub *Publication) IsConnected() bool {
	return !pub.IsClosed() && pub.conductor.isPublicationConnected(pub.metaData.TimeOfLastStatusMsg.Get())
}

// IsClosed returns whether this Publication has been closed
func (pub *Publication) IsClosed() bool {
	return pub.isClosed.Get()
}

// Close will close this publication with the driver. This is a blocking call.
func (pub *Publication) Close() error {
	if pub.isClosed.CompareAndSet(false, true) {
		<-pub.conductor.releasePublication(pub.registrationID)
	}

	return nil
}

// Offer is the primary send mechanism on Publication
func (pub *Publication) Offer(buffer *atomic.Buffer, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {

	newPosition := PublicationClosed

	if !pub.IsClosed() {
		limit := pub.publicationLimit.get()
		partitionIndex := pub.metaData.ActivePartitionIx.Get()
		termAppender := pub.appenders[partitionIndex]
		rawTail := termAppender.RawTail()
		termOffset := rawTail & 0xFFFFFFFF
		position := computeTermBeginPosition(logbuffer.TermID(rawTail), pub.positionBitsToShift, pub.initialTermID) + termOffset

		if logger.IsEnabledFor(logging.DEBUG) {
			logger.Debugf("Offering at %d of %d (pubLmt: %v)", position, limit, pub.publicationLimit)
		}
		if position < limit {
			var appendResult term.AppenderResult
			resValSupplier := term.DefaultReservedValueSupplier
			if nil != reservedValueSupplier {
				resValSupplier = reservedValueSupplier
			}
			if length <= pub.maxPayloadLength {
				termAppender.AppendUnfragmentedMessage(&appendResult, buffer, offset, length, resValSupplier)
			} else {
				pub.checkForMaxMessageLength(length)
				termAppender.AppendFragmentedMessage(&appendResult, buffer, offset, length, pub.maxPayloadLength, resValSupplier)
			}

			if appendResult.TermOffset() > 0 {
				newPosition = (position - termOffset) + appendResult.TermOffset()
			} else {
				newPosition = AdminAction
				if appendResult.TermOffset() == term.AppenderTripped {
					nextIndex := nextPartitionIndex(partitionIndex)

					pub.appenders[nextIndex].SetTailTermID(appendResult.TermID() + 1)
					pub.metaData.ActivePartitionIx.Set(nextIndex)
				}
			}
		} else if pub.conductor.isPublicationConnected(pub.metaData.TimeOfLastStatusMsg.Get()) {
			newPosition = BackPressured
		} else {
			newPosition = NotConnected
		}
	}

	return newPosition
}

func (pub *Publication) checkForMaxMessageLength(length int32) {
	if length > pub.maxMessageLength {
		panic(fmt.Sprintf("Encoded message exceeds maxMessageLength of %d, length=%d", pub.maxPayloadLength, length))
	}
}

func computeTermBeginPosition(activeTermID, positionBitsToShift, initialTermID int32) int64 {
	termCount := int64(activeTermID - initialTermID)

	return termCount << uint32(positionBitsToShift)
}

func nextPartitionIndex(currentIndex int32) int32 {
	return util.FastMod3(uint64(currentIndex) + 1)
}
