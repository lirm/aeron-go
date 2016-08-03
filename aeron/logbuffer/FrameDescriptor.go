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

package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"unsafe"
)

var FrameDescriptor = struct {
	FRAME_ALIGNMENT int32

	BEGIN_FRAG   uint8
	END_FRAG     uint8
	UNFRAGMENTED uint8

	ALIGNED_HEADER_LENGTH int32

	VERSION_OFFSET int32
	FLAGS_OFFSET   int32
	TYPE_OFFSET    int32
	LENGTH_OFFSET  int32
	TERM_OFFSET    int32
}{
	32,
	0x80,
	0x40,
	0x80 | 0x40,
	32,
	4,
	5,
	6,
	0,
	8,
}

func flagsOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.FLAGS_FIELD_OFFSET
}

func lengthOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.FRAME_LENGTH_FIELD_OFFSET
}

func typeOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.TYPE_FIELD_OFFSET
}

func FrameLengthVolatile(logBuffer *atomic.Buffer, frameOffset int32) int32 {
	offset := lengthOffset(frameOffset)
	return logBuffer.GetInt32Volatile(offset)
}

func FrameLengthOrdered(logBuffer *atomic.Buffer, frameOffset int32, frameLength int32) {
	logBuffer.PutInt32Ordered(lengthOffset(frameOffset), frameLength)
}

func IsPaddingFrame(logBuffer *atomic.Buffer, frameOffset int32) bool {
	return logBuffer.GetUInt16(typeOffset(frameOffset)) == DataFrameHeader.HDR_TYPE_PAD
}

func SetFrameType(logBuffer *atomic.Buffer, frameOffset int32, typ uint16) {
	logBuffer.PutUInt16(typeOffset(frameOffset), typ)
}

func FrameFlags(logBuffer *atomic.Buffer, frameOffset int32, flags uint8) {
	logBuffer.PutUInt8(flagsOffset(frameOffset), flags)
}

func DefaultFrameHeader(logMetaDataBuffer *atomic.Buffer) *atomic.Buffer {
	headerPtr := unsafe.Pointer(uintptr(logMetaDataBuffer.Ptr()) + Descriptor.LOG_DEFAULT_FRAME_HEADER_OFFSET)

	return atomic.MakeBuffer(headerPtr, DataFrameHeader.LENGTH)
}
