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

/* LogBufferMetaData is the flyweight for LogBuffer meta data
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

	TailCounter         []flyweight.Int64Field // 0, 8, 16
	ActivePartitionIx   flyweight.Int32Field // 24
	padding0            flyweight.Padding //28
	TimeOfLastStatusMsg flyweight.Int64Field //128
	padding1            flyweight.Padding
	RegID               flyweight.Int64Field // 256
	InitTermID          flyweight.Int32Field
	DefaultFrameHdrLen  flyweight.Int32Field
	MTULen              flyweight.Int32Field
	padding2            flyweight.Padding
	DefaultFrameHeader  flyweight.RawDataField // 384
	padding3            flyweight.Padding
}

func (m *LogBufferMetaData) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	m.TailCounter = make([]flyweight.Int64Field, 3)
	pos += m.TailCounter[0].Wrap(buf, pos)
	pos += m.TailCounter[1].Wrap(buf, pos)
	pos += m.TailCounter[2].Wrap(buf, pos)
	pos += m.ActivePartitionIx.Wrap(buf, pos)
	pos += m.padding0.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)
	//fmt.Printf(" -> pos %d\n", pos)
	pos += m.TimeOfLastStatusMsg.Wrap(buf, pos)
	pos += m.padding1.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)
	//fmt.Printf(" -> pos %d\n", pos)
	pos += m.RegID.Wrap(buf, pos)
	pos += m.InitTermID.Wrap(buf, pos)
	pos += m.DefaultFrameHdrLen.Wrap(buf, pos)
	pos += m.MTULen.Wrap(buf, pos)
	//fmt.Printf(" -> pos %d\n", pos)
	pos += m.padding2.Wrap(buf, pos, util.CacheLineLength, util.CacheLineLength)
	//fmt.Printf(" -> pos %d\n", pos)
	pos += m.DefaultFrameHeader.Wrap(buf, pos, DataFrameHeader.Length)
	//fmt.Printf(" -> pos %d\n", pos)
	pos += m.padding3.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)
	//fmt.Printf(" -> pos %d\n", pos)

	m.SetSize(pos - offset)
	return m
}

var Descriptor = struct {
	TermTailCounterOffset int32

	logDefaultFrameHeaderOffset uintptr
	logMetaDataLength           int32
}{
	0,

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

func TermID(rawTail int64) int32 {
	return int32(rawTail >> 32)
}

func NextPartitionIndex(currentIndex int32) int32 {
	return util.FastMod3(uint64(currentIndex) + 1)
}
