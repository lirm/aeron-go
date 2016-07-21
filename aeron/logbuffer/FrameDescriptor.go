package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/buffers"
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

func FrameLengthVolatile(logBuffer *buffers.Atomic, frameOffset int32) int32 {
	offset := lengthOffset(frameOffset)
	return logBuffer.GetInt32Volatile(offset)
}

func FrameLengthOrdered(logBuffer *buffers.Atomic, frameOffset int32, frameLength int32) {
	logBuffer.PutInt32Ordered(lengthOffset(frameOffset), frameLength)
}

func IsPaddingFrame(logBuffer *buffers.Atomic, frameOffset int32) bool {
	return logBuffer.GetUInt16(typeOffset(frameOffset)) == DataFrameHeader.HDR_TYPE_PAD
}

func SetFrameType(logBuffer *buffers.Atomic, frameOffset int32, typ uint16) {
	logBuffer.PutUInt16(typeOffset(frameOffset), typ)
}

func FrameFlags(logBuffer *buffers.Atomic, frameOffset int32, flags uint8) {
	logBuffer.PutUInt8(flagsOffset(frameOffset), flags)
}

func DefaultFrameHeader(logMetaDataBuffer *buffers.Atomic) *buffers.Atomic {
	headerPtr := unsafe.Pointer(uintptr(logMetaDataBuffer.Ptr()) + Descriptor.LOG_DEFAULT_FRAME_HEADER_OFFSET)

	return buffers.MakeAtomic(headerPtr, DataFrameHeader.LENGTH)
}
