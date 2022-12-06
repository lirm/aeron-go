// Copyright 2022 Steven Stern
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package counters

import (
	"fmt"
)

// Aeron.NullValue, but using it directly creates a circular dependency
const nullValue = -1

type ReadableCounter struct {
	Reader         *Reader
	RegistrationId int64
	CounterId      int32
	metaDataOffset int32
	valueOffset    int32
	isClosed       bool
}

func NewReadableRegisteredCounter(reader *Reader, registrationId int64, counterId int32) (*ReadableCounter, error) {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return nil, fmt.Errorf("counterId=%d maxCounterId=%d", counterId, reader.maxCounterID)
	}
	metaDataOffset := counterId * MetadataLength
	valueOffset := counterId * CounterLength
	return &ReadableCounter{reader, registrationId, counterId, metaDataOffset, valueOffset, false}, nil
}

func NewReadableCounter(reader *Reader, counterId int32) (*ReadableCounter, error) {
	return NewReadableRegisteredCounter(reader, nullValue, counterId)
}

func (rc *ReadableCounter) State() int32 {
	return rc.Reader.metaData.GetInt32Volatile(rc.metaDataOffset)
}

func (rc *ReadableCounter) Label() string {
	return rc.Reader.labelValue(rc.metaDataOffset)
}

func (rc *ReadableCounter) Get() int64 {
	return rc.Reader.values.GetInt64Volatile(rc.valueOffset)
}

func (rc *ReadableCounter) GetWeak() int64 {
	return rc.Reader.values.GetInt64(rc.valueOffset)
}

func (rc *ReadableCounter) Close() {
	rc.isClosed = true
}

func (rc *ReadableCounter) IsClosed() bool {
	return rc.isClosed
}
