/*
Copyright 2016-2018 Stanislav Liberman

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

package counters

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
)

const COUNTER_LENGTH = util.CacheLineLength * 2

const FULL_LABEL_LENGTH = util.CacheLineLength * 6
const KEY_OFFSET = 16
const LABEL_OFFSET = util.CacheLineLength * 2
const METADATA_LENGTH = LABEL_OFFSET + FULL_LABEL_LENGTH
const MAX_KEY_LENGTH = (util.CacheLineLength * 2) - (util.SizeOfInt32 * 2) - util.SizeOfInt64

const RECORD_UNUSED int32 = 0
const RECORD_ALLOCATED int32 = 1

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
			typeId := reader.metaData.GetInt32(i + 4)
			label := reader.labelValue(i)

			// TODO Get the key buffer

			value := reader.values.GetInt64Volatile(id * COUNTER_LENGTH)

			//fmt.Printf("Reading at offset %d; counterState=%d; typeId=%d\n", i, recordStatus, typeId)

			cb(Counter{id, typeId, value, label})
		}
		id++
	}
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
func (reader *Reader) GetKeyInt64(counterId int32, offset int32) (int64, error) {
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
