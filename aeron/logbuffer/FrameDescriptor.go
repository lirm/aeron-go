package logbuffer

import (
	"fmt"
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

func termOffsetOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.TERM_OFFSET_FIELD_OFFSET
}

func typeOffset(frameOffset int32) int32 {
	return frameOffset + DataFrameHeader.TYPE_FIELD_OFFSET
}

func FrameLengthVolatile(logBuffer *buffers.Atomic, frameOffset int32) int32 {
	// TODO: need to byte order to LITTLE_ENDIAN
	offset := lengthOffset(frameOffset)
	// log.Printf("lengthOffset: %d", offset)
	return logBuffer.GetInt32Volatile(offset)
}

func FrameLengthOrdered(logBuffer *buffers.Atomic, frameOffset int32, frameLength int32) {
	// TODO: need to byte order to LITTLE_ENDIAN
	logBuffer.PutInt32Ordered(lengthOffset(frameOffset), frameLength)
}

func IsPaddingFrame(logBuffer *buffers.Atomic, frameOffset int32) bool {
	return logBuffer.GetUInt16(typeOffset(frameOffset)) == DataFrameHeader.HDR_TYPE_PAD
}

func checkHeaderLength(length int32) {
	if length != DataFrameHeader.LENGTH {
		panic(
			fmt.Sprintf("Frame header length %d must be equal to %d", length, DataFrameHeader.LENGTH))
	}
}

func checkMaxFrameLength(length int32) {
	if (length & (FrameDescriptor.FRAME_ALIGNMENT - 1)) != 0 {
		panic(
			fmt.Sprintf("Max frame length must be a multiple of %d, length=%d", FrameDescriptor.FRAME_ALIGNMENT, length))
	}
}

func computeMaxMessageLength(capacity int32) int32 {
	return capacity / 8
}

func setFrameType(logBuffer *buffers.Atomic, frameOffset int32, typ uint16) {
	logBuffer.PutUInt16(typeOffset(frameOffset), typ)
}

func FrameType(logBuffer *buffers.Atomic, frameOffset int32) uint16 {
	return logBuffer.GetUInt16(frameOffset)
}

func SetFrameType(logBuffer *buffers.Atomic, frameOffset int32, typ uint16) {
	logBuffer.PutUInt16(typeOffset(frameOffset), typ)
}

func frameFlags(logBuffer *buffers.Atomic, frameOffset int32, flags uint8) {
	logBuffer.PutUInt8(flagsOffset(frameOffset), flags)
}

func frameTermOffset(logBuffer *buffers.Atomic, frameOffset int32, termOffset int32) {
	logBuffer.PutInt32(termOffsetOffset(frameOffset), termOffset)
}

func frameVersion(logBuffer *buffers.Atomic, frameOffset int32) uint8 {
	return logBuffer.GetUInt8(frameOffset)
}

func DefaultFrameHeader(logMetaDataBuffer *buffers.Atomic) *buffers.Atomic {
	headerPtr := unsafe.Pointer(uintptr(logMetaDataBuffer.Ptr()) + uintptr(Descriptor.LOG_DEFAULT_FRAME_HEADER_OFFSET))

	return buffers.MakeAtomic(headerPtr, DataFrameHeader.LENGTH)
}
