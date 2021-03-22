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
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
	"time"
)

// Control contains everything required for the archive subscription/response side
type Control struct {
	Context      *ArchiveContext
	Subscription *aeron.Subscription
	Listeners    *ArchiveListeners
	State        ControlState

	// Polling results
	Results ControlResults

	// FIXME: auth
	EncodedChallenge   []byte
	challengeSessionId int64 // FIXME: Todo

	marshaller *codecs.SbeGoMarshaller // FIXME: sort out
}

// The polling mechanism is not parameterizsed so we need to set state for the results as we go
// These pieces are filled out by various ResponsePollers which will set IsPollComplete to true
type ControlResults struct {
	CorrelationId        int64                         // FIXME: we may want overlapping operations
	ControlResponse      *codecs.ControlResponse       // FIXME: We may want a queue here
	RecordingDescriptors []*codecs.RecordingDescriptor // FIXME: ditto
	IsPollComplete       bool
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
func NewControl(context *ArchiveContext) *Control {
	control := new(Control)
	control.Context = context

	return control
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control
func ControlFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)
	var recordingSignalEvent = new(codecs.RecordingSignalEvent)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		// FIXME: We could use an ErrorHandler/Listener
		logger.Error("Failed to decode control message header", err)

	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			logger.Infof("Uncorrelated control response correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse) // Not much to be done here as we can't correlate
			return
		}

		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

	case recordingSignalEvent.SbeTemplateId():
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode recording signal", err)
		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// The connection handling specific fragment handler.
// This mechanism only alows us to pass results back via global state which we do in control.State
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
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			logger.Infof("Uncorrelated control response correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse) // Not much to be done here as we can't correlate
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

		// Looking good, so update state and store the SessionId
		control.State.state = ControlStateConnected
		control.State.err = nil
		control.Context.SessionId = controlResponse.ControlSessionId

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

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.Context.Options.RangeChecking

	control.Results.ControlResponse = nil  // Clear old results
	control.Results.IsPollComplete = false // Clear completion flag

	// FIXME: check what controlledPoll might do instead
	return control.Subscription.Poll(handler, fragmentLimit)
}

// Poll for a matching ControlReponse
func (control *Control) PollNextResponse(correlationId int64) error {
	logger.Debugf("PollNextResponse(%d) start", correlationId)

	start := time.Now()
	fragmentsWanted := 10 // FIXME: check is this safe to assume

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

		if time.Since(start) > control.Context.Options.Timeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.Context.Options.IdleStrategy.Idle(0)
	}
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
		if control.Results.ControlResponse.ControlSessionId != control.Context.SessionId {
			// FIXME: Other than log?
			control.ErrorHandler(fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.Context.SessionId))
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
		logger.Debugf("PollForResponse(%d) complete", correlationId)
		return nil
	}
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control
func DescriptorFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	// logger.Debugf("DescriptorFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var recordingDescriptor = new(codecs.RecordingDescriptor)
	var controlResponse = new(codecs.ControlResponse)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		// FIXME: We could use an ErrorHandler/Listener
		logger.Error("Failed to decode control message header", err)

	}

	switch hdr.TemplateId {
	case recordingDescriptor.SbeTemplateId():
		// logger.Debugf("Received RecordingDescriptorResponse: length %d", buf.Len())
		if err := recordingDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode RecordingDescriptor", err)
		}
		// logger.Debugf("RecordingDescriptor: %#v\n", recordingDescriptor)

		// Look it up
		control, ok := correlationsMap[recordingDescriptor.CorrelationId]
		if !ok {
			logger.Infof("Uncorrelated control response correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse) // Not much to be done here as we can't correlate
			return
		}

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)

	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode control response", err)
		}

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			logger.Infof("Uncorrelated control response correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse) // Not much to be done here as we can't correlate
			return
		}

		// We're basically finished but let's log some info if we can
		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

		// Expected
		if controlResponse.Code == codecs.ControlResponseCode.RECORDING_UNKNOWN {
			logger.Debugf("ControlResponse error UNKNOWN: %s\n", controlResponse.ErrorMessage)
			return
		}

		// Unexpected
		if controlResponse.Code == codecs.ControlResponseCode.ERROR {
			logger.Debugf("ControlResponse error ERROR: %s\n%#v", controlResponse.ErrorMessage, controlResponse)
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

		if time.Since(start) > control.Context.Options.Timeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.Context.Options.IdleStrategy.Idle(0)
	}

	return nil
}

// Poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationId int64, fragmentsWanted int32) error {
	logger.Debugf("PollForDescriptors(%d)", correlationId)

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.Context.Options.RangeChecking

	control.Results.ControlResponse = nil      // Clear old results
	control.Results.IsPollComplete = false     // Clear completion flag
	control.Results.RecordingDescriptors = nil // Clear previous search

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
		if control.Results.ControlResponse.ControlSessionId != control.Context.SessionId {
			// FIXME: Other than log?1
			control.ErrorHandler(fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.Context.SessionId))
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
