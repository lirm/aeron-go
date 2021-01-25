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

// FIXME: Reentrancy options: a) giant lock, b) parameterise, c) allocate on fly
var marshaller *codecs.SbeGoMarshaller = codecs.NewSbeGoMarshaller()

// Create and marshal a connect request
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
	if err := request.Encode(marshaller, buffer, defaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Create and marshal a StartRecordingRequest
func StartRecordingRequestPacket(channel string, stream int32, correlationID int64, sourceLocation codecs.SourceLocationEnum) ([]byte, error) {
	var request codecs.StartRecordingRequest

	request.ControlSessionId = controlSessionID.Inc()
	request.CorrelationId = correlationID
	request.StreamId = stream
	request.SourceLocation = sourceLocation
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, defaults.RangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}