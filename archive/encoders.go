// Copyright (C) 2021 Talos, Inc.
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

package archive

import (
	"bytes"
	"github.com/lirm/aeron-go/archive/codecs"
)

// Encoders for all the protocol packets Each of these functions
// creates a []byte suitable for sending over the wire by using the
// generated encoders created using the simple-binary-encosing

// FIXME: Reentrancy options: a) giant lock, b) parameterise, c) allocate on fly

var marshaller *codecs.SbeGoMarshaller = codecs.NewSbeGoMarshaller()

func ConnectRequestPacket(responseChannel string, responseStream int32, correlationID int64) ([]byte, error) {
	var request codecs.ConnectRequest

	request.CorrelationId = correlationID
	request.Version = SemanticVersion()
	request.ResponseStreamId = responseStream
	request.ResponseChannel = []uint8(responseChannel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// deprecated
func StartRecordingRequestPacket(sessionId int64, correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, channel string) ([]byte, error) {
	return StartRecordingRequest2Packet(sessionId, correlationId, stream, sourceLocation, false, channel)
}

func StartRecordingRequest2Packet(sessionId int64, correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) ([]byte, error) {
	var request codecs.StartRecordingRequest2

	request.Channel = []uint8(channel)
	request.StreamId = stream
	request.SourceLocation = sourceLocation
	if autoStop {
		request.AutoStop = codecs.BooleanType.TRUE
	} // else FALSE by default
	request.CorrelationId = correlationId
	request.ControlSessionId = sessionId

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingRequestPacket(sessionId int64, correlationId int64, stream int32, channel string) ([]byte, error) {
	var request codecs.StopRecordingRequest

	request.ControlSessionId = sessionId
	request.CorrelationId = correlationId
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingsForUriRequestPacket(sessionId int64, correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) ([]byte, error) {
	var request codecs.ListRecordingsForUriRequest
	request.ControlSessionId = sessionId
	request.CorrelationId = correlationId
	request.FromRecordingId = fromRecordingId
	request.RecordCount = recordCount
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplayRequestPacket(sessionId int64, correlationId int64, recordingId int64, position int64, length int64, replayStream int32, replayChannel string) ([]byte, error) {
	var request codecs.ReplayRequest
	request.ControlSessionId = sessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position
	request.ReplayStreamId = replayStream
	request.ReplayChannel = []uint8(replayChannel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopReplayRequestPacket(sessionId int64, correlationId int64, replaySessionId int64) ([]byte, error) {
	var request codecs.StopReplayRequest
	request.ControlSessionId = sessionId
	request.CorrelationId = correlationId
	request.ReplaySessionId = replaySessionId

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, ArchiveDefaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
