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
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
	"time"
)

// The polling mechanism is not parameterizsed so we need to set state for the results as we go
// These pieces are filled out by various ResponsePollers all of which will set IsPollComplete to true
type ControlResults struct {
	FragmentLimit        int
	CorrelationId        int64                         // FIXME: we may want overlapping operations
	ControlResponse      *codecs.ControlResponse       // FIXME: We may want a queue here
	RecordingDescriptors []*codecs.RecordingDescriptor // FIXME: ditto
	IsPollComplete       bool                          // FIXME: Consider making this atomic
}

// Control contains everything required for the archive control pub/sub request/response pair
type Control struct {
	ResponseChannel    string
	ResponseStream     int32
	RequestChannel     string
	RequestStream      int32
	Subscription       *aeron.Subscription
	Publication        *aeron.Publication
	State              ControlState
	IdleStrategy       idlestrategy.Idler
	marshaller         *codecs.SbeGoMarshaller
	RangeChecking      bool
	challengeSessionId int64 // FIXME: Todo
	SessionId          int64

	// Polling results
	Results ControlResults

	// FIXME: auth
	EncodedChallenge []byte
	rsa              RecordingSignalAdapter
}

// An archive "connection" involves some to and fro
const ControlStateError = -1
const ControlStateNew = 0
const ControlStateConnectRequestSent = 1
const ControlStateConnectRequestOk = 2
const ControlStateConnected = 3
const ControlStateTimedOut = 4

type ControlState struct {
	state int
	err   error
}

// Create a new initialized control. Note that a control does require inititializtion for it's channels
func NewControl() *Control {
	control := new(Control)
	control.RangeChecking = ArchiveDefaults.RangeChecking
	control.IdleStrategy = ArchiveDefaults.ControlIdleStrategy
	control.ResponseChannel = ArchiveDefaults.ResponseChannel
	control.ResponseStream = ArchiveDefaults.ResponseStream
	control.RequestChannel = ArchiveDefaults.RequestChannel
	control.RequestStream = ArchiveDefaults.RequestStream
	control.marshaller = codecs.NewSbeGoMarshaller()

	return control
}

// useful to see in debug mode
func ArchiveNewSubscriptionHandler(string, int32, int64) {
	logger.Debugf("Archive NewSubscriptionandler\n")
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control
func ControlFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		// FIXME: Wef could use an ErrorHandler/Listener
		logger.Error("Failed to decode control message header", err)

	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			logger.Info("Uncorrelated control response correlationId=", controlResponse.CorrelationId) // Not much to be done here as we can't correlate
			return
		}

		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via a channel
// FIXME: This is what the connection establishment currently uses, switch it over ro non connection specific one
func ConnectionControlFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlSubscriptionHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: We should use an ErrorHandler/Listener
		// Not much to be done here as we can't correlate
		logger.Error("Failed to decode control message header", err)

	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			logger.Info("Uncorrelated control response correlationId=", controlResponse.CorrelationId) // Not much to be done here as we can't correlate
			return
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
		control.SessionId = controlResponse.ControlSessionId

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// The control response poller uses local state to pass back information from the underlying subscription
func (control *Control) ErrorHandler(err error) {
	// FIXME: for now I'm just logging
	logger.Errorf(err.Error())
}

// The control response poller uses local state to pass back information from the underlying subscription
func (control *Control) Poll(handler term.FragmentHandler, fragmentLimit int) int {
	control.Results.ControlResponse = nil
	control.Results.IsPollComplete = false

	// FIXME: check what controlledPoll might do instead
	return control.Subscription.Poll(handler, fragmentLimit)
}

// Poll for a matching ControlReponse
func (control *Control) PollNextResponse(correlationId int64) error {
	logger.Debugf("PollNextResponse(%d) start", correlationId)

	start := time.Now()
	fragmentsWanted := 1 // FIXME: check is this safe to assume

	for {
		fragmentsWanted -= control.Poll(ControlFragmentHandler, fragmentsWanted)

		// Check result
		if control.Results.IsPollComplete {
			logger.Debugf("PollNextResponse(%d) complete", correlationId)
			if control.Results.ControlResponse.Code != codecs.ControlResponseCode.OK {
				err := fmt.Errorf("Control Response failure: %s", control.Results.ControlResponse.ErrorMessage)
				logger.Debug(err)
				return err
			} else {
				return nil
			}
		}

		if control.Subscription.IsClosed() {
			return fmt.Errorf("response channel from archive is not connected")
		}

		if time.Since(start) > ArchiveDefaults.ControlTimeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.IdleStrategy.Idle(0)
	}

	return nil
}

// Poll for a specific correlationId
// Returns nil on success, error otherwise, with detail passed back via Control.ControlResponse
func (control *Control) PollForResponse(correlationId int64) error {
	logger.Debugf("PollForResponse(%d)", correlationId)

	for {
		// Check for error
		if err := control.PollNextResponse(correlationId); err != nil {
			return err
		}

		// Check we're on the right session
		if control.Results.ControlResponse.ControlSessionId != control.SessionId {
			// FIXME: Other than log?
			control.ErrorHandler(fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.SessionId))
			control.Results.IsPollComplete = true
			return nil
		}

		// Check we've got the right correlationId
		if control.Results.ControlResponse.CorrelationId != correlationId {
			// FIXME: Other than log?
			control.ErrorHandler(fmt.Errorf("Control Response expected CorrelationId %d, received %d", correlationId, control.Results.ControlResponse.CorrelationId))
			control.Results.IsPollComplete = true
			return nil
		}

		control.Results.IsPollComplete = true
		return nil
	}
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control
func DescriptorFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("DescriptorFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var recordingDescriptor = new(codecs.RecordingDescriptor)
	var controlResponse = new(codecs.ControlResponse)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		// FIXME: Wef could use an ErrorHandler/Listener
		logger.Error("Failed to decode control message header", err)

	}

	switch hdr.TemplateId {
	case recordingDescriptor.SbeTemplateId():
		logger.Debugf("Received RecordingDescriptorlResponse: length %d", buf.Len())
		if err := recordingDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode RecordingDescriptor", err)
		}
		logger.Debugf("RecordingDescriptor: %#v\n", recordingDescriptor)

		// Look it up
		control, ok := correlationsMap[recordingDescriptor.CorrelationId]
		if !ok {
			logger.Info("Uncorrelated control response correlationId=", controlResponse.CorrelationId) // Not much to be done here as we can't correlate
			return
		}

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)

	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			// Not much to be done here as we can't correlate
			logger.Info("Uncorrelated control response correlationId=", controlResponse.CorrelationId)
			return
		}
		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

		if controlResponse.Code != codecs.ControlResponseCode.RECORDING_UNKNOWN {
			logger.Debugf("ControlResponse error UNKNOWN: %s\n", controlResponse.ErrorMessage)
			return
		}
		if controlResponse.Code != codecs.ControlResponseCode.ERROR {
			logger.Debugf("ControlResponse error ERROR: %s\n", controlResponse.ErrorMessage)
			return
		}

	default:
		logger.Debugf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// Poll for a fragmentLimit of Descriptors
func (control *Control) PollNextDescriptor(correlationId int64, fragmentsWanted int) error {
	logger.Debugf("PollNextDescriptor(%d) start", correlationId)
	start := time.Now()

	for !control.Results.IsPollComplete {
		fragmentsWanted -= control.Poll(DescriptorFragmentHandler, fragmentsWanted)

		if fragmentsWanted <= 0 {
			logger.Debugf("PollNextDescriptor(%d) all fragments received", fragmentsWanted)
			control.Results.IsPollComplete = true
		}

		if control.Subscription.IsClosed() {
			return fmt.Errorf("response channel from archive is not connected")
		}

		if control.Results.IsPollComplete {
			logger.Debugf("PollNextDescriptor(%d) complete", correlationId)
			return nil
		}

		if fragmentsWanted > 0 {
			continue
		}

		// FIXME: maybe a separate timeout, this could legitimately be a while
		if time.Since(start) > ArchiveDefaults.ControlTimeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.IdleStrategy.Idle(0)
	}

	return nil
}

// Poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationId int64, fragmentsWanted int32) error {
	logger.Debugf("PollForDescriptors(%d)", correlationId)

	for {
		// Check for error
		if err := control.PollNextDescriptor(correlationId, int(fragmentsWanted)); err != nil {
			control.Results.IsPollComplete = true
			return err
		}

		if control.Results.IsPollComplete {
			logger.Debugf("PollForDescriptors(%d) complete", correlationId)
			return nil
		}

		// Check we're on the right session
		if control.Results.ControlResponse.ControlSessionId != control.SessionId {
			// FIXME: Other than log?1
			control.ErrorHandler(fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.SessionId))
			control.Results.IsPollComplete = true
			return nil
		}

		// Check we've got the right correlationId
		if control.Results.ControlResponse.CorrelationId != correlationId {
			// FIXME: Other than log?
			control.ErrorHandler(fmt.Errorf("Control Response expected CorrelationId %d, received %d", correlationId, control.Results.ControlResponse.CorrelationId))
			control.Results.IsPollComplete = true
			return nil
		}
	}
}
