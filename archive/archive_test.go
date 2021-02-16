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
	"github.com/lirm/aeron-go/archive/codecs"
	logging "github.com/op/go-logging"
	"testing"
)

// Rather than mock or spawn an archive-media-driver we're just seeing
// if we can connect to one and if we can we'll run some tests. If the
// init fails to connect then we'll skip the tests
// FIXME: this plan fails as aeron-go calls log.Fatalf() !!!
var context *ArchiveContext
var archive *Archive
var connectionError error

type TestCases struct {
	sampleStream  int32
	sampleChannel string
	replayStream  int32
	replayChannel string
}

var testCases = []TestCases{
	{int32(*TestConfig.SampleStream), *TestConfig.SampleChannel, int32(*TestConfig.ReplayStream), *TestConfig.ReplayChannel},
}

func init() {
	context = NewArchiveContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	archive, connectionError = ArchiveConnect(context)
}

// This should always pass
func TestConnection(t *testing.T) {
	if connectionError != nil || archive == nil {
		t.Log("Skipping as not connected to archive-media-driver")
		return
	}
}

// Test adding a recording
func TestStartRecording(t *testing.T) {
	if connectionError != nil || archive == nil {
		t.Log("Skipping as not connected to archive-media-driver")
		return
	}

	pub, err := archive.StartRecording(testCases[0].sampleChannel, testCases[0].sampleStream, codecs.SourceLocation.LOCAL, true)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("pub:%#v", pub)
}

// Test adding a recording
func TestListRecordingsForUri(t *testing.T) {
	if connectionError != nil || archive == nil {
		t.Log("Skipping as not connected to archive-media-driver")
		return
	}

	if testing.Verbose() {
		logging.SetLevel(logging.DEBUG, "archive")
	}
	count, err := archive.ListRecordingsForUri(0, 100, "aeron", testCases[0].sampleStream)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("count:%d", count)
}
