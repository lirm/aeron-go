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

var channel = *TestConfig.RequestChannel
var stream = int32(*TestConfig.RequestStream)

// Execute all the encoders as a sanity check
func TestEncoders(t *testing.T) {
	packet, err := ConnectRequestPacket(channel, stream, 1)
	if err != nil {
		t.Log("ConnectRequestPacket() failed")
		t.Fail()
	}
	if len(packet) != 61 { // Look for unexepected change
		t.Logf("ConnectRequestPacket failed length check: %d", len(packet))
		t.Fail()
	}

	packet, err = StartRecordingRequestPacket(channel, stream, 1, codecs.SourceLocation.LOCAL)
	if err != nil {
		t.Log("StartRecordingRequestPacket() failed")
		t.Fail()
	}
	if len(packet) != 69 { // Look for unexepected change
		t.Logf("StartRecordingRequestPacket() failed length check: %d", len(packet))
		t.Fail()
	}
}
