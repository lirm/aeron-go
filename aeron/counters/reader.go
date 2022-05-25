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

// See: SubscriptionReadyFlyweights.java
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                        Correlation ID                         |
//  |                                                               |
//  +---------------------------------------------------------------+
//  |                  Channel Status Indicator ID                  |
//  +---------------------------------------------------------------+

package counters

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/command"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/util"

)

const COUNTER_LENGTH = util.CacheLineLength * 2

const FULL_LABEL_LENGTH = util.CacheLineLength * 6
const LABEL_OFFSET = util.CacheLineLength * 2
const METADATA_LENGTH = LABEL_OFFSET + FULL_LABEL_LENGTH
const MAX_KEY_LENGTH = (util.CacheLineLength * 2) - (util.SizeOfInt32 * 2) - util.SizeOfInt64

const TYPE_ID_OFFSET = 4
const FREE_FOR_REUSE_DEADLINE_OFFSET = TYPE_ID_OFFSET + 4
const KEY_OFFSET = FREE_FOR_REUSE_DEADLINE_OFFSET + 8

const RECORD_RECLAIMED int32 = -1
const RECORD_UNUSED int32 = 0
const RECORD_ALLOCATED int32 = 1

const NullCounterId = int32(-1)

const RECORDING_ID_OFFSET = 0
const SESSION_ID_OFFSET = RECORDING_ID_OFFSET + 8
const SOURCE_IDENTITY_LENGTH_OFFSET = SESSION_ID_OFFSET + 4
const SOURCE_IDENTITY_OFFSET = SOURCE_IDENTITY_LENGTH_OFFSET + 4

const NULL_RECORDING_ID = -1

// From LocalSocketAddressStatus.Java
const CHANNEL_STATUS_ID_OFFSET = 0
const LOCAL_SOCKET_ADDRESS_LENGTH_OFFSET = CHANNEL_STATUS_ID_OFFSET + 4
const LOCAL_SOCKET_ADDRESS_STRING_OFFSET = LOCAL_SOCKET_ADDRESS_LENGTH_OFFSET + 4
const MAX_IPV6_LENGTH = len("[ffff:ffff:ffff:ffff:ffff:ffff:255.255.255.255]:65536")
const LOCAL_SOCKET_ADDRESS_STATUS_TYPE_ID = 14

// Channel Status is set when the listeneradapter callbacks are invoked
const (
	ChannelStatusErrored      = -1 // Channel has errored. Check logs for information
	ChannelStatusInitializing = 0  // Channel is being initialized
	ChannelStatusActive       = 1  // Channel has finished initialization and is active
	ChannelStatusClosing      = 2  // Channel is being closed
)

// StatusString provides a convenience method for logging and error handling
func ChannelStatusString(channelStatus int) string {
	switch channelStatus {
	case ChannelStatusErrored:
		return "ChannelStatusErrored"
	case ChannelStatusInitializing:
		return "ChannelStatusInitializing"
	case ChannelStatusActive:
		return "ChannelStatusActive"
	case ChannelStatusClosing:
		return "ChannelStatusClosing"
	default:
		return "Unknown"
	}
}

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
			label := reader.labelValue(i)

			// TODO Get the key buffer

			value := reader.values.GetInt64Volatile(id * COUNTER_LENGTH)

			//fmt.Printf("Reading at offset %d; counterState=%d; typeId=%d\n", i, recordStatus, typeId)

			cb(Counter{id, typeId, value, label})
		}
		id++
	}
}

func (reader *Reader) FindCounter(typeId int32, keyFilter func(keyBuffer *atomic.Buffer) bool) int32 {
	var keyBuf atomic.Buffer
	for id := 0; id < reader.maxCounterID; id++ {
		metaDataOffset := int32(id) * METADATA_LENGTH
		recordStatus := reader.metaData.GetInt32Volatile(metaDataOffset)
		if recordStatus == RECORD_UNUSED {
			break
		} else if RECORD_ALLOCATED == recordStatus {
			thisTypeId := reader.metaData.GetInt32(metaDataOffset + 4)
			if thisTypeId == typeId {
				// requires Go 1.17: keyPtr := unsafe.Add(reader.metaData.Ptr(), metaDataOffset+KEY_OFFSET)
				keyPtr := unsafe.Pointer(uintptr(reader.metaData.Ptr()) + uintptr(metaDataOffset+KEY_OFFSET))
				keyBuf.Wrap(keyPtr, MAX_KEY_LENGTH)
				if keyFilter == nil || keyFilter(&keyBuf) {
					return int32(id)
				}
			}
		}
	}
	return NullCounterId
}

// GetKeyPartInt32 returns an int32 portion of the key at the specified offset
func (reader *Reader) GetKeyPartInt32(counterId int32, offset int32) (int32, error) {
	if err := reader.validateCounterIdAndOffset(counterId, offset+util.SizeOfInt32); err != nil {
		return 0, err
	}
	metaDataOffset := counterId * METADATA_LENGTH
	recordStatus := reader.metaData.GetInt32Volatile(metaDataOffset)
	if recordStatus != RECORD_ALLOCATED {
		return 0, fmt.Errorf("counterId=%d recordStatus=%d", counterId, recordStatus)
	}
	return reader.metaData.GetInt32(metaDataOffset + KEY_OFFSET + offset), nil
}

// GetKeyPartInt64 returns an int64 portion of the key at the specified offset
func (reader *Reader) GetKeyPartInt64(counterId int32, offset int32) (int64, error) {
	if err := reader.validateCounterIdAndOffset(counterId, offset+util.SizeOfInt64); err != nil {
		return 0, err
	}
	metaDataOffset := counterId * METADATA_LENGTH
	recordStatus := reader.metaData.GetInt32Volatile(metaDataOffset)
	if recordStatus != RECORD_ALLOCATED {
		return 0, fmt.Errorf("counterId=%d recordStatus=%d", counterId, recordStatus)
	}
	return reader.metaData.GetInt64(metaDataOffset + KEY_OFFSET + offset), nil
}

func (reader *Reader) validateCounterIdAndOffset(counterId int32, offset int32) error {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return fmt.Errorf("counterId=%d maxCounterId=%d", counterId, reader.maxCounterID)
	}
	if offset < 0 || offset >= MAX_KEY_LENGTH {
		return fmt.Errorf("counterId=%d offset=%d maxKeyLength=%d", counterId, offset, MAX_KEY_LENGTH)
	}
	return nil
}

func (reader *Reader) labelValue(metaDataOffset int32) string {
	labelSize := reader.metaData.GetInt32(metaDataOffset + LABEL_OFFSET)
	return string(reader.metaData.GetBytesArray(metaDataOffset+LABEL_OFFSET+4, labelSize))
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
// Returns the counterID if found otherwise NullCounterId
func (reader *Reader) FindCounterIDByRecording(recordingID int64) int32 {
	buffer := reader.MetaData()

	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		if counterState == RECORD_ALLOCATED && reader.CounterTypeID(i) == command.RecordingPosition {
			// fmt.Printf("recordingID=%d, lookup=%d\n", recordingID, buffer.GetInt64((i*METADATA_LENGTH)+KEY_OFFSET+RECORDING_ID_OFFSET))
			if recordingID == buffer.GetInt64((i*METADATA_LENGTH)+KEY_OFFSET+RECORDING_ID_OFFSET) {
				return i
			}
		} else if RECORD_UNUSED == counterState {
			break
		}
	}
	return NullCounterId
}

// FindCounterIDBySession finds the active counterID for a stream based on the sessionID
// Returns the counterID if found otherwise NullCounterId
func (reader *Reader) FindCounterIDBySession(sessionID int32) int32 {
	buffer := reader.MetaData()
	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		counterTypeID := reader.CounterTypeID(i)
		if counterState == RECORD_ALLOCATED && counterTypeID == command.RecordingPosition {
			// fmt.Printf("i=%d, counterState=%d, counterTypeID=%d, %d=%d\n", i, counterState, counterTypeID, sessionID, buffer.GetInt32((i*METADATA_LENGTH)+KEY_OFFSET+SESSION_ID_OFFSET))
			if sessionID == buffer.GetInt32((i*METADATA_LENGTH)+KEY_OFFSET+SESSION_ID_OFFSET) {
				return i
			}
		} else if RECORD_UNUSED == counterState {
			break
		}
	}
	return NullCounterId
}

// RecordingID finds the active counterID for a stream based on the sessionID
func (reader *Reader) RecordingID(counterID int32) int64 {
	buffer := reader.MetaData()
	if reader.CounterState(counterID) == RECORD_ALLOCATED && reader.CounterTypeID(counterID) == command.RecordingPosition {
		return buffer.GetInt64((counterID * METADATA_LENGTH) + KEY_OFFSET + RECORDING_ID_OFFSET)
	}
	return NULL_RECORDING_ID
}

// SourceIdentity returns the source identity for the recording
func (reader *Reader) GetSourceIdentity(counterID int32) string {
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

// AwaitRecordingCounterID waits for a specific counterID to appear and return it
// On timeout return (NullCounterId, error)
func (reader *Reader) AwaitRecordingCounterID(sessionID int32) (int32, error) {
	start := time.Now()
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 50} // FIXME: tune, // FIXME: parameterize?

	for {
		counterID := reader.FindCounterIDBySession(sessionID)
		//fmt.Printf("counterID:%d, sessionID:%d\n", counterID, sessionID)
		if counterID != NullCounterId {
			return counterID, nil
		} else {
			// Idle and check timeout
			idler.Idle(0)
			if time.Since(start) > time.Second*15 { // FIXME: tune, // FIXME: parameterize?
				return NullCounterId, fmt.Errorf("Timeout waiting for counter:%d", counterID)
			}
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

// FindAddresses returns the list of currently bound local sockets.
// Based on LocalSocketAddressStatus.java
// channelStatusId is the identity of the counter for the channel which aggregates the transports.
// returns a list of active bound local socket addresses.
func (reader *Reader) FindAddresses(channelStatusID int) (addresses [][]byte) {

	buffer := reader.MetaData()

	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		counterTypeID := reader.CounterTypeID(i)
		if counterState == RECORD_ALLOCATED && counterTypeID == LOCAL_SOCKET_ADDRESS_STATUS_TYPE_ID {

			// fmt.Printf("FindAddreses i=%d, counterState=%d, counterTypeID=%d\n", i, counterState, counterTypeID)
			keyOffset := i*METADATA_LENGTH + KEY_OFFSET
			counterChannelStatusID := int(buffer.GetInt32(keyOffset + CHANNEL_STATUS_ID_OFFSET))
			counterValue := reader.CounterValue(i)

			if channelStatusID == counterChannelStatusID && counterValue == ChannelStatusActive {
				length := buffer.GetInt32(keyOffset + LOCAL_SOCKET_ADDRESS_LENGTH_OFFSET)
				if length > 0 {
					b := make([]byte, length)
					buffer.GetBytes(keyOffset+LOCAL_SOCKET_ADDRESS_STRING_OFFSET, b)
					addresses = append(addresses, b)
					// fmt.Printf("FindAddreses: this one is %s\n", b)
				}
			}
		} else if counterState == RECORD_UNUSED {
			break
		}
	}
	return addresses
}

// FindAddress returns the first currently bound local sockets.
// Based on LocalSocketAddressStatus.java
//   channelStatus   is the value for the channel which aggregates the transports.
//   channelStatusID is the identity of the counter for the channel which aggregates the transports.
// returns a list of active bound local socket addresses.
func (reader *Reader) FindAddress(channelStatus int, channelStatusID int) (address []byte) {

	buffer := reader.MetaData()

	if channelStatus != ChannelStatusActive {
		return nil
	}

	for i, size := int32(0), reader.MaxCounterID(); i < int32(size); i++ {
		counterState := reader.CounterState(i)
		counterTypeID := reader.CounterTypeID(i)

		if counterState == RECORD_ALLOCATED && counterTypeID == LOCAL_SOCKET_ADDRESS_STATUS_TYPE_ID {

			// fmt.Printf("FindAddreses(%d, %d)) i=%d, counterState=%d, counterTypeID=%d\n", channelStatus, channelStatusID, i, counterState, counterTypeID)
			keyOffset := i*METADATA_LENGTH + KEY_OFFSET
			counterChannelStatusID := buffer.GetInt32(keyOffset + CHANNEL_STATUS_ID_OFFSET)
			counterValue := reader.CounterValue(i)

			// fmt.Printf("\tFindAddreses channelStatusID(%d) == int(counterChannelStatusID)(%d)\n", channelStatusID, int(counterChannelStatusID))
			// fmt.Printf("\tFindAddreses counterValue(%d) == int64(ChannelStatusActive)(%d)\n", counterValue, ChannelStatusActive)
			if channelStatusID == int(counterChannelStatusID) && counterValue == int64(ChannelStatusActive) {
				fmt.Printf("\t\tFindAddreses XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX grabbing value\n")
				length := buffer.GetInt32(keyOffset + LOCAL_SOCKET_ADDRESS_LENGTH_OFFSET)
				if length > 0 {
					b := make([]byte, length)
					buffer.GetBytes(keyOffset+LOCAL_SOCKET_ADDRESS_STRING_OFFSET, b)
					// fmt.Printf("FindAddreses: this one is %s\n", b)
					return b
				}
			}
		} else if counterState == RECORD_UNUSED {
			break
		}
	}
	return nil
}
