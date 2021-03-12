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
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
	"time"
)

// Proxy class for encapsulating encoding and sending of control protocol messages to an archive
type Proxy struct {
	Publication  *aeron.Publication
	Marshaller   *codecs.SbeGoMarshaller
	SessionId    int64
	IdleStrategy idlestrategy.Idler
	Timeout      time.Duration
	Retries      int
}

// Create a proxy with default settings
func NewProxy(publication *aeron.Publication, idleStrategy idlestrategy.Idler, sessionId int64) *Proxy {
	proxy := new(Proxy)
	proxy.Publication = publication
	proxy.IdleStrategy = idleStrategy
	proxy.Marshaller = codecs.NewSbeGoMarshaller()
	proxy.Timeout = ArchiveDefaults.ControlTimeout
	proxy.Retries = ArchiveDefaults.ControlRetries
	proxy.SessionId = sessionId

	return proxy
}

// Registered for logging
func ArchiveProxyNewPublicationHandler(channel string, stream int32, session int32, regId int64) {
	logger.Debugf("ArchiveProxyNewPublicationHandler channel:%s stream:%d, session:%d, regId:%d", channel, stream, session, regId)
}

// Offer to our request publication with a retry to allow time for the image establishment, some back pressure etc
func (proxy *Proxy) Offer(buffer *atomic.Buffer, offset int32, length int32, reservedValueSupplier term.ReservedValueSupplier) int64 {
	start := time.Now()
	var ret int64
	for time.Since(start) < ArchiveDefaults.ControlTimeout {
		ret = proxy.Publication.Offer(buffer, offset, length, reservedValueSupplier)
		switch ret {
		// Retry on these
		case aeron.NotConnected, aeron.BackPressured, aeron.AdminAction:
			proxy.IdleStrategy.Idle(0)

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
// publication

func (proxy *Proxy) Connect(responseChannel string, responseStream int32, correlationId int64) error {

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

func (proxy *Proxy) CloseSession() error {
	// Create a packet and send it
	bytes, err := CloseSessionRequestPacket(proxy.SessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// StartRecording
// Uses the more recent protocol addition StartdRecordingRequest2 which added autoStop
func (proxy *Proxy) StartRecording(correlationId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) error {

	bytes, err := StartRecordingRequest2Packet(proxy.SessionId, correlationId, stream, sourceLocation, autoStop, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecording(correlationId int64, stream int32, channel string) error {
	// Create a packet and send it
	bytes, err := StopRecordingRequestPacket(proxy.SessionId, correlationId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) Replay(correlationId int64, recordingId int64, position int64, length int64, replayChannel string, replayStream int32) error {

	// Create a packet and send it
	bytes, err := ReplayRequestPacket(proxy.SessionId, correlationId, recordingId, position, length, replayStream, replayChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopReplay(correlationId int64, replaySessionId int64) error {
	// Create a packet and send it
	bytes, err := StopReplayRequestPacket(proxy.SessionId, correlationId, replaySessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ListRecordings(correlationId int64, fromRecordingId int64, recordCount int32) error {
	// Create a packet and send it
	bytes, err := ListRecordingsRequestPacket(proxy.SessionId, correlationId, fromRecordingId, recordCount)
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
func (proxy *Proxy) ListRecordingsForUri(correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) error {

	bytes, err := ListRecordingsForUriRequestPacket(proxy.SessionId, correlationId, fromRecordingId, recordCount, stream, channel)

	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ListRecording(correlationId int64, fromRecordingId int64) error {
	// Create a packet and send it
	bytes, err := ListRecordingRequestPacket(proxy.SessionId, correlationId, fromRecordingId)
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
func (proxy *Proxy) ExtendRecording(correlationId int64, recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) error {
	// Create a packet and send it
	bytes, err := ExtendRecordingRequest2Packet(proxy.SessionId, correlationId, recordingId, stream, sourceLocation, autoStop, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

// RecordingPosition
func (proxy *Proxy) RecordingPosition(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := RecordingPositionRequestPacket(proxy.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) TruncateRecording(correlationId int64, recordingId int64, position int64) error {
	// Create a packet and send it
	bytes, err := TruncateRecordingPacket(proxy.SessionId, correlationId, recordingId, position)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecordingBySubscriptionId(correlationId int64, subscriptionId int64) error {
	// Create a packet and send it
	bytes, err := StopRecordingSubscriptionPacket(proxy.SessionId, correlationId, subscriptionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopRecordingByIdentity(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := StopRecordingByIdentityPacket(proxy.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopPosition(correlationId int64, recordingId int64) error {
	// Create a packet and send it
	bytes, err := StopPositionPacket(proxy.SessionId, correlationId, recordingId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) FindLastMatchingRecording(correlationId int64, minRecordingId int64, sessionId int32, stream int32, channel string) error {

	// Create a packet and send it
	bytes, err := FindLastMatchingRecordingPacket(proxy.SessionId, correlationId, minRecordingId, sessionId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) ListRecordingSubscriptionsRequest(controlSessionId int64, correlationId int64, pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channel string) error {

	// Create a packet and send it
	bytes, err := ListRecordingSubscriptionsPacket(proxy.SessionId, correlationId, pseudoIndex, subscriptionCount, applyStreamId, stream, channel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) BoundedReplayRequest(controlSessionId int64, correlationId int64, recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) error {

	// Create a packet and send it
	bytes, err := BoundedReplayPacket(proxy.SessionId, correlationId, recordingId, position, length, limitCounterId, replayStream, replayChannel)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}

func (proxy *Proxy) StopAllReplaysRequest(controlSessionId int64, correlationId int64, recordingId int64) error {

	// Create a packet and send it
	bytes, err := StopAllReplaysPacket(proxy.SessionId, correlationId, recordingId)
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

func (proxy *Proxy) PurgeRecordingRequest(correlationId int64, replaySessionId int64) error {
	// Create a packet and send it
	bytes, err := PurgeRecordingRequestPacket(proxy.SessionId, correlationId, replaySessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Offer failed: %d", ret)
	}

	return nil
}
