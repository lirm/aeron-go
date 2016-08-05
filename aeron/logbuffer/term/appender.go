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
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
	"math"
	"unsafe"
)

const (
	APPENDER_TRIPPED int32 = -1
	APPENDER_FAILED  int32 = -2
)

var DEFAULT_RESERVED_VALUE_SUPPLIER ReservedValueSupplier = func(termBuffer *atomic.Buffer, termOffset int32, length int32) int64 { return 0 }

type ReservedValueSupplier func(termBuffer *atomic.Buffer, termOffset int32, length int32) int64

type HeaderWriter struct {
	sessionId int32
	streamId  int32
	buffer atomic.Buffer
}

func (header *HeaderWriter) Fill(defaultHdr *atomic.Buffer) {
	header.sessionId = defaultHdr.GetInt32(logbuffer.DataFrameHeader.SessionIDFieldOffset)
	header.streamId = defaultHdr.GetInt32(logbuffer.DataFrameHeader.StreamIDFieldOffset)
}

func (header *HeaderWriter) write(termBuffer *atomic.Buffer, offset, length, termId int32) {
	termBuffer.PutInt32Ordered(offset, -length)

	headerPtr := uintptr(termBuffer.Ptr()) + uintptr(offset)
	header.buffer.Wrap(unsafe.Pointer(headerPtr), logbuffer.DataFrameHeader.Length)

	header.buffer.PutInt8(logbuffer.DataFrameHeader.VersionFieldOffset, logbuffer.DataFrameHeader.CurrentVersion)
	header.buffer.PutUInt8(logbuffer.DataFrameHeader.FlagsFieldOffset, logbuffer.FrameDescriptor.BeginFrag|logbuffer.FrameDescriptor.EndFrag)
	header.buffer.PutUInt16(logbuffer.DataFrameHeader.TypeFieldOffset, logbuffer.DataFrameHeader.TypeData)
	header.buffer.PutInt32(logbuffer.DataFrameHeader.TermOffsetFieldOffset, offset)
	header.buffer.PutInt32(logbuffer.DataFrameHeader.SessionIDFieldOffset, header.sessionId)
	header.buffer.PutInt32(logbuffer.DataFrameHeader.StreamIDFieldOffset, header.streamId)
	header.buffer.PutInt32(logbuffer.DataFrameHeader.TermIDFieldOffset, termId)
}

type Appender struct {
	termBuffer *atomic.Buffer
	tailBuffer *atomic.Buffer
	tailOffset int32
}

type AppenderResult struct {
	termOffset int64
	termId     int32
}

func (result *AppenderResult) TermOffset() int64 {
	return result.termOffset
}

func (result *AppenderResult) TermId() int32 {
	return result.termId
}

func MakeAppender(termBuffer *atomic.Buffer, metaDataBuffer *atomic.Buffer, partitionIndex int) *Appender {
	appender := new(Appender)
	appender.termBuffer = termBuffer
	appender.tailBuffer = metaDataBuffer
	appender.tailOffset = logbuffer.Descriptor.TermTailCounterOffset + (int32(partitionIndex) * util.SizeOfInt64)

	return appender
}

func (appender *Appender) RawTail() int64 {
	return appender.tailBuffer.GetInt64Volatile(appender.tailOffset)
}

func (appender *Appender) getAndAddRawTail(alignedLength int32) int64 {
	return appender.tailBuffer.GetAndAddInt64(appender.tailOffset, int64(alignedLength))
}

func (appender *Appender) Claim(result *AppenderResult, header *HeaderWriter, length int32, claim *logbuffer.Claim) {

	frameLength := length + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.Alignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermID(rawTail)
	result.termOffset = termOffset + int64(alignedLength)
	if result.termOffset > int64(termLength) {
		appender.handleEndOfLogCondition(result, appender.termBuffer, int32(termOffset), header, termLength)
	} else {
		offset := int32(termOffset)
		header.write(appender.termBuffer, offset, frameLength, result.termId)
		claim.Wrap(appender.termBuffer, offset, frameLength)
	}
}

func (appender *Appender) AppendUnfragmentedMessage(result *AppenderResult, header *HeaderWriter,
	srcBuffer *atomic.Buffer, srcOffset int32, length int32, reservedValueSupplier ReservedValueSupplier) {

	frameLength := length + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.Alignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermID(rawTail)
	result.termOffset = termOffset + int64(alignedLength)
	if result.termOffset > int64(termLength) {
		appender.handleEndOfLogCondition(result, appender.termBuffer, int32(termOffset), header, termLength)
	} else {
		offset := int32(termOffset)
		header.write(appender.termBuffer, offset, frameLength, logbuffer.TermID(rawTail))
		appender.termBuffer.PutBytes(offset+logbuffer.DataFrameHeader.Length, srcBuffer, srcOffset, length)

		if nil != reservedValueSupplier {
			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)
		}

		logbuffer.FrameLengthOrdered(appender.termBuffer, offset, frameLength)
	}
}

func (appender *Appender) AppendFragmentedMessage(result *AppenderResult, header *HeaderWriter,
	srcBuffer *atomic.Buffer, srcOffset int32, length int32, maxPayloadLength int32, reservedValueSupplier ReservedValueSupplier) {

	numMaxPayloads := length / maxPayloadLength
	remainingPayload := length % maxPayloadLength
	var lastFrameLength int32 = 0
	if remainingPayload > 0 {
		lastFrameLength = util.AlignInt32(remainingPayload+logbuffer.DataFrameHeader.Length, logbuffer.FrameDescriptor.Alignment)
	}
	requiredLength := (numMaxPayloads * (maxPayloadLength + logbuffer.DataFrameHeader.Length)) + lastFrameLength
	rawTail := appender.getAndAddRawTail(requiredLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermID(rawTail)
	result.termOffset = termOffset + int64(requiredLength)
	if result.termOffset > int64(termLength) {
		appender.handleEndOfLogCondition(result, appender.termBuffer, int32(termOffset), header, termLength)
	} else {
		flags := logbuffer.FrameDescriptor.BeginFrag
		remaining := length
		offset := int32(termOffset)

		for remaining > 0 {
			bytesToWrite := int32(math.Min(float64(remaining), float64(maxPayloadLength)))
			frameLength := bytesToWrite + logbuffer.DataFrameHeader.Length
			alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.Alignment)

			header.write(appender.termBuffer, offset, frameLength, result.termId)
			appender.termBuffer.PutBytes(
				offset+logbuffer.DataFrameHeader.Length, srcBuffer, srcOffset+(length-remaining), bytesToWrite)

			if remaining <= maxPayloadLength {
				flags |= logbuffer.FrameDescriptor.EndFrag
			}

			logbuffer.FrameFlags(appender.termBuffer, offset, flags)

			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)

			logbuffer.FrameLengthOrdered(appender.termBuffer, offset, frameLength)

			flags = 0
			offset += alignedLength
			remaining -= bytesToWrite
		}
	}
}

func (appender *Appender) handleEndOfLogCondition(result *AppenderResult,
	termBuffer *atomic.Buffer, termOffset int32, header *HeaderWriter, termLength int32) {
	result.termOffset = int64(APPENDER_FAILED)

	if termOffset <= termLength {
		result.termOffset = int64(APPENDER_TRIPPED)

		if termOffset < termLength {
			paddingLength := termLength - termOffset
			header.write(termBuffer, termOffset, paddingLength, result.termId)
			logbuffer.SetFrameType(termBuffer, termOffset, logbuffer.DataFrameHeader.TypePad)
			logbuffer.FrameLengthOrdered(termBuffer, termOffset, paddingLength)
		}
	}
}

func (appender *Appender) SetTailTermId(termId int32) {
	appender.tailBuffer.PutInt64(appender.tailOffset, int64(termId)<<32)
}
