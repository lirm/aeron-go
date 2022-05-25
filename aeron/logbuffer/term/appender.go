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

package term

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/flyweight"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
)

const (
	// AppenderTripped is returned when the end of the term has been reached and buffer roll was done
	AppenderTripped int64 = -1

	// AppenderFailed is returned when appending is not possible due to position being outside of the term. ??
	AppenderFailed int64 = -2

	beginFrag    uint8 = 0x80
	endFrag      uint8 = 0x40
	unfragmented uint8 = 0x80 | 0x40
)

// DefaultReservedValueSupplier is the default reserved value provider
var DefaultReservedValueSupplier ReservedValueSupplier = func(termBuffer *atomic.Buffer, termOffset int32, length int32) int64 {
	return 0
}

// ReservedValueSupplier is the type definition for a provider of user supplied header data
type ReservedValueSupplier func(termBuffer *atomic.Buffer, termOffset int32, length int32) int64

// HeaderWriter is a helper class for writing frame header to the term
type headerWriter struct {
	sessionID int32
	streamID  int32
}

func (header *headerWriter) fill(defaultHdr *atomic.Buffer) {
	header.sessionID = defaultHdr.GetInt32(logbuffer.DataFrameHeader.SessionIDFieldOffset)
	header.streamID = defaultHdr.GetInt32(logbuffer.DataFrameHeader.StreamIDFieldOffset)
}

func (header *headerWriter) write(termBuffer *atomic.Buffer, offset, length, termID int32) {
	termBuffer.PutInt32Ordered(offset, -length)

	termBuffer.PutInt8(offset+logbuffer.DataFrameHeader.VersionFieldOffset, logbuffer.DataFrameHeader.CurrentVersion)
	termBuffer.PutUInt8(offset+logbuffer.DataFrameHeader.FlagsFieldOffset, unfragmented)
	termBuffer.PutUInt16(offset+logbuffer.DataFrameHeader.TypeFieldOffset, logbuffer.DataFrameHeader.TypeData)
	termBuffer.PutInt32(offset+logbuffer.DataFrameHeader.TermOffsetFieldOffset, offset)
	termBuffer.PutInt32(offset+logbuffer.DataFrameHeader.SessionIDFieldOffset, header.sessionID)
	termBuffer.PutInt32(offset+logbuffer.DataFrameHeader.StreamIDFieldOffset, header.streamID)
	termBuffer.PutInt32(offset+logbuffer.DataFrameHeader.TermIDFieldOffset, termID)
}

// Appender type is the term writer
type Appender struct {
	termBuffer   *atomic.Buffer
	tailCounter  flyweight.Int64Field
	headerWriter headerWriter
}

// MakeAppender is the factory function for term Appenders
func MakeAppender(logBuffers *logbuffer.LogBuffers, partitionIndex int) *Appender {

	appender := new(Appender)
	appender.termBuffer = logBuffers.Buffer(partitionIndex)
	appender.tailCounter = logBuffers.Meta().TailCounter[partitionIndex]

	header := logBuffers.Meta().DefaultFrameHeader.Get()
	appender.headerWriter.fill(header)

	return appender
}

// RawTail is the accessor to the raw value of the tail offset used by Publication
func (appender *Appender) RawTail() int64 {
	return appender.tailCounter.Get()
}

// SetRawTail sets the raw value of the tail. It should not be used outside of testing
func (appender *Appender) SetRawTail(v int64) {
	appender.tailCounter.Set(v)
}

func (appender *Appender) getAndAddRawTail(alignedLength int32) int64 {
	return appender.tailCounter.GetAndAddInt64(int64(alignedLength))
}

// Claim is the interface for using Buffer Claims for zero copy sends
func (appender *Appender) Claim(length int32, claim *logbuffer.Claim) (termOffset int64, termID int32) {

	frameLength := length + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termOffset = rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset += int64(alignedLength)
	if termOffset > int64(termLength) {
		termOffset = handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength)
	} else {
		offset := int32(termOffset)
		appender.headerWriter.write(appender.termBuffer, offset, frameLength, termID)
		claim.Wrap(appender.termBuffer, offset, frameLength)
	}

	return termOffset, termID
}

// AppendUnfragmentedMessage appends an unfragmented message in a single frame to the term
func (appender *Appender) AppendUnfragmentedMessage(srcBuffer *atomic.Buffer, srcOffset int32, length int32,
	reservedValueSupplier ReservedValueSupplier) (resultingOffset int64, termID int32) {

	frameLength := length + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset := rawTail & 0xFFFFFFFF
	resultingOffset = termOffset + int64(alignedLength)
	if resultingOffset > int64(termLength) {
		resultingOffset = handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength)
	} else {
		offset := int32(termOffset)
		appender.headerWriter.write(appender.termBuffer, offset, frameLength, logbuffer.TermID(rawTail))
		appender.termBuffer.PutBytes(offset+logbuffer.DataFrameHeader.Length, srcBuffer, srcOffset, length)

		if nil != reservedValueSupplier {
			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)
		}

		logbuffer.SetFrameLength(appender.termBuffer, offset, frameLength)
	}

	return resultingOffset, termID
}

// AppendUnfragmentedMessage2 appends the given pair of buffers as an unfragmented message in a single frame to the term
func (appender *Appender) AppendUnfragmentedMessage2(
	srcBufferOne *atomic.Buffer, srcOffsetOne int32, lengthOne int32,
	srcBufferTwo *atomic.Buffer, srcOffsetTwo int32, lengthTwo int32,
	reservedValueSupplier ReservedValueSupplier,
) (resultingOffset int64, termID int32) {

	frameLength := lengthOne + lengthTwo + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset := rawTail & 0xFFFFFFFF
	resultingOffset = termOffset + int64(alignedLength)
	if resultingOffset > int64(termLength) {
		resultingOffset = handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength)
	} else {
		offset := int32(termOffset)
		dataOffset := offset + logbuffer.DataFrameHeader.Length
		appender.headerWriter.write(appender.termBuffer, offset, frameLength, logbuffer.TermID(rawTail))
		appender.termBuffer.PutBytes(dataOffset, srcBufferOne, srcOffsetOne, lengthOne)
		appender.termBuffer.PutBytes(dataOffset+lengthOne, srcBufferTwo, srcOffsetTwo, lengthTwo)

		if nil != reservedValueSupplier {
			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)
		}

		logbuffer.SetFrameLength(appender.termBuffer, offset, frameLength)
	}

	return resultingOffset, termID
}

// AppendFragmentedMessage appends a message greater than frame length as a batch of fragments
func (appender *Appender) AppendFragmentedMessage(srcBuffer *atomic.Buffer, srcOffset int32, length int32,
	maxPayloadLength int32, reservedValueSupplier ReservedValueSupplier) (resultingOffset int64, termID int32) {

	numMaxPayloads := length / maxPayloadLength
	remainingPayload := length % maxPayloadLength
	var lastFrameLength int32
	if remainingPayload > 0 {
		lastFrameLength = util.AlignInt32(remainingPayload+logbuffer.DataFrameHeader.Length, logbuffer.FrameAlignment)
	}
	requiredLength := (numMaxPayloads * (maxPayloadLength + logbuffer.DataFrameHeader.Length)) + lastFrameLength
	rawTail := appender.getAndAddRawTail(requiredLength)

	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset := rawTail & 0xFFFFFFFF
	resultingOffset = termOffset + int64(requiredLength)
	if resultingOffset > int64(termLength) {
		resultingOffset = handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength)
	} else {
		flags := beginFrag
		remaining := length
		offset := int32(termOffset)

		for remaining > 0 {
			bytesToWrite := minInt32(remaining, maxPayloadLength)
			frameLength := bytesToWrite + logbuffer.DataFrameHeader.Length
			alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)

			appender.headerWriter.write(appender.termBuffer, offset, frameLength, termID)
			appender.termBuffer.PutBytes(
				offset+logbuffer.DataFrameHeader.Length, srcBuffer, srcOffset+(length-remaining), bytesToWrite)

			if remaining <= maxPayloadLength {
				flags |= endFrag
			}

			logbuffer.FrameFlags(appender.termBuffer, offset, flags)

			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)

			logbuffer.SetFrameLength(appender.termBuffer, offset, frameLength)

			flags = 0
			offset += alignedLength
			remaining -= bytesToWrite
		}
	}

	return resultingOffset, termID
}

// AppendFragmentedMessage2 appends the given pair of buffers (with combined length greater than max frame length)
//  as a batch of fragments
func (appender *Appender) AppendFragmentedMessage2(
	srcBufferOne *atomic.Buffer, srcOffsetOne int32, lengthOne int32,
	srcBufferTwo *atomic.Buffer, srcOffsetTwo int32, lengthTwo int32,
	maxPayloadLength int32, reservedValueSupplier ReservedValueSupplier,
) (resultingOffset int64, termID int32) {
	length := lengthOne + lengthTwo
	numMaxPayloads := length / maxPayloadLength
	remainingPayload := length % maxPayloadLength
	var lastFrameLength int32
	if remainingPayload > 0 {
		lastFrameLength = util.AlignInt32(remainingPayload+logbuffer.DataFrameHeader.Length, logbuffer.FrameAlignment)
	}
	requiredLength := (numMaxPayloads * (maxPayloadLength + logbuffer.DataFrameHeader.Length)) + lastFrameLength
	rawTail := appender.getAndAddRawTail(requiredLength)

	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset := rawTail & 0xFFFFFFFF
	resultingOffset = termOffset + int64(requiredLength)
	if resultingOffset > int64(termLength) {
		resultingOffset = handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength)
	} else {
		flags := beginFrag
		remaining := length
		frameOffset := int32(termOffset)
		var posOne, posTwo int32

		for remaining > 0 {
			bytesToWrite := minInt32(remaining, maxPayloadLength)
			frameLength := bytesToWrite + logbuffer.DataFrameHeader.Length
			alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)

			appender.headerWriter.write(appender.termBuffer, frameOffset, frameLength, termID)

			var bytesWritten int32
			payloadOffset := frameOffset + logbuffer.DataFrameHeader.Length
			for bytesWritten < bytesToWrite {
				remainingOne := lengthOne - posOne
				if remainingOne > 0 {
					numBytes := minInt32(bytesToWrite-bytesWritten, remainingOne)
					appender.termBuffer.PutBytes(payloadOffset, srcBufferOne, srcOffsetOne+posOne, numBytes)
					bytesWritten += numBytes
					payloadOffset += numBytes
					posOne += numBytes
				} else {
					numBytes := minInt32(bytesToWrite-bytesWritten, lengthTwo-posTwo)
					appender.termBuffer.PutBytes(payloadOffset, srcBufferTwo, srcOffsetTwo+posTwo, numBytes)
					bytesWritten += numBytes
					payloadOffset += numBytes
					posTwo += numBytes
				}
			}

			if remaining <= maxPayloadLength {
				flags |= endFrag
			}
			logbuffer.FrameFlags(appender.termBuffer, frameOffset, flags)

			reservedValue := reservedValueSupplier(appender.termBuffer, frameOffset, frameLength)
			appender.termBuffer.PutInt64(frameOffset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)

			logbuffer.SetFrameLength(appender.termBuffer, frameOffset, frameLength)

			flags = 0
			frameOffset += alignedLength
			remaining -= bytesToWrite
		}
	}

	return resultingOffset, termID
}

func handleEndOfLogCondition(termID int32, termBuffer *atomic.Buffer, termOffset int32, header *headerWriter,
	termLength int32) int64 {
	newOffset := AppenderFailed

	if termOffset <= termLength {
		newOffset = AppenderTripped

		if termOffset < termLength {
			paddingLength := termLength - termOffset
			header.write(termBuffer, termOffset, paddingLength, termID)
			logbuffer.SetFrameType(termBuffer, termOffset, logbuffer.DataFrameHeader.TypePad)
			logbuffer.SetFrameLength(termBuffer, termOffset, paddingLength)
		}
	}

	return newOffset
}

func (appender *Appender) SetTailTermID(termID int32) {
	appender.tailCounter.Set(int64(termID) << 32)
}

func minInt32(v1, v2 int32) int32 {
	if v1 < v2 {
		return v1
	} else {
		return v2
	}
}
