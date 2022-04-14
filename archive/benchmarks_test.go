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
	"bytes"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/archive/codecs"
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

// Benchmark the Descriptor Fragment Handler using a recording descriptor as the example.
// Used to evaluate whether optimizations are worthwhile here
func BenchmarkDescriptorFragmentHandler(b *testing.B) {
	rd := codecs.RecordingDescriptor{
		ControlSessionId:  1,
		CorrelationId:     2,
		RecordingId:       3,
		StartTimestamp:    4,
		StopTimestamp:     5,
		StartPosition:     6,
		StopPosition:      7,
		InitialTermId:     8,
		SegmentFileLength: 9,
		TermBufferLength:  10,
		MtuLength:         11,
		SessionId:         12,
		StreamId:          13,
		StrippedChannel:   []uint8("stripped"),
		OriginalChannel:   []uint8("original"),
		SourceIdentity:    []uint8("source")}

	marshaller := codecs.NewSbeGoMarshaller()
	buffer := new(bytes.Buffer)

	header := codecs.MessageHeader{BlockLength: rd.SbeBlockLength(), TemplateId: rd.SbeTemplateId(), SchemaId: rd.SbeSchemaId(), Version: rd.SbeSchemaVersion()}

	if err := header.Encode(marshaller, buffer); err != nil {
		b.Logf("header encode failed")
		b.FailNow()
	}
	if err := rd.Encode(marshaller, buffer, rangeChecking); err != nil {
		b.Logf("header encode failed")
		b.FailNow()
	}

	bytes := buffer.Bytes()
	length := int32(len(bytes))
	atomicbuffer := atomic.MakeBuffer(bytes, len(bytes))

	// Mock the correlationId to Control map
	correlations.Store(rd.CorrelationId, new(Control))

	for n := 0; n < b.N; n++ {
		DescriptorFragmentHandler(archive.SessionID, atomicbuffer, 0, length, nil)
	}

}
