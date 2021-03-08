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
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/archive/codecs"
	logging "github.com/op/go-logging"
	"time"
)

// Archive is the primary interface to the media driver for managing archiving
type Archive struct {
	aeron   *aeron.Aeron
	context *ArchiveContext
	Proxy   *Proxy
	Control *Control
}

// Globals
var logger = logging.MustGetLogger("archive")
var _correlationId atomic.Long
var sessionsMap map[int64]*Control     // Used for recording signal sessionId lookup
var correlationsMap map[int64]*Control // Used for connection establishment and commands, correlationId lookup
var recordingsMap map[int64]*Control   // Used for recordings, recordingId lookup

// Inititialization
func init() {
	_correlationId.Set(time.Now().UnixNano())

	sessionsMap = make(map[int64]*Control)
	correlationsMap = make(map[int64]*Control)
	recordingsMap = make(map[int64]*Control)

	// FIXME: Provide options
	logging.SetLevel(ArchiveDefaults.ArchiveLoglevel, "archive")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "aeron")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "memmap")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "driver")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "counters")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "logbuffers")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "buffer")
	logging.SetLevel(ArchiveDefaults.AeronLoglevel, "rb")
}

func ArchiveAvailableImageHandler(image *aeron.Image) {
	logger.Infof("Archive NewAvailableImageHandler\n")
}

func ArchiveUnavailableImageHandler(image *aeron.Image) {
	logger.Infof("Archive NewUnavalableImageHandler\n")
}

// Utility function to convert a ReplaySessionId into a streamId
// It's actually just the least significant 32 bits
func ReplaySessionIdToStreamId(replaySessionId int64) int32 {
	return int32(replaySessionId)
}

// Utility function to add a session to a channel URI
// On failure it will return the original
func AddReplaySessionIdToChannel(channel string, replaySessionId int32) (string, error) {
	uri, err := aeron.ParseChannelUri(channel)
	if err != nil {
		return channel, err
	}
	uri.Set("session-id", fmt.Sprint(replaySessionId))
	return uri.String(), nil
}

// ArchiveConnect factory method to create a Archive instance from the ArchiveContext settings
func ArchiveConnect(context *ArchiveContext) (*Archive, error) {
	var err error

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
	archive.Proxy = NewProxy(control.Publication, context.IdleStrategy, control.SessionId)

	// FIXME: Java can somehow use an ephemeral port looked up here ...
	// FIXME: Java and C++ use AUTH and Challenge/Response

	// And intitiate the connection
	control.State.state = ControlStateConnectRequestSent
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = control  // Add it to our map so we can find it
	defer correlationsMapClean(correlationId) // Clear the lookup

	// Send the request and poll for the reply, giving up if we hit our timeout
	// FIXME: This can return but the image is not yet properly established so delay a little
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
	idler.Idle(0)

	if err := archive.Proxy.Connect(control.ResponseChannel, control.ResponseStream, correlationId); err != nil {
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

	// Store the SessionId in the proxy as well
	archive.Proxy.SessionId = control.SessionId
	logger.Infof("Archive connection established for sessionId:%d\n", control.SessionId)
	sessionsMap[archive.Control.SessionId] = archive.Control // Add it to our map so we can look it up

	// FIXME: Return the archive with the control intact, not sure if this the right thing to do on failure
	return archive, control.State.err
}

// Close will terminate client conductor and remove all publications and subscriptions from the media driver
func (archive *Archive) Close() error {
	archive.Proxy.CloseSession()
	archive.Control.Publication.Close()
	archive.Control.Subscription.Close()
	return archive.aeron.Close()
}

// AddSubscription will add a new subscription to the driver.
// Returns a channel, which can be used for either blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddSubscription(channel string, streamId int32) chan *aeron.Subscription {
	return archive.aeron.AddSubscription(channel, streamId)
}

// AddPublication will add a new publication to the driver. If such
// publication already exists within ClientConductor the same instance
// will be returned.  Returns a channel, which can be used for either
// blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddPublication(channel string, streamId int32) chan *aeron.Publication {
	return archive.aeron.AddPublication(channel, streamId)
}

// AddExclusivePublication will add a new exclusive publication to the driver. If such publication already
// exists within ClientConductor the same instance will be returned.
// Returns a channel, which can be used for either blocking or non-blocking want for media driver confirmation
func (archive *Archive) AddExclusivePublication(channel string, streamId int32) chan *aeron.Publication {
	return archive.aeron.AddExclusivePublication(channel, streamId)
}

// ClientId returns the client identity that has been allocated for communicating with the media driver.
func (archive *Archive) ClientId() int64 {
	return archive.aeron.ClientID()
}

// Clear the connections map of a correlationId. Done by a function so it can defer'ed
func correlationsMapClean(correlationId int64) {
	delete(correlationsMap, correlationId)
}

// Start recording a channel/stream
// Returns recordingId on success, explanatory error otherwise. Also See RecordingIdToSessionId()
func (archive *Archive) StartRecording(channel string, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool) (int64, error) {

	logger.Debugf("StartRecording(%s:%d)\n", channel, stream)
	// FIXME: locking
	// FIXME: check open

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StartRecording(correlationId, stream, sourceLocation, autoStop, channel); err != nil {
		return 0, err
	}
	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil
}

// StopRecording can be performed by RecordingId, by SubscriptionId, by Publication, or by a channel/stream pairing

// StopRecording by Channel and Stream
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
func (archive *Archive) StopRecording(channel string, stream int32) (int64, error) {
	logger.Debugf("StopRecording(%s:%d)\n", channel, stream)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecording(correlationId, stream, channel); err != nil {
		return 0, err
	}
	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil
}

// StopRecording by RecordingId as returned by StartRecording
func (archive *Archive) StopRecordingByRecordingId(recordingId int64) (int64, error) {
	logger.Debugf("StopRecordingByRecordingId(%d)\n", recordingId)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingByIdentity(correlationId, recordingId); err != nil {
		return 0, err
	}
	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil
}

// StopRecording by SubscriptionId
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
func (archive *Archive) StopRecordingBySubscriptionId(subscriptionId int64) (int64, error) {
	logger.Debugf("StopRecordingByRecordingId(%d)\n", subscriptionId)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingBySubscriptionId(correlationId, subscriptionId); err != nil {
		return 0, err
	}
	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil
}

// StopRecording by Publication
// Stop recording a sessionId specific recording that pertains to the given Publication
func (archive *Archive) StopRecordingByPublication(publication aeron.Publication) (int64, error) {
	channel, err := AddReplaySessionIdToChannel(publication.Channel(), publication.SessionID())
	if err != nil {
		return 0, err
	}
	return archive.StopRecording(channel, publication.StreamID())
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

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup
	fmt.Printf("Start recording correlationId:%d\n", correlationId)
	// FIXME: semantics of autoStop here?
	if err := archive.Proxy.StartRecording(correlationId, stream, codecs.SourceLocation.LOCAL, false, channel); err != nil {
		// FIXME: cleanup
		return nil, err
	}

	return publication, nil
}

// Get a new correlation Id
func NextCorrelationId() int64 {
	return _correlationId.Inc()
}

// List up to recordCount recording descriptors from fromRecordingId
// with a limit of recordCount for a given channel and stream
// returning the number of descriptors consumed.  If fromRecordingId
// is greater than we return 0.
func (archive *Archive) ListRecordingsForUri(fromRecordingId int64, recordCount int32, channelFragment string, stream int32) (int, error) {

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ListRecordingsForUri(correlationId, fromRecordingId, recordCount, stream, channelFragment); err != nil {
		return 0, err
	}

	if err := archive.Control.PollForDescriptors(correlationId, recordCount); err != nil {
		return 0, err
	}

	// If there's a ControlResponse let's see what transpired
	response := archive.Control.Results.ControlResponse
	if response != nil {
		switch response.Code {
		case codecs.ControlResponseCode.ERROR:
			return 0, fmt.Errorf("Response for correlationId %d (relevantId %d) failed %s", response.CorrelationId, response.CorrelationId, response.ErrorMessage)

		case codecs.ControlResponseCode.RECORDING_UNKNOWN:
			return len(archive.Control.Results.RecordingDescriptors), nil

		}
	}

	// Otherwise we can return our results
	return len(archive.Control.Results.RecordingDescriptors), nil

}

// Start a replay for a length in bytes of a recording from a position.
//
// If the position is FIXME: Java NULL_POSITION (-1) then the stream will
// be replayed from the start.
//
// If the length is FIXME: Java MAX_VALUE (2^31-1) the replay will follow a
// live recording.
//
// If the length is FIXME: Java NULL_LENGTH (-1) the replay will
// replay the whole stream of unknown length.
//
// The lower 32-bits of the returned value contains the ImageSessionId() of the received replay. All
// 64-bits are required to uniquely identify the replay when calling StopReplay(). The lower 32-bits
// can be obtained by casting the int64 value to an int32. (FIXME: provide a mechanism)
//
// Returns a ReplaySessionId - the id of the replay session which will be the same as the Image sessionId
// of the received replay for correlation with the matching channel and stream id in the lower 32 bits
func (archive *Archive) StartReplay(recordingId int64, position int64, length int64, replayChannel string, replayStream int32) (int64, error) {

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.Replay(correlationId, recordingId, position, length, replayChannel, replayStream); err != nil {
		return 0, err
	}

	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil

}

// StopReplay
func (archive *Archive) StopReplay(replaySessionId int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopReplay(correlationId, replaySessionId); err != nil {
		return 0, err
	}

	if err := archive.Control.PollForResponse(correlationId); err != nil {
		return 0, err
	}

	return archive.Control.Results.ControlResponse.RelevantId, nil
}
