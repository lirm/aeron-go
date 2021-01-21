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

// FIXME: optimize: decide where buffer allocation happens
// FIXME: reentrancy?
// FIXME: optimize: common header?

// Create and marshal a StartRecordingRequest
func (archive *Archive) StartRecordingRequestPacket(channel string, stream int32) (*bytes.Buffer, error) {
	var request codecs.StartRecordingRequest

	request.ControlSessionId = controlSessionID.Inc()
	request.CorrelationId = archive.aeron.NextCorrelationID()
	request.StreamId = stream
	request.SourceLocation = codecs.SourceLocation.LOCAL
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(archive.Control.marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(archive.Control.marshaller, buffer, archive.Control.rangeCheck); err != nil {
		return nil, err
	}

	return buffer, nil
}

// Create and marshal a connect request
func (archive *Archive) ConnectRequestPacket(responseChannel string, responseStream int32) (*bytes.Buffer, error) {
	var request codecs.ConnectRequest

	request.CorrelationId = archive.aeron.NextCorrelationID()
	request.Version = codecs.SemanticVersion()
	request.ResponseStreamId = responseStream
	request.ResponseChannel = []uint8(responseChannel)

	// Marshal it
	header := codecs.MessageHeader{request.SbeBlockLength(), request.SbeTemplateId(), request.SbeSchemaId(), request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(archive.Control.marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(archive.Control.marshaller, buffer, archive.Control.rangeCheck); err != nil {
		return nil, err
	}

	return buffer, nil
}
