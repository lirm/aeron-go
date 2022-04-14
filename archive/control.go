// Copyright (C) 2021-2022 Talos, Inc.
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
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer/term"
	"github.com/corymonroe-coinbase/aeron-go/archive/codecs"
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
	ErrorResponse                    error // Used by PollForErrorResponse
}

// PollContext contains the information we'll need in the image Poll()
// callback to match against our request or for async events to invoke
// the appropriate listener
type PollContext struct {
	control       *Control
	correlationID int64
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

// The current subscription handler doesn't provide a mechanism for passing context
// so we return data via the control's Results
func controlFragmentHandler(context interface{}, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	pollContext, ok := context.(*PollContext)
	if !ok {
		logger.Errorf("context conversion failed")
		return
	}

	logger.Debugf("controlFragmentHandler: correlationID:%d offset:%d length:%d header:%#v", pollContext.correlationID, offset, length, header)
	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't really tell what went wrong
		err2 := fmt.Errorf("controlFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if pollContext.control.archive.Listeners.ErrorListener != nil {
			pollContext.control.archive.Listeners.ErrorListener(err2)
		}
		return
	}

	// Look up our control
	c, ok := correlations.Load(pollContext.correlationID)
	if !ok {
		// something has gone horribly wrong and we can't correlate
		if pollContext.control.archive.Listeners.ErrorListener != nil {
			pollContext.control.archive.Listeners.ErrorListener(fmt.Errorf("failed to locate control via correlationID %d", pollContext.correlationID))
		}
		logger.Debugf("failed to locate control via correlationID %d", pollContext.correlationID)
		return
	}
	control := c.(*Control)

	switch hdr.TemplateId {
	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("controlFragmentHandler/controlResponse: Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't see what's gone wrong
			err2 := fmt.Errorf("controlFragmentHandler failed to decode control response:%w", err)
			// Call the global error handler, ugly but it's all we've got
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return
		}

		// Check this was for us
		if controlResponse.ControlSessionId == control.archive.SessionID && controlResponse.CorrelationId == pollContext.correlationID {
			// Set our state to let the caller of Poll() which triggered this know they have something
			// We're basically finished so prepare our OOB return values and log some info if we can
			logger.Debugf("controlFragmentHandler/controlResponse: received for sessionID:%d, correlationID:%d", controlResponse.ControlSessionId, controlResponse.CorrelationId)
			control.Results.ControlResponse = controlResponse
			control.Results.IsPollComplete = true
		} else {
			logger.Debugf("controlFragmentHandler/controlResponse ignoring sessionID:%d, correlationID:%d", controlResponse.ControlSessionId, controlResponse.CorrelationId)
		}

	case codecIds.recordingSignalEvent:
		var recordingSignalEvent = new(codecs.RecordingSignalEvent)

		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't really tell what went wrong
			err2 := fmt.Errorf("ControlFragmentHandler failed to decode recording signal: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
		}
		if pollContext.control.archive.Listeners.RecordingSignalListener != nil {
			pollContext.control.archive.Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	// These can happen when testing/reconnecting or if multiple clients are on the same channel/stream
	case codecIds.recordingDescriptor:
		logger.Debugf("controlFragmentHandler: ignoring RecordingDescriptor type %d", hdr.TemplateId)
	case codecIds.recordingSubscriptionDescriptor:
		logger.Debugf("controlFragmentHandler: ignoring RecordingSubscriptionDescriptor type %d", hdr.TemplateId)

	default:
		// This can happen when testing/adding new functionality
		fmt.Printf("controlFragmentHandler: Unexpected message type %d\n", hdr.TemplateId)
	}
}

// ConnectionControlFragmentHandler is the connection handling specific fragment handler.
// This mechanism only alows us to pass results back via global state which we do in control.State
func ConnectionControlFragmentHandler(context *PollContext, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ConnectionControlFragmentHandler: correlationID:%d offset:%d length: %d header: %#v", context.correlationID, offset, length, header)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("ConnectionControlFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if context.control.archive.Listeners.ErrorListener != nil {
			context.control.archive.Listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ConnectionControlFragmentHandler failed to decode control response: %w", err)
			if context.control.archive.Listeners.ErrorListener != nil {
				context.control.archive.Listeners.ErrorListener(err2)
			}
			logger.Debugf("ConnectionControlFragmentHandler failed to decode control response: %w", err)
			return
		}

		// Look this message up
		c, ok := correlations.Load(controlResponse.CorrelationId)
		if !ok {
			logger.Debugf("connectionControlFragmentHandler/controlResponse: ignoring correlationID=%d [%s]\n%#v", controlResponse.CorrelationId, string(controlResponse.ErrorMessage), controlResponse)
			return
		}
		control := c.(*Control)

		// Check this was for us
		if controlResponse.CorrelationId == context.correlationID {
			// Check result
			if controlResponse.Code != codecs.ControlResponseCode.OK {
				control.State.state = ControlStateError
				control.State.err = fmt.Errorf("Control Response failure: %s", controlResponse.ErrorMessage)
				if context.control.archive.Listeners.ErrorListener != nil {
					context.control.archive.Listeners.ErrorListener(control.State.err)
				}
				return
			}

			// assert state change
			if control.State.state != ControlStateConnectRequestSent {
				control.State.state = ControlStateError
				control.State.err = fmt.Errorf("Control Response not expecting response")
				if context.control.archive.Listeners.ErrorListener != nil {
					context.control.archive.Listeners.ErrorListener(control.State.err)
				}
			}

			// Looking good, so update state and store the SessionID
			control.State.state = ControlStateConnected
			control.State.err = nil
			control.archive.SessionID = controlResponse.ControlSessionId
		} else {
			// It's conceivable if the same application is making concurrent connection attempts using
			// the same channel/stream that we can reach here which is our parent's problem
			logger.Debugf("connectionControlFragmentHandler/controlResponse: ignoring correlationID=%d", controlResponse.CorrelationId)
		}

	case codecIds.challenge:
		var challenge = new(codecs.Challenge)

		if err := challenge.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("ControlFragmentHandler failed to decode challenge: %w", err)
			if context.control.archive.Listeners.ErrorListener != nil {
				context.control.archive.Listeners.ErrorListener(err2)
			}
		}

		logger.Infof("ControlFragmentHandler: challenge:%s, session:%d, correlationID:%d", challenge.EncodedChallenge, challenge.ControlSessionId, challenge.CorrelationId)

		// Look this message up
		c, ok := correlations.Load(challenge.CorrelationId)
		if !ok {
			logger.Debugf("connectionControlFragmentHandler/controlResponse: ignoring correlationID=%d", challenge.CorrelationId)
			return
		}
		control := c.(*Control)

		// Check this was for us
		if challenge.CorrelationId == context.correlationID {

			// Check the challenge is expected iff our option for this is not nil
			if control.archive.Options.AuthChallenge != nil {
				if !bytes.Equal(control.archive.Options.AuthChallenge, challenge.EncodedChallenge) {
					control.State.err = fmt.Errorf("ChallengeResponse Unexpected: expected:%v received:%v", control.archive.Options.AuthChallenge, challenge.EncodedChallenge)
					return
				}
			}

			// Send the response
			// Looking good, so update state and store the SessionID
			control.State.state = ControlStateChallenged
			control.State.err = nil
			control.archive.SessionID = challenge.ControlSessionId
			control.archive.Proxy.ChallengeResponse(challenge.CorrelationId, control.archive.Options.AuthResponse)
		} else {
			// It's conceivable if the same application is making concurrent connection attempts using
			// the same channel/stream that we can reach here which is our parent's problem
			logger.Debugf("connectionControlFragmentHandler/challengr: ignoring correlationID=%d", challenge.CorrelationId)
		}

	// These can happen when testing/reconnecting or if multiple clients are on the same channel/stream
	case codecIds.recordingDescriptor:
		logger.Debugf("connectionControlFragmentHandler: ignoring RecordingDescriptor type %d", hdr.TemplateId)
	case codecIds.recordingSubscriptionDescriptor:
		logger.Debugf("connectionControlFragmentHandler: ignoring RecordingSubscriptionDescriptor type %d", hdr.TemplateId)
	case codecIds.recordingSignalEvent:
		logger.Debugf("connectionControlFragmentHandler: ignoring recordingSignalEvent type %d", hdr.TemplateId)

	default:
		fmt.Printf("ConnectionControlFragmentHandler: Insert decoder for type: %d", hdr.TemplateId)
	}
}

// PollForErrorResponse polls the response stream for errors or async events.
//
// It will continue until it either receives an error or the queue is empty.
//
// If any control messages are present then they will be discarded so this
// call should not be used unless there are no outstanding operations.
//
// This may be used to check for errors, to dispatch async events, and
// to catch up on messages not for this session if for example the
// same channel and stream are in use by other sessions.
//
// Returns an error if we detect an archive operation in progress
// and a count of how many messages were consumed
func (control *Control) PollForErrorResponse() (int, error) {

	logger.Debugf("PollForErrorResponse(%d)", control.archive.SessionID)
	context := PollContext{control, 0}
	received := 0

	// Poll for async events, errors etc until the queue is drained
	for {
		ret := control.PollWithContext(
			func(buf *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
				errorResponseFragmentHandler(&context, buf, offset, length, header)
			}, 1)
		received += ret

		// If we received a response with an error then return it
		if control.Results.ErrorResponse != nil {
			return received, control.Results.ErrorResponse
		}

		// If we polled and did nothing then return
		if ret == 0 {
			return received, nil
		}
	}

	return 0, nil // Should not happen
}

// errorResponseFragmentHandler is used to check for errors and async events on an idle control
// session. Essentially:
//    ignore messages not on our session ID
//    process recordingSignalEvents
//    Log a warning if we have interrupted a synchronous event
func errorResponseFragmentHandler(context interface{}, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	pollContext, ok := context.(*PollContext)
	if !ok {
		logger.Errorf("context conversion failed")
		return
	}

	logger.Debugf("errorResponseFragmentHandler: offset:%d length: %d", offset, length)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("ConnectionControlFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly, but it's all we've got
		if pollContext.control.archive.Listeners.ErrorListener != nil {
			pollContext.control.archive.Listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("controlFragmentHandler/controlResponse: Received controlResponse")
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't see what's gone wrong
			err2 := fmt.Errorf("errorResponseFragmentHandler failed to decode control response:%w", err)
			// Call the global error handler, ugly, but it's all we've got
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return
		}

		// If this was for us then check for errors
		if controlResponse.ControlSessionId == pollContext.control.archive.SessionID {
			if controlResponse.Code == codecs.ControlResponseCode.ERROR {
				pollContext.control.Results.ErrorResponse = fmt.Errorf("PollForErrorResponse received a ControlResponse (correlationId:%d Code:ERROR error=\"%s\"", controlResponse.CorrelationId, controlResponse.ErrorMessage)
			}
		}
		return

	case codecIds.challenge:
		var challenge = new(codecs.Challenge)

		if err := challenge.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("errorResponseFragmentHandler failed to decode challenge: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
		}

		// If this was for us then that's bad
		if challenge.ControlSessionId == pollContext.control.archive.SessionID {
			pollContext.control.Results.ErrorResponse = fmt.Errorf("Received and ignoring challenge  (correlationID:%d). EerrorResponse should not be called on in parallel with sync operations", challenge.CorrelationId)
			logger.Warning(pollContext.control.Results.ErrorResponse)
			return
		}

	case codecIds.recordingDescriptor:
		var rd = new(codecs.RecordingDescriptor)

		if err := rd.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("errorResponseFragmentHandler failed to decode recordingSubscription: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
		}

		// If this was for us then that's bad
		if rd.ControlSessionId == pollContext.control.archive.SessionID {
			pollContext.control.Results.ErrorResponse = fmt.Errorf("Received and ignoring recordingDescriptor (correlationID:%d). ErrorResponse should not be called on in parallel with sync operations", rd.CorrelationId)
			logger.Warning(pollContext.control.Results.ErrorResponse)
			return
		}

	case codecIds.recordingSubscriptionDescriptor:
		var rsd = new(codecs.RecordingSubscriptionDescriptor)

		if err := rsd.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("errorResponseFragmentHandler failed to decode recordingSubscription: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
		}

		// If this was for us then that's bad
		if rsd.ControlSessionId == pollContext.control.archive.SessionID {
			pollContext.control.Results.ErrorResponse = fmt.Errorf("Received and ignoring recordingsubscriptionDescriptor (correlationID:%d). ErrorResponse should not be called on in parallel with sync operations", rsd.CorrelationId)
			logger.Warning(pollContext.control.Results.ErrorResponse)
		}

	case codecIds.recordingSignalEvent:
		var rse = new(codecs.RecordingSignalEvent)

		if err := rse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't really tell what went wrong
			err2 := fmt.Errorf("errorResponseFragmentHandler failed to decode recording signal: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
		}
		if pollContext.control.archive.Listeners.RecordingSignalListener != nil {
			pollContext.control.archive.Listeners.RecordingSignalListener(rse)
		}

		if rse.ControlSessionId == pollContext.control.archive.SessionID {
			// We can call the async callback if it exists
			if pollContext.control.archive.Listeners.RecordingSignalListener != nil {
				pollContext.control.archive.Listeners.RecordingSignalListener(rse)
			}
		}

	default:
		fmt.Printf("errorResponseFragmentHandler: Insert decoder for type: %d", hdr.TemplateId)
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

// PollWithContext provides a Poll() with a context argument
func (control *Control) PollWithContext(handler term.FragmentHandler, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil  // Clear old results
	control.Results.IsPollComplete = false // Clear completion flag
	control.Results.ErrorResponse = nil    // Clear ErrorResponse

	return control.Subscription.PollWithContext(handler, fragmentLimit)
}

// PollForResponse polls for a specific correlationID
// Returns (relevantId, nil) on success, (0 or relevantId, error) on failure
// More complex responses are contained in Control.ControlResponse after the call
func (control *Control) PollForResponse(correlationID int64, sessionID int64) (int64, error) {
	logger.Debugf("PollForResponse(%d:%d)", correlationID, sessionID)

	// Poll for events.
	//
	// As we can get async events we receive a fairly arbitrary
	// number of messages here.
	//
	// Additionally if two clients use the same channel/stream for
	// responses they will see each others messages so we just
	// continually poll until we get our response or timeout
	// without having completed and treat that as an error
	start := time.Now()
	context := PollContext{control, correlationID}

	for {
		ret := control.PollWithContext(
			func(buf *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
				controlFragmentHandler(&context, buf, offset, length, header)
			}, 10)

		// Check result
		if control.Results.IsPollComplete {
			logger.Debugf("PollForResponse(%d:%d) complete, result is %d", correlationID, sessionID, control.Results.ControlResponse.Code)
			if control.Results.ControlResponse.Code != codecs.ControlResponseCode.OK {
				err := fmt.Errorf("Control Response failure: %s", control.Results.ControlResponse.ErrorMessage)
				logger.Debug(err) // log it in debug mode as an aid to diagnosis
				return control.Results.ControlResponse.RelevantId, err
			}
			// logger.Debugf("PollForResponse(%d:%d) success", correlationID, sessionID)
			return control.Results.ControlResponse.RelevantId, nil
		}

		if control.Subscription.IsClosed() {
			return 0, fmt.Errorf("response channel from archive is not connected")
		}

		if time.Since(start) > control.archive.Options.Timeout {
			return 0, fmt.Errorf("timeout waiting for correlationID %d", correlationID)
		}

		// Idle as configured if there was nothing there
		// logger.Debugf("PollForResponse(%d:%d) idle", correlationID, sessionID)
		if ret == 0 {
			control.archive.Options.IdleStrategy.Idle(0)
		}
	}

	return 0, fmt.Errorf("PollForResponse out of loop") // can't happen
}

// DescriptorFragmentHandler is used to poll for descriptors (both recording and subscription)
// The current subscription handler doesn't provide a mechanism for passing a context
// so we return data via the control's Results
func DescriptorFragmentHandler(context interface{}, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	pollContext, ok := context.(*PollContext)
	if !ok {
		logger.Errorf("context conversion failed")
		return
	}

	// logger.Debugf("DescriptorFragmentHandler: correlationID:%d offset:%d length: %d header: %#v\n", pollContext.correlationID, offset, length, header)

	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("DescriptorFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if pollContext.control.archive.Listeners.ErrorListener != nil {
			pollContext.control.archive.Listeners.ErrorListener(err2)
		}
		return
	}

	// Look up our control
	c, ok := correlations.Load(pollContext.correlationID)
	if !ok {
		// something has gone horribly wrong and we can't correlate
		if pollContext.control.archive.Listeners.ErrorListener != nil {
			pollContext.control.archive.Listeners.ErrorListener(fmt.Errorf("failed to locate control via correlationID %d", pollContext.correlationID))
		}
		logger.Debugf("failed to locate control via correlationID %d", pollContext.correlationID)
		return
	}
	control := c.(*Control)

	switch hdr.TemplateId {
	case codecIds.recordingDescriptor:
		var recordingDescriptor = new(codecs.RecordingDescriptor)
		logger.Debugf("Received RecordingDescriptor: length %d", buf.Len())
		if err := recordingDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode RecordingDescriptor: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return
		}
		logger.Debugf("RecordingDescriptor: %#v", recordingDescriptor)

		// Check this was for us
		if recordingDescriptor.ControlSessionId == control.archive.SessionID && recordingDescriptor.CorrelationId == pollContext.correlationID {
			// Set our state to let the caller of Poll() which triggered this know they have something
			control.Results.RecordingDescriptors = append(control.Results.RecordingDescriptors, recordingDescriptor)
			control.Results.FragmentsReceived++
		} else {
			logger.Debugf("descriptorFragmentHandler/recordingDescriptor ignoring sessionID:%d, pollContext.correlationID:%d", recordingDescriptor.ControlSessionId, recordingDescriptor.CorrelationId)
		}

	case codecIds.recordingSubscriptionDescriptor:
		logger.Debugf("Received RecordingSubscriptionDescriptor: length %d", buf.Len())
		var recordingSubscriptionDescriptor = new(codecs.RecordingSubscriptionDescriptor)
		if err := recordingSubscriptionDescriptor.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode RecordingSubscriptioDescriptor: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return
		}

		// Check this was for us
		if recordingSubscriptionDescriptor.ControlSessionId == control.archive.SessionID && recordingSubscriptionDescriptor.CorrelationId == pollContext.correlationID {
			// Set our state to let the caller of Poll() which triggered this know they have something
			control.Results.RecordingSubscriptionDescriptors = append(control.Results.RecordingSubscriptionDescriptors, recordingSubscriptionDescriptor)
			control.Results.FragmentsReceived++
		} else {
			logger.Debugf("descriptorFragmentHandler/recordingSubscriptionDescriptor ignoring sessionID:%d, correlationID:%d", recordingSubscriptionDescriptor.ControlSessionId, recordingSubscriptionDescriptor.CorrelationId)
		}

	case codecIds.controlResponse:
		var controlResponse = new(codecs.ControlResponse)
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode control response: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return
		}

		// Check this was for us
		if controlResponse.ControlSessionId == control.archive.SessionID && controlResponse.CorrelationId == pollContext.correlationID {
			// Set our state to let the caller of Poll() which triggered this know they have something
			// We're basically finished so prepare our OOB return values and log some info if we can
			logger.Debugf("descriptorFragmentHandler/controlResponse: received for sessionID:%d, correlationID:%d", controlResponse.ControlSessionId, controlResponse.CorrelationId)
			control.Results.ControlResponse = controlResponse
			control.Results.IsPollComplete = true
		} else {
			logger.Debugf("descriptorFragmentHandler/controlResponse ignoring sessionID:%d, correlationID:%d", controlResponse.ControlSessionId, controlResponse.CorrelationId)
		}

		// Expected if there are no results
		if controlResponse.Code == codecs.ControlResponseCode.RECORDING_UNKNOWN {
			logger.Debugf("ControlResponse error UNKNOWN: %s", controlResponse.ErrorMessage)
			return
		}

		// Unexpected so log but we deal with it in the parent
		if controlResponse.Code == codecs.ControlResponseCode.ERROR {
			logger.Debugf("ControlResponse error ERROR: %s\n%#v", controlResponse.ErrorMessage, controlResponse)
			return
		}

	case codecIds.recordingSignalEvent:
		var recordingSignalEvent = new(codecs.RecordingSignalEvent)
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			err2 := fmt.Errorf("failed to decode recording signal: %w", err)
			if pollContext.control.archive.Listeners.ErrorListener != nil {
				pollContext.control.archive.Listeners.ErrorListener(err2)
			}
			return

		}
		if pollContext.control.archive.Listeners.RecordingSignalListener != nil {
			pollContext.control.archive.Listeners.RecordingSignalListener(recordingSignalEvent)
		}

	default:
		logger.Debug("descriptorFragmentHandler: Insert decoder for type: %d", hdr.TemplateId)
	}
}

// PollForDescriptors to poll for recording descriptors, adding them to the set in the control
func (control *Control) PollForDescriptors(correlationID int64, sessionID int64, fragmentsWanted int32) error {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = control.archive.Options.RangeChecking

	control.Results.ControlResponse = nil                  // Clear old results
	control.Results.IsPollComplete = false                 // Clear completion flag
	control.Results.RecordingDescriptors = nil             // Clear previous results
	control.Results.RecordingSubscriptionDescriptors = nil // Clear previous results
	control.Results.FragmentsReceived = 0                  // Reset our loop results count

	start := time.Now()
	descriptorCount := 0
	pollContext := PollContext{control, correlationID}

	for !control.Results.IsPollComplete {
		logger.Debugf("PollForDescriptors(%d:%d, %d)", correlationID, sessionID, int(fragmentsWanted)-descriptorCount)
		fragments := control.PollWithContext(
			func(buf *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
				DescriptorFragmentHandler(&pollContext, buf, offset, length, header)
			}, int(fragmentsWanted)-descriptorCount)
		logger.Debugf("PollWithContext(%d:%d) returned %d fragments", correlationID, sessionID, fragments)
		descriptorCount = len(control.Results.RecordingDescriptors) + len(control.Results.RecordingSubscriptionDescriptors)

		// A control response may have told us we're complete or we may have all we asked for
		if control.Results.IsPollComplete || descriptorCount >= int(fragmentsWanted) {

			logger.Debugf("PollNextDescriptor(%d:%d) complete", correlationID, sessionID)
			return nil
		}

		// Check wer're live
		if control.Subscription.IsClosed() {
			return fmt.Errorf("response channel from archive is not connected")
		}

		// Check timeout
		if time.Since(start) > control.archive.Options.Timeout {
			return fmt.Errorf("PollNextDescriptor timeout waiting for correlationID %d", correlationID)
		}

		// If we received something then loop straight away
		if fragments > 0 {
			logger.Debugf("PollForDescriptors(%d:%d) looping with %d of %d", correlationID, sessionID, control.Results.FragmentsReceived, fragmentsWanted)
			continue
		}

		// If we are yet to receive anything then idle
		if descriptorCount == 0 {
			logger.Debugf("PollForDescriptors(%d:%d) idling with %d of %d", correlationID, sessionID, control.Results.FragmentsReceived, fragmentsWanted)
			control.archive.Options.IdleStrategy.Idle(0)
		} else {
			// Otherwise we have received everything and we're done
			logger.Debugf("PollForDescriptors(%d/%d) complete at %d of %d (fragments:%d)", correlationID, sessionID, control.Results.FragmentsReceived, fragmentsWanted, fragments)
			return nil
		}

	}
	return nil
}
