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
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/flyweight"
	"github.com/lirm/aeron-go/aeron/util"
)

const (
	Clean                   int32 = 0
	NeedsCleaning           int32 = 1
	PartitionCount          int   = 3
	LogMetaDataSectionIndex int   = PartitionCount

	termMinLength        int32 = 64 * 1024
	maxSingleMappingSize int64 = 0x7FFFFFFF

	sizeofLogMetadata int32 = util.CacheLineLength * 2
)

/**
 * <pre>
 *   0                   1                   2                   3
 *   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |                       Tail Counter 0                          |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                       Tail Counter 1                          |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                       Tail Counter 2                          |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                   Active Partition Index                      |
 *  +---------------------------------------------------------------+
 *  |                      Cache Line Padding                      ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 *  |                 Time of Last Status Message                   |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                      Cache Line Padding                      ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 *  |                 Registration / Correlation ID                 |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                        Initial Term Id                        |
 *  +---------------------------------------------------------------+
 *  |                  Default Frame Header Length                  |
 *  +---------------------------------------------------------------+
 *  |                          MTU Length                           |
 *  +---------------------------------------------------------------+
 *  |                      Cache Line Padding                      ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 *  |                    Default Frame Header                      ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 * </pre>
 */

type LogBufferMetaData struct {
	flyweight.FWBase

	TailCounter0        flyweight.Int64Field
	TailCounter1        flyweight.Int64Field
	TailCounter2        flyweight.Int64Field
	ActivePartitionIx   flyweight.Int32Field
	padding0            flyweight.Padding
	TimeOfLastStatusMsg flyweight.Int64Field
	padding1            flyweight.Padding
	RegID               flyweight.Int64Field
	InitTermID          flyweight.Int32Field
	DefaultFrameHdrLen  flyweight.Int32Field
	MTULen              flyweight.Int32Field
	padding2            flyweight.Padding
	DefaultFrameHeader  flyweight.RawDataField
}

func (m *LogBufferMetaData) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.TailCounter0.Wrap(buf, pos)
	pos += m.TailCounter1.Wrap(buf, pos)
	pos += m.TailCounter2.Wrap(buf, pos)
	pos += m.ActivePartitionIx.Wrap(buf, pos)
	pos += m.padding0.Wrap(buf, pos, util.CacheLineLength*2)
	pos += m.TimeOfLastStatusMsg.Wrap(buf, pos)
	pos += m.padding1.Wrap(buf, pos, util.CacheLineLength*2)
	pos += m.RegID.Wrap(buf, pos)
	pos += m.InitTermID.Wrap(buf, pos)
	pos += m.DefaultFrameHdrLen.Wrap(buf, pos)
	pos += m.MTULen.Wrap(buf, pos)
	pos += m.padding1.Wrap(buf, pos, util.CacheLineLength*2)
	pos += m.DefaultFrameHeader.Wrap(buf, pos, m.DefaultFrameHdrLen.Get())

	m.SetSize(pos - offset)
	return m
}

var Descriptor = struct {
	TermTailCounterOffset int32

	logActivePartitionIndexOffset    int32
	logTimeOfLastStatusMessageOffset int32
	logInitialTermIDOffset           int32
	logDefaultFrameHeaderLenOffset   int32
	logMtuLengthOffset               int32

	logDefaultFrameHeaderOffset uintptr
	logMetaDataLength           int32
}{
	0,

	util.SizeOfInt64 * int32(PartitionCount),
	util.CacheLineLength * 2,
	util.CacheLineLength*4 + 8,
	util.CacheLineLength*4 + 12,
	util.CacheLineLength*4 + 16,

	uintptr(util.CacheLineLength * 5),
	util.CacheLineLength * 7,
}

func checkTermLength(termLength int64) {
	if termLength < int64(termMinLength) {
		panic(fmt.Sprintf("Term length less than min size of %d, length=%d",
			termMinLength, termLength))
	}

	if (termLength & (int64(FrameAlignment) - 1)) != 0 {
		panic(fmt.Sprintf("Term length not a multiple of %d, length=%d",
			FrameAlignment, termLength))
	}
}

func ComputeTermBeginPosition(activeTermID, positionBitsToShift, initialTermID int32) int64 {
	termCount := int64(activeTermID - initialTermID)

	return termCount << uint32(positionBitsToShift)
}

func computeTermLength(logLength int64) int64 {
	return (logLength - int64(Descriptor.logMetaDataLength)) / int64(PartitionCount)
}

func IndexByPosition(position int64, positionBitsToShift uint8) int32 {
	term := uint64(position) >> positionBitsToShift
	return util.FastMod3(term)
}

func InitialTermID(logMetaDataBuffer *atomic.Buffer) int32 {
	return logMetaDataBuffer.GetInt32(Descriptor.logInitialTermIDOffset)
}

func MtuLength(logMetaDataBuffer *atomic.Buffer) int32 {
	return logMetaDataBuffer.GetInt32(Descriptor.logMtuLengthOffset)
}

func ActivePartitionIndex(logMetaDataBuffer *atomic.Buffer) int32 {
	return logMetaDataBuffer.GetInt32Volatile(Descriptor.logActivePartitionIndexOffset)
}

func SetActivePartitionIndex(logMetaDataBuffer *atomic.Buffer, index int32) {
	logMetaDataBuffer.PutInt32Ordered(Descriptor.logActivePartitionIndexOffset, index)
}

func TimeOfLastStatusMessage(logMetaDataBuffer *atomic.Buffer) int64 {
	return logMetaDataBuffer.GetInt64Volatile(Descriptor.logTimeOfLastStatusMessageOffset)
}

func TermID(rawTail int64) int32 {
	return int32(rawTail >> 32)
}

func NextPartitionIndex(currentIndex int32) int32 {
	return util.FastMod3(uint64(currentIndex) + 1)
}
