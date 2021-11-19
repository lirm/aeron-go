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
	State        controlState

	// Polling results
	Results ControlResults

	archive *Archive // link to parent
}

// ControlResults for holding state over a Control request/response
// The polling mechanism is not parameterizsed so we need to set state for the results as we go
// These pieces are filled out by various ResponsePollers which will set IsPollComplete to true
type ControlResults struct {
	CorrelationId                    int64
	ControlResponse                  *codecs.ControlResponse
	RecordingDescriptors             []*codecs.RecordingDescriptor
	RecordingSubscriptionDescriptors []*codecs.RecordingSubscriptionDescriptor
	IsPollComplete                   bool
	FragmentsReceived                int
}

// Archive Connection State used internally for connection establishment
const (
	ControlStateError              = -1
	ControlStateNew                = iota
	ControlStateConnectRequestSent = iota
	ControlStateChallenged         = iota
	ControlStateConnectRequestOk   = iota
	ControlStateConnected          = iota
	ControlStateTimedOut           = iota
)

// Used internally to handle connection state
type controlState struct {
	state int
	err   error
}

// CodecIds stops us allocating every object when we need only one
// Arguably SBE should give us a static value
type CodecIds struct {
	controlResponse                 uint16
	challenge                       uint16
	recordingDescriptor             uint16
	recordingSubscriptionDescriptor uint16
	recordingSignalEvent            uint16
	recordingStarted                uint16
	recordingProgress               uint16
	recordingStopped                uint16
}

var codecIds CodecIds

func init() {
	var controlResponse codecs.ControlResponse
	var challenge codecs.Challenge
	var recordingDescriptor codecs.RecordingDescriptor
	var recordingSubscriptionDescriptor codecs.RecordingSubscriptionDescriptor
	var recordingSignalEvent codecs.RecordingSignalEvent
	var recordingStarted = new(codecs.RecordingStarted)
	var recordingProgress = new(codecs.RecordingProgress)
	var recordingStopped = new(codecs.RecordingStopped)

	codecIds.controlResponse = controlResponse.SbeTemplateId()
	codecIds.challenge = challenge.SbeTemplateId()
	codecIds.recordingDescriptor = recordingDescriptor.SbeTemplateId()
	codecIds.recordingSubscriptionDescriptor = recordingSubscriptionDescriptor.SbeTemplateId()
	codecIds.recordingSignalEvent = recordingSignalEvent.SbeTemplateId()
	codecIds.recordingStarted = recordingStarted.SbeTemplateId()
	codecIds.recordingProgress = recordingProgress.SbeTemplateId()
	codecIds.recordingStopped = recordingStopped.SbeTemplateId()
}

// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control's Results
func controlFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("controlFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't really tell what went wrong
		err2 := fmt.Errorf("controlFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if Listeners.ErrorListener != nil {
			Listeners.ErrorListener(err2)
		}
		return
	}

	switch hdr.TemplateId {
	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't see what's gone wrong
			err2 := fmt.Errorf("controlFragmentHandler failed to decode control response:%w", err)
			// Call the global error handler, ugly but it's all we've got
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Check this is for our session by looking it up
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			// Must have been for someone else which can happen if two or more clients have
			// use the same channel/stream
			logger.Debugf("controlFragmentHandler/controlResponse: ignoring correlationID=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			return
		}
		control := c.(*Control)
		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

	case codecIds.recordingSignalEvent:
		var recordingSignalEvent = new(codecs.RecordingSignalEvent)

		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't really tell what went wrong
			err2 := fmt.Errorf("ControlFargamentHandler failed to decode recording signal: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		// This can happen when testing/adding new functionality
		fmt.Printf("controlFragmentHandler: Unexpected message type %d\n", hdr.TemplateId)
	}
}

// ConnectionControlFragmentHandler is the connection handling specific fragment handler.
// This mechanism only alows us to pass results back via global state which we do in control.State
func ConnectionControlFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlSubscriptionHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("ConnectionControlFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if Listeners.ErrorListener != nil {
			Listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
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
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			// Not much to be done here as we can't know what went wrong
			err := fmt.Errorf("connectionControlFragmentHandler: ignoring uncorrelated correlationID=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

		// Check result
		if controlResponse.Code != codecs.ControlResponseCode.OK {
			control.State.state = ControlStateError
			control.State.err = fmt.Errorf("Control Response failure: %s", controlResponse.ErrorMessage)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(control.State.err)
			}
			return
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

	case codecIds.challenge:
		var challenge = new(codecs.Challenge)

		if err := challenge.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ControlFragmentHandler failed to decode challenge: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}

		logger.Infof("ControlFragmentHandler: challenge:%s, session:%d, correlationID:%d", challenge.EncodedChallenge, challenge.ControlSessionId, challenge.CorrelationId)

		// Look it up
		c, ok := correlations.Load(challenge.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("connectionControlFragmentHandler: ignoring uncorrelated correlationID=%d", challenge.CorrelationId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

		// Check the challenge is expected iff our option for this is not nil
		if control.archive.Options.AuthChallenge != nil {
			if !bytes.Equal(control.archive.Options.AuthChallenge, challenge.EncodedChallenge) {
				control.State.err = fmt.Errorf("ChallengeResponse Unexpected: expected:%v received:%v", control.archive.Options.AuthChallenge, challenge.EncodedChallenge)
				return
			}
		}

		// Send the response
		// Looking good, so update state and store the SessionId
		control.State.state = ControlStateChallenged
		control.State.err = nil
		control.archive.SessionId = challenge.ControlSessionId
		control.archive.Proxy.ChallengeResponse(challenge.CorrelationId, control.archive.Options.AuthResponse)

	default:
		fmt.Printf("ConnectionControlFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
	}
}

// Poll provides the control response poller using local state to pass
// back data from the underlying subscription
func (control *Control) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil  // Clear old results
	control.Results.IsPollComplete = false // Clear completion flag

	return control.Subscription.Poll(handler, fragmentLimit)
}

// PollNextResponse for a matching ControlReponse
func (control *Control) PollNextResponse(correlationID int64) error {
	logger.Debugf("PollNextResponse(%d) start", correlationID)

	start := time.Now()

	// Poll for events As we can get async events we receive a
	// fairly arbitrary number of messages here Additionally if
	// two clients use the same channel/stream for responses they
	// will see each others messages.
	// So we just continually poll until we get our response or
	// timeout without having completed and treat that as an error
	for {
		control.Poll(controlFragmentHandler, 1)

		// Check result
		if control.Results.IsPollComplete {
			logger.Debugf("PollNextResponse(%d) complete", correlationID)
			if control.Results.ControlResponse.Code != codecs.ControlResponseCode.OK {
				err := fmt.Errorf("Control Response failure: %s", control.Results.ControlResponse.ErrorMessage)
				logger.Debug(err)
				return err
			}
			return nil
		}

		if control.Subscription.IsClosed() {
			return fmt.Errorf("response channel from archive is not connected")
		}

		if time.Since(start) > control.archive.Options.Timeout {
			return fmt.Errorf("timeout waiting for correlationID %d", correlationID)
		}

		// Idle as configured
		control.archive.Options.IdleStrategy.Idle(0)
	}
}

// PollForResponse polls for a specific correlationID
// Returns nil, relevantId on success, error, 0 failure
// More complex responses are contained in Control.ControlResponse after the call
func (control *Control) PollForResponse(correlationID int64) (int64, error) {
	logger.Debugf("PollForResponse(%d)", correlationID)

	for {
		// Check for error
		if err := control.PollNextResponse(correlationID); err != nil {
			return 0, err
		}

		// Check we're on the right session
		if control.Results.ControlResponse.ControlSessionId != control.archive.SessionId {
			err := fmt.Errorf("PollForResponse() expected SessionId %d, received %d", control.Results.ControlResponse.ControlSessionId, control.archive.SessionId)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			control.Results.IsPollComplete = true
			return 0, err
		}

		// Sanity Check we've got the right correlationID which in theory can't happen
		if control.Results.ControlResponse.CorrelationId != correlationID {
			logger.Errorf("PollForResponse() unexpected CorrelationID expected:%d, received:%d", correlationID, control.Results.ControlResponse.CorrelationId)
		} else {
			control.Results.IsPollComplete = true
			logger.Debugf("PollForResponse(%d) complete", correlationID)
			return control.Results.ControlResponse.RelevantId, nil
		}
	}
}

// DescriptorFragmentHandler is used to poll for descriptors (both recording and subscription)
// The current subscription handler doesn't provide a mechanism for passing a rock
// so we return data via the control's Results
// FIXME:Bug Need to adjust fragment counts in case something async happens
func DescriptorFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	// logger.Debugf("DescriptorFragmentHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("DescriptorFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if Listeners.ErrorListener != nil {
			Listeners.ErrorListener(err2)
		}
		return
	}

	switch hdr.TemplateId {
	case codecIds.recordingDescriptor:
		var recordingDescriptor = new(codecs.RecordingDescriptor)
		// logger.Debugf("Received RecordingDescriptor: length %d", buf.Len())
		if err := recordingDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode RecordingDescriptor: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}
		// logger.Debugf("RecordingDescriptor: %#v\n", recordingDescriptor)

		// Look it up
		c, ok := correlations.Load(recordingDescriptor.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			logger.Debugf("DescriptorFragmentHandler ignoring correlationID=%d\n%#v", recordingDescriptor.CorrelationId)
			return
		}
		control := c.(*Control)

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)
		control.Results.FragmentsReceived++

	case codecIds.recordingSubscriptionDescriptor:
		// logger.Debugf("Received RecordingSubscriptionDescriptor: length %d", buf.Len())
		var recordingSubscriptionDescriptor = new(codecs.RecordingSubscriptionDescriptor)
		if err := recordingSubscriptionDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode RecordingSubscriptioDescriptor: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		c, ok := correlations.Load(recordingSubscriptionDescriptor.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("uncorrelated recordingSubscriptionDescriptor correlationID=%d\n%#v", recordingSubscriptionDescriptor.CorrelationId, recordingSubscriptionDescriptor)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingSubscriptionDescriptors = append(control.Results.RecordingSubscriptionDescriptors, recordingSubscriptionDescriptor)
		control.Results.FragmentsReceived++

	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode control response: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("DescriptorFragmentHandler: Uncorrelated control response correlationID=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

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

	case codecIds.recordingSignalEvent:
		var recordingSignalEvent = new(codecs.RecordingSignalEvent)
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode recording signal: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return

		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		err := fmt.Errorf("descriptorFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
		logger.Debug("Error: %s\n", err.Error())
		fmt.Printf("Error: %s\n", err.Error())
	}
}

// PollNextDescriptor to poll for a fragmentLimit of Descriptors
func (control *Control) PollNextDescriptor(correlationID int64, fragmentsWanted int) error {
	logger.Debugf("PollNextDescriptor(%d) start", correlationID)
	start := time.Now()

	// Poll for descriptors
	// We should end up with a maximum of fragmentsWanted or finish with a ControlResults
	for !control.Results.IsPollComplete {
		control.Poll(DescriptorFragmentHandler, fragmentsWanted-control.Results.FragmentsReceived)

		if fragmentsWanted <= control.Results.FragmentsReceived {
			logger.Debugf("PollNextDescriptor(%d) all fragments received", fragmentsWanted)
			control.Results.IsPollComplete = true
		}

		if control.Subscription.IsClosed() {
			return fmt.Errorf("response channel from archive is not connected")
		}

		if control.Results.IsPollComplete {
			logger.Debugf("PollNextDescriptor(%d) complete", correlationID)
			return nil
		}

		if time.Since(start) > control.archive.Options.Timeout {
			return fmt.Errorf("PollNextDescriptor timeout waiting for correlationID %d", correlationID)
		}

		control.archive.Options.IdleStrategy.Idle(0)
	}

	return nil
}

// PollForDescriptors to poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationID int64, fragmentsWanted int32) error {
	logger.Debugf("PollForDescriptors(%d)", correlationID)

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil                  // Clear old results
	control.Results.IsPollComplete = false                 // Clear completion flag
	control.Results.RecordingDescriptors = nil             // Clear previous results
	control.Results.RecordingSubscriptionDescriptors = nil // Clear previous results
	control.Results.FragmentsReceived = 0                  // Reset our results count

	for {
		// Check for error
		if err := control.PollNextDescriptor(correlationID, int(fragmentsWanted)); err != nil {
			control.Results.IsPollComplete = true
			return err
		}

		if control.Results.IsPollComplete {
			logger.Debugf("PollForDescriptors(%d) complete", correlationID)
			return nil
		}
	}
}
