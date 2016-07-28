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
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/aeron/util"
	"sync/atomic"
)

const (
	NOT_CONNECTED      int64 = -1
	BACK_PRESSURED     int64 = -2
	ADMIN_ACTION       int64 = -3
	PUBLICATION_CLOSED int64 = -4
)

type Publication struct {
	conductor           *ClientConductor
	logMetaDataBuffer   *buffers.Atomic
	channel             string
	registrationId      int64
	streamId            int32
	sessionId           int32
	initialTermId       int32
	maxPayloadLength    int32
	positionBitsToShift int32
	publicationLimit    buffers.Position

	isClosed atomic.Value

	appenders    [logbuffer.PARTITION_COUNT]*term.Appender
	headerWriter term.HeaderWriter
}

func NewPublication(logBuffers *logbuffer.LogBuffers) *Publication {
	pub := new(Publication)
	pub.logMetaDataBuffer = logBuffers.Buffer(logbuffer.Descriptor.LOG_META_DATA_SECTION_INDEX)
	pub.initialTermId = logbuffer.InitialTermId(pub.logMetaDataBuffer)
	pub.maxPayloadLength = logbuffer.MtuLength(pub.logMetaDataBuffer) - logbuffer.DataFrameHeader.LENGTH
	pub.positionBitsToShift = int32(util.NumberOfTrailingZeroes(logBuffers.Buffer(0).Capacity()))
	header := logbuffer.DefaultFrameHeader(pub.logMetaDataBuffer)
	pub.headerWriter.Fill(header)
	pub.isClosed.Store(false)

	for i := 0; i < logbuffer.PARTITION_COUNT; i++ {
		appender := term.MakeAppender(logBuffers.Buffer(i), pub.logMetaDataBuffer, i)
		logger.Debugf("TermAppender[%d]: %v", i, appender)
		pub.appenders[i] = appender
	}

	return pub
}

func (pub *Publication) IsClosed() bool {
	return pub.isClosed.Load().(bool)
}

func (pub *Publication) Close() error {
	if pub != nil && !pub.IsClosed() {
		pub.isClosed.Store(true)
		<-pub.conductor.releasePublication(pub.registrationId)
	}

	return nil
}

func (pub *Publication) checkForMaxMessageLength(length int32) {
	if length > pub.maxPayloadLength {
		panic(fmt.Sprintf("Encoded message exceeds maxMessageLength of %d, length=%d", pub.maxPayloadLength, length))
	}
}

func (pub *Publication) IsConnected() bool {
	return !pub.IsClosed() && pub.conductor.isPublicationConnected(logbuffer.TimeOfLastStatusMessage(pub.logMetaDataBuffer))
}

func (pub *Publication) newPosition(index int32, currentTail int32, position int64, result *term.AppenderResult) int64 {
	newPosition := ADMIN_ACTION

	if result.TermOffset() > 0 {
		newPosition = (position - int64(currentTail)) + result.TermOffset()
	} else if result.TermOffset() == int64(term.APPENDER_TRIPPED) {
		nextIndex := logbuffer.NextPartitionIndex(index)

		pub.appenders[nextIndex].SetTailTermId(result.TermId() + 1)
		logbuffer.SetActivePartitionIndex(pub.logMetaDataBuffer, nextIndex)
	}

	return newPosition
}

func (pub *Publication) Offer(buffer *buffers.Atomic, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {

	newPosition := PUBLICATION_CLOSED

	if !pub.IsClosed() {
		limit := pub.publicationLimit.Get()
		partitionIndex := logbuffer.ActivePartitionIndex(pub.logMetaDataBuffer)
		termAppender := pub.appenders[partitionIndex]
		rawTail := termAppender.RawTail()
		termOffset := rawTail & 0xFFFFFFFF
		position := logbuffer.ComputeTermBeginPosition(logbuffer.TermId(rawTail), pub.positionBitsToShift, pub.initialTermId) + termOffset

		logger.Debugf("Offering at %d of %d (pubLmt: %v)", position, limit, pub.publicationLimit)
		if position < limit {
			var appendResult term.AppenderResult
			var resValSupplier term.ReservedValueSupplier = term.DEFAULT_RESERVED_VALUE_SUPPLIER
			if nil != reservedValueSupplier {
				resValSupplier = reservedValueSupplier
			}
			if length <= pub.maxPayloadLength {
				termAppender.AppendUnfragmentedMessage(&appendResult, &pub.headerWriter, buffer, offset, length, resValSupplier)
			} else {
				pub.checkForMaxMessageLength(length)
				termAppender.AppendFragmentedMessage(&appendResult, &pub.headerWriter, buffer, offset, length, pub.maxPayloadLength, resValSupplier)
			}

			newPosition = pub.newPosition(partitionIndex, int32(termOffset), position, &appendResult)
			logger.Debugf("publication new position: %d, term offset: %d", newPosition, termOffset)
		} else if pub.conductor.isPublicationConnected(logbuffer.TimeOfLastStatusMessage(pub.logMetaDataBuffer)) {
			newPosition = BACK_PRESSURED
		} else {
			newPosition = NOT_CONNECTED
		}
	}

	return newPosition
}
