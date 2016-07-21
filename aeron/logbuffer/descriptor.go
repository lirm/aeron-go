package logbuffer

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/util"
)

const (
	CLEAN           int32 = 0
	NEEDS_CLEANING  int32 = 1
	PARTITION_COUNT int   = 3
)

var Descriptor = struct {
	TERM_MIN_LENGTH         int32
	MAX_SINGLE_MAPPING_SIZE int64

	LOG_META_DATA_SECTION_INDEX         int
	LOG_DEFAULT_FRAME_HEADER_MAX_LENGTH int32

	SIZEOF_LOG_METADATA int32

	TERM_TAIL_COUNTER_OFFSET               int32
	LOG_ACTIVE_PARTITION_INDEX_OFFSET      int32
	LOG_TIME_OF_LAST_STATUS_MESSAGE_OFFSET int32
	LOG_INITIAL_TERM_ID_OFFSET             int32
	LOG_DEFAULT_FRAME_HEADER_LENGTH_OFFSET int32
	LOG_MTU_LENGTH_OFFSET                  int32

	LOG_DEFAULT_FRAME_HEADER_OFFSET uintptr
	LOG_META_DATA_LENGTH            int32
}{
	64 * 1024,
	0x7FFFFFFF,

	PARTITION_COUNT,
	util.CACHE_LINE_LENGTH * 2,

	util.CACHE_LINE_LENGTH * 2,

	0,
	util.SIZEOF_INT64 * int32(PARTITION_COUNT),
	util.CACHE_LINE_LENGTH * 2,
	util.CACHE_LINE_LENGTH*4 + 8,
	util.CACHE_LINE_LENGTH*4 + 12,
	util.CACHE_LINE_LENGTH*4 + 16,

	uintptr(util.CACHE_LINE_LENGTH * 5),
	util.CACHE_LINE_LENGTH * 7,
}

func checkTermLength(termLength int64) {
	if termLength < int64(Descriptor.TERM_MIN_LENGTH) {
		panic(fmt.Sprintf("Term length less than min size of %d, length=%d",
			Descriptor.TERM_MIN_LENGTH, termLength))
	}

	if (termLength & (int64(FrameDescriptor.FRAME_ALIGNMENT) - 1)) != 0 {
		panic(fmt.Sprintf("Term length not a multiple of %d, length=%d",
			FrameDescriptor.FRAME_ALIGNMENT, termLength))
	}
}

func ComputeTermBeginPosition(activeTermId, positionBitsToShift, initialTermId int32) int64 {
	termCount := int64(activeTermId - initialTermId)

	return termCount << uint32(positionBitsToShift)
}

func computeTermLength(logLength int64) int64 {
	return (logLength - int64(Descriptor.LOG_META_DATA_LENGTH)) / int64(PARTITION_COUNT)
}

func IndexByPosition(position int64, positionBitsToShift uint8) int32 {
	term := uint64(position) >> positionBitsToShift
	return util.FastMod3(term)
}

func InitialTermId(logMetaDataBuffer *buffers.Atomic) int32 {
	return logMetaDataBuffer.GetInt32(Descriptor.LOG_INITIAL_TERM_ID_OFFSET)
}

func MtuLength(logMetaDataBuffer *buffers.Atomic) int32 {
	return logMetaDataBuffer.GetInt32(Descriptor.LOG_MTU_LENGTH_OFFSET)
}

func ActivePartitionIndex(logMetaDataBuffer *buffers.Atomic) int32 {
	return logMetaDataBuffer.GetInt32Volatile(Descriptor.LOG_ACTIVE_PARTITION_INDEX_OFFSET)
}

func SetActivePartitionIndex(logMetaDataBuffer *buffers.Atomic, index int32) {
	logMetaDataBuffer.PutInt32Ordered(Descriptor.LOG_ACTIVE_PARTITION_INDEX_OFFSET, index)
}

func TimeOfLastStatusMessage(logMetaDataBuffer *buffers.Atomic) int64 {
	return logMetaDataBuffer.GetInt64Volatile(Descriptor.LOG_TIME_OF_LAST_STATUS_MESSAGE_OFFSET)
}

func TermId(rawTail int64) int32 {
	return int32(rawTail >> 32)
}

func NextPartitionIndex(currentIndex int32) int32 {
	return util.FastMod3(uint64(currentIndex) + 1)
}
