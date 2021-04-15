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
	"testing"
)

// check to see if this is worth avoiding in the fragment handlers
// which are called quite frequently
func BenchmarkCodecObjectCreation(b *testing.B) {
	var recordingStarted *codecs.RecordingStarted
	for n := 0; n < b.N; n++ {
		recordingStarted = new(codecs.RecordingStarted)
	}
	// declared but not used
	if recordingStarted == nil {
		b.Logf("all done")
	}
}
