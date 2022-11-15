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

package archive

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
)

const RECORDING_POSITION_COUNTER_TYPE_ID = 100
const NULL_RECORDING_ID = -1

const (
	RECORDING_ID_OFFSET           = 0
	SESSION_ID_OFFSET             = RECORDING_ID_OFFSET + 8
	SOURCE_IDENTITY_LENGTH_OFFSET = SESSION_ID_OFFSET + 4
	SOURCE_IDENTITY_OFFSET        = SOURCE_IDENTITY_LENGTH_OFFSET + 4
)

// FindCounterIdByRecording finds the active counter id for a stream based on the recording id.
// Returns the counter id if found otherwise NullCounterId
func FindCounterIdByRecording(countersReader *counters.Reader, recordingId int64) (int32, error) {
	return countersReader.FindCounter(RECORDING_POSITION_COUNTER_TYPE_ID, func(keyBuffer *atomic.Buffer) bool {
		return keyBuffer.GetInt64(RECORDING_ID_OFFSET) == recordingId
	})
}

// FindCounterIdBySession finds the active counterID for a stream based on the session id.
// Returns the counter id if found otherwise NullCounterId
func FindCounterIdBySession(countersReader *counters.Reader, sessionId int32) (int32, error) {
	return countersReader.FindCounter(RECORDING_POSITION_COUNTER_TYPE_ID, func(keyBuffer *atomic.Buffer) bool {
		return keyBuffer.GetInt32(SESSION_ID_OFFSET) == sessionId
	})
}

func GetRecordingId(countersReader *counters.Reader, counterId int32) int64 {
	if countersReader.GetCounterTypeId(counterId) != RECORDING_POSITION_COUNTER_TYPE_ID {
		return NULL_RECORDING_ID
	}
	recordingId, err := countersReader.GetKeyPartInt64(counterId, RECORDING_ID_OFFSET)
	if err != nil {
		return NULL_RECORDING_ID
	}
	return recordingId
}

// GetSourceIdentity returns the source identity for the recording
func GetSourceIdentity(countersReader *counters.Reader, counterId int32) string {
	if countersReader.GetCounterTypeId(counterId) != RECORDING_POSITION_COUNTER_TYPE_ID {
		return ""
	}
	identity, err := countersReader.GetKeyPartString(counterId, SOURCE_IDENTITY_LENGTH_OFFSET)
	if err != nil {
		return ""
	}
	return identity
}

// IsRecordingActive tells us if the recording counter is still active?
func IsRecordingActive(countersReader *counters.Reader, counterId int32, recordingId int64) bool {
	recId, err := countersReader.GetKeyPartInt64(counterId, RECORDING_ID_OFFSET)
	return err == nil && recId == recordingId
}
