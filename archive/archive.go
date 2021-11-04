// Copyright (C) 2021 Talos, Inc.
// Copyright (C) 2014-2021 Real Logic Limited.
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
	aeron        *aeron.Aeron            // Embedded aeron
	aeronContext *aeron.Context          // Embedded aeron context, see context.go for available wrapper functions
	Options      *Options                // Configuration options
	SessionId    int64                   // Allocated by the archiving media driver
	Proxy        *Proxy                  // For outgoing protocol messages (publish/request)
	Control      *Control                // For incoming protocol messages (subscribe/reponse)
	Events       *RecordingEventsAdapter // For async recording events (must be enabled)
}

// Constant values used to control behaviour of StartReplay
const (
	RecordingPositionNull = int64(-1)        // Replay a stream from the start.
	RecordingLengthNull   = int64(-1)        // Replay will follow a live recording
	RecordingLengthMax    = int64(2<<31 - 1) // Replay the whole stream
)

// replication flag used fir duplication instead of extension, see Replicate and variants
const (
	RecordingIdNullValue = int32(-1)
)

// Listeners may be set to get callbacks on various operations.  This
// is a global as internally the aeron library provides no context
// within the the FragmentAssemblers without any user data (or other
// context). Listeners.ErrorListener() if set will be called if for
// example protocol unmarshalling goes wrong.
var Listeners *ArchiveListeners

// ArchiveListeners contains all the callbacks
// By default only the ErrorListener is set to a logging listener.  If
// "archive" is at loglevel DEBUG then logging listeners are set for
// all listeners.
//
// The signal listener will be called in normal operation if set.
//
// The image listeners will be be called in normal operation if set.
//
// The RecordingEvent listeners require RecordingEventEnable() to be called
// as well as having the RecordingEvent Poll() called by user code.
type ArchiveListeners struct {
	// Called on errors for things like uncorrelated control messages
	ErrorListener func(error)

	// Async protocol events if enabled
	RecordingEventStartedListener  func(*codecs.RecordingStarted)
	RecordingEventProgressListener func(*codecs.RecordingProgress)
	RecordingEventStoppedListener  func(*codecs.RecordingStopped)

	// Async protocol event
	RecordingSignalListener func(*codecs.RecordingSignalEvent)

	// Async events from the underlying Aeron instance
	NewSubscriptionListener  func(string, int32, int64)
	NewPublicationListener   func(string, int32, int32, int64)
	AvailableImageListener   func(*aeron.Image)
	UnavailableImageListener func(*aeron.Image)
}

// LoggingErrorListener is set by default and will report internal failures when
// returning an error is not possible
func LoggingErrorListener(err error) {
	logger.Errorf("Error: %s\n", err.Error())
}

// LoggingRecordingSignalListener (called by default only in DEBUG)
func LoggingRecordingSignalListener(rse *codecs.RecordingSignalEvent) {
	logger.Infof("RecordingSignalListener, signal event is %#v\n", rse)
}

// LoggingRecordingEventStartedListener (called by default only in DEBUG)
func LoggingRecordingEventStartedListener(rs *codecs.RecordingStarted) {
	logger.Infof("RecordingEventStartedListener: %#v\n", rs)
}

// LoggingRecordingEventProgressListener (called by default only in DEBUG)
func LoggingRecordingEventProgressListener(rp *codecs.RecordingProgress) {
	logger.Infof("RecordingEventProgressListener, event is %#v\n", rp)
}

// LoggingRecordingEventStoppedListener (called by default only in DEBUG)
func LoggingRecordingEventStoppedListener(rs *codecs.RecordingStopped) {
	logger.Infof("RecordingEventStoppedListener, event is %#v\n", rs)
}

// LoggingNewSubscriptionListener from underlying aeron (called by default only in DEBUG)
func LoggingNewSubscriptionListener(channel string, stream int32, correlationId int64) {
	logger.Infof("NewSubscriptionListener(channel:%s stream:%d correlationId:%d)\n", channel, stream, correlationId)
}

// LoggingNewPublicationListener from underlying aeron (called by default only in DEBUG)
func LoggingNewPublicationListener(channel string, stream int32, session int32, regId int64) {
	logger.Infof("NewPublicationListener(channel:%s stream:%d, session:%d, regId:%d)", channel, stream, session, regId)
}

// LoggingAvailableImageListener from underlying aeron (called by default only in DEBUG)
func LoggingAvailableImageListener(image *aeron.Image) {
	logger.Infof("NewAvailableImageListener, sessionId is %d\n", image.SessionID())
}

// LoggingUnavailableImageListener from underlying aeron (called by default only in DEBUG)
func LoggingUnavailableImageListener(image *aeron.Image) {
	logger.Infof("NewUnavalableImageListener, sessionId is %d\n", image.SessionID())
}

// Also set globally (and set via the Options) is the protocol
// marshalling checks. When unmarshalling we lack context so we need a
// global copy of the options values which we set before calling
// Poll() to ensure it's current
//
// Use the Options structure to set this
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

// Utility to create a new correlation Id
func nextCorrelationId() int64 {
	return _correlationId.Inc()
}

// ReplaySessionIdToStreamId utility function to convert a ReplaySessionId into a streamId
func ReplaySessionIdToStreamId(replaySessionId int64) int32 {
	// It's actually just the least significant 32 bits
	return int32(replaySessionId)
}

// AddSessionIdToChannel utility function to add a session to a channel URI
// On failure it will return the original and an error
func AddSessionIdToChannel(channel string, sessionId int32) (string, error) {
	uri, err := aeron.ParseChannelUri(channel)
	if err != nil {
		return channel, err
	}
	uri.Set("session-id", fmt.Sprint(sessionId))
	return uri.String(), nil
}

// NewArchive factory method to create an Archive instance
// You may provide your own archive Options or otherwise one will be created from defaults
// You may provide your own aeron Context or otherwise one will be created from defaults
func NewArchive(options *Options, context *aeron.Context) (*Archive, error) {
	var err error

	archive := new(Archive)
	archive.aeron = new(aeron.Aeron)
	archive.aeronContext = context

	// Use the provided options or use our defaults
	if options != nil {
		archive.Options = options
	} else {
		if archive.Options == nil {
			// Create a new set
			archive.Options = DefaultOptions()
		}
	}

	// Set the logging levels
	logging.SetLevel(archive.Options.ArchiveLoglevel, "archive")
	logging.SetLevel(archive.Options.AeronLoglevel, "aeron")
	logging.SetLevel(archive.Options.AeronLoglevel, "memmap")
	logging.SetLevel(archive.Options.AeronLoglevel, "driver")
	logging.SetLevel(archive.Options.AeronLoglevel, "counters")
	logging.SetLevel(archive.Options.AeronLoglevel, "logbuffers")
	logging.SetLevel(archive.Options.AeronLoglevel, "buffer")
	logging.SetLevel(archive.Options.AeronLoglevel, "rb")

	// Setup the Control (subscriber/response)
	archive.Control = new(Control)
	archive.Control.archive = archive

	// Setup the Proxy (publisher/request)
	archive.Proxy = new(Proxy)
	archive.Proxy.archive = archive
	archive.Proxy.marshaller = codecs.NewSbeGoMarshaller()

	// Setup Recording Events (although it's not enabled by default)
	archive.Events = new(RecordingEventsAdapter)
	archive.Events.archive = archive

	// Create the listeners and populate
	Listeners = new(ArchiveListeners)
	Listeners.ErrorListener = LoggingErrorListener
	archive.SetAeronErrorHandler(LoggingErrorListener)

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

		archive.aeronContext.NewSubscriptionHandler(Listeners.NewSubscriptionListener)
		archive.aeronContext.NewPublicationHandler(Listeners.NewPublicationListener)
	}

	// Connect the underlying aeron
	archive.aeron, err = aeron.Connect(archive.aeronContext)
	if err != nil {
		return archive, err
	}

	// and then the subscription, it's poller and initiate a connection
	archive.Control.Subscription = <-archive.aeron.AddSubscription(archive.Options.ResponseChannel, archive.Options.ResponseStream)
	logger.Debugf("Control response subscription: %#v", archive.Control.Subscription)

	// Create the publication half for the proxy that looks after sending requests on that
	archive.Proxy.Publication = <-archive.aeron.AddExclusivePublication(archive.Options.RequestChannel, archive.Options.RequestStream)
	logger.Debugf("Proxy request publication: %#v", archive.Proxy.Publication)

	// And intitiate the connection
	archive.Control.State.state = ControlStateConnectRequestSent
	correlationId := nextCorrelationId()
	logger.Debugf("NewArchive correlationId is %d", correlationId)
	correlationsMap[correlationId] = archive.Control // Add it to our map so we can find it
	defer correlationsMapClean(correlationId)        // Clear the lookup

	// Use Auth if requested
	if archive.Options.AuthEnabled {
		if err := archive.Proxy.AuthConnectRequest(correlationId, archive.Options.ResponseStream, archive.Options.ResponseChannel, archive.Options.AuthCredentials); err != nil {
			logger.Errorf("AuthConnectRequest failed: %s\n", err)
			return nil, err
		}
	} else {
		if err := archive.Proxy.ConnectRequest(correlationId, archive.Options.ResponseStream, archive.Options.ResponseChannel); err != nil {
			logger.Errorf("ConnectRequest failed: %s\n", err)
			return nil, err
		}
	}

	start := time.Now()
	for archive.Control.State.state != ControlStateConnected && archive.Control.State.err == nil {
		fragments := archive.Control.Poll(ConnectionControlFragmentHandler, 1)
		if fragments > 0 {
			logger.Debugf("Read %d fragment(s)\n", fragments)
		}

		// Check for timeout
		if time.Since(start) > archive.Options.Timeout {
			archive.Control.State.state = ControlStateTimedOut
			archive.Control.State.err = fmt.Errorf("Operation timed out")
			break
		} else {
			archive.Options.IdleStrategy.Idle(0)
		}
	}

	if archive.Control.State.err != nil {
		logger.Errorf("Connect failed: %s\n", archive.Control.State.err)
	} else if archive.Control.State.state != ControlStateConnected {
		logger.Error("Connect failed\n")
	} else {
		logger.Infof("Archive connection established for sessionId:%d\n", archive.SessionId)

		// Store the SessionId
		sessionsMap[archive.SessionId] = archive.Control // Add it to our map so we can look it up
	}

	return archive, archive.Control.State.err
}

// Close will terminate client conductor and remove all publications and subscriptions from the media driver
func (archive *Archive) Close() error {
	archive.Proxy.CloseSessionRequest()
	archive.Proxy.Publication.Close()
	archive.Control.Subscription.Close()
	delete(sessionsMap, archive.SessionId)
	return archive.aeron.Close()
}

// SetAeronErrorHandler sets the aeron error handler
func (archive *Archive) SetAeronErrorHandler(handler func(error)) {
	archive.aeronContext.ErrorHandler(handler)
}

// SetAeronDir sets the root directory for media driver files
func (archive *Archive) SetAeronDir(dir string) {
	archive.aeronContext.AeronDir(dir)
}

// SetAeronMediaDriverTimeout sets the timeout for keep alives to media driver
func (archive *Archive) SetAeronMediaDriverTimeout(timeout time.Duration) {
	archive.aeronContext.MediaDriverTimeout(timeout)
}

// SetAeronResourceLingerTimeout sets the timeout for resource cleanup after they're released
func (archive *Archive) SetAeronResourceLingerTimeout(timeout time.Duration) {
	archive.aeronContext.ResourceLingerTimeout(timeout)
}

// SetAeronInterServiceTimeout sets the timeout for client heartbeat
func (archive *Archive) SetAeronInterServiceTimeout(timeout time.Duration) {
	archive.aeronContext.InterServiceTimeout(timeout)
}

// SetAeronPublicationConnectionTimeout sets the timeout for publications
func (archive *Archive) SetAeronPublicationConnectionTimeout(timeout time.Duration) {
	archive.aeronContext.PublicationConnectionTimeout(timeout)
}

// AeronCncFileName returns the name of the Counters file
func (archive *Archive) AeronCncFileName() string {
	return archive.aeronContext.CncFileName()
}

// EnableRecordingEvents starts recording events flowing
// Events are returned via the three callbacks which should be
// overridden from the default logging listeners defined in the Listeners
func (archive *Archive) EnableRecordingEvents() {
	archive.Events.Subscription = <-archive.aeron.AddSubscription(archive.Options.RecordingEventsChannel, archive.Options.RecordingEventsStream)
	archive.Events.Enabled = true
	logger.Debugf("RecordingEvents subscription: %#v", archive.Events.Subscription)
}

// DisableRecordingEvents stops recording events flowing
func (archive *Archive) DisableRecordingEvents() {
	archive.Events.Subscription.Close()
	archive.Events.Enabled = false
	logger.Debugf("RecordingEvents subscription closed")
}

// RecordingEventsPoll is used to poll for recording events
func (archive *Archive) RecordingEventsPoll() int {
	return archive.Events.Poll(nil, 1)
}

// AddSubscription will add a new subscription to the driver.
//
// Returns a channel, which can be used for either blocking or non-blocking wait for media driver confirmation
func (archive *Archive) AddSubscription(channel string, streamId int32) chan *aeron.Subscription {
	return archive.aeron.AddSubscription(channel, streamId)
}

// AddPublication will add a new publication to the driver.
//
// If such a publication already exists within ClientConductor the same instance
// will be returned.
//
// Returns a channel, which can be used for either blocking or
// non-blocking want for media driver confirmation
func (archive *Archive) AddPublication(channel string, streamId int32) chan *aeron.Publication {
	return archive.aeron.AddPublication(channel, streamId)
}

// AddExclusivePublication will add a new exclusive publication to the driver.
//
// If such a publication already exists within ClientConductor the
// same instance will be returned.
//
// Returns a channel, which can be used for either blocking or
// non-blocking want for media driver confirmation
func (archive *Archive) AddExclusivePublication(channel string, streamId int32) chan *aeron.Publication {
	return archive.aeron.AddExclusivePublication(channel, streamId)
}

// ClientId returns the client identity that has been allocated for communicating with the media driver.
func (archive *Archive) ClientId() int64 {
	return archive.aeron.ClientID()
}

// Clear the connections map of a correlationId. Done by a function so it can defer()ed.
func correlationsMapClean(correlationId int64) {
	delete(correlationsMap, correlationId)
}

// StartRecording a channel/stream
//
// Channels that include sessionId parameters are considered different
// than channels without sessionIds. If a publication matches both a
// sessionId specific channel recording and a non-sessionId specific
// recording, it will be recorded twice.
//
// Returns (subscriptionId, nil) or (0, error) on failure.  The
// SubscriptionId can be used in StopRecordingBySubscription()
func (archive *Archive) StartRecording(channel string, stream int32, isLocal bool, autoStop bool) (int64, error) {

	logger.Debugf("StartRecording(%s:%d)\n", channel, stream)

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StartRecordingRequest(correlationId, stream, isLocal, autoStop, channel); err != nil {
		return 0, err
	}
	return archive.Control.PollForResponse(correlationId)
}

// StopRecording can be performed by RecordingId, by SubscriptionId,
// by Publication, or by a channel/stream pairing (default)

// StopRecording by Channel and Stream
//
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
func (archive *Archive) StopRecording(channel string, stream int32) error {
	logger.Debugf("StopRecording(%s:%d)\n", channel, stream)

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingRequest(correlationId, stream, channel); err != nil {
		return err
	}
	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// StopRecordingByIdentity for the RecordingIdentity looked up in ListRecording*()
//
// Returns True if the recording was stopped or false if the recording is not currently active
// and (false, error) if something went wrong
func (archive *Archive) StopRecordingByIdentity(recordingId int64) (bool, error) {
	logger.Debugf("StopRecordingByIdentity(%d)\n", recordingId)

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingByIdentityRequest(correlationId, recordingId); err != nil {
		return false, err
	}
	res, err := archive.Control.PollForResponse(correlationId)
	if res != 0 {
		logger.Debugf("StopRecordingByIdentity result was %d\n", res)
	}

	return res == 0, err
}

// StopRecordingBySubscriptionId as returned by StartRecording
//
// Channels that include sessionId parameters are considered different than channels without sessionIds. Stopping
// recording on a channel without a sessionId parameter will not stop the recording of any sessionId specific
// recordings that use the same channel and streamId.
//
// Returns error on failure, nil on success
func (archive *Archive) StopRecordingBySubscriptionId(subscriptionId int64) error {
	logger.Debugf("StopRecordingBySubscriptionId(%d)\n", subscriptionId)

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopRecordingSubscriptionRequest(correlationId, subscriptionId); err != nil {
		return err
	}
	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// StopRecordingByPublication to stop recording a sessionId specific
// recording that pertains to the given Publication
//
// Returns error on failure, nil on success
func (archive *Archive) StopRecordingByPublication(publication aeron.Publication) error {
	channel, err := AddSessionIdToChannel(publication.Channel(), publication.SessionID())
	if err != nil {
		return err
	}
	return archive.StopRecording(channel, publication.StreamID())
}

// AddRecordedPublication to set it up forrecording.
//
// This creates a per-session recording which can fail if:
// Sending the request fails - see error for detail
// Publication.IsOriginal() is false // FIXME: check semantics
func (archive *Archive) AddRecordedPublication(channel string, stream int32) (*aeron.Publication, error) {

	// This can fail in aeron via log.Fatalf(), not much we can do
	publication := <-archive.AddPublication(channel, stream)
	if !publication.IsOriginal() {
		return nil, fmt.Errorf("publication already added for channel=%s stream=%d", channel, stream)
	}

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	sessionChannel, err := AddSessionIdToChannel(publication.Channel(), publication.SessionID())
	if err != nil {
		publication.Close()
		return nil, err
	}

	if err := archive.Proxy.StartRecordingRequest(correlationId, stream, true, false, sessionChannel); err != nil {
		publication.Close()
		return nil, err
	}

	if _, err := archive.Control.PollForResponse(correlationId); err != nil {
		publication.Close()
		return nil, err
	}

	return publication, nil
}

// ListRecordings up to recordCount recording descriptors
func (archive *Archive) ListRecordings(fromRecordingId int64, recordCount int32) ([]*codecs.RecordingDescriptor, error) {
	correlationId := nextCorrelationId()
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

// ListRecordingsForUri will list up to recordCount recording descriptors from fromRecordingId
// with a limit of recordCount for a given channel and stream.
//
// Returns the number of descriptors consumed. If fromRecordingId is
// greater than the largest known we return 0.
func (archive *Archive) ListRecordingsForUri(fromRecordingId int64, recordCount int32, channelFragment string, stream int32) ([]*codecs.RecordingDescriptor, error) {

	correlationId := nextCorrelationId()
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

// ListRecording will fetch the recording descriptor for a recordingId
//
// Returns a single recording descriptor or nil if there was no match
func (archive *Archive) ListRecording(recordingId int64) (*codecs.RecordingDescriptor, error) {
	correlationId := nextCorrelationId()
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

// StartReplay for a length in bytes of a recording from a position.
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

	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplayRequest(correlationId, recordingId, position, length, replayChannel, replayStream); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// BoundedReplay to start a replay for a length in bytes of a
// recording from a position bounded by a position counter.
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
func (archive *Archive) BoundedReplay(recordingId int64, position int64, length int64, limitCounterId int32, replayStream int32, replayChannel string) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.BoundedReplayRequest(correlationId, recordingId, position, length, limitCounterId, replayStream, replayChannel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// StopReplay for a  session.
//
// Returns error on failure, nil on success
func (archive *Archive) StopReplay(replaySessionId int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopReplayRequest(correlationId, replaySessionId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// StopAllReplays for a given recordingId
//
// Returns error on failure, nil on success
func (archive *Archive) StopAllReplays(recordingId int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopAllReplaysRequest(correlationId, recordingId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// ExtendRecording to extend an existing non-active recording of a channel and stream pairing.
//
// The channel must be configured for the initial position from which it will be extended.
//
// Returns the subscriptionId of the recording that can be passed to StopRecording()
func (archive *Archive) ExtendRecording(recordingId int64, stream int32, sourceLocation codecs.SourceLocationEnum, autoStop bool, channel string) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ExtendRecordingRequest(correlationId, recordingId, stream, sourceLocation, autoStop, channel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// GetRecordingPosition of the position recorded for an active recording.
//
// Returns the recording position or if there are no active
// recordings then RecordingPositionNull.
func (archive *Archive) GetRecordingPosition(recordingId int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.RecordingPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// TruncateRecording of a stopped recording to a given position that
// is less than the stopped position. The provided position must be on
// a fragment boundary. Truncating a recording to the start position
// effectively deletes the recording. If the truncate operation will
// result in deleting segments then this will occur
// asynchronously. Before extending a truncated recording which has
// segments being asynchronously being deleted then you should await
// completion via the RecordingSignal Delete
//
// Returns nil on success, error on failre
func (archive *Archive) TruncateRecording(recordingId int64, position int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.TruncateRecordingRequest(correlationId, recordingId, position); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// GetStartPosition for a recording.
//
// Return the start position of the recording or (0, error) on failure
func (archive *Archive) GetStartPosition(recordingId int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StartPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// GetStopPosition for a recording.
//
// Return the stop position, or RecordingPositionNull if still active.
func (archive *Archive) GetStopPosition(recordingId int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopPositionRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// FindLastMatchingRecording that matches the given criteria.
//
// Returns the RecordingId or RecordingIdNullValue if no match
func (archive *Archive) FindLastMatchingRecording(minRecordingId int64, sessionId int32, stream int32, channel string) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.FindLastMatchingRecordingRequest(correlationId, minRecordingId, sessionId, stream, channel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// ListRecordingSubscriptions to list the active recording
// subscriptions in the archive create via StartRecording or
// ExtendRecording.
//
// Returns a (possibly empty) list of RecordingSubscriptionDescriptors
func (archive *Archive) ListRecordingSubscriptions(pseudoIndex int32, subscriptionCount int32, applyStreamId bool, stream int32, channelFragment string) ([]*codecs.RecordingSubscriptionDescriptor, error) {
	correlationId := nextCorrelationId()
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

// DetachSegments from the beginning of a recording up to the
// provided new start position. The new start position must be first
// byte position of a segment after the existing start position.  It
// is not possible to detach segments which are active for recording
// or being replayed.
//
// Returns error on failure, nil on success
func (archive *Archive) DetachSegments(recordingId int64, newStartPosition int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.DetachSegmentsRequest(correlationId, recordingId, newStartPosition); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// DeleteDetachedSegments which have been previously detached from a recording.
//
// Returns the count of deleted segment files.
func (archive *Archive) DeleteDetachedSegments(recordingId int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.DeleteDetachedSegmentsRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// PurgeSegments (detach and delete) to segments from the beginning of
// a recording up to the provided new start position. The new start
// position must be first byte position of a segment after the
// existing start position. It is not possible to detach segments
// which are active for recording or being replayed.
//
// Returns the count of deleted segment files.
func (archive *Archive) PurgeSegments(recordingId int64, newStartPosition int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.PurgeSegmentsRequest(correlationId, recordingId, newStartPosition); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// AttachSegments to the beginning of a recording to restore history
// that was previously detached.
// Segment files must match the existing recording and join exactly to
// the start position of the recording they are being attached to.
//
// Returns the count of attached segment files.
func (archive *Archive) AttachSegments(recordingId int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.AttachSegmentsRequest(correlationId, recordingId); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// MigrateSegments from a source recording and attach them to the
// beginning of a destination recording.
//
// The source recording must match the destination recording for
// segment length, term length, mtu length, stream id, plus the stop
// position and term id of the source must join with the start
// position of the destination and be on a segment boundary.
//
// The source recording will be effectively truncated back to its
// start position after the migration.  Returns the count of attached
// segment files.
//
// Returns the count of attached segment files.
func (archive *Archive) MigrateSegments(recordingId int64, position int64) (int64, error) {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.MigrateSegmentsRequest(correlationId, recordingId, position); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// KeepAlive will perform a simple packet exchange with the media-driver
//
// Returns error on failure, nil on success
func (archive *Archive) KeepAlive() error {
	correlationId := nextCorrelationId()
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
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplicateRequest(correlationId, srcRecordingId, dstRecordingId, srcControlStreamId, srcControlChannel, liveDestination); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// Replicate2 will replicate a recording from a source archive to a
// destination which can be considered a backup for a primary
// archive. The source recording will be replayed via the provided
// replay channel and use the original stream id.  If the destination
// recording id is RecordingIdNullValue (-1) then a new destination
// recording is created, otherwise the provided destination recording
// id will be extended. The details of the source recording descriptor
// will be replicated.
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
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.ReplicateRequest2(correlationId, srcRecordingId, dstRecordingId, stopPosition, channelTagId, srcControlStreamId, srcControlChannel, liveDestination, replicationChannel); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// TaggedReplicate to replicate a recording from a source archive to a
// destination which can be considered a backup for a primary
// archive. The source recording will be replayed via the provided
// replay channel and use the original stream id.  If the destination
// recording id is RecordingIdNullValue (-1) then a new destination
// recording is created, otherwise the provided destination recording
// id will be extended. The details of the source recording descriptor
// will be replicated.
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
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.TaggedReplicateRequest(correlationId, srcRecordingId, dstRecordingId, channelTagId, subscriptionTagId, srcControlStreamId, srcControlChannel, liveDestination); err != nil {
		return 0, err
	}

	return archive.Control.PollForResponse(correlationId)
}

// StopReplication of a replication request
//
// Returns error on failure, nil on success
func (archive *Archive) StopReplication(replicationId int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.StopReplicationRequest(correlationId, replicationId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}

// PurgeRecording of a stopped recording, i.e. mark recording as
// Invalid and delete the corresponding segment files. The space in
// the Catalog will be reclaimed upon compaction.
//
// Returns error on failure, nil on success
func (archive *Archive) PurgeRecording(recordingId int64) error {
	correlationId := nextCorrelationId()
	correlationsMap[correlationId] = archive.Control // Set the lookup
	defer correlationsMapClean(correlationId)        // Clear the lookup

	if err := archive.Proxy.PurgeRecordingRequest(correlationId, recordingId); err != nil {
		return err
	}

	_, err := archive.Control.PollForResponse(correlationId)
	return err
}
