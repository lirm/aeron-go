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
	"fmt"
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/archive/codecs"
)

// RecordingEventsAdapter is used to poll for recording events on a subscription.
type RecordingEventsAdapter struct {
	Subscription *aeron.Subscription
	Enabled      bool
	archive      *Archive // link to parent
}

// FragmentHandlerWithListeners provides a FragmentHandler with ArchiveListeners
type FragmentHandlerWithListeners func(listeners *ArchiveListeners, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header)

// PollWithContext the aeron subscription handler.
// If you pass it a nil handler it will use the builtin and call the Listeners
// If you ask for 0 fragments it will only return one fragment (if available)
func (rea *RecordingEventsAdapter) PollWithContext(handler FragmentHandlerWithListeners, fragmentLimit int) int {

	// Update our globals in case they've changed so we use the current state in our callback
	rangeChecking = rea.archive.Options.RangeChecking

	if handler == nil {
		handler = reFragmentHandler
	}
	if fragmentLimit == 0 {
		fragmentLimit = 1
	}
	return rea.Subscription.PollWithContext(
		func(buf *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
			handler(rea.archive.Listeners, buf, offset, length, header)
		}, fragmentLimit)
}

func reFragmentHandler(listeners *ArchiveListeners, buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	var hdr codecs.SbeGoMessageHeader

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()

	if err := hdr.Decode(marshaller, buf); err != nil {
		// Not much to be done here as we can't correlate
		err2 := fmt.Errorf("reFragmentHandler() failed to decode control message header: %w", err)
		// Call the global error handler, ugly but it's all we've got
		if listeners.ErrorListener != nil {
			listeners.ErrorListener(err2)
		}
	}

	switch hdr.TemplateId {
	case codecIds.recordingStarted:
		var recordingStarted = new(codecs.RecordingStarted)
		logger.Debugf("Received RecordingStarted: length %d", buf.Len())
		if err := recordingStarted.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			err2 := fmt.Errorf("Decode() of RecordingStarted failed: %w", err)
			if listeners.ErrorListener != nil {
				listeners.ErrorListener(err2)
			}
		} else {
			// logger.Debugf("RecordingStarted: %#v\n", recordingStarted)
			// Call the Listener
			if listeners.RecordingEventStartedListener != nil {
				listeners.RecordingEventStartedListener(recordingStarted)
			}
		}

	case codecIds.recordingProgress:
		var recordingProgress = new(codecs.RecordingProgress)
		logger.Debugf("Received RecordingProgress: length %d", buf.Len())
		if err := recordingProgress.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			err2 := fmt.Errorf("Decode() of RecordingProgress failed: %w", err)
			if listeners.ErrorListener != nil {
				listeners.ErrorListener(err2)
			}
		} else {
			logger.Debugf("RecordingProgress: %#v\n", recordingProgress)
			// Call the Listener
			if listeners.RecordingEventProgressListener != nil {
				listeners.RecordingEventProgressListener(recordingProgress)
			}
		}

	case codecIds.recordingStopped:
		var recordingStopped = new(codecs.RecordingStopped)
		logger.Debugf("Received RecordingStopped: length %d", buf.Len())
		if err := recordingStopped.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, rangeChecking); err != nil {
			err2 := fmt.Errorf("Decode() of RecordingStopped failed: %w", err)
			if listeners.ErrorListener != nil {
				listeners.ErrorListener(err2)
			}
		} else {
			logger.Debugf("RecordingStopped: %#v\n", recordingStopped)
			// Call the Listener
			if listeners.RecordingEventStoppedListener != nil {
				listeners.RecordingEventStoppedListener(recordingStopped)
			}
		}

	default:
		logger.Errorf("Insert decoder for type: %d\n", hdr.TemplateId)
	}
}
