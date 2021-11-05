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
	ExtraFragments                   int
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
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("controlFragmentHandler() failed to decode control message header: %w", err)
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
			err2 := fmt.Errorf("controlFragmentHandler failed to decode control response:%w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}

		// Check this is for our session
		// Look it up
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("controlFragmentHandler uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			logger.Infof("controlFragmentHandler/controlResponse: Uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse) // Not much to be done here as we can't correlate
			return
		}
		control := c.(*Control)
		control.Results.ControlResponse = controlResponse
		control.Results.IsPollComplete = true

	case codecIds.recordingSignalEvent:
		var recordingSignalEvent = new(codecs.RecordingSignalEvent)

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

		// If we can locate this correlationId then we can let our parent know we
		// will want an extra fragment
		c, ok := correlations.Load(recordingSignalEvent.CorrelationId)
		if !ok {
			logger.Infof("controlFragmentHandler/recordingSignalEvent: Uncorrelated recordingSignalEvent correlationId=%d\n%#v", recordingSignalEvent.CorrelationId, recordingSignalEvent) // Not much to be done here as we can't correlate
			return
		}
		control := c.(*Control)
		control.Results.ExtraFragments++

	default:
		// This can happen when testing
		fmt.Printf("controlFragmentHandler: Unexpected message type %d\n", hdr.TemplateId)
	}

	return
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
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("ConnectionControlFragmentHandler: Uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
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
			err2 := fmt.Errorf("ControlFargamentHandler failed to decode challenge: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
		}

		logger.Infof("ControlFragmentHandler: challenge:%s, session:%d, correlationId:%d", challenge.EncodedChallenge, challenge.ControlSessionId, challenge.CorrelationId)

		// Look it up
		c, ok := correlations.Load(challenge.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("ConnectionControlFragmentHandler: Uncorrelated challenge correlationId=%d", challenge.CorrelationId)
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

	return
}

// Poll provides rhe control response poller uses local state to pass
// back data from the underlying subscription
func (control *Control) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil  // Clear old results
	control.Results.IsPollComplete = false // Clear completion flag
	control.Results.ExtraFragments = 0     // Clear extra fragment count

	return control.Subscription.Poll(handler, fragmentLimit)
}

// PollNextResponse for a matching ControlReponse
func (control *Control) PollNextResponse(correlationId int64) error {
	logger.Debugf("PollNextResponse(%d) start", correlationId)

	start := time.Now()

	// Poll for events
	// As we can get async events we need to track how many/ extra events we might be wanting
	// but we're also generous in how many we request as here we only want one
	fragmentsWanted := 10
	control.Results.ExtraFragments = 0
	for {
		fragmentsReceived := control.Poll(controlFragmentHandler, fragmentsWanted)
		fragmentsWanted = fragmentsWanted - fragmentsReceived + control.Results.ExtraFragments
		control.Results.ExtraFragments = 0

		// Check result
		if control.Results.IsPollComplete {
			logger.Debugf("PollNextResponse(%d) complete", correlationId)
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
			return fmt.Errorf("timeout waiting for correlationId %d", correlationId)
		}
		control.archive.Options.IdleStrategy.Idle(0)
	}
}

// PollForResponse polls for a specific correlationId
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

// PollForErrorResponse polls the response stream once for an
// error. If another message is present then it will be skipped over
// so only call when not expecting another response.
func (control *Control) PollForErrorResponse() error {

	// Poll for events
	// As we can get async events we need to track how many/ extra events we might be wanting
	// but we're also generous in how many we request as here we only want one
	fragmentsWanted := 10
	control.Results.ExtraFragments = 0
	for {
		fragmentsReceived := control.Poll(controlFragmentHandler, 1)
		fragmentsWanted = fragmentsWanted - fragmentsReceived + control.Results.ExtraFragments
		control.Results.ExtraFragments = 0

		if fragmentsWanted >= 1 {
			continue
		}

		// If we received a response with an error then return it
		if fragmentsWanted == 0 {
			if control.Results.ControlResponse.Code == codecs.ControlResponseCode.ERROR {
				return fmt.Errorf(string(control.Results.ControlResponse.ErrorMessage))
			}
		} else {
			// Nothing there yet
			return nil
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
	}

	switch hdr.TemplateId {
	case codecIds.recordingDescriptor:
		var recordingDescriptor = new(codecs.RecordingDescriptor)
		// logger.Debugf("Received RecordingDescriptorResponse: length %d", buf.Len())
		// Not much to be done here as we can't correlate
		if err := recordingDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			err2 := fmt.Errorf("Failed to decode RecordingDescriptor: %w", err)
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
			err := fmt.Errorf("Uncorrelated recordingDesciptor correlationId=%d\n%#v", recordingDescriptor.CorrelationId, recordingDescriptor)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)

	case codecIds.recordingSubscriptionDescriptor:
		var recordingSubscriptionDescriptor = new(codecs.RecordingSubscriptionDescriptor)
		if err := recordingSubscriptionDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("Failed to decode RecordingSubscriptioDescriptor: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return
		}

		// Look it up
		c, ok := correlations.Load(recordingSubscriptionDescriptor.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("Uncorrelated recordingSubscriptionDescriptor correlationId=%d\n%#v", recordingSubscriptionDescriptor.CorrelationId, recordingSubscriptionDescriptor)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err)
			}
			return
		}
		control := c.(*Control)

		// Set our state to let the caller of Poll() which triggered this know they have something
		control.Results.RecordingSubscriptionDescriptors = append(control.Results.RecordingSubscriptionDescriptors, recordingSubscriptionDescriptor)

	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
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
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			// Not much to be done here as we can't correlate
			err := fmt.Errorf("DescriptorFragmentHandler: Uncorrelated control response correlationId=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
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
			err2 := fmt.Errorf("Failed to decode recording signal: %w", err)
			if Listeners.ErrorListener != nil {
				Listeners.ErrorListener(err2)
			}
			return

		}
		if Listeners.RecordingSignalListener != nil {
			Listeners.RecordingSignalListener(recordingSignalEvent)
		}

		// If we can locate this correlationId then we can let our parent know we
		// will want an extra fragment
		c, ok := correlations.Load(recordingSignalEvent.CorrelationId)
		if !ok {
			logger.Infof("DescriptorFragmentHandler: Uncorrelated control response correlationId=%d\n%#v", recordingSignalEvent.CorrelationId, recordingSignalEvent)
			return
		}
		control := c.(*Control)
		control.Results.ExtraFragments++

	default:
		logger.Debugf("DescriptorFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
		fmt.Printf("DescriptorFragmentHandler: Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// PollNextDescriptor to poll for a fragmentLimit of Descriptors
func (control *Control) PollNextDescriptor(correlationId int64, fragmentsWanted int) error {
	logger.Debugf("PollNextDescriptor(%d) start", correlationId)
	start := time.Now()

	// Poll for events
	// As we can get async events we need to track how many/ extra events we might be wanting
	// but we're also generous in how many we request as here we only want one
	for !control.Results.IsPollComplete {
		fragmentsWanted -= control.Poll(DescriptorFragmentHandler, fragmentsWanted)

		// Adjust the fragment count for oob things like recordingsignals
		fragmentsWanted += control.Results.ExtraFragments
		control.Results.ExtraFragments = 0

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

// PollForDescriptors to poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationId int64, fragmentsWanted int32) error {
	logger.Debugf("PollForDescriptors(%d)", correlationId)

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil                  // Clear old results
	control.Results.IsPollComplete = false                 // Clear completion flag
	control.Results.ExtraFragments = 0                     // Clear extra fragment count
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
