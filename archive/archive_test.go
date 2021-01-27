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
	"testing"
	"time"
)

// Rather than mock or spawn an archive-media-driver we're just seeing
// if we can connect to one and if we can we'll run some tests. If the
// init fails to connect then we'll skip the tests
var context *ArchiveContext
var archive *Archive

func init() {
	context = NewArchiveContext()
	context.AeronDir(*TestConfig.AeronPrefix)
	archive, _ = ArchiveConnect(context)
}

// This should always pass
func TestConnection(t *testing.T) {
	if archive == nil {
		t.Log("Skipping test connection")
	}
}

// Test adding a recording
func TestRecordedPublication(t *testing.T) {
	if archive == nil {
		t.Log("Skipping test connection")
	}

	time.Sleep(time.Second) // FIXME: delay

	pub, err := archive.AddRecordedPublication(*TestConfig.RecordingChannel, int32(*TestConfig.RecordingStream))
	if err != nil {
		t.Fail()
	}
	t.Logf("pub:%#v", pub)

	time.Sleep(time.Second) // FIXME: delay
}
