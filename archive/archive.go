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
	"github.com/lirm/aeron-go/archive/codecs"
	logging "github.com/op/go-logging"
	"time"
)

// Archive is the primary interface to the media driver for managing archiving
type Archive struct {
	aeron   *aeron.Aeron
	context *ArchiveContext
	proxy   *Proxy
	Control *Control // FIXME: one-to-one and merge in?
}

// Control contains everything required for the archive control pub/sub request/response pair
// FIXME: What does the control structure gain us here if it's one-to-one
type Control struct {
	Subscription                *aeron.Subscription
	Publication                 *aeron.Publication
	State                       ControlState
	correlationID               int64
	terminateSubscriptionPoller chan bool
	marshaller                  *codecs.SbeGoMarshaller
	sessionID                   int64
	challengeSessionID          int64 // FIXME: Todo
	subscriptionID              int64
	publicationID               int64
	isClosed                    bool // FIXME: make this a method and state
	rangeCheck                  bool
}

// An archive "connection" involves some to and fro
const ConnectionStateError = -1
const ConnectionStateNew = 0
const ConnectionStateConnectRequestSent = 1
const ConnectionStateConnectRequestOk = 2
const ConnectionStateConnected = 3

type ControlState struct {
	state int
	err   error
}

// Globals
var logger = logging.MustGetLogger("archive")

// Internal
var controlSessionID atomic.Long        // FIXME: copy driver and correlationID and hide
var correlationID atomic.Long           // FIXME: copy driver and correlationID and hide
var controls = make(map[int64]*Control) // FIXME: we should maybe have a map of correlation IDs

func ArchiveAvailableImageHandler(*aeron.Image) {
	logger.Infof("Archive NewAvailableImageHandler\n")
}

func ArchiveUnavailableImageHandler(*aeron.Image) {
	logger.Infof("Archive NewUnavalableImageHandler\n")
}

// FIXME: move
func ArchiveNewSubscriptionHandler(string, int32, int64) {
	logger.Debugf("Archive NewSubscriptionandler\n")
}

// FIXME: This will all move elsewhere into a distinct listener/poller
func ControlSubscriptionNext(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	logger.Debugf("ControlSubscriptionHandler: offset:%d length: %d header: %#v\n", offset, length, header)

	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)

	data := buffer.GetBytesArray(offset, length)

	buf := bytes.NewBuffer(data)
	marshaller := codecs.NewSbeGoMarshaller()
	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: Should we use an ErrorHandler?
		logger.Error("Failed to decode control message header", err) // Not much to be done here as we can't correlate

	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, defaults.RangeChecking); err != nil {
			logger.Error("Failed to decode control response", err) // Not much to be done here as we can't correlate
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := controls[controlResponse.CorrelationId]
		if !ok {
			logger.Error("Failed to correlate control response") // Not much to be done here as we can't correlate
		}

		// Check result
		if controlResponse.Code != codecs.ControlResponseCode.OK {
			control.State.state = ConnectionStateError
			control.State.err = fmt.Errorf("Control Response failure: %s", controlResponse.ErrorMessage)
			logger.Warning(control.State.err)
		}

		// assert state change
		if control.State.state != ConnectionStateConnectRequestSent {
			control.State.state = ConnectionStateError
			control.State.err = fmt.Errorf("Control Response not expecting response")
			logger.Error(control.State.err)
		}

		// Looking good
		control.State.state = ConnectionStateConnected

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

// ArchiveConnect factory method to create a Archive instance from the ArchiveContext settings
func ArchiveConnect(context *ArchiveContext) (*Archive, error) {
	var err error

	// FIXME: Provide options
	logging.SetLevel(logging.INFO, "aeron")
	logging.SetLevel(logging.INFO, "archive")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "driver")
	logging.SetLevel(logging.INFO, "counters")
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "buffer")
	logging.SetLevel(logging.INFO, "rb")

	archive := new(Archive)
	archive.aeron = new(aeron.Aeron)
	if context != nil {
		archive.context = NewArchiveContext()
	} else {
		archive.context = context
	}
	archive.Control = new(Control)

	// Setup the Control
	archive.Control.marshaller = codecs.NewSbeGoMarshaller()
	archive.Control.sessionID = controlSessionID.Inc()
	controls[archive.Control.sessionID] = archive.Control

	// Connect the underlying aeron
	logger.Debugf("Archive connecting with context: %v", context.aeronContext)
	archive.aeron, err = aeron.Connect(archive.context.aeronContext)
	if err != nil {
		return archive, err
	}

	// and then the subscription, it's poller and initiate a connection
	archive.Control.Subscription = <-archive.aeron.AddSubscription(archive.context.ResponseChannel, archive.context.ResponseStream)
	logger.Debugf("Control response subscription: %#v", archive.Control.Subscription)
	//  FIXME: archive.StartControlSubscriptionPoller() after connect

	// Create the publication half and the proxy that looks after sending requests on that
	archive.Control.Publication = <-archive.aeron.AddPublication(archive.context.RequestChannel, archive.context.RequestStream)
	logger.Debugf("Control request publication: %#v", archive.Control.Publication)
	archive.proxy = NewProxy(*archive.Control.Publication, context.IdleStrategy)

	// FIXME: And now might begin the connection state machine dance but we'll be cheap and simple as starting point
	// FIXME: Java can somehow use an ephemeral port looked up here ...
	// FIXME: Java and C++ use AUTH and Challenge/Response

	// Set up and send a connection request
	archive.Control.State.state = ConnectionStateConnectRequestSent
	archive.Control.correlationID = correlationID.Inc()
	if err := archive.proxy.ConnectRequest(archive.context.ResponseChannel, archive.context.ResponseStream, archive.Control.correlationID); err != nil {
		fmt.Printf("ConnectRequest failed: %s\n", err)
	}

	// FIXME: this moves
	for fragments := 0; fragments == 0; {
		fragments = archive.Control.Subscription.Poll(ControlSubscriptionNext, 10)
		if fragments > 0 {
			fmt.Printf("Read %d fragments\n", fragments)
		} else {
			// FIXME: Add default connection timeout
			// archive.context.IdleStrategy.Idle()
			idleStrategy := idlestrategy.Sleeping{SleepFor: time.Second}
			idleStrategy.Idle(1)
		}
	}

	return archive, nil
}

// Close will terminate client conductor and remove all publications and subscriptions from the media driver
func (archive *Archive) Close() error {
	return archive.aeron.Close()
}

// AddSubscription will add a new subscription to the driver.
// Returns a channel, which can be used for either blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddSubscription(channel string, streamID int32) chan *aeron.Subscription {

	return archive.aeron.AddSubscription(channel, streamID)
}

// AddPublication will add a new publication to the driver. If such
// publication already exists within ClientConductor the same instance
// will be returned.  Returns a channel, which can be used for either
// blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddPublication(channel string, streamID int32) chan *aeron.Publication {

	return archive.aeron.AddPublication(channel, streamID)
}

// AddExclusivePublication will add a new exclusive publication to the driver. If such publication already
// exists within ClientConductor the same instance will be returned.
// Returns a channel, which can be used for either blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddExclusivePublication(channel string, streamID int32) chan *aeron.Publication {
	return archive.aeron.AddExclusivePublication(channel, streamID)
}

// NextCorrelationID generates the next correlation id that is unique for the connected Media Driver.
// This is useful generating correlation identifiers for pairing requests with responses in a clients own
// application protocol.
//
// This method is thread safe and will work across processes that all use the same media driver.
func (archive *Archive) NextCorrelationID() int64 {
	return archive.aeron.NextCorrelationID()
}

// ClientID returns the client identity that has been allocated for communicating with the media driver.
func (archive *Archive) ClientID() int64 {
	return archive.aeron.ClientID()
}

// Add a Recorded Publication and set it up to be recorded.

// This can fail if:
//   Publication.IsOriginal() is false // FIXME: check semantics
//   Sending the request fails - see error for detail
//
// FIXME: Formalize the error handling
func (archive *Archive) AddRecordedPublication(channel string, stream int32) (*aeron.Publication, error) {

	// FIXME: can that fail?
	publication := <-archive.AddPublication(channel, stream)
	if !publication.IsOriginal() {
		// FIXME: cleanup
		return nil, fmt.Errorf("publication already added for channel=%s stream=%d", channel, stream)
	}

	if err := archive.proxy.StartRecordingRequest(channel, stream, archive.NextCorrelationID(), codecs.SourceLocation.LOCAL); err != nil {
		// FIXME: cleanup
		return nil, err
	}

	return publication, nil
}
