/*
Copyright 2016-2017 Stanislav Liberman

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
	NeedsCleaning                 = 1
	PartitionCount                = 3
	LogMetaDataSectionIndex       = PartitionCount

	termMinLength        int32 = 64 * 1024
	termMaxLength        int32 = 1024 * 1024 * 1024
	pageMinSize          int32 = 4 * 1024
	pageMaxSize          int32 = 1024 * 1024 * 1024
	maxSingleMappingSize int64 = 0x7FFFFFFF
	logMetaDataLength          = pageMinSize
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
 *  |                      Active Term Count                        |
 *  +---------------------------------------------------------------+
 *  |                     Cache Line Padding                       ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 *  |                    End of Stream Position                     |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                        Is Connected                           |
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
 *  |                         Term Length                           |
 *  +---------------------------------------------------------------+
 *  |                          Page Size                            |
 *  +---------------------------------------------------------------+
 *  |                      Cache Line Padding                      ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 *  |                     Default Frame Header                     ...
 * ...                                                              |
 *  +---------------------------------------------------------------+
 * </pre>
 */
type LogBufferMetaData struct {
	flyweight.FWBase

	TailCounter        []flyweight.Int64Field // 0, 8, 16
	ActiveTermCountOff flyweight.Int32Field   // 24
	padding0           flyweight.Padding      // 28
	EndOfStreamPosOff  flyweight.Int64Field   // 128
	IsConnected        flyweight.Int32Field   // 136
	padding1           flyweight.Padding      // 140
	CorrelationId      flyweight.Int64Field   // 256
	InitTermID         flyweight.Int32Field   // 264
	DefaultFrameHdrLen flyweight.Int32Field   // 270
	MTULen             flyweight.Int32Field   // 274
	TermLen            flyweight.Int32Field   // 278
	PageSize           flyweight.Int32Field   // 282
	padding2           flyweight.Padding      // 286
	DefaultFrameHeader flyweight.RawDataField // 384
	padding3           flyweight.Padding
}

func (m *LogBufferMetaData) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	m.TailCounter = make([]flyweight.Int64Field, 3)
	pos += m.TailCounter[0].Wrap(buf, pos)
	pos += m.TailCounter[1].Wrap(buf, pos)
	pos += m.TailCounter[2].Wrap(buf, pos)
	pos += m.ActiveTermCountOff.Wrap(buf, pos)
	pos += m.padding0.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)
	pos += m.EndOfStreamPosOff.Wrap(buf, pos)
	pos += m.IsConnected.Wrap(buf, pos)
	pos += m.padding1.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.InitTermID.Wrap(buf, pos)
	pos += m.DefaultFrameHdrLen.Wrap(buf, pos)
	pos += m.MTULen.Wrap(buf, pos)
	pos += m.TermLen.Wrap(buf, pos)
	pos += m.PageSize.Wrap(buf, pos)
	pos += m.padding2.Wrap(buf, pos, util.CacheLineLength, util.CacheLineLength)
	pos += m.DefaultFrameHeader.Wrap(buf, pos, DataFrameHeader.Length)
	pos += m.padding3.Wrap(buf, pos, util.CacheLineLength*2, util.CacheLineLength)

	m.SetSize(pos - offset)
	return m
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

func computeTermLength(logLength int64) int64 {
	return (logLength - int64(logMetaDataLength)) / int64(PartitionCount)
}

func TermID(rawTail int64) int32 {
	return int32(rawTail >> 32)
}

func RotateLog(logMetaDataBuffer *LogBufferMetaData, currentTermCount int32, currentTermId int32) {
	nextTermId := currentTermId + 1
	nextTermCount := currentTermCount + 1
	nextIndex := nextTermCount % PartitionCount
	expectedTermId := nextTermId - PartitionCount

	tail := logMetaDataBuffer.TailCounter[nextIndex]
	rawTail := tail.Get()
	if expectedTermId == TermID(rawTail) {
		for !tail.CAS(rawTail, int64(nextTermId)<<32) {

		}
	}

	// This should be CAS
	// LogBufferDescriptor::casActiveTermCount(logMetaDataBuffer, currentTermCount, nextTermCount)
	logMetaDataBuffer.ActiveTermCountOff.Set(nextTermCount)
}
