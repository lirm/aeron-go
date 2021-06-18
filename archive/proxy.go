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
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
	"time"
)

// Proxy class for encapsulating encoding and sending of control protocol messages to an archive
type Proxy struct {
	Publication *aeron.Publication
	archive     *Archive                // link to parent
	marshaller  *codecs.SbeGoMarshaller // currently shared as we're not reentrant (but could be here)
}

// Offer to our request publication with a retry to allow time for the image establishment, some back pressure etc
func (proxy *Proxy) Offer(buffer *atomic.Buffer, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {
	start := time.Now()
	var ret int64
	for time.Since(start) < proxy.archive.Options.Timeout {
		ret = proxy.Publication.Offer(buffer, offset, length, reservedValueSupplier)
		switch ret {
		// Retry on these
		case aeron.NotConnected, aeron.BackPressured, aeron.AdminAction:
			proxy.archive.Options.IdleStrategy.Idle(0)

		// Fail or succeed on other values
		default:
			return ret
		}
	}

	// Give up, returning the last failure
	logger.Debugf("Proxy.Offer timing out [%d]", ret)
	return ret

}

// From here we have all the functions that create a data packet and send it on the
// publication. Responses will be processed on the control

func (proxy *Proxy) ConnectRequest(correlationId int64, responseStream int32, responseChannel string) error {

	// Create a packet and send it
	bytes, err := codecs.ConnectRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, correlationId, responseStream, responseChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) CloseSessionRequest() error {
	// Create a packet and send it
	bytes, err := codecs.CloseSessionRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// StartRecordingRequest
// Uses the more recent protocol addition StartdRecordingRequest2 which added autoStop
func (proxy *Proxy) StartRecordingRequest(correlationId int64, stream int32, isLocal bool, autoStop bool, channel string) error {

	bytes, err := codecs.StartRecordingRequest2Packet(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, stream, isLocal, autoStop, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecordingRequest(correlationId int64, stream int32, channel string) error {
	// Create a packet and send it
	bytes, err := codecs.StopRecordingRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ReplayRequest(correlationId int64, recordingId int64, position int64, length int64, replayChannel string, replayStream int32) error {

	// Create a packet and send it
	bytes, err := codecs.ReplayRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, position, length, replayStream, replayChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopReplayRequest(correlationId int64, replaySessionId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StopReplayRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, replaySessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// Lists up to recordCount recordings starting at fromRecordingId
func (proxy *Proxy) ListRecordingsRequest(correlationId int64, fromRecordingId int64, recordCount int32) error {
	// Create a packet and send it
	bytes, err := codecs.ListRecordingsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, fromRecordingId, recordCount)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// ListRecordingsForUri
// Lists up to recordCount recordings that match the channel and stream
func (proxy *Proxy) ListRecordingsForUriRequest(correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) error {

	bytes, err := codecs.ListRecordingsForUriRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, fromRecordingId, recordCount, stream, channel)

	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// Retrieves a recording descriptor for a specific recordingId
func (proxy *Proxy) ListRecordingRequest(correlationId int64, fromRecordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.ListRecordingRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, fromRecordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// ExtendRecording
// Uses the more recent protocol addition ExtendRecordingRequest2 which added autoStop
func (proxy *Proxy) ExtendRecordingRequest(correlationId int64, recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) error {
	// Create a packet and send it
	bytes, err := codecs.ExtendRecordingRequest2Packet(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, stream, sourceLocation, autoStop, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// RecordingPosition
func (proxy *Proxy) RecordingPositionRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.RecordingPositionRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) TruncateRecordingRequest(correlationId int64, recordingId int64, position int64) error {
	// Create a packet and send it
	bytes, err := codecs.TruncateRecordingRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, position)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecordingSubscriptionRequest(correlationId int64, subscriptionId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StopRecordingSubscriptionPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, subscriptionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecordingByIdentityRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StopRecordingByIdentityPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopPositionRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StopPositionPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) FindLastMatchingRecordingRequest(correlationId int64, minRecordingId int64, sessionId int32, stream int32, channel string) error {

	// Create a packet and send it
	bytes, err := codecs.FindLastMatchingRecordingPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, minRecordingId, sessionId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ListRecordingSubscriptionsRequest(correlationId int64, pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channel string) error {

	// Create a packet and send it
	bytes, err := codecs.ListRecordingSubscriptionsPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, pseudoIndex, subscriptionCount, applyStreamId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) BoundedReplayRequest(correlationId int64, recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) error {

	// Create a packet and send it
	bytes, err := codecs.BoundedReplayPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, position, length, limitCounterId, replayStream, replayChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopAllReplaysRequest(correlationId int64, recordingId int64) error {

	// Create a packet and send it
	bytes, err := codecs.StopAllReplaysPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) CatalogHeaderRequest(version int32, length int32, nextRecordingId int64, alignment int32) error {

	// Create a packet and send it
	bytes, err := codecs.CatalogHeaderPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, version, length, nextRecordingId, alignment)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ReplicateRequest(correlationId int64, srcRecordingId int64, dstRecordingId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) error {
	// Create a packet and send it
	bytes, err := codecs.ReplicateRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, srcRecordingId, dstRecordingId, srcControlStreamId, srcControlChannel, liveDestination)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ReplicateRequest2(correlationId int64, srcRecordingId int64, dstRecordingId int64, stopPosition int64, channelTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string, replicationChannel string) error {
	// Create a packet and send it
	bytes, err := codecs.ReplicateRequest2Packet(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, srcRecordingId, dstRecordingId, stopPosition, channelTagId, srcControlStreamId, srcControlChannel, liveDestination, replicationChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopReplicationRequest(correlationId int64, replicationId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StopReplicationRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, replicationId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StartPositionRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.StartPositionRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) DetachSegmentsRequest(correlationId int64, recordingId int64, newStartPosition int64) error {
	// Create a packet and send it
	bytes, err := codecs.DetachSegmentsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, newStartPosition)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) DeleteDetachedSegmentsRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.DeleteDetachedSegmentsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) PurgeSegmentsRequest(correlationId int64, recordingId int64, newStartPosition int64) error {
	// Create a packet and send it
	bytes, err := codecs.PurgeSegmentsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId, newStartPosition)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}
func (proxy *Proxy) AttachSegmentsRequest(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.AttachSegmentsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) MigrateSegmentsRequest(correlationId int64, srcRecordingId int64, destRecordingId int64) error {
	// Create a packet and send it
	bytes, err := codecs.MigrateSegmentsRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, srcRecordingId, destRecordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) AuthConnectRequest(correlationId int64, responseStream int32, responseChannel string, encodedCredentials []uint8) error {

	// Create a packet and send it
	bytes, err := codecs.AuthConnectRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, correlationId, responseStream, responseChannel, encodedCredentials)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ChallengeResponse(correlationId int64, encodedCredentials []uint8) error {
	// Create a packet and send it
	bytes, err := codecs.ChallengeResponsePacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, encodedCredentials)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) KeepAliveRequest(correlationId int64) error {
	// Create a packet and send it
	bytes, err := codecs.KeepAliveRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) TaggedReplicateRequest(correlationId int64, srcRecordingId int64, dstRecordingId int64, channelTagId int64, subscriptionTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) error {
	// Create a packet and send it
	bytes, err := codecs.TaggedReplicateRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, srcRecordingId, dstRecordingId, channelTagId, subscriptionTagId, srcControlStreamId, srcControlChannel, liveDestination)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) PurgeRecordingRequest(correlationId int64, replaySessionId int64) error {
	// Create a packet and send it
	bytes, err := codecs.PurgeRecordingRequestPacket(proxy.marshaller, proxy.archive.Options.RangeChecking, proxy.archive.SessionId, correlationId, replaySessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}
