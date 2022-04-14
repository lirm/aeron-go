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
	"flag"
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logging"
	"github.com/corymonroe-coinbase/aeron-go/archive/codecs"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// Rather than mock or spawn an archive-media-driver we're just seeing
// if we can connect to one and if we can we'll run some tests. If the
// init fails to connect then we'll skip the tests
// FIXME:BiggerPicture this plan fails as aeron-go calls log.Fatalf() if the media driver is not running!
var archive *Archive
var haveArchive = false
var DEBUG = false

type TestCases struct {
	sampleStream  int32
	sampleChannel string
	replayStream  int32
	replayChannel string
}

var testCases = []TestCases{
	{int32(*TestConfig.SampleStream), *TestConfig.SampleChannel, int32(*TestConfig.ReplayStream), *TestConfig.ReplayChannel},
}

// For testing async events
type TestCounters struct {
	recordingSignalCount        int
	recordingEventStartedCount  int
	recordingEventProgressCount int
	recordingEventStoppedCount  int
}

var testCounters TestCounters

func RecordingSignalListener(rse *codecs.RecordingSignalEvent) {
	testCounters.recordingSignalCount++
}

func RecordingEventStartedListener(rs *codecs.RecordingStarted) {
	testCounters.recordingEventStartedCount++
}

func RecordingEventProgressListener(rp *codecs.RecordingProgress) {
	testCounters.recordingEventProgressCount++
}

func RecordingEventStoppedListener(rs *codecs.RecordingStopped) {
	testCounters.recordingEventStoppedCount++
}

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	context := aeron.NewContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	options := DefaultOptions()

	// Cleaning up after test runs can take a little time so we
	// randomize the streams in use to make that less likely
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	testCases[0].sampleStream += int32(r.Intn(1000))
	testCases[0].replayStream += int32(r.Intn(1000))
	if testCases[0].sampleStream == testCases[0].replayStream {
		testCases[0].replayStream++
	}

	if *TestConfig.Debug {
		log.Printf("Setting verbose logging")
		log.Printf("Using %s/%d and %s/%d", testCases[0].sampleChannel, testCases[0].sampleStream, testCases[0].replayChannel, testCases[0].replayStream)
		options.ArchiveLoglevel = logging.DEBUG
		DEBUG = true
	}

	archive, err = NewArchive(options, context)
	if err != nil || archive == nil {
		log.Printf("archive-media-driver connection failed, skipping all archive_tests:%s", err.Error())
		return
	}
	haveArchive = true

	result := m.Run()
	if result != 0 {
		archive.Close()
		os.Exit(result)
	}

	// FIXME disable auth testing
	archive.Close()
	os.Exit(result)

	// Test auth
	options.AuthEnabled = true
	options.AuthCredentials = []uint8(*TestConfig.AuthCredentials)
	options.AuthChallenge = []uint8(*TestConfig.AuthChallenge)
	options.AuthResponse = []uint8(*TestConfig.AuthResponse)

	testCases[0].sampleStream += int32(r.Intn(1000))
	testCases[0].replayStream += int32(r.Intn(1000))
	if testCases[0].sampleStream == testCases[0].replayStream {
		testCases[0].replayStream++
	}

	archive, err = NewArchive(options, context)
	if err != nil || archive == nil {
		log.Printf("secure-archive-media-driver connection failed, skipping allsecure  archive_tests:%s", err.Error())
		haveArchive = false
		return
	}

	haveArchive = true
	result = m.Run()

	archive.Close()
	os.Exit(result)
}

// This should always pass
func TestConnection(t *testing.T) {
	if !haveArchive {
		return
	}

	// PollForErrorEvents should be safe
	for i := 0; i < 10; i++ {
		count, err := archive.PollForErrorResponse()
		if err != nil {
			t.Logf("PollforErrorRespose() recieved %d responses, err is %s", count, err)
			t.FailNow()
		}
		idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
		idler.Idle(0)
	}

}

// Test KeepAlive
func TestKeepAlive(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	if err := archive.KeepAlive(); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

// Helper to check values of counters
func CounterValuesMatch(c TestCounters, signals int, started int, progress int, stopped int, t *testing.T) bool {
	if testCounters.recordingSignalCount != signals {
		t.Logf("testCounters.recordingSignalCount[%d] != signals[%d]", testCounters.recordingSignalCount, signals)
		return false
	}
	if testCounters.recordingEventStartedCount != started {
		t.Logf("testCounters.recordingEventStartedCount[%d] != started[%d]", testCounters.recordingEventStartedCount, started)
		return false
	}
	if testCounters.recordingEventProgressCount != progress {
		t.Logf("testCounters.recordingEventProgressCount[%d] != progress[%d]", testCounters.recordingEventProgressCount, progress)
		return false
	}
	if testCounters.recordingEventStoppedCount != stopped {
		t.Logf("testCounters.recordingEventStoppedCount[%d] != stopped[%d]", testCounters.recordingEventStoppedCount, stopped)
		return false
	}
	return true
}

// Test that Archive RPCs will fail correctly
func TestRPCFailure(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	// Ask to stop a bogus recording
	res, err := archive.StopRecordingByIdentity(0xdeadbeef)
	if err == nil || res {
		t.Logf("RPC succeeded and should have failed")
		t.FailNow()
	}

}

// Test the recording event signals appear
func TestAsyncEvents(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	archive.Listeners.RecordingSignalListener = RecordingSignalListener
	archive.Listeners.RecordingEventStartedListener = RecordingEventStartedListener
	archive.Listeners.RecordingEventProgressListener = RecordingEventProgressListener
	archive.Listeners.RecordingEventStoppedListener = RecordingEventStoppedListener

	testCounters = TestCounters{0, 0, 0, 0}
	if !CounterValuesMatch(testCounters, 0, 0, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	archive.EnableRecordingEvents()
	archive.RecordingEventsPoll()

	if !CounterValuesMatch(testCounters, 0, 0, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	publication, err := archive.AddRecordedPublication(testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Delay a little to get the publication is established
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
	idler.Idle(0)

	archive.RecordingEventsPoll()
	if !CounterValuesMatch(testCounters, 1, 1, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	if err := archive.StopRecordingByPublication(*publication); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !CounterValuesMatch(testCounters, 2, 1, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	archive.RecordingEventsPoll()
	if !CounterValuesMatch(testCounters, 2, 1, 0, 1, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	// Cleanup
	archive.DisableRecordingEvents()
	archive.Listeners.RecordingSignalListener = nil
	archive.Listeners.RecordingEventStartedListener = nil
	archive.Listeners.RecordingEventProgressListener = nil
	archive.Listeners.RecordingEventStoppedListener = nil
	testCounters = TestCounters{0, 0, 0, 0}
	if !CounterValuesMatch(testCounters, 0, 0, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	publication.Close()
}

// Test PollForErrorEvents
func TestPollForErrorEvents(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	archive.Listeners.RecordingSignalListener = RecordingSignalListener

	testCounters = TestCounters{0, 0, 0, 0}
	if !CounterValuesMatch(testCounters, 0, 0, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	publication, err := archive.AddRecordedPublication(testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// PollForErrorEvents should simply return successfully with the recording signal event having arrived
	_, err = archive.PollForErrorResponse()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !CounterValuesMatch(testCounters, 1, 0, 0, 0, t) {
		t.Log("Async event counters mismatch")
		t.FailNow()
	}

	// Delay a little to get the publication established
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 500}
	idler.Idle(0)

	if err := archive.StopRecordingByPublication(*publication); err != nil {
		t.Log(err)
		t.FailNow()
	}
	publication.Close()

	// Now we'll reach inside the archive a little to leave an outstanding request in the queue
	// We know a StopRecording of a non-existent recording should fail but this call will succeed
	// as it's only the request half
	err = archive.Proxy.StopRecordingSubscriptionRequest(12345, 54321)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// So now PollForErrorResponse should get the reply to that
	// and fail because overlapping I/O is very bad. Note that if
	// normal archive calls are made then we have locking to
	// prevent this
	idler.Idle(0)
	count, err := archive.Control.PollForErrorResponse()
	if err == nil {
		t.Logf("PollForErrorResponse succeeded and should have failed: count is %d", count)
		t.FailNow()
	}
	if count != 1 {
		t.Logf("PollForErrorResponse failed correctly but count is %d and should have been 1", count)
		t.FailNow()
	}

	// Then PollForErrorResponse should see no further messages
	idler.Idle(0)
	count, err = archive.Control.PollForErrorResponse()
	if err != nil {
		t.Logf("PollForErrorResponse failed")
		t.FailNow()
	}
	if count != 0 {
		t.Logf("PollForErrorResponse succeeded but count is %d and should have been 0", count)
		t.FailNow()
	}
}

// Test adding a recording and then removing it - by Publication (session specific)
func TestStartStopRecordingByPublication(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	publication, err := archive.AddRecordedPublication(testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Delay a little to get the publication is established
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
	idler.Idle(0)

	if err := archive.StopRecordingByPublication(*publication); err != nil {
		t.Log(err)
		t.FailNow()
	}
	publication.Close()

}

// Test adding a recording and then removing it - by Subscription
func TestStartStopRecordingBySubscription(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	// Start snd stop by subscription
	subscriptionID, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = archive.StopRecordingBySubscriptionId(subscriptionID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

// Test adding a recording and then removing it - by Channel and Stream
func TestStartStopRecordingByChannelAndStream(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	// Start snd stop by channel&stream
	_, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = archive.StopRecording(testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Start snd stop by identity is done in other tests
}

// Test adding a recording and then removing it, checking the listing counts between times
func TestListRecordings(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	recordings, err := archive.ListRecordingsForUri(0, 100, testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	initial := len(recordings)
	t.Logf("Initial count is %d", initial)

	// Add a recording
	subscriptionID, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("SubscriptionID is %d", subscriptionID)

	// Add a publication on that
	publication := <-archive.AddPublication(testCases[0].sampleChannel, testCases[0].sampleStream)
	if testing.Verbose() && DEBUG {
		t.Logf("Publication is %#v", publication)
	}

	recordings, err = archive.ListRecordingsForUri(0, 100, testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(recordings) == 0 {
		t.Log("No Recordings!")
		t.FailNow()
	}

	//  Grab the recordingID
	recordingID := recordings[len(recordings)-1].RecordingId
	t.Logf("Working count is %d, recordingID is %d", len(recordings), recordingID)

	// Cleanup
	res, err := archive.StopRecordingByIdentity(recordingID)
	if err != nil {
		t.Logf("StopRecordingByIdentity(%d) failed: %s", recordingID, err.Error())
	} else if !res {
		t.Logf("StopRecordingByIdentity(%d) failed", recordingID)
	}
	if err := archive.PurgeRecording(recordingID); err != nil {
		t.Logf("PurgeRecording(%d) failed: %s", recordingID, err.Error())
	}
	publication.Close()

	recordings, err = archive.ListRecordingsForUri(0, 100, testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	final := len(recordings)
	t.Logf("Final count is %d", final)

	if initial != final {
		t.Logf("Number of recordings changed from %d to %d", initial, final)
		t.Fail()
	}
}

// Test starting a replay
func TestStartStopReplay(t *testing.T) {
	if !haveArchive {
		return
	}

	// Add a recording to make sure there is one
	subscriptionID, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("SubscriptionID is %d", subscriptionID)

	// Add a publication on that
	publication := <-archive.AddPublication(testCases[0].sampleChannel, testCases[0].sampleStream)
	t.Logf("Publication found %v", publication)

	recordings, err := archive.ListRecordingsForUri(0, 100, "aeron", testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(recordings) == 0 {
		t.Log("No recordings!")
		t.FailNow()
	}

	// That should give us a recordingID
	recordingID := recordings[len(recordings)-1].RecordingId

	replayID, err := archive.StartReplay(recordingID, 0, RecordingLengthNull, testCases[0].replayChannel, testCases[0].replayStream)
	if err != nil {
		t.Logf("StartReplay failed: %d, %s", replayID, err.Error())
		t.FailNow()
	}
	if err := archive.StopReplay(replayID); err != nil {
		t.Logf("StopReplay(%d) failed: %s", replayID, err.Error())
	}

	// So ListRecordingsForUri should find something
	recordings, err = archive.ListRecordingsForUri(0, 100, testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Working count is %d, recordingID is %d", len(recordings), recordingID)

	// And ListRecordings should also find something
	recordings, err = archive.ListRecordings(0, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(recordings) == 0 {
		t.Log("No recordings!")
		t.FailNow()
	}
	recordingID = recordings[len(recordings)-1].RecordingId
	t.Logf("Working count is %d, recordingID is %d", len(recordings), recordingID)

	// ListRecording should find one by the above Id
	recording, err := archive.ListRecording(recordingID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if recordingID != recording.RecordingId {
		t.Log("ListRecording did not return the correct record descriptor")
		t.FailNow()
	}
	t.Logf("ListRecording(%d) returned %#v", recordingID, *recording)

	// ListRecording should not find one with a bad ID
	badID := int64(-127)
	recording, err = archive.ListRecording(badID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if recording != nil {
		t.Log("ListRecording returned a record descriptor and should not have")
		t.FailNow()
	}
	t.Logf("ListRecording(%d) correctly returned nil", badID)

	// While we're here, check ListRecordingSubscription is working
	descriptors, err := archive.ListRecordingSubscriptions(0, 10, false, 0, "aeron")
	if err != nil {
		t.Logf("ListRecordingSubscriptions() returned error:%s", err.Error())
		t.FailNow()
	}
	if descriptors == nil {
		t.Logf("ListRecordingSubscriptions() returned no descriptors")
		t.FailNow()
	}
	t.Logf("ListRecordingSubscriptions() returned %d descriptor(s)", len(descriptors))

	// Cleanup
	res, err := archive.StopRecordingByIdentity(recordingID)
	if err != nil {
		t.Logf("StopRecordingByIdentity(%d) failed: %s", recordingID, err.Error())
		t.FailNow()
	}
	if !res {
		t.Logf("StopRecordingByIdentity(%d) failed", recordingID)
	}
}

// Test starting a replay
// FIXME:BiggerPicture Disabled as aeron calls log.Fatalf()
func DisabledTestAddRecordedPublicationFailure(t *testing.T) {
	if !haveArchive {
		return
	}

	pub, err := archive.AddRecordedPublication("bogus", 99)
	if err != nil || pub != nil {
		t.Logf("Add recorded publication succeeded and should have failed. error:%s, pub%#v", err, pub)
		t.FailNow()
	}
}

// Test concurrency
func DisabledTestConcurrentConnections(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go ConcurrentSimple(&wg, i, t)
	}
	wg.Wait()

}

func ConcurrentSimple(wg *sync.WaitGroup, n int, t *testing.T) {

	var err error
	context := aeron.NewContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	options := DefaultOptions()
	// options.ArchiveLoglevel = logging.DEBUG

	defer wg.Done()
	t.Logf("Worker %d starting", n)

	// Randomize our stream
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	testCases[0].sampleStream += int32(r.Intn(10000))

	archive, err = NewArchive(options, context)
	if err != nil || archive == nil {
		t.Logf("archive-media-driver connection failed, skipping all archive_tests:%s", err.Error())
		return
	}

	// Thump out some Start and Stop RecordingRequests. If we do too many we'tll timeout, or need to backoff
	if false {
		for i := 0; i < 5; i++ {
			_, err := archive.StartRecording(testCases[0].sampleChannel, int32(20000+n*100+i), true, true)
			if err != nil {
				t.Logf("StartRecording failed for worker %d, attempt %d: %s", n, i, err.Error())
				return
			}
		}
	}

	archive.Close()
	t.Logf("Worker %d exiting", n)
}
