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

// Package archive provides API access to Aeron's archive-media-driver
package archive

import (
	"bytes"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/archive/codecs"
	"time"
)

// Control contains everything required for the archive control pub/sub request/response pair
type Control struct {
	ResponseChannel             string
	ResponseStream              int32
	RequestChannel              string
	RequestStream               int32
	Subscription                *aeron.Subscription
	Publication                 *aeron.Publication
	State                       ControlState
	ConnectionChange            chan ControlState // Invoked on connect, disconnect, and error
	PolledBytes                 chan []byte       // Public for returning data from poller
	terminateSubscriptionPoller chan bool         // Internal for stopping the poller
	marshaller                  *codecs.SbeGoMarshaller
	RangeChecking               bool
	challengeSessionID          int64 // FIXME: Todo
	SessionID                   int64
	CorrelationID               int64 // FIXME: we'll want overlapping operations
}

// An archive "connection" involves some to and fro
const ControlStateError = -1
const ControlStateNew = 0
const ControlStateConnectRequestSent = 1
const ControlStateConnectRequestOk = 2
const ControlStateConnected = 3
const ControlStateDisconnected = 4

type ControlState struct {
	state int
	err   error
}

// Globals
var correlations = make(map[int64]*Control) // Map of correlationID so we can correlate responses

// Create a new initialized control. Note that a control does require inititializtion for it's channels
func NewControl() *Control {
	control := new(Control)
	control.RangeChecking = ArchiveDefaults.RangeChecking
	control.ConnectionChange = make(chan ControlState)
	control.PolledBytes = make(chan []byte)
	control.terminateSubscriptionPoller = make(chan bool)
	control.marshaller = codecs.NewSbeGoMarshaller()
	control.ResponseChannel = ArchiveDefaults.ResponseChannel
	control.ResponseStream = ArchiveDefaults.ResponseStream
	control.RequestChannel = ArchiveDefaults.RequestChannel
	control.RequestStream = ArchiveDefaults.RequestStream

	return control
}

// useful to see in debug mode
func ArchiveNewSubscriptionHandler(string, int32, int64) {
	logger.Debugf("Archive NewSubscriptionandler\n")
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via a channel
func ControlSubscriptionHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlSubscriptionHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)

	buf := new(bytes.Buffer)
	_ = buffer.GetBytesArray(offset, length)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: Should we use an ErrorHandler?
		logger.Error("Failed to decode control message header", err) // Not much to be done here as we can't correlate

	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			logger.Error("Failed to decode control response", err) // Not much to be done here as we can't correlate
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := correlations[controlResponse.CorrelationId]
		if !ok {
			logger.Error("Failed to correlate control response correlationID=", controlResponse.CorrelationId) // Not much to be done here as we can't correlate
			fmt.Printf("correlations:\n%#v\n", correlations)
		}

		// Check result
		if controlResponse.Code != codecs.ControlResponseCode.OK {
			control.State.state = ControlStateError
			control.State.err = fmt.Errorf("Control Response failure: %s", controlResponse.ErrorMessage)
			logger.Warning(control.State.err)
		}

		// assert state change
		if control.State.state != ControlStateConnectRequestSent {
			control.State.state = ControlStateError
			control.State.err = fmt.Errorf("Control Response not expecting response")
			logger.Error(control.State.err)
		}

		// Looking good
		control.State.state = ControlStateConnected
		control.State.err = nil
		control.SessionID = controlResponse.ControlSessionId
		control.ConnectionChange <- control.State

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// Start the control response subscription poller.
func (control *Control) StartControlSubscriptionPoller() {
	control.terminateSubscriptionPoller = make(chan bool, 1)
	go control.ControlSubscriptionPoller()
}

// Stop the control response subscription poller.
func (control *Control) StopControlSubscriptionPoller() {
	close(control.terminateSubscriptionPoller)
}

// The control response poller itself which invokes the handler when data is available
func (control *Control) ControlSubscriptionPoller() {
	fmt.Printf("ControlSubscriptionPoller starting\n")
	for {
		fragments := control.Subscription.Poll(ControlSubscriptionHandler, 10)
		if fragments > 0 {
			logger.Debugf("Read %d fragments\n", fragments)
		}

		// Check for exit
		select {
		case res, ok := <-control.terminateSubscriptionPoller:
			if res || !ok {
				fmt.Printf("Exiting ControlSubscriptionPoller\n")
				return
			}
		case <-time.After(ArchiveDefaults.ControlIdleTime):
			// timed out
		}
	}
}
