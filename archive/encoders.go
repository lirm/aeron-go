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
// FIXME: rangeCheckign needs to be global in the fragment assemblers but not here. Parameterize!

var marshaller *codecs.SbeGoMarshaller = codecs.NewSbeGoMarshaller()

func ConnectRequestPacket(responseChannel string, responseStream int32, correlationID int64) ([]byte, error) {
	var request codecs.ConnectRequest

	request.CorrelationId = correlationID
	request.Version = SemanticVersion()
	request.ResponseStreamId = responseStream
	request.ResponseChannel = []uint8(responseChannel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CloseSessionRequestPacket(controlSessionId int64) ([]byte, error) {
	var request codecs.CloseSessionRequest

	request.ControlSessionId = controlSessionId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
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
func StartRecordingRequestPacket(controlSessionId int64, correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, channel string) ([]byte, error) {
	return StartRecordingRequest2Packet(controlSessionId, correlationId, stream, sourceLocation, false, channel)
}

func StartRecordingRequest2Packet(controlSessionId int64, correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) ([]byte, error) {
	var request codecs.StartRecordingRequest2

	request.Channel = []uint8(channel)
	request.StreamId = stream
	request.SourceLocation = sourceLocation
	if autoStop {
		request.AutoStop = codecs.BooleanType.TRUE
	} // else FALSE by default
	request.CorrelationId = correlationId
	request.ControlSessionId = controlSessionId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingRequestPacket(controlSessionId int64, correlationId int64, stream int32, channel string) ([]byte, error) {
	var request codecs.StopRecordingRequest

	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplayRequestPacket(controlSessionId int64, correlationId int64, recordingId int64, position int64, length int64, replayStream int32, replayChannel string) ([]byte, error) {
	var request codecs.ReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position
	request.Length = length
	request.ReplayStreamId = replayStream
	request.ReplayChannel = []uint8(replayChannel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopReplayRequestPacket(controlSessionId int64, correlationId int64, replaySessionId int64) ([]byte, error) {
	var request codecs.StopReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.ReplaySessionId = replaySessionId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingsRequestPacket(controlSessionId int64, correlationId int64, fromRecordingId int64, recordCount int32) ([]byte, error) {
	var request codecs.ListRecordingsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.FromRecordingId = fromRecordingId
	request.RecordCount = recordCount

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingsForUriRequestPacket(controlSessionId int64, correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) ([]byte, error) {
	var request codecs.ListRecordingsForUriRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.FromRecordingId = fromRecordingId
	request.RecordCount = recordCount
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.ListRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
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
func ExtendRecordingRequestPacket(controlSessionId int64, correlationId int64, recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, channel string) ([]byte, error) {
	return ExtendRecordingRequest2Packet(controlSessionId, correlationId, recordingId, stream, sourceLocation, false, channel)
}

func ExtendRecordingRequest2Packet(controlSessionId int64, correlationId int64, recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) ([]byte, error) {
	var request codecs.ExtendRecordingRequest2
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.StreamId = stream
	request.SourceLocation = sourceLocation
	if autoStop {
		request.AutoStop = codecs.BooleanType.TRUE
	} // else FALSE by default
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func RecordingPositionRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.RecordingPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func TruncateRecordingRequestPacket(controlSessionId int64, correlationId int64, recordingId int64, position int64) ([]byte, error) {
	var request codecs.TruncateRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingSubscriptionPacket(controlSessionId int64, correlationId int64, subscriptionId int64) ([]byte, error) {
	var request codecs.StopRecordingSubscriptionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SubscriptionId = subscriptionId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopRecordingByIdentityPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.StopRecordingByIdentityRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopPositionPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.StopPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func FindLastMatchingRecordingPacket(controlSessionId int64, correlationId int64, minRecordingId int64, sessionId int32, stream int32, channel string) ([]byte, error) {
	var request codecs.FindLastMatchingRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.MinRecordingId = minRecordingId
	request.SessionId = sessionId
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListRecordingSubscriptionsPacket(controlSessionId int64, correlationId int64, pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channel string) ([]byte, error) {
	var request codecs.ListRecordingSubscriptionsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.PseudoIndex = pseudoIndex
	request.SubscriptionCount = subscriptionCount
	if applyStreamId {
		request.ApplyStreamId = codecs.BooleanType.TRUE
	} // else FALSE by default
	request.StreamId = stream
	request.Channel = []uint8(channel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func BoundedReplayPacket(controlSessionId int64, correlationId int64, recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) ([]byte, error) {
	var request codecs.BoundedReplayRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.Position = position
	request.Length = length
	request.LimitCounterId = limitCounterId
	request.ReplayStreamId = replayStream
	request.ReplayChannel = []uint8(replayChannel)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopAllReplaysPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.StopAllReplaysRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CatalogHeaderPacket(version int32, length int32, nextRecordingId int64, alignment int32) ([]byte, error) {
	var request codecs.CatalogHeader
	request.Version = version
	request.Length = length
	request.NextRecordingId = nextRecordingId
	request.Alignment = alignment

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ReplicateRequestPacket(controlSessionId int64, correlationId int64, srcRecordingId int64, dstRecordingId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) ([]byte, error) {
	var request codecs.ReplicateRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = dstRecordingId
	request.SrcControlStreamId = srcControlStreamId
	request.SrcControlChannel = []uint8(srcControlChannel)
	request.LiveDestination = []uint8(liveDestination)

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StopReplicationRequestPacket(controlSessionId int64, correlationId int64, replicationId int64) ([]byte, error) {
	var request codecs.StopReplicationRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.ReplicationId = replicationId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func StartPositionRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.StartPositionRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DetachSegmentsRequestPacket(controlSessionId int64, correlationId int64, recordingId int64, newStartPosition int64) ([]byte, error) {
	var request codecs.DetachSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.NewStartPosition = newStartPosition

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DeleteDetachedSegmentsRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.DeleteDetachedSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PurgeSegmentsRequestPacket(controlSessionId int64, correlationId int64, recordingId int64, newStartPosition int64) ([]byte, error) {
	var request codecs.PurgeSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId
	request.NewStartPosition = newStartPosition

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func AttachSegmentsRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.AttachSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func MigrateSegmentsRequestPacket(controlSessionId int64, correlationId int64, srcRecordingId int64, destRecordingId int64) ([]byte, error) {
	var request codecs.MigrateSegmentsRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.SrcRecordingId = srcRecordingId
	request.DstRecordingId = destRecordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PurgeRecordingRequestPacket(controlSessionId int64, correlationId int64, recordingId int64) ([]byte, error) {
	var request codecs.PurgeRecordingRequest
	request.ControlSessionId = controlSessionId
	request.CorrelationId = correlationId
	request.RecordingId = recordingId

	// Marshal it
	header := codecs.MessageHeader{BlockLength: request.SbeBlockLength(), TemplateId: request.SbeTemplateId(), SchemaId: request.SbeSchemaId(), Version: request.SbeSchemaVersion()}
	buffer := new(bytes.Buffer)
	if err := header.Encode(marshaller, buffer); err != nil {
		return nil, err
	}
	if err := request.Encode(marshaller, buffer, rangeChecking); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// FIXME: Todo (although some are incoming)

/*

RecordingDescriptorHeader
RecordingDescriptor
RecordingSubscriptionDescriptor

[rcv] RecordingSignalEvent

AuthConnectRequest
Challenge
ChallengeResponse
KeepAliveRequest
TaggedReplicateRequest
* StartRecordingRequest2
// StopRecordingByIdentityRequest
[rcv] RecordingStarted
[rcv] RecordingProgress
[rcv] RecordingStopped
* PurgeRecordingRequest

*/
