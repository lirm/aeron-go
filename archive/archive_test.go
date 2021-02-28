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
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/archive/codecs"
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
var context *ArchiveContext
var archive *Archive
var haveArchive bool = false
var DEBUG = true

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
	context = NewArchiveContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	archive, err = ArchiveConnect(context)
	if err != nil || archive == nil {
		log.Printf("archive-media-driver connection failed, skipping all archive_tests:%s", err.Error())
		return
	} else {
		haveArchive = true
	}

	result := m.Run()
	archive.Close()
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 500}
	idler.Idle(0)
	os.Exit(result)
}

// This should always pass
func TestConnection(t *testing.T) {
	if !haveArchive {
		return
	}
}

// Test adding a recording
func TestStartStopRecording(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}
	recordingId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, codecs.SourceLocation.LOCAL, true)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("id:%#v", recordingId)
	StopRecording(recordingId)

}

// Test adding a recording
func TestListRecordingsForUri(t *testing.T) {
	if !haveArchive {
		return
	}

	if testing.Verbose() && DEBUG {
		logging.SetLevel(logging.DEBUG, "archive")
	}

	count, err := archive.ListRecordingsForUri(0, 100, "aeron", testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if count == 0 {
		// Add a recording to make sure there is one
		recordingId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, codecs.SourceLocation.LOCAL, true)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		t.Logf("Added new recording: id:%d", recordingId)
		defer StopRecording(recordingId)
		// FIXME: Again a delay is a bogus thing to do here.
		idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 500}
		idler.Idle(0)
	}

	count, err = archive.ListRecordingsForUri(0, 100, "aeron", testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("count:%d", count)
}

// Test starting a replay
func TestStartStopReplay(t *testing.T) {
	if !haveArchive {
		return
	}

	// Add a recording to make sure there is one
	recordingId, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, codecs.SourceLocation.LOCAL, true)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("recordingId:%#v", recordingId)
	defer StopRecording(recordingId)

	// FIXME: Delay a little to get that established
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 500}
	idler.Idle(0)

	count, err := archive.ListRecordingsForUri(0, 100, "aeron", testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if count == 0 {
		t.Log("FIXME:No recordings to start")
		t.FailNow()

	}

	recordingId = archive.Control.Results.RecordingDescriptors[count-1].RecordingId
	t.Logf("id:%#v", recordingId)
	replayId, err := archive.StartReplay(recordingId, 0, -1, testCases[0].replayChannel, testCases[0].replayStream)
	if err != nil {
		t.Logf("StartReplay failed: %d, %s", replayId, err.Error())
		t.FailNow()
	}
	defer StopReplay(replayId)

	return

}

// Defer functions to keep the tests tidy

func StopRecording(recordingId int64) {
	log.Printf("StopRecordingBySubscriptionId(%d)", recordingId)
	res, err := archive.StopRecordingBySubscriptionId(recordingId)
	if err != nil {
		log.Printf("StopRecordingBySubscriptionId(%d) failed:%d %s", recordingId, res, err.Error())
	}
}

func StopReplay(replayId int64) {
	log.Printf("StopReplay(%d)", replayId)
	res, err := archive.StopReplay(replayId)
	if err != nil {
		log.Printf("StopReplay(%d) failed:%d %s", replayId, res, err.Error())
	}
}
