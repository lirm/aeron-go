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
	Context     *ArchiveContext
	Publication *aeron.Publication
}

// Create a proxy with default settings
func NewProxy(context *ArchiveContext) *Proxy {
	proxy := new(Proxy)
	proxy.Context = context

	return proxy
}

// Offer to our request publication with a retry to allow time for the image establishment, some back pressure etc
func (proxy *Proxy) Offer(buffer *atomic.Buffer, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {
	start := time.Now()
	var ret int64
	for time.Since(start) < proxy.Context.Options.Timeout {
		ret = proxy.Publication.Offer(buffer, offset, length, reservedValueSupplier)
		switch ret {
		// Retry on these
		case aeron.NotConnected, aeron.BackPressured, aeron.AdminAction:
			proxy.Context.Options.IdleStrategy.Idle(0)

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

func (proxy *Proxy) ConnectRequest(responseChannel string, responseStream int32, correlationId int64) error {

	// Create a packet and send it
	bytes, err := ConnectRequestPacket(responseChannel, responseStream, correlationId)
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
	bytes, err := CloseSessionRequestPacket(proxy.Context.SessionId)
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
func (proxy *Proxy) StartRecordingRequest(correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) error {

	bytes, err := StartRecordingRequest2Packet(proxy.Context.SessionId, correlationId, stream, sourceLocation, autoStop, channel)
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
	bytes, err := StopRecordingRequestPacket(proxy.Context.SessionId, correlationId, stream, channel)
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
	bytes, err := ReplayRequestPacket(proxy.Context.SessionId, correlationId, recordingId, position, length, replayStream, replayChannel)
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
	bytes, err := StopReplayRequestPacket(proxy.Context.SessionId, correlationId, replaySessionId)
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
	bytes, err := ListRecordingsRequestPacket(proxy.Context.SessionId, correlationId, fromRecordingId, recordCount)
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

	bytes, err := ListRecordingsForUriRequestPacket(proxy.Context.SessionId, correlationId, fromRecordingId, recordCount, stream, channel)

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
	bytes, err := ListRecordingRequestPacket(proxy.Context.SessionId, correlationId, fromRecordingId)
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
	bytes, err := ExtendRecordingRequest2Packet(proxy.Context.SessionId, correlationId, recordingId, stream, sourceLocation, autoStop, channel)
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
	bytes, err := RecordingPositionRequestPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := TruncateRecordingRequestPacket(proxy.Context.SessionId, correlationId, recordingId, position)
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
	bytes, err := StopRecordingSubscriptionPacket(proxy.Context.SessionId, correlationId, subscriptionId)
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
	bytes, err := StopRecordingByIdentityPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := StopPositionPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := FindLastMatchingRecordingPacket(proxy.Context.SessionId, correlationId, minRecordingId, sessionId, stream, channel)
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
	bytes, err := ListRecordingSubscriptionsPacket(proxy.Context.SessionId, correlationId, pseudoIndex, subscriptionCount, applyStreamId, stream, channel)
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
	bytes, err := BoundedReplayPacket(proxy.Context.SessionId, correlationId, recordingId, position, length, limitCounterId, replayStream, replayChannel)
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
	bytes, err := StopAllReplaysPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := CatalogHeaderPacket(version, length, nextRecordingId, alignment)
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
	bytes, err := ReplicateRequestPacket(proxy.Context.SessionId, correlationId, srcRecordingId, dstRecordingId, srcControlStreamId, srcControlChannel, liveDestination)
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
	bytes, err := ReplicateRequest2Packet(proxy.Context.SessionId, correlationId, srcRecordingId, dstRecordingId, stopPosition, channelTagId, srcControlStreamId, srcControlChannel, liveDestination, replicationChannel)
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
	bytes, err := StopReplicationRequestPacket(proxy.Context.SessionId, correlationId, replicationId)
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
	bytes, err := StartPositionRequestPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := DetachSegmentsRequestPacket(proxy.Context.SessionId, correlationId, recordingId, newStartPosition)
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
	bytes, err := DeleteDetachedSegmentsRequestPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := PurgeSegmentsRequestPacket(proxy.Context.SessionId, correlationId, recordingId, newStartPosition)
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
	bytes, err := AttachSegmentsRequestPacket(proxy.Context.SessionId, correlationId, recordingId)
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
	bytes, err := MigrateSegmentsRequestPacket(proxy.Context.SessionId, correlationId, srcRecordingId, destRecordingId)
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
	bytes, err := TaggedReplicateRequestPacket(proxy.Context.SessionId, correlationId, srcRecordingId, dstRecordingId, channelTagId, subscriptionTagId, srcControlStreamId, srcControlChannel, liveDestination)
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
	bytes, err := PurgeRecordingRequestPacket(proxy.Context.SessionId, correlationId, replaySessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}
