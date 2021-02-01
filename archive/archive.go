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
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
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

// Globals
var logger = logging.MustGetLogger("archive")
var _correlationID atomic.Long

// Inititialization
func init() {
	_correlationID.Set(time.Now().UnixNano())
}

func ArchiveAvailableImageHandler(*aeron.Image) {
	logger.Infof("Archive NewAvailableImageHandler\n")
}

func ArchiveUnavailableImageHandler(*aeron.Image) {
	logger.Infof("Archive NewUnavalableImageHandler\n")
}

// FIXME: move
// ArchiveConnect factory method to create a Archive instance from the ArchiveContext settings
func ArchiveConnect(context *ArchiveContext) (*Archive, error) {
	var err error

	// FIXME: Provide options
	logging.SetLevel(logging.INFO, "aeron")
	logging.SetLevel(logging.DEBUG, "archive")
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

	// Setup the Control
	control := NewControl()
	archive.Control = control

	// Connect the underlying aeron
	logger.Debugf("Archive connecting with context: %v", context.aeronContext)
	archive.aeron, err = aeron.Connect(archive.context.aeronContext)
	if err != nil {
		return archive, err
	}

	// and then the subscription, it's poller and initiate a connection
	control.Subscription = <-archive.aeron.AddSubscription(control.ResponseChannel, control.ResponseStream)
	logger.Debugf("Control response subscription: %#v", control.Subscription)

	// Start the control response subscription poller
	control.StartControlSubscriptionPoller()

	// Create the publication half and the proxy that looks after sending requests on that
	control.Publication = <-archive.aeron.AddExclusivePublication(control.RequestChannel, control.RequestStream)
	logger.Debugf("Control request publication: %#v", control.Publication)
	archive.proxy = NewProxy(control.Publication, context.IdleStrategy, control.SessionID)

	// FIXME: For now we're just fobbing this off to the control handler but we'll want the original exhange
	// to be special
	// FIXME: Java can somehow use an ephemeral port looked up here ...
	// FIXME: Java and C++ use AUTH and Challenge/Response

	// And intitiate the connection
	control.State.state = ControlStateConnectRequestSent
	correlationID := NextCorrelationID()
	correlations[correlationID] = control // Add it to our map so we can find it

	// FIXME: we should retry instead of waiting
	time.Sleep(time.Second)

	if err := archive.proxy.ConnectRequest(control.ResponseChannel, control.ResponseStream, correlationID); err != nil {
		logger.Errorf("ConnectRequest failed: %s\n", err)
		return nil, err
	}
	// Wait for the response
	ControlState := <-control.ConnectionChange // FIXME: time this out
	if ControlState.err != nil || ControlState.state != ControlStateConnected {
		logger.Error("Connect failed: %s\n", err)
		return nil, err
	}
	// FIXME: delete(correlations, correlationID) // remove it from map

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

	correlationID := NextCorrelationID()
	correlations[correlationID] = archive.Control // Add it to our map so we can find it
	fmt.Printf("Start recording correlationId:%d\n", correlationID)
	if err := archive.proxy.StartRecordingRequest(channel, stream, correlationID, codecs.SourceLocation.LOCAL); err != nil {
		// FIXME: cleanup
		return nil, err
	}

	return publication, nil
}

// Get a new correlation ID
func NextCorrelationID() int64 {
	return _correlationID.Inc()
}
