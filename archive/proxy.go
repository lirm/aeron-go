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
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
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

// Offer to our request publication
func (proxy *Proxy) Offer(buf bytes.Buffer) int64 {
	bytes := buf.Bytes()
	length := int32(buf.Len())
	buffer := atomic.MakeBuffer(bytes, length)

	var ret int64
	for retries := proxy.Retries; retries > 0; retries-- {
		ret = proxy.Publication.Offer(buffer, 0, length, nil)
		switch ret {
		case aeron.NotConnected:
			return ret // Fail immediately
		case aeron.PublicationClosed:
			return ret // Fail immediately
		case aeron.MaxPositionExceeded:
			return ret // Fail immediately
		default:
			if ret > 0 {
				return ret // Succeed
			} else {
				// Retry (aeron.BackPressured or aeron.AdminAction)
				proxy.IdleStrategy.Idle(0)
			}
		}
	}

	// Give up, returning the last failure
	logger.Debugf("Proxy.Offer giving up [%d]", ret)
	return ret

}

// Make a Connect Request
func (proxy *Proxy) ConnectRequest(responseChannel string, responseStream int32, correlationId int64) error {

	// Create a packet and send it
	bytes, err := ConnectRequestPacket(responseChannel, responseStream, correlationId)
	if err != nil {
		return err
	}

	if ret := proxy.Publication.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Publication.Offer failed: %d", ret)
	}

	return nil
}

// Start a Recorded Publication
func (proxy *Proxy) StartRecording(channel string, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, correlationId int64, sessionId int64) error {

	bytes, err := StartRecordingRequest2Packet(channel, stream, sourceLocation, autoStop, correlationId, sessionId)
	if err != nil {
		return err
	}

	if ret := proxy.Publication.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Publication.Offer failed: %d", ret)
	}

	return nil
}

// Start a Recorded Publication
func (proxy *Proxy) ListRecordingsForUri(sessionId int64, correlationId int64, fromRecordingId int64, recordCount int32, stream int32, channel string) error {

	bytes, err := ListRecordingsForUriRequestPacket(sessionId, correlationId, fromRecordingId, recordCount, stream, channel)

	if err != nil {
		return err
	}

	if ret := proxy.Publication.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
		return fmt.Errorf("Publication.Offer failed: %d", ret)
	}

	return nil
}
