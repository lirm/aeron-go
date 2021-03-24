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
	"bytes"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
)

type RecordingEventsAdapter struct {
	Context      *ArchiveContext
	Subscription *aeron.Subscription
	Enabled      bool
}

// Create a new recording event adapter
func NewRecordingEventsAdapter(context *ArchiveContext) *RecordingEventsAdapter {
	rea := new(RecordingEventsAdapter)
	rea.Context = context

	return rea
}

// The response poller wraps the aeron subscription handler.
// If you pass it a nil handler it will use the builtin and call the Listeners
// If you ask for 0 fragments it will only return one fragment (if available)
func (rea *RecordingEventsAdapter) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = rea.Context.Options.RangeChecking

	if handler == nil {
		handler = ReFragmentHandler
	}
	if fragmentLimit == 0 {
		fragmentLimit = 1
	}
	return rea.Subscription.Poll(handler, fragmentLimit)
}

func ReFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	var hdr codecs.SbeGoMessageHeader
	var recordingStarted = new(codecs.RecordingStarted)
	var recordingProgress = new(codecs.RecordingProgress)
	var recordingStopped = new(codecs.RecordingStopped)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()

	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: Should we use an ErrorHandler?
		logger.Error("Failed to decode message header", err)
		return
	}

	switch hdr.TemplateId {
	case recordingStarted.SbeTemplateId():
		logger.Debugf("Received RecordingStarted: length %d", buf.Len())
		if err := recordingStarted.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			// FIXME: Listeners.ErrorListener()
			logger.Error("Failed to decode RecordingStarted", err)

		} else {
			// Call the Listener
			if Listeners.RecordingEventStartedListener != nil {
				Listeners.RecordingEventStartedListener(recordingStarted)
			}
		}

	case recordingProgress.SbeTemplateId():
		logger.Debugf("Received RecordingProgress: length %d", buf.Len())
		if err := recordingProgress.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			logger.Error("Failed to decode RecordingProgress", err)
		} else {
			logger.Debugf("RecordingProgress: %#v\n", recordingProgress)
			// Call the Listener
			if Listeners.RecordingEventProgressListener != nil {
				Listeners.RecordingEventProgressListener(recordingProgress)
			}
		}

	case recordingStopped.SbeTemplateId():
		logger.Debugf("Received RecordingStopped: length %d", buf.Len())
		if err := recordingStopped.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			logger.Error("Failed to decode RecordingStopped", err)
		} else {
			logger.Debugf("RecordingStopped: %#v\n", recordingStopped)
			// Call the Listener
			if Listeners.RecordingEventStoppedListener != nil {
				Listeners.RecordingEventStoppedListener(recordingStopped)
			}
		}

	default:
		logger.Errorf("Insert decoder for type: %d\n", hdr.TemplateId)
	}
}
