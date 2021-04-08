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
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	logging "github.com/op/go-logging"
	"log"
	"os"
	"testing"
	"time"
)

// Rather than mock or spawn an archive-media-driver we're just seeing
// if we can connect to one and if we can we'll run some tests. If the
// init fails to connect then we'll skip the tests
// FIXME: this plan fails as aeron-go calls log.Fatalf() if the media driver is not running !!!
var archive *Archive
var haveArchive bool = false
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

func TestMain(m *testing.M) {

	var err error
	context := aeron.NewContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	options := DefaultOptions()
	if *TestConfig.Verbose {
		log.Printf("Setting verbose logging")
		options.ArchiveLoglevel = logging.DEBUG
	}

	archive, err = NewArchive(options, context)
	if err != nil || archive == nil {
		log.Printf("archive-media-driver connection failed, skipping all archive_tests:%s", err.Error())
		return
	} else {
		haveArchive = true
	}

	result := m.Run()
	archive.Close()
	os.Exit(result)
}

// This should always pass
func TestConnection(t *testing.T) {
	if !haveArchive {
		return
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
	subscriptionId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = archive.StopRecordingBySubscriptionId(subscriptionId)
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
	subscriptionId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("SubscriptionId is %d", subscriptionId)

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
	//  Grab the recordingId
	recordingId := recordings[len(recordings)-1].RecordingId
	t.Logf("Working count is %d, recordingId is %d", len(recordings), recordingId)

	// Cleanup
	if res, err := archive.StopRecordingByIdentity(recordingId); err != nil || !res {
		t.Logf("StopRecordingByIdentity(%d) failed: %s", recordingId, err.Error())
	}
	if err := archive.PurgeRecording(recordingId); err != nil {
		t.Logf("PurgeRecording(%d) failed: %s", recordingId, err.Error())
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
	subscriptionId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, true, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("SubscriptionId is %d", subscriptionId)

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

	// That should give us a recordingId
	recordingId := recordings[len(recordings)-1].RecordingId

	replayId, err := archive.StartReplay(recordingId, 0, RecordingLengthNull, testCases[0].replayChannel, testCases[0].replayStream)
	if err != nil {
		t.Logf("StartReplay failed: %d, %s", replayId, err.Error())
		t.FailNow()
	}
	if err := archive.StopReplay(replayId); err != nil {
		t.Logf("StopReplay(%d) failed: %s", replayId, err.Error())
	}

	// So ListRecordingsForUri should find something
	recordings, err = archive.ListRecordingsForUri(0, 100, testCases[0].sampleChannel, testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Working count is %d, recordingId is %d", len(recordings), recordingId)

	// And ListRecordings should also find something
	recordings, err = archive.ListRecordings(0, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	recordingId = recordings[len(recordings)-1].RecordingId
	t.Logf("Working count is %d, recordingId is %d", len(recordings), recordingId)

	// ListRecording should find one by the above Id
	recording, err := archive.ListRecording(recordingId)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if recordingId != recording.RecordingId {
		t.Log("ListRecording did not return the correct record descriptor")
		t.FailNow()
	}
	t.Logf("ListRecording(%d) returned %#v", recordingId, *recording)

	// ListRecording should not find one with a bad Id
	badId := int64(-127)
	recording, err = archive.ListRecording(badId)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if recording != nil {
		t.Log("ListRecording returned a record descriptor and should not have")
		t.FailNow()
	}
	t.Logf("ListRecording(%d) correctly returned nil", badId)

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
	res, err := archive.StopRecordingByIdentity(recordingId)
	if err != nil {
		t.Logf("StopRecordingByIdentity(%d) failed: %s", recordingId, err.Error())
		t.FailNow()
	}
	if !res {
		t.Logf("StopRecordingByIdentity(%d) failed", recordingId)
	}

	return

}

// Test starting a replay
// FIXME: Disabled as aeron calls log.Fatalf()
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
