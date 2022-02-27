/*
Copyright 2016-2018 Stanislav Liberman
Copyright 2022 Talos, Inc.

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

// See  ~agrona/agrona/src/main/java/org/agrona/concurrent/status/CountersReader.java
//
// Values Buffer
//
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                        Counter Value                          |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                       Registration Id                         |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                          Owner Id                             |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                     104 bytes of padding                     ...
// ...                                                              |
//  +---------------------------------------------------------------+
//  |                   Repeats to end of buffer                   ...
//  |                                                               |
// ...                                                              |
//  +---------------------------------------------------------------+
//
// Meta Data Buffer
//
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                        Record State                           |
//  +---------------------------------------------------------------+
//  |                          Type Id                              |
//  +---------------------------------------------------------------+
//  |                  Free-for-reuse Deadline (ms)                 |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                      112 bytes for key                       ...
// ...                                                              |
//  +-+-------------------------------------------------------------+
//  |R|                      Label Length                           |
//  +-+-------------------------------------------------------------+
//  |                     380 bytes of Label                       ...
// ...                                                              |
//  +---------------------------------------------------------------+
//  |                   Repeats to end of buffer                   ...
//  |                                                               |
// ...                                                              |
//  +---------------------------------------------------------------+

// See: ~aeron/aeron-archive/src/main/java/io/aeron/archive/status/RecordingPos.Java
// for which the below is a subset
//
// The position a recording has reached when being archived.
//
// Key has the following layout:
//
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                        Recording ID                           |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                         Session ID                            |
//  +---------------------------------------------------------------+
//  |                Source Identity for the Image                  |
//  |                                                              ...
// ...                                                              |
//  +---------------------------------------------------------------+

package counters

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/command"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/util"

	"fmt"
	"time"
)

const COUNTER_LENGTH = util.CacheLineLength * 2

const FULL_LABEL_LENGTH = util.CacheLineLength * 6
const LABEL_OFFSET = util.CacheLineLength * 2
const METADATA_LENGTH = LABEL_OFFSET + FULL_LABEL_LENGTH

const TYPE_ID_OFFSET = 4
const FREE_FOR_REUSE_DEADLINE_OFFSET = TYPE_ID_OFFSET + 4
const KEY_OFFSET = FREE_FOR_REUSE_DEADLINE_OFFSET + 8

const RECORD_RECLAIMED int32 = -1
const RECORD_UNUSED int32 = 0
const RECORD_ALLOCATED int32 = 1

const RECORDING_ID_OFFSET = 0
const SESSION_ID_OFFSET = RECORDING_ID_OFFSET + 64
const SOURCE_IDENTITY_LENGTH_OFFSET = SESSION_ID_OFFSET + 32
const SOURCE_IDENTITY_OFFSET = SOURCE_IDENTITY_LENGTH_OFFSET + 32

const NULL_RECORDING_ID = -1
const NULL_COUNTER_ID = -1

type Reader struct {
	metaData *atomic.Buffer
	values   *atomic.Buffer

	maxCounterID int
}

type Counter struct {
	Id     int32
	TypeId int32
	Value  int64
	Label  string
}

func NewReader(values, metaData *atomic.Buffer) *Reader {

	reader := Reader{metaData: metaData, values: values}
	reader.maxCounterID = int(values.Capacity() / COUNTER_LENGTH)

	return &reader
}

func (reader *Reader) Scan(cb func(Counter)) {

	var id int32 = 0
	var i int32 = 0

	for capacity := reader.metaData.Capacity(); i < capacity; i += METADATA_LENGTH {
		recordStatus := reader.metaData.GetInt32Volatile(i)
		if RECORD_UNUSED == recordStatus {
			break
		} else if RECORD_ALLOCATED == recordStatus {
			typeId := reader.metaData.GetInt32(i + TYPE_ID_OFFSET)
			labelSize := reader.metaData.GetInt32(i + LABEL_OFFSET)
			label := string(reader.metaData.GetBytesArray(i+4+LABEL_OFFSET, labelSize))

			// TODO Get the key buffer

			value := reader.values.GetInt64Volatile(id * COUNTER_LENGTH)

			//fmt.Printf("Reading at offset %d; counterState=%d; typeId=%d\n", i, recordStatus, typeId)

			cb(Counter{id, typeId, value, label})

			id++
		}
	}
}

// MaxCounterID returns the reader's maximum CounterID
func (reader *Reader) MaxCounterID() int {
	return reader.maxCounterID
}

// MetaData returns the reader's MetaData
func (reader *Reader) MetaData() *atomic.Buffer {
	return reader.metaData
}

// CounterState returns  the state for a counter
func (reader *Reader) CounterState(counterID int32) int32 {
	return reader.metaData.GetInt32Volatile(counterID * METADATA_LENGTH)
}

// CounterTypeID returns the typeID for a counter
func (reader *Reader) CounterTypeID(counterID int32) int32 {
	return reader.metaData.GetInt32(counterID*METADATA_LENGTH + TYPE_ID_OFFSET)
}

// FindCounterIDByRecording finds the active counter id for a stream
// based on the recording id
// Returns the counterID if found otherwise NULL_COUNTER_ID
func (reader *Reader) FindCounterIDByRecording(recordingID int64) int32 {
	buffer := reader.MetaData()

	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		if counterState == RECORD_ALLOCATED && reader.CounterTypeID(i) == command.RecordingPosition {
			if recordingID == buffer.GetInt64((i*METADATA_LENGTH)+KEY_OFFSET+RECORDING_ID_OFFSET) {
				return i
			}
		} else if RECORD_UNUSED == counterState {
			break
		}
	}
	return NULL_COUNTER_ID
}

// FindCounterIDBySession finds the active counterID for a stream based on the sessionID
// Returns the counterID if found otherwise NULL_COUNTER_ID
func (reader *Reader) FindCounterIDBySession(sessionID int32) int32 {
	buffer := reader.MetaData()
	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		counterTypeID := reader.CounterTypeID(i)
		if counterState == RECORD_ALLOCATED && counterTypeID == command.RecordingPosition {
			if sessionID == buffer.GetInt32((i*METADATA_LENGTH)+KEY_OFFSET+SESSION_ID_OFFSET) {
				return i
			}
		} else if RECORD_UNUSED == counterState {
			break
		}
	}
	return NULL_COUNTER_ID
}

// RecordingID finds the active counterID for a stream based on the sessionID
func (reader *Reader) RecordingID(counterID int32) int64 {
	buffer := reader.MetaData()
	if reader.CounterState(counterID) == RECORD_ALLOCATED && reader.CounterTypeID(counterID) == command.RecordingPosition {
		return buffer.GetInt64((counterID * METADATA_LENGTH) + KEY_OFFSET + RECORDING_ID_OFFSET)
	}
	return NULL_RECORDING_ID
}

// SourceIdentity retturns the source identity for the recording
func (reader *Reader) SourceIdentity(counterID int32) string {
	buffer := reader.MetaData()
	if reader.CounterState(counterID) == RECORD_ALLOCATED && reader.CounterTypeID(counterID) == command.RecordingPosition {
		offset := (counterID * METADATA_LENGTH) + KEY_OFFSET
		length := buffer.GetInt32(offset + SOURCE_IDENTITY_LENGTH_OFFSET)
		byteArray := buffer.GetBytesArray(offset+SOURCE_IDENTITY_OFFSET, length)
		return string(byteArray)
	}

	return ""
}

// IsActive tells us if the recording counter is still active?
func (reader *Reader) IsActive(counterID int32, recordingID int64) bool {
	buffer := reader.MetaData()
	offset := counterID * METADATA_LENGTH
	return reader.CounterTypeID(counterID) == command.RecordingPosition && buffer.GetInt64(offset+KEY_OFFSET+RECORDING_ID_OFFSET) == recordingID && reader.CounterState(counterID) == RECORD_ALLOCATED
}

// CounterValue returns the value for a given counterID
func (reader *Reader) CounterValue(counterID int32) int64 {
	buffer := reader.MetaData()
	offset := counterID * METADATA_LENGTH // The value offset is zero as it's first up
	return buffer.GetInt64Volatile(offset)
}

// AwaitRecordingCounterID waits for a specific counterID to appear
// and return it On timeout return (NULLCOUNTER_ID, error)
func (reader *Reader) AwaitRecordingCounterID(sessionID int32) (int32, error) {
	start := time.Now()
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 10}

	for {
		counterID := reader.FindCounterIDBySession(sessionID)
		if counterID != NULL_COUNTER_ID {
			return counterID, nil
		}

		// Idle and check timeout
		idler.Idle(0)
		if time.Since(start) > time.Second*5 { // FIXME: tune, // FIXME: paramterize?
			return NULL_COUNTER_ID, fmt.Errorf("Timeout waiting for counter:%d", counterID)
		}
	}
}

// AwaitPosition waits for a counter to reach (or exceed) a position
func (reader *Reader) AwaitPosition(counterID int32, position int64) error {
	start := time.Now()
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 10}

	for reader.CounterValue(counterID) < position {
		if reader.CounterState(counterID) != RECORD_ALLOCATED {
			return fmt.Errorf("counter %d not active", counterID)
		}

		// Idle and check timeout
		idler.Idle(0)
		if time.Since(start) > time.Second*5 { // FIXME: tune, // FIXME: paramterize?
			return fmt.Errorf("Timeout waiting for counter:%d position", counterID)
		}
	}

	return nil
}
