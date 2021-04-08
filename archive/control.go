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
	Subscription *aeron.Subscription
	State        ControlState

	// Polling results
	Results ControlResults

	// FIXME: auth
	EncodedChallenge   []byte
	challengeSessionId int64 // FIXME: Todo

	marshaller *codecs.SbeGoMarshaller // FIXME: sort out

	archive *Archive // link to parent
}

// The polling mechanism is not parameterizsed so we need to set state for the results as we go
// These pieces are filled out by various ResponsePollers which will set IsPollComplete to true
type ControlResults struct {
	CorrelationId                    int64
	ControlResponse                  *codecs.ControlResponse
	RecordingDescriptors             []*codecs.RecordingDescriptor
	RecordingSubscriptionDescriptors []*codecs.RecordingSubscriptionDescriptor
	IsPollComplete                   bool
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

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control's Results
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
		err2 := fmt.Errorf("ControlFragmentHandler failed to decode control message header: %w", err)
		if Listeners.ErrorListener != nil {
			Listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ControlFragmentHandler failed to decode control response:%w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}

		// Check this is for our session
		_, ok := sessionsMap[controlResponse.ControlSessionId]
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("ControlFragmentHandler: Unexpected sessionId in control response: %d", controlResponse.ControlSessionId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("ControlFragmentHandler uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

	case recordingSignalEvent.SbeTemplateId():
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ControlFargamentHandler failed to decode recording signal: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		fmt.Printf("ControlFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
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
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("ConnectionControlFragmentHandler failed to decode control message header: %w", err)
		if Listeners.ErrorListener != nil {
			Listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ConnectionControlFragmentHandler failed to decode control response: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("ConnectionControlFragmentHandler: Uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		// Check result
		if controlResponse.Code != codecs.ControlResponseCode.OK {
			control.State.state = ControlStateError
			control.State.err = fmt.Errorf("Control Response failure: %s", controlResponse.ErrorMessage)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(control.State.err)
			}
		}

		// assert state change
		if control.State.state != ControlStateConnectRequestSent {
			control.State.state = ControlStateError
			control.State.err = fmt.Errorf("Control Response not expecting response")
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(control.State.err)
			}
		}

		// Looking good, so update state and store the SessionId
		control.State.state = ControlStateConnected
		control.State.err = nil
		control.archive.SessionId = controlResponse.ControlSessionId

	default:
		fmt.Printf("ConnectionControlFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// The control response poller uses local state to pass back information from the underlying subscription
func (control *Control) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

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

		if time.Since(start) > control.archive.Options.Timeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.archive.Options.IdleStrategy.Idle(0)
	}
}

// Poll for a specific correlationId
// Returns nil, relevantId on success, error, 0 failure
// More complex responses are contained in Control.ControlResponse after the call
func (control *Control) PollForResponse(correlationId int64) (int64, error) {
	logger.Debugf("PollForResponse(%d)", correlationId)

	for {
		// Check for error
		if err := control.PollNextResponse(correlationId); err != nil {
			return 0, err
		}

		// Check we're on the right session
		if control.Results.ControlResponse.ControlSessionId != control.archive.SessionId {
			err := fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.archive.SessionId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			control.Results.IsPollComplete = true
			return 0, err
		}

		// Check we've got the right correlationId
		// This is usually a sign of a logic error in handling the protocol so we'll log it and move on
		// This can also happen on reconnection to same control subscription
		if control.Results.ControlResponse.CorrelationId != correlationId {
			err := fmt.Errorf("Control Response expected CorrelationId %d, received %d", correlationId, control.Results.ControlResponse.CorrelationId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
		}

		control.Results.IsPollComplete = true
		logger.Debugf("PollForResponse(%d) complete", correlationId)
		return control.Results.ControlResponse.RelevantId, nil
	}
}

// Poll the response stream once for an error. If another message is
// present then it will be skipped over so only call when not
// expecting another response.
//
// FIXME: Need to adjust fragment counts in case something else async happens
func (control *Control) PollForErrorResponse() error {

	fragments := control.Poll(ControlFragmentHandler, 1)
	if fragments == 0 {
		return nil
	}
	if control.Results.ControlResponse.Code == codecs.ControlResponseCode.ERROR {
		return fmt.Errorf(string(control.Results.ControlResponse.ErrorMessage))
	}

	return nil
}

// Poll for descriptors (both recording and subscription)
// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control's Results
func DescriptorFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	// logger.Debugf("DescriptorFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var recordingDescriptor = new(codecs.RecordingDescriptor)
	var recordingSubscriptionDescriptor = new(codecs.RecordingSubscriptionDescriptor)
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
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("Uncorrelated recordingDesciptor correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)

	case recordingSubscriptionDescriptor.SbeTemplateId():
		if err := recordingSubscriptionDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("Failed to decode RecordingSubscriptioDescriptor: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		control, ok := correlationsMap[recordingSubscriptionDescriptor.CorrelationId]
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("Uncorrelated recordingSubscriptionDescriptor correlationId=%d\n%#v", controlResponse.CorrelationId, controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingSubscriptionDescriptors = append(control.Results.RecordingSubscriptionDescriptors, recordingSubscriptionDescriptor)

	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("Failed to decode control response: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		control, ok := correlationsMap[controlResponse.CorrelationId]
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("DescriptorFragmentHandler: Uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}

		// We're basically finished so prepare our OOB return values and log some info if we can
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

	case recordingSignalEvent.SbeTemplateId():
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("Failed to decode recording signal: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return

		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		logger.Debugf("DescriptorFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
		fmt.Printf("DescriptorFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
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

		if time.Since(start) > control.archive.Options.Timeout {
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.archive.Options.IdleStrategy.Idle(0)
	}

	return nil
}

// Poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationId int64, fragmentsWanted int32) error {
	logger.Debugf("PollForDescriptors(%d)", correlationId)

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil                  // Clear old results
	control.Results.IsPollComplete = false                 // Clear completion flag
	control.Results.RecordingDescriptors = nil             // Clear previous search
	control.Results.RecordingSubscriptionDescriptors = nil // Clear previous search

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
		if control.Results.ControlResponse.ControlSessionId != control.archive.SessionId {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("Control Response expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.archive.SessionId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}

			control.Results.IsPollComplete = true
			return nil
		}

		// Check we've got the right correlationId
		if control.Results.ControlResponse.CorrelationId != correlationId {
			err := fmt.Errorf("Control Response expected CorrelationId %d, received %d", correlationId, control.Results.ControlResponse.CorrelationId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			control.Results.IsPollComplete = true
			return nil
		}
	}
}
