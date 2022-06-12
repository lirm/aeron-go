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

package counters

import (
	"fmt"
	"unsafe"

	"github.com/lirm/aeron-go/aeron/atomic"
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

			// fmt.Printf("Reading at offset %d; counterState=%d; typeId=%d\n", i, recordStatus, typeId)

			cb(Counter{id, typeId, value, label})
		}
		id++
	}
}

func (reader *Reader) ScanForType(typeId int32, callback func(counterId int32, keyBuffer *atomic.Buffer) bool) {
	var keyBuf atomic.Buffer
	for id := 0; id < reader.maxCounterID; id++ {
		counterId := int32(id)
		metaDataOffset := counterId * METADATA_LENGTH
		recordStatus := reader.metaData.GetInt32Volatile(metaDataOffset)
		if recordStatus == RECORD_UNUSED {
			break
		} else if RECORD_ALLOCATED == recordStatus {
			thisTypeId := reader.metaData.GetInt32(metaDataOffset + 4)
			if thisTypeId == typeId {
				// requires Go 1.17: keyPtr := unsafe.Add(reader.metaData.Ptr(), metaDataOffset+KEY_OFFSET)
				keyPtr := unsafe.Pointer(uintptr(reader.metaData.Ptr()) + uintptr(metaDataOffset+KEY_OFFSET))
				keyBuf.Wrap(keyPtr, MAX_KEY_LENGTH)
				if !callback(counterId, &keyBuf) {
					break
				}
			}
		}
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

// GetKeyPartString returns a string portion of the key, assuming the string is prefixed by its length
// (as an 32-bit int) at the specified offset.
func (reader *Reader) GetKeyPartString(counterId int32, offset int32) (string, error) {
	if err := reader.validateCounterIdAndOffset(counterId, offset+util.SizeOfInt32); err != nil {
		return "", err
	}
	metaDataOffset := counterId * METADATA_LENGTH
	recordStatus := reader.metaData.GetInt32Volatile(metaDataOffset)
	if recordStatus != RECORD_ALLOCATED {
		return "", fmt.Errorf("counterId=%d recordStatus=%d", counterId, recordStatus)
	}
	lengthOffset := metaDataOffset + KEY_OFFSET + offset
	length := reader.metaData.GetInt32(lengthOffset)
	if length < 0 || (offset+length) > MAX_KEY_LENGTH {
		return "", fmt.Errorf("counterId=%d offset=%d length=%d", counterId, offset, length)
	}
	return string(reader.metaData.GetBytesArray(lengthOffset+4, length)), nil
}

// GetCounterValue returns the value of the given counter id (as a volatile read).
func (reader *Reader) GetCounterValue(counterId int32) int64 {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return 0
	}
	return reader.values.GetInt64Volatile(counterId * COUNTER_LENGTH)
}

// GetCounterTypeId returns the type id for a counter.
func (reader *Reader) GetCounterTypeId(counterId int32) int32 {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return -1
	}
	return reader.metaData.GetInt32(counterId*METADATA_LENGTH + TYPE_ID_OFFSET)
}

func (reader *Reader) IsCounterAllocated(counterId int32) bool {
	return counterId >= 0 && counterId < int32(reader.maxCounterID) &&
		reader.metaData.GetInt32Volatile(counterId*METADATA_LENGTH) == RECORD_ALLOCATED
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
