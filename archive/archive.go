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
	Control *Control
}

// Globals
var logger = logging.MustGetLogger("archive")
var _correlationID atomic.Long
var sessionsMap map[int64]*Control    // Used for recording signal sessionId lookup
var connectionsMap map[int64]*Control // Used for connection establishment and commands, correlationId lookup
var recordingsMap map[int64]*Control  // Used for recordings, recordingId lookup

// Inititialization
func init() {
	_correlationID.Set(time.Now().UnixNano())

	sessionsMap = make(map[int64]*Control)
	connectionsMap = make(map[int64]*Control)
	recordingsMap = make(map[int64]*Control)
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
	logging.SetLevel(ArchiveDefaults.ArchiveLoglevel, "archive")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "aeron")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "memmap")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "driver")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "counters")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "logbuffers")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "buffer")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "rb")

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

	// Create the publication half and the proxy that looks after sending requests on that
	control.Publication = <-archive.aeron.AddExclusivePublication(control.RequestChannel, control.RequestStream)
	logger.Debugf("Control request publication: %#v", control.Publication)
	archive.proxy = NewProxy(control.Publication, context.IdleStrategy, control.SessionId)

	// FIXME: Java can somehow use an ephemeral port looked up here ...
	// FIXME: Java and C++ use AUTH and Challenge/Response

	// And intitiate the connection
	control.State.state = ControlStateConnectRequestSent
	correlationID := NextCorrelationID()
	connectionsMap[correlationID] = control // Add it to our map so we can find it

	// Send the request and poll for the reply, giving up if we hit our timeout
	if err := archive.proxy.ConnectRequest(control.ResponseChannel, control.ResponseStream, correlationID); err != nil {
		logger.Errorf("ConnectRequest failed: %s\n", err)
		return nil, err
	}

	start := time.Now()
	for control.State.state != ControlStateConnected && control.State.err == nil {
		fragments := archive.Control.Poll(ConnectionControlFragmentHandler, 1)
		if fragments > 0 {
			logger.Debugf("Read %d fragments\n", fragments)
		}

		// Check for timeout
		if time.Since(start) > ArchiveDefaults.ControlTimeout {
			control.State.state = ControlStateTimedOut
			delete(connectionsMap, correlationID) // remove it from map

		} else {
			control.IdleStrategy.Idle(0)
		}
	}

	if control.State.err != nil {
		logger.Errorf("Connect failed: %s\n", err)
	}
	if control.State.state != ControlStateConnected {
		logger.Error("Connect failed\n")
	}

	logger.Infof("Archive connection established for sessionID:%d\n", control.SessionId)
	sessionsMap[archive.Control.SessionId] = archive.Control // Add it to our map so we can look it up

	// FIXME: Return the archive with the control intact, not sure if this the right thing to do on failure
	return archive, control.State.err
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

	// FIXME: check failure
	publication := <-archive.AddPublication(channel, stream)
	if !publication.IsOriginal() {
		// FIXME: cleanup
		return nil, fmt.Errorf("publication already added for channel=%s stream=%d", channel, stream)
	}

	correlationID := NextCorrelationID()
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
