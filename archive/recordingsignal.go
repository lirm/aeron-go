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

// Package archive provides API access to Aeron's archive-media-driver
package archive

import (
	"bytes"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	_ "github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
)

// The Recording event listener types
type ControlResponseListener func(*codecs.ControlResponse)
type RecordingStartedListener func(*codecs.RecordingStarted)
type RecordingProgressListener func(*codecs.RecordingProgress)
type RecordingStoppedListener func(*codecs.RecordingStopped)
type RecordingSignalListener func(*codecs.RecordingSignalEvent)

// Recording signals implement the asynchronous receiving and dispatching of of recording events
type RecordingSignalAdapter struct {
	SessionID               int64
	Consumer                int64 // FIXME?
	Subscription            aeron.Subscription
	controlResponseListener ControlResponseListener
	StartListener           RecordingStartedListener
	ProgressListener        RecordingProgressListener
	StoppedListener         RecordingStoppedListener
	SignalListener          RecordingSignalListener
	fragmentLimit           int
	isDone                  bool
}

// The control response poller wraps the aeron subscription handler
func (rsa *RecordingSignalAdapter) Poll(handler term.FragmentHandler, fragmentLimit int) int {
	rsa.isDone = false
	return rsa.Subscription.Poll(handler, fragmentLimit)
}

func RecordingEventFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	var hdr codecs.SbeGoMessageHeader
	var controlResponse = new(codecs.ControlResponse)
	var recordingSignalEvent = new(codecs.RecordingSignalEvent)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()

	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: Should we use an ErrorHandler?
		logger.Error("Failed to decode control message header", err) // Not much to be done here as we can't correlate
	}

	switch hdr.TemplateId {
	case controlResponse.SbeTemplateId():
		logger.Debugf("Received controlResponse: length %d", buf.Len())
		if err := controlResponse.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			logger.Error("Failed to decode control response", err) // Not much to be done here as we can't correlate
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := sessionsMap[controlResponse.ControlSessionId]
		if !ok {
			logger.Info("Uncorrelated control response SessionID=", controlResponse.ControlSessionId) // Not much to be done here as we can't correlate
			return
		}

		// Call the Listener
		control.rsa.controlResponseListener(controlResponse)
		control.rsa.isDone = true

	case recordingSignalEvent.SbeTemplateId():
		logger.Debugf("Received RecordingSignalEvent: length %d", buf.Len())
		if err := recordingSignalEvent.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			logger.Error("Failed to decode control response", err) // Not much to be done here as we can't correlate
		}
		logger.Debugf("ControlResponse: %#v\n", controlResponse)

		// Look it up
		control, ok := sessionsMap[recordingSignalEvent.ControlSessionId]
		// Not much to be done here as we can't correlate
		if !ok {
			logger.Infof("Uncorrelated control response SessionID=%d, CorrelationID=%d", controlResponse.ControlSessionId, controlResponse.CorrelationId)
			return
		}

		// Call the Listener
		control.rsa.SignalListener(recordingSignalEvent)
		control.rsa.isDone = true

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}

	return
}

func OnFragmentHandler(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	var hdr codecs.SbeGoMessageHeader
	var recordingStarted = new(codecs.RecordingStarted)
	var recordingProgress = new(codecs.RecordingProgress)
	var recordingStopped = new(codecs.RecordingStopped)

	buf := new(bytes.Buffer)
	buffer.WriteBytes(buf, offset, length)

	marshaller := codecs.NewSbeGoMarshaller()

	if err := hdr.Decode(marshaller, buf); err != nil {
		// FIXME: Should we use an ErrorHandler?
		logger.Error("Failed to decode control message header", err) // Not much to be done here as we can't correlate
	}

	switch hdr.TemplateId {
	case recordingStarted.SbeTemplateId():
		logger.Debugf("Received RecordingStarted: length %d", buf.Len())
		if err := recordingStarted.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode RecordingStarted", err)
		}
		logger.Debugf("RecordingStarted: %#v\n", recordingStarted)

		// Look it up
		control, ok := sessionsMap[int64(recordingStarted.SessionId)] // FIXME: cast means my assumption is wrong
		// Not much to be done here as we can't correlate
		if !ok {
			logger.Infof("Uncorrelated control response SessionID=%d\n", recordingStarted.SessionId)
			return
		}

		// Call the Listener
		control.rsa.StartListener(recordingStarted)

	case recordingProgress.SbeTemplateId():
		logger.Debugf("Received RecordingProgress: length %d", buf.Len())
		if err := recordingProgress.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode RecordingProgress", err)
		}
		logger.Debugf("RecordingProgress: %#v\n", recordingProgress)

		// Look it up
		control, ok := recordingsMap[recordingProgress.RecordingId]
		// Not much to be done here as we can't correlate
		if !ok {
			logger.Infof("Uncorrelated recording event RecordingId=%d\n", recordingProgress.RecordingId)
			return
		}

		// Call the Listener
		control.rsa.ProgressListener(recordingProgress)

	case recordingStopped.SbeTemplateId():
		logger.Debugf("Received RecordingStopped: length %d", buf.Len())
		if err := recordingStopped.Decode(marshaller, buf, hdr.Version, hdr.BlockLength, ArchiveDefaults.RangeChecking); err != nil {
			// Not much to be done here as we can't correlate
			logger.Error("Failed to decode RecordingStopped", err)
		}
		logger.Debugf("RecordingStopped: %#v\n", recordingStopped)

		// Look it up
		control, ok := recordingsMap[recordingStopped.RecordingId]
		// Not much to be done here as we can't correlate
		if !ok {
			logger.Infof("Uncorrelated recording event RecordingId=%d\n", recordingStopped.RecordingId)
			return
		}

		// Call the Listener
		control.rsa.StoppedListener(recordingStopped)

	default:
		fmt.Printf("Insert decoder for type: %d\n", hdr.TemplateId)
	}
}
