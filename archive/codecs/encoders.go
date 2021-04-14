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

// Package codecs contains the archive protocol packet encoding and decoding
package codecs

import (
	"bytes"
)

// Encoders for all the protocol packets Each of these functions
// creates a []byte suitable for sending over the wire by using the
// generated encoders created using simple-binary-encoding.
//
// All the packet specificationss are defined in the aeron-archive protocol
// maintained at:
// http://github.com/real-logic/aeron/blob/master/aeron-archive/src/main/resources/archive/aeron-archive-codecs.xml)
//
// The codecs are generated from that specification using Simple
// Binary Encoding (SBE) from https://github.com/real-logic/simple-binary-encoding
func ConnectRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, correlationID int64, responseStream int32, responseChannel string) ([]byte, error) {
	var request ConnectRequest

	request.CorrelationId = correlationID
	request.Version = SemanticVersion()
	request.ResponseStreamId = responseStream
	request.ResponseChannel = []uint8(responseChannel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CloseSessionRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64) ([]byte, error) {
	var request CloseSessionRequest

	request.ControlSessionId = controlSessionId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// deprecated
func StartRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, stream int32, isLocal bool, channel string) ([]byte, error) {
	return StartRecordingRequest2Packet(marshaller, rangeChecking, controlSessionId, correlationId, stream, isLocal, false, channel)
}

func StartRecordingRequest2Packet(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, stream int32, isLocal bool, autoStop bool, channel string) ([]byte, error) {
	var request StartRecordingRequest2

	request.Channel = []uint8(channel)
	request.StreamId = stream
	if isLocal {
		request.SourceLocation = SourceLocation.LOCAL
	} else {
		request.SourceLocation = SourceLocation.REMOTE
	}
	if autoStop {
		request.AutoStop = BooleanType.TRUE
	} // else FALSE by default
	request.CorrelationId = correlationId
	request.ControlSessionId = controlSessionId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, stream int32, channel string) ([]byte, error) {
	var request StopRecordingRequest

	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplayRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, position int64, length int64, replayStream int32, replayChannel string) ([]byte, error) {
	var request ReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position
	request.Length = length
	request.ReplayStreamId = replayStream
	request.ReplayChannel = []uint8(replayChannel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopReplayRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, replaySessionId int64) ([]byte, error) {
	var request StopReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.ReplaySessionId = replaySessionId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, fromRecordingId int64, recordCount int32) ([]byte, error) {
	var request ListRecordingsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.FromRecordingId = fromRecordingId
	request.RecordCount = recordCount

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingsForUriRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) ([]byte, error) {
	var request ListRecordingsForUriRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.FromRecordingId = fromRecordingId
	request.RecordCount = recordCount
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request ListRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// deprecated
func ExtendRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, stream int32, sourceLocation SourceLocationEnum, channel string) ([]byte, error) {
	return ExtendRecordingRequest2Packet(marshaller, rangeChecking, controlSessionId, correlationId, recordingId, stream, sourceLocation, false, channel)
}

func ExtendRecordingRequest2Packet(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, stream int32, sourceLocation SourceLocationEnum, autoStop bool, channel string) ([]byte, error) {
	var request ExtendRecordingRequest2
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.StreamId = stream
	request.SourceLocation = sourceLocation
	if autoStop {
		request.AutoStop = BooleanType.TRUE
	} // else FALSE by default
	request.Channel = []uint8(channel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func RecordingPositionRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request RecordingPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func TruncateRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, position int64) ([]byte, error) {
	var request TruncateRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingSubscriptionPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, subscriptionId int64) ([]byte, error) {
	var request StopRecordingSubscriptionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SubscriptionId = subscriptionId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingByIdentityPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request StopRecordingByIdentityRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopPositionPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request StopPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func FindLastMatchingRecordingPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, minRecordingId int64, sessionId int32, stream int32, channel string) ([]byte, error) {
	var request FindLastMatchingRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.MinRecordingId = minRecordingId
	request.SessionId = sessionId
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingSubscriptionsPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channel string) ([]byte, error) {
	var request ListRecordingSubscriptionsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.PseudoIndex = pseudoIndex
	request.SubscriptionCount = subscriptionCount
	if applyStreamId {
		request.ApplyStreamId = BooleanType.TRUE
	} // else FALSE by default
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func BoundedReplayPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) ([]byte, error) {
	var request BoundedReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position
	request.Length = length
	request.LimitCounterId = limitCounterId
	request.ReplayStreamId = replayStream
	request.ReplayChannel = []uint8(replayChannel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopAllReplaysPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request StopAllReplaysRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CatalogHeaderPacket(marshaller *SbeGoMarshaller, rangeChecking bool, version int32, length int32, nextRecordingId int64, alignment int32) ([]byte, error) {
	var request CatalogHeader
	request.Version = version
	request.Length = length
	request.NextRecordingId = nextRecordingId
	request.Alignment = alignment

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplicateRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, srcRecordingId int64, dstRecordingId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) ([]byte, error) {
	var request ReplicateRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = dstRecordingId
	request.SrcControlStreamId = srcControlStreamId
	request.SrcControlChannel = []uint8(srcControlChannel)
	request.LiveDestination = []uint8(liveDestination)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopReplicationRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, replicationId int64) ([]byte, error) {
	var request StopReplicationRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.ReplicationId = replicationId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StartPositionRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request StartPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DetachSegmentsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, newStartPosition int64) ([]byte, error) {
	var request DetachSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.NewStartPosition = newStartPosition

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DeleteDetachedSegmentsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request DeleteDetachedSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PurgeSegmentsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64, newStartPosition int64) ([]byte, error) {
	var request PurgeSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.NewStartPosition = newStartPosition

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func AttachSegmentsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request AttachSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func AuthConnectRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, correlationID int64, responseStream int32, responseChannel string, encodedCredentials []uint8) ([]byte, error) {
	var request AuthConnectRequest

	request.CorrelationId = correlationID
	request.Version = SemanticVersion()
	request.ResponseStreamId = responseStream
	request.ResponseChannel = []uint8(responseChannel)
	request.EncodedCredentials = encodedCredentials

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func MigrateSegmentsRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, srcRecordingId int64, destRecordingId int64) ([]byte, error) {
	var request MigrateSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = destRecordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func KeepAliveRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64) ([]byte, error) {
	var request KeepAliveRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func TaggedReplicateRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, srcRecordingId int64, dstRecordingId int64, channelTagId int64, subscriptionTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) ([]byte, error) {
	var request TaggedReplicateRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = dstRecordingId
	request.ChannelTagId = channelTagId
	request.SubscriptionTagId = subscriptionTagId
	request.SrcControlStreamId = srcControlStreamId
	request.SrcControlChannel = []uint8(srcControlChannel)
	request.LiveDestination = []uint8(liveDestination)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplicateRequest2Packet(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, srcRecordingId int64, dstRecordingId int64, stopPosition int64, channelTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string, replicationChannel string) ([]byte, error) {
	var request ReplicateRequest2
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = dstRecordingId
	request.StopPosition = stopPosition
	request.ChannelTagId = channelTagId
	request.SrcControlStreamId = srcControlStreamId
	request.SrcControlChannel = []uint8(srcControlChannel)
	request.LiveDestination = []uint8(liveDestination)
	request.ReplicationChannel = []uint8(replicationChannel)

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PurgeRecordingRequestPacket(marshaller *SbeGoMarshaller, rangeChecking bool, controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request PurgeRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
