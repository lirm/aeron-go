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
	Context *ArchiveContext
	Proxy   *Proxy
	Control *Control
	Events  *RecordingEventsAdapter
}

// constants relating to StartReplay
const RecordingPositionNull = int64(-1)     // Replay the stream from the start.
const RecordingLengthNull = int64(-1)       // Replay will follow a live recording
const RecordingLengthMax = int64(2<<31 - 1) // Replay the whole stream

// constants used elsewhere
const RecordingIdNullValue = int32(-1) // Java's io.aeron.Aeron#NULL_VALUE

// By default all but one of these callbacks are active, and all need to be
// set to user functions to be invoked. This can be done at any time
//
// If the loglevel is set to DEBUG, then all of the default listeners
// will be set to logging listeners.
//
// The Signal Listener if set will be called in normal operation
//
// The Image listeners will be be called in normal operation if set
//
// The ReccordingEvent listeners require RecordingEventEnable() to be called
// as well as having the RecordingEvent Poll() called by user code
type ArchiveListeners struct {
	RecordingEventStartedListener  func(*codecs.RecordingStarted)
	RecordingEventProgressListener func(*codecs.RecordingProgress)
	RecordingEventStoppedListener  func(*codecs.RecordingStopped)

	RecordingSignalListener func(*codecs.RecordingSignalEvent)

	AvailableImageListener   func(*aeron.Image)
	UnavailableImageListener func(*aeron.Image)

	NewSubscriptionListener func(string, int32, int64)
	NewPublicationListener  func(string, int32, int32, int64)
}

// Some Listeners that log for convenience/debug
func LoggingAvailableImageListener(image *aeron.Image) {
	logger.Infof("NewAvailableImageListener, sessionId is %d\n", image.SessionID())
}

func LoggingUnavailableImageListener(image *aeron.Image) {
	logger.Infof("NewUnavalableImageListener, sessionId is %d\n", image.SessionID())
}

func LoggingRecordingSignalListener(rse *codecs.RecordingSignalEvent) {
	logger.Infof("RecordingSignalListener, signal event is %#v\n", rse)
}

func LoggingRecordingEventStartedListener(rs *codecs.RecordingStarted) {
	logger.Infof("RecordingEventStartedListener: %#v\n", rs)
}

func LoggingRecordingEventProgressListener(rp *codecs.RecordingProgress) {
	logger.Infof("RecordingEventProgressListener, event is %#v\n", rp)
}

func LoggingRecordingEventStoppedListener(rs *codecs.RecordingStopped) {
	logger.Infof("RecordingEventStoppedListener, event is %#v\n", rs)
}

func LoggingNewSubscriptionListener(channel string, stream int32, correlationId int64) {
	logger.Infof("NewSubscriptionListener(channel:%s stream:%d correlationId:%d)\n", channel, stream, correlationId)
}

func LoggingNewPublicationListener(channel string, stream int32, session int32, regId int64) {
	logger.Infof("NewPublicationListener(channel:%s stream:%d, session:%d, regId:%d)", channel, stream, session, regId)
}

// Listeners may be set to get callbacks on various operations.
// This global as the aeron library calls the FragmentAssemblers without any user
// data (or other context).
var Listeners *ArchiveListeners

// Also set globally (and set via the Options) is the protocol
// marshalling checks. If protocol marshaling goes wrong we lack
// context so it needs to be global.
var rangeChecking bool

// Other globals used internally
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
}

// Get a new correlation Id
func NextCorrelationId() int64 {
	return _correlationId.Inc()
}

// Utility function to convert a ReplaySessionId into a streamId
func ReplaySessionIdToStreamId(replaySessionId int64) int32 {
	// It's actually just the least significant 32 bits
	return int32(replaySessionId)
}

// Utility function to add a session to a channel URI
// On failure it will return the original and an error
func AddSessionIdToChannel(channel string, sessionId int32) (string, error) {
	uri, err := aeron.ParseChannelUri(channel)
	if err != nil {
		return channel, err
	}
	uri.Set("session-id", fmt.Sprint(sessionId))
	return uri.String(), nil
}

// ArchiveConnect factory method to create a Archive instance from the ArchiveContext settings
// You may provide your own archive context which may include an aeron context
func NewArchive(context *ArchiveContext, options *Options) (*Archive, error) {
	var err error

	archive := new(Archive)
	archive.aeron = new(aeron.Aeron)

	// Use they're context or allocate a default one for them
	if context != nil {
		archive.Context = context
	} else {
		archive.Context = NewArchiveContext()
	}

	// Use the provided options or use our defaults
	if options != nil {
		archive.Context.Options = options
	} else {
		if archive.Context.Options == nil {
			// Create a new set
			archive.Context.Options = DefaultOptions()
		}
	}

	// FIXME: strip once development complete
	archive.Context.Options.RangeChecking = true

	logging.SetLevel(archive.Context.Options.ArchiveLoglevel, "archive")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "aeron")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "memmap")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "driver")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "counters")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "logbuffers")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "buffer")
	logging.SetLevel(archive.Context.Options.AeronLoglevel, "rb")

	// Setup the Control (subscriber/response)
	archive.Control = NewControl(context)

	// Setup the Proxy (publisher/request)
	archive.Proxy = NewProxy(context)

	// Setup Recording Events (although it's not enabled by default)
	archive.Events = NewRecordingEventsAdapter(context)

	Listeners = new(ArchiveListeners)
	// In Debug mode initialize our listeners with simple loggers
	// Note that these actually log at INFO so you can do this manually for INFO if you like
	if logging.GetLevel("archive") >= logging.DEBUG {
		logger.Debugf("Setting logging listeners")

		Listeners.RecordingEventStartedListener = LoggingRecordingEventStartedListener
		Listeners.RecordingEventProgressListener = LoggingRecordingEventProgressListener
		Listeners.RecordingEventStoppedListener = LoggingRecordingEventStoppedListener

		Listeners.RecordingSignalListener = LoggingRecordingSignalListener

		Listeners.AvailableImageListener = LoggingAvailableImageListener
		Listeners.UnavailableImageListener = LoggingUnavailableImageListener

		Listeners.NewSubscriptionListener = LoggingNewSubscriptionListener
		Listeners.NewPublicationListener = LoggingNewPublicationListener

		archive.Context.aeronContext.NewSubscriptionHandler(Listeners.NewSubscriptionListener)
		archive.Context.aeronContext.NewPublicationHandler(Listeners.NewPublicationListener)
	}

	// Connect the underlying aeron
	logger.Debugf("Archive connecting with context: %v", context.aeronContext)
	archive.aeron, err = aeron.Connect(archive.Context.aeronContext)
	if err != nil {
		return archive, err
	}

	// and then the subscription, it's poller and initiate a connection
	archive.Control.Subscription = <-archive.aeron.AddSubscription(archive.Context.Options.ResponseChannel, archive.Context.Options.ResponseStream)
	logger.Debugf("Control response subscription: %#v", archive.Control.Subscription)

	// Create the publication half for the proxy that looks after sending requests on that
	archive.Proxy.Publication = <-archive.aeron.AddExclusivePublication(archive.Context.Options.RequestChannel, archive.Context.Options.RequestStream)
	logger.Debugf("Proxy request publication: %#v", archive.Proxy.Publication)

	// FIXME: Java can somehow use an ephemeral port looked up here ...
	// FIXME: Java and C++ use AUTH and Challenge/Response

	// And intitiate the connection
	archive.Control.State.state = ControlStateConnectRequestSent
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Add it to our map so we can find it
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ConnectRequest(archive.Context.Options.ResponseChannel, archive.Context.Options.ResponseStream, correlationId); err != nil {
		logger.Errorf("ConnectRequest failed: %s\n", err)
		return nil, err
	}

	start := time.Now()
	for archive.Control.State.state != ControlStateConnected && archive.Control.State.err == nil {
		fragments := archive.Control.Poll(ConnectionControlFragmentHandler, 1)
		if fragments > 0 {
			logger.Debugf("Read %d fragments\n", fragments)
		}

		// Check for timeout
		if time.Since(start) > archive.Context.Options.Timeout {
			archive.Control.State.state = ControlStateTimedOut
		} else {
			archive.Context.Options.IdleStrategy.Idle(0)
		}
	}

	if archive.Control.State.err != nil {
		logger.Errorf("Connect failed: %s\n", err)
	}
	if archive.Control.State.state != ControlStateConnected {
		logger.Error("Connect failed\n")
	}

	// Store the SessionId in the proxy as well
	logger.Infof("Archive connection established for sessionId:%d\n", archive.Context.SessionId)
	sessionsMap[archive.Context.SessionId] = archive.Control // Add it to our map so we can look it up

	// FIXME: Return the archive with the control intact, not sure if this the right thing to do on failure
	return archive, archive.Control.State.err
}

// Close will terminate client conductor and remove all publications and subscriptions from the media driver
func (archive *Archive) Close() error {
	archive.Proxy.CloseSessionRequest()
	archive.Proxy.Publication.Close()
	archive.Control.Subscription.Close()
	delete(sessionsMap, archive.Context.SessionId)
	return archive.aeron.Close()
}

// Start recording events flowing
// Events are returned via the three callbacks which should be
// overridden from the default logging listeners defined in the Listeners
func (archive *Archive) EnableRecordingEvents() {
	archive.Events.Subscription = <-archive.aeron.AddSubscription(archive.Context.Options.RecordingEventsChannel, archive.Context.Options.RecordingEventsStream)
	archive.Events.Enabled = true
	logger.Debugf("RecordingEvents subscription: %#v", archive.Events.Subscription)
}

// Stop recording events flowing
func (archive *Archive) DisableRecordingEvents() {
	archive.Events.Subscription.Close()
	archive.Events.Enabled = false
	logger.Debugf("RecordingEvents subscription closed")
}

// Poll for recording events
func (archive *Archive) RecordingEventsPoll() int {
	return archive.Events.Poll(nil, 1)
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
//
// Channels that include sessionId parameters are considered different
// than channels without sessionIds. If a publication matches both a
// sessionId specific channel recording and a non-sessionId specific
// recording, it will be recorded twice.
// Returns (subscriptionId, nil) or (0, error) on failure.
// The SubscriptionId can be used in StopRecordingBySubscription()
// FIXME: test StopRecordingBySubscription()
func (archive *Archive) StartRecording(channel string, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool) (int64, error) {

	logger.Debugf("StartRecording(%s:%d)\n", channel, stream)
	// FIXME: locking
	// FIXME: check open

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StartRecordingRequest(correlationId, stream, sourceLocation, autoStop, channel); err != nil {
		return 0, err
	}
	return archive.Control.PollForResponse(correlationId)
}

// StopRecording can be performed by RecordingId, by SubscriptionId, by Publication, or by a channel/stream pairing (default)

// StopRecording by Channel and Stream
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
func (archive *Archive) StopRecording(channel string, stream int32) error {
	logger.Debugf("StopRecording(%s:%d)\n", channel, stream)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingRequest(correlationId, stream, channel); err != nil {
		return err
	}
	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// StopRecording by RecordingId as looked up in ListRecording*()
// Returns True if the recording was stopped or false if the recording is not currently active
// and (false, error) if something went wrong
func (archive *Archive) StopRecordingByIdentity(recordingId int64) (bool, error) {
	logger.Debugf("StopRecordingByIdentity(%d)\n", recordingId)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingByIdentityRequest(correlationId, recordingId); err != nil {
		return false, err
	}
	res, err := archive.Control.PollForResponse(correlationId)
	return res != 0, err
}

// StopRecording by SubscriptionId
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
// Returns error on failure, nil on success
func (archive *Archive) StopRecordingBySubscriptionId(subscriptionId int64) error {
	logger.Debugf("StopRecordingBySubscriptionId(%d)\n", subscriptionId)

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingSubscriptionRequest(correlationId, subscriptionId); err != nil {
		return err
	}
	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// StopRecording by Publication
// Stop recording a sessionId specific recording that pertains to the given Publication
// Returns error on failure, nil on success
func (archive *Archive) StopRecordingByPublication(publication aeron.Publication) error {
	channel, err := AddSessionIdToChannel(publication.Channel(), publication.SessionID())
	if err != nil {
		return err
	}
	return archive.StopRecording(channel, publication.StreamID())
}

// Add a Recorded Publication and set it up to be recorded.
// This creates a per-session recording
//
// This can fail if:
//   Publication.IsOriginal() is false // FIXME: check semantics
//   Sending the request fails - see error for detail
func (archive *Archive) AddRecordedPublication(channel string, stream int32) (*aeron.Publication, error) {

	// FIXME: check failure
	publication := <-archive.AddPublication(channel, stream)
	if !publication.IsOriginal() {
		return nil, fmt.Errorf("publication already added for channel=%s stream=%d", channel, stream)
	}

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup
	fmt.Printf("Start recording correlationId:%d\n", correlationId)

	sessionChannel, err := AddSessionIdToChannel(publication.Channel(), publication.SessionID())
	if err != nil {
		publication.Close()
		return nil, err
	}

	if err := archive.Proxy.StartRecordingRequest(correlationId, stream, codecs.SourceLocation.LOCAL, false, sessionChannel); err != nil {
		publication.Close()
		return nil, err
	}

	if _, err := archive.Control.PollForResponse(correlationId); err != nil {
		publication.Close()
		return nil, err
	}

	return publication, nil
}

// List up to recordCount recording descriptors
func (archive *Archive) ListRecordings(fromRecordingId int64, recordCount int32) ([]*codecs.RecordingDescriptor, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ListRecordingsRequest(correlationId, fromRecordingId, recordCount); err != nil {
		return nil, err
	}
	if err := archive.Control.PollForDescriptors(correlationId, recordCount); err != nil {
		return nil, err
	}

	// If there's a ControlResponse let's see what transpired
	response := archive.Control.Results.ControlResponse
	if response != nil {
		switch response.Code {
		case codecs.ControlResponseCode.ERROR:
			return nil, fmt.Errorf("Response for correlationId %d (relevantId %d) failed %s", response.CorrelationId, response.RelevantId, response.ErrorMessage)

		case codecs.ControlResponseCode.RECORDING_UNKNOWN:
			return archive.Control.Results.RecordingDescriptors, nil
		}
	}

	// Otherwise we can return our results
	return archive.Control.Results.RecordingDescriptors, nil
}

// List up to recordCount recording descriptors from fromRecordingId
// with a limit of recordCount for a given channel and stream
// returning the number of descriptors consumed.  If fromRecordingId
// is greater than we return 0.
func (archive *Archive) ListRecordingsForUri(fromRecordingId int64, recordCount int32, channelFragment string, stream int32) ([]*codecs.RecordingDescriptor, error) {

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ListRecordingsForUriRequest(correlationId, fromRecordingId, recordCount, stream, channelFragment); err != nil {
		return nil, err
	}

	if err := archive.Control.PollForDescriptors(correlationId, recordCount); err != nil {
		return nil, err
	}

	// If there's a ControlResponse let's see what transpired
	response := archive.Control.Results.ControlResponse
	if response != nil {
		switch response.Code {
		case codecs.ControlResponseCode.ERROR:
			return nil, fmt.Errorf("Response for correlationId %d (relevantId %d) failed %s", response.CorrelationId, response.RelevantId, response.ErrorMessage)

		case codecs.ControlResponseCode.RECORDING_UNKNOWN:
			return archive.Control.Results.RecordingDescriptors, nil
		}
	}

	// Otherwise we can return our results
	return archive.Control.Results.RecordingDescriptors, nil
}

// Grab the recording descriptor for a recordingId
// Returns a single recording descriptor or nil if there was no match
func (archive *Archive) ListRecording(recordingId int64) (*codecs.RecordingDescriptor, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ListRecordingRequest(correlationId, recordingId); err != nil {
		return nil, err
	}
	if err := archive.Control.PollForDescriptors(correlationId, 1); err != nil {
		return nil, err
	}

	// If there's a ControlResponse let's see what transpired
	response := archive.Control.Results.ControlResponse
	if response != nil {
		switch response.Code {
		case codecs.ControlResponseCode.ERROR:
			return nil, fmt.Errorf("Response for correlationId %d (relevantId %d) failed %s", response.CorrelationId, response.RelevantId, response.ErrorMessage)

		case codecs.ControlResponseCode.RECORDING_UNKNOWN:
			return nil, nil
		}
	}

	// Otherwise we can return our results
	return archive.Control.Results.RecordingDescriptors[0], nil
}

// Start a replay for a length in bytes of a recording from a position.
//
// If the position is RecordingPositionNull (-1) then the stream will
// be replayed from the start.
//
// If the length is RecordingLengthMax (2^31-1) the replay will follow
// a live recording.
//
// If the length is RecordingLengthNull (-1) the replay will
// replay the whole stream of unknown length.
//
// The lower 32-bits of the returned value contains the ImageSessionId() of the received replay. All
// 64-bits are required to uniquely identify the replay when calling StopReplay(). The lower 32-bits
// can be obtained by casting the int64 value to an int32. See ReplaySessionIdToStreamId() helper.
//
// Returns a ReplaySessionId - the id of the replay session which will be the same as the Image sessionId
// of the received replay for correlation with the matching channel and stream id in the lower 32 bits
func (archive *Archive) StartReplay(recordingId int64, position int64, length int64, replayChannel string, replayStream int32) (int64, error) {

	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplayRequest(correlationId, recordingId, position, length, replayChannel, replayStream); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Start a replay for a length in bytes of a recording from a position
// bounded by a position counter. If the position is
// RecordingPositionNull then the stream will be replayed from the
// start.  The lower 32-bits of the returned value contains the
// sessionId of the received replay. All 64-bits are required to
// uniquely identify the replay when calling StopReplay The lower
// 32-bits can be obtained by casting casting the int64 value to an
// int32. See ReplaySessionIdToStreamId() helper.
// Returns a ReplaySessionId - the id of the replay session which will
// be the same as the Image sessionId of the received replay for
// correlation with the matching channel and stream id in the lower 32
// bits.
func (archive *Archive) BoundedReplay(recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.BoundedReplayRequest(correlationId, recordingId, position, length, limitCounterId, replayStream, replayChannel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Stop a Replay session
// Returns error on failure, nil on success
func (archive *Archive) StopReplay(replaySessionId int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopReplayRequest(correlationId, replaySessionId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// Stop all Replays for a given recordingId
func (archive *Archive) StopAllReplays(recordingId int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopAllReplaysRequest(correlationId, recordingId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// Extend an existing nont-active recording of a channel and stream
// pairing. The channel must be configured for the initial position
// from which it will be extended. This can be done with FIXME:
// ChannelUriStringBuilder#initialPosition(long, int, int). The
// details required to initialise can be found by calling FIXME:
// listRecording(long, RecordingDescriptorConsumer).
//
// Returns the subscriptionId of the recording that can be passed to
// StopRecording()
func (archive *Archive) ExtendRecording(recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ExtendRecordingRequest(correlationId, recordingId, stream, sourceLocation, autoStop, channel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Get the position recorded for an active recording. If no active
// recording then return RecordingPositionNull.
func (archive *Archive) GetRecordingPosition(recordingId int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.RecordingPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Truncate a stopped recording to a given position that is less than
// the stopped position. The provided position must be on a fragment
// boundary. Truncating a recording to the start position effectively
// deletes the recording.  If the truncate operation will result in
// deleting segments then this will occur asynchronously. Before
// extending a truncated recording which has segments being
// asynchronously being deleted then you should await completion
// via the RecordingSignal Delete
func (archive *Archive) TruncateRecording(recordingId int64, position int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.TruncateRecordingRequest(correlationId, recordingId, position); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// Get the start position for a recording.
// Return the start position of the recording
func (archive *Archive) GetStartPosition(recordingId int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StartPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Get the stop position for a recording.
// Return the stop position, or RecordingPositionNull if still active.
func (archive *Archive) GetStopPosition(recordingId int64, position int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Find the last recording that matches the given criteria.
func (archive *Archive) FindLastMatchingRecording(minRecordingId int64, sessionId int32, stream int32, channel string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.FindLastMatchingRecordingRequest(correlationId, minRecordingId, sessionId, stream, channel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// List active recording subscriptions in the archive create via StartRecording or ExtendRecording.
// Returns a (possibly empty) list of RecordingSubscriptionDescriptors
func (archive *Archive) ListRecordingSubscriptions(pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channelFragment string) ([]*codecs.RecordingSubscriptionDescriptor, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ListRecordingSubscriptionsRequest(correlationId, pseudoIndex, subscriptionCount, applyStreamId, stream, channelFragment); err != nil {
		return nil, err
	}

	if err := archive.Control.PollForDescriptors(correlationId, subscriptionCount); err != nil {
		return nil, err
	}

	// If there's a ControlResponse let's see what transpired
	response := archive.Control.Results.ControlResponse
	if response != nil {
		switch response.Code {
		case codecs.ControlResponseCode.ERROR:
			return nil, fmt.Errorf("Response for correlationId %d (relevantId %d) failed %s", response.CorrelationId, response.RelevantId, response.ErrorMessage)

		case codecs.ControlResponseCode.SUBSCRIPTION_UNKNOWN:
			return archive.Control.Results.RecordingSubscriptionDescriptors, nil
		}
	}

	// Otherwise we can return our results
	return archive.Control.Results.RecordingSubscriptionDescriptors, nil

}

// Detach segments from the beginning of a recording up to the
// provided new start position. The new start position must be first
// byte position of a segment after the existing start position.  It
// is not possible to detach segments which are active for recording
// or being replayed.
func (archive *Archive) DetachSegments(recordingId int64, newStartPosition int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.DetachSegmentsRequest(correlationId, recordingId, newStartPosition); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// Delete segments which have been previously detached from a recording.
// Returns the count of deleted segment files.
func (archive *Archive) DeleteDetachedSegments(recordingId int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.DeleteDetachedSegmentsRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Purge (detach and delete) segments from the beginning of a
// recording up to the provided new start position. The new start
// position must be first byte position of a segment after the
// existing start position. It is not possible to detach segments
// which are active for recording or being replayed.
// Returns the count of deleted segment files.
func (archive *Archive) PurgeSegments(recordingId int64, newStartPosition int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.PurgeSegmentsRequest(correlationId, recordingId, newStartPosition); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Attach segments to the beginning of a recording to restore history
// that was previously detached.
// Segment files must match the existing recording and join exactly to
// the start position of the recording they are being attached to.
// Returns the count of attached segment files.
func (archive *Archive) AttachSegments(recordingId int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.AttachSegmentsRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Migrate segments from a source recording and attach them to the
// beginning of a destination recording.
// The source recording must match the destination recording for
// segment length, term length, mtu length, stream id, plus the stop
// position and term id of the source must join with the start
// position of the destination and be on a segment boundary.
// The source recording will be effectively truncated back to its
// start position after the migration.  Returns the count of attached
// segment files.
// Returns the count of attached segment files.
func (archive *Archive) MigrateSegments(recordingId int64, position int64) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.MigrateSegmentsRequest(correlationId, recordingId, position); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// KeepAlive
func (archive *Archive) KeepAlive() error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	return archive.Proxy.KeepAliveRequest(correlationId)
}

// Replicate a recording from a source archive to a destination which
// can be considered a backup for a primary archive. The source
// recording will be replayed via the provided replay channel and use
// the original stream id.  If the destination recording id is
// RecordingIdNullValue (-1) then a new destination recording is
// created, otherwise the provided destination recording id will be
// extended. The details of the source recording descriptor will be
// replicated.
//
// For a source recording that is still active the replay can merge
// with the live stream and then follow it directly and no longer
// require the replay from the source. This would require a multicast
// live destination.
//
// Errors will be reported asynchronously and can be checked for with
// PollForErrorResponse() Follow progress with the RecordingSignalAdapter
//
// srcRecordingId     recording id which must exist in the source archive.
// dstRecordingId     recording to extend in the destination, otherwise {@link io.aeron.Aeron#NULL_VALUE}.
// srcControlStreamId remote control stream id for the source archive to instruct the replay on.
// srcControlChannel  remote control channel for the source archive to instruct the replay on.
// liveDestination    destination for the live stream if merge is required. nil for no merge.
//
// Returns the replication session id which can be passed StopReplication()
func (archive *Archive) Replicate(srcRecordingId int64, dstRecordingId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplicateRequest(correlationId, srcRecordingId, dstRecordingId, srcControlStreamId, srcControlChannel, liveDestination); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Replicate a recording from a source archive to a destination which
// can be considered a backup for a primary archive. The source
// recording will be replayed via the provided replay channel and use
// the original stream id.  If the destination recording id is
// RecordingIdNullValue (-1) then a new destination recording is
// created, otherwise the provided destination recording id will be
// extended. The details of the source recording descriptor will be
// replicated.
//
// For a source recording that is still active the replay can merge
// with the live stream and then follow it directly and no longer
// require the replay from the source. This would require a multicast
// live destination.
//
// Errors will be reported asynchronously and can be checked for with
// PollForErrorResponse() Follow progress with the RecordingSignalAdapter
//
// srcRecordingId     recording id which must exist in the source archive.
// dstRecordingId     recording to extend in the destination, otherwise {@link io.aeron.Aeron#NULL_VALUE}.
// stopPosition       position to stop the replication. RecordingPositionNull to stop at end of current recording.
// srcControlStreamId remote control stream id for the source archive to instruct the replay on.
// srcControlChannel  remote control channel for the source archive to instruct the replay on.
// liveDestination    destination for the live stream if merge is required. nil for no merge.
// replicationChannel channel over which the replication will occur. Empty or null for default channel.
//
// Returns the replication session id which can be passed StopReplication()
func (archive *Archive) Replicate2(srcRecordingId int64, dstRecordingId int64, stopPosition int64, channelTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string, replicationChannel string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplicateRequest2(correlationId, srcRecordingId, dstRecordingId, stopPosition, channelTagId, srcControlStreamId, srcControlChannel, liveDestination, replicationChannel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Replicate a recording from a source archive to a destination which
// can be considered a backup for a primary archive. The source
// recording will be replayed via the provided replay channel and use
// the original stream id.  If the destination recording id is
// RecordingIdNullValue (-1) then a new destination recording is
// created, otherwise the provided destination recording id will be
// extended. The details of the source recording descriptor will be
// replicated.
//
// The subscription used in the archive will be tagged
// with the provided tags. For a source recording that is still active
// the replay can merge with the live stream and then follow it
// directly and no longer require the replay from the source. This
// would require a multicast live destination.
//
// Errors will be reported asynchronously and can be checked for with
// PollForErrorResponse() Follow progress with the RecordingSignalAdapter
//
// srcRecordingId     recording id which must exist in the source archive.
// dstRecordingId     recording to extend in the destination, otherwise {@link io.aeron.Aeron#NULL_VALUE}.
// channelTagId       used to tag the replication subscription.
// subscriptionTagId  used to tag the replication subscription.
// srcControlStreamId remote control stream id for the source archive to instruct the replay on.
// srcControlChannel  remote control channel for the source archive to instruct the replay on.
// liveDestination    destination for the live stream if merge is required. nil for no merge.
//
// Returns the replication session id which can be passed StopReplication()
func (archive *Archive) TaggedReplicate(srcRecordingId int64, dstRecordingId int64, channelTagId int64, subscriptionTagId int64, srcControlStreamId int32, srcControlChannel string, liveDestination string) (int64, error) {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.TaggedReplicateRequest(correlationId, srcRecordingId, dstRecordingId, channelTagId, subscriptionTagId, srcControlStreamId, srcControlChannel, liveDestination); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Stop a replication request
func (archive *Archive) StopReplication(replicationId int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopReplicationRequest(correlationId, replicationId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// Purge a stopped recording, i.e. mark recording as Invalid and
// delete the corresponding segment files. The space in the Catalog
// will be reclaimed upon compaction.
func (archive *Archive) PurgeRecording(recordingId int64) error {
	correlationId := NextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.PurgeRecordingRequest(correlationId, recordingId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}
