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
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
	"math"
	"unsafe"
	"github.com/lirm/aeron-go/aeron/atomic"
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
}

func (header *HeaderWriter) Fill(defaultHdr *atomic.Buffer) {
	header.sessionId = defaultHdr.GetInt32(logbuffer.DataFrameHeader.SESSION_ID_FIELD_OFFSET)
	header.streamId = defaultHdr.GetInt32(logbuffer.DataFrameHeader.STREAM_ID_FIELD_OFFSET)
}

func (header *HeaderWriter) write(termBuffer *atomic.Buffer, offset, length, termId int32) {
	termBuffer.PutInt32Ordered(offset, -length)

	headerPtr := uintptr(termBuffer.Ptr()) + uintptr(offset)
	headerBuffer := atomic.MakeBuffer(unsafe.Pointer(headerPtr), logbuffer.DataFrameHeader.LENGTH)

	headerBuffer.PutInt8(logbuffer.DataFrameHeader.VERSION_FIELD_OFFSET, logbuffer.DataFrameHeader.CURRENT_VERSION)
	headerBuffer.PutUInt8(logbuffer.DataFrameHeader.FLAGS_FIELD_OFFSET, logbuffer.FrameDescriptor.BEGIN_FRAG|logbuffer.FrameDescriptor.END_FRAG)
	headerBuffer.PutUInt16(logbuffer.DataFrameHeader.TYPE_FIELD_OFFSET, logbuffer.DataFrameHeader.HDR_TYPE_DATA)
	headerBuffer.PutInt32(logbuffer.DataFrameHeader.TERM_OFFSET_FIELD_OFFSET, offset)
	headerBuffer.PutInt32(logbuffer.DataFrameHeader.SESSION_ID_FIELD_OFFSET, header.sessionId)
	headerBuffer.PutInt32(logbuffer.DataFrameHeader.STREAM_ID_FIELD_OFFSET, header.streamId)
	headerBuffer.PutInt32(logbuffer.DataFrameHeader.TERM_ID_FIELD_OFFSET, termId)
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
	appender.tailOffset = logbuffer.Descriptor.TERM_TAIL_COUNTER_OFFSET + (int32(partitionIndex) * util.SIZEOF_INT64)

	return appender
}

func (appender *Appender) RawTail() int64 {
	return appender.tailBuffer.GetInt64Volatile(appender.tailOffset)
}

func (appender *Appender) getAndAddRawTail(alignedLength int32) int64 {
	return appender.tailBuffer.GetAndAddInt64(appender.tailOffset, int64(alignedLength))
}

func (appender *Appender) Claim(result *AppenderResult, header *HeaderWriter, length int32, claim *logbuffer.Claim) {

	frameLength := length + logbuffer.DataFrameHeader.LENGTH
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.FRAME_ALIGNMENT)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermId(rawTail)
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

	frameLength := length + logbuffer.DataFrameHeader.LENGTH
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.FRAME_ALIGNMENT)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermId(rawTail)
	result.termOffset = termOffset + int64(alignedLength)
	if result.termOffset > int64(termLength) {
		appender.handleEndOfLogCondition(result, appender.termBuffer, int32(termOffset), header, termLength)
	} else {
		offset := int32(termOffset)
		header.write(appender.termBuffer, offset, frameLength, logbuffer.TermId(rawTail))
		appender.termBuffer.PutBytes(offset+logbuffer.DataFrameHeader.LENGTH, srcBuffer, srcOffset, length)

		if nil != reservedValueSupplier {
			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.RESERVED_VALUE_FIELD_OFFSET, reservedValue)
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
		lastFrameLength = util.AlignInt32(remainingPayload+logbuffer.DataFrameHeader.LENGTH, logbuffer.FrameDescriptor.FRAME_ALIGNMENT)
	}
	requiredLength := (numMaxPayloads * (maxPayloadLength + logbuffer.DataFrameHeader.LENGTH)) + lastFrameLength
	rawTail := appender.getAndAddRawTail(requiredLength)
	termOffset := rawTail & 0xFFFFFFFF

	termLength := appender.termBuffer.Capacity()

	result.termId = logbuffer.TermId(rawTail)
	result.termOffset = termOffset + int64(requiredLength)
	if result.termOffset > int64(termLength) {
		appender.handleEndOfLogCondition(result, appender.termBuffer, int32(termOffset), header, termLength)
	} else {
		flags := logbuffer.FrameDescriptor.BEGIN_FRAG
		remaining := length
		offset := int32(termOffset)

		for remaining > 0 {
			bytesToWrite := int32(math.Min(float64(remaining), float64(maxPayloadLength)))
			frameLength := bytesToWrite + logbuffer.DataFrameHeader.LENGTH
			alignedLength := util.AlignInt32(frameLength, logbuffer.FrameDescriptor.FRAME_ALIGNMENT)

			header.write(appender.termBuffer, offset, frameLength, result.termId)
			appender.termBuffer.PutBytes(
				offset+logbuffer.DataFrameHeader.LENGTH, srcBuffer, srcOffset+(length-remaining), bytesToWrite)

			if remaining <= maxPayloadLength {
				flags |= logbuffer.FrameDescriptor.END_FRAG
			}

			logbuffer.FrameFlags(appender.termBuffer, offset, flags)

			reservedValue := reservedValueSupplier(appender.termBuffer, offset, frameLength)
			appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.RESERVED_VALUE_FIELD_OFFSET, reservedValue)

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
			logbuffer.SetFrameType(termBuffer, termOffset, logbuffer.DataFrameHeader.HDR_TYPE_PAD)
			logbuffer.FrameLengthOrdered(termBuffer, termOffset, paddingLength)
		}
	}
}

func (appender *Appender) SetTailTermId(termId int32) {
	appender.tailBuffer.PutInt64(appender.tailOffset, int64(termId)<<32)
}
