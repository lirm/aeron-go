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
)

const (
	// FrameAlignment frame alignment
	FrameAlignment int32 = 32
)

func flagsOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.FlagsFieldOffset
}

func lengthOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.FrameLengthFieldOffset
}

func typeOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.TypeFieldOffset
}

func FrameLengthVolatile(logBuffer *atomic.Buffer, frameOffset int32) int32 {
	offset := lengthOffset(frameOffset)
	return logBuffer.GetInt32Volatile(offset)
}

func FrameLengthOrdered(logBuffer *atomic.Buffer, frameOffset int32, frameLength int32) {
	logBuffer.PutInt32Ordered(lengthOffset(frameOffset), frameLength)
}

func IsPaddingFrame(logBuffer *atomic.Buffer, frameOffset int32) bool {
	return logBuffer.GetUInt16(typeOffset(frameOffset)) == DataFrameHeader.TypePad
}

func SetFrameType(logBuffer *atomic.Buffer, frameOffset int32, typ uint16) {
	logBuffer.PutUInt16(typeOffset(frameOffset), typ)
}

func FrameFlags(logBuffer *atomic.Buffer, frameOffset int32, flags uint8) {
	logBuffer.PutUInt8(flagsOffset(frameOffset), flags)
}

func DefaultFrameHeader(logMetaDataBuffer *atomic.Buffer) *atomic.Buffer {
	headerPtr := uintptr(logMetaDataBuffer.Ptr()) + Descriptor.logDefaultFrameHeaderOffset

	return atomic.MakeBuffer(headerPtr, DataFrameHeader.Length)
}
