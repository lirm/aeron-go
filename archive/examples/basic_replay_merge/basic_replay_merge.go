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

// An example replayed subscriber
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/lirm/aeron-go/archive"
	"github.com/lirm/aeron-go/archive/codecs"
	"github.com/lirm/aeron-go/archive/examples"
	"github.com/lirm/aeron-go/archive/replaymerge"
)

var logID = "basic_replay_merge"
var logger = logging.MustGetLogger(logID)

func main() {
	flag.Parse()

	sampleChannel := *examples.Config.MdcChannel
	sampleStream := int32(*examples.Config.SampleStream)
	responseStream := sampleStream + 2

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.Config.DriverTimeout)
	context := aeron.NewContext().
		AeronDir(*examples.Config.AeronPrefix).
		MediaDriverTimeout(timeout).
		AvailableImageHandler(func(image aeron.Image) {
			logger.Infof("Image available: %+v", image)
		}).
		UnavailableImageHandler(func(image aeron.Image) {
			logger.Infof("Image unavailable: %+v", image)
		}).
		ErrorHandler(func(err error) {
			logger.Warning(err)
		})

	options := archive.DefaultOptions()
	options.RequestChannel = *examples.Config.RequestChannel
	options.RequestStream = int32(*examples.Config.RequestStream)
	options.ResponseChannel = *examples.Config.ResponseChannel
	options.ResponseStream = responseStream

	t := false
	examples.Config.Verbose = &t
	if *examples.Config.Verbose {
		fmt.Printf("Setting loglevel: archive.DEBUG/aeron.INFO\n")
		options.ArchiveLoglevel = logging.DEBUG
		options.AeronLoglevel = logging.DEBUG
		logging.SetLevel(logging.DEBUG, logID)
	} else {
		logging.SetLevel(logging.NOTICE, logID)
	}

	arch, err := archive.NewArchive(options, context)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer arch.Close()

	// Enable recording events although the defaults will only log in debug mode
	arch.EnableRecordingEvents()

	recording, err := FindLatestRecording(arch, sampleChannel, sampleStream)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	subChannelUri, _ := aeron.ParseChannelUri("aeron:udp")
	subChannelUri.SetControlMode(aeron.MdcControlModeManual)
	subChannelUri.SetSessionID(recording.SessionId)
	subChannel := subChannelUri.String()

	logger.Infof("Subscribing to channel:%s, stream:%d", subChannel, sampleStream)
	subscription, err := arch.AddSubscription(subChannel, sampleStream)
	if err != nil {
		logger.Fatal(err)
	}
	defer subscription.Close()

	counter := 0
	printHandler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)
		logger.Infof("Message: %s", bytes)
		counter++
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}

	liveDestination := sampleChannel

	replayDestinationUri, _ := aeron.ParseChannelUri("aeron:udp?endpoint=localhost:0")
	replayDestination := replayDestinationUri.String()

	replayChannelUri, _ := aeron.ParseChannelUri(replayDestination)
	replayChannelUri.SetSessionID(recording.SessionId)
	replayChannel := replayChannelUri.String()

	logger.Infof(
		"ReplayMerge(subChannel:%s,%d,replayChannel:%s,replayDestination:%s,liveDestination:%s)",
		subscription.Channel(),
		subscription.StreamID(),
		replayChannel,
		replayDestination,
		liveDestination,
	)
	merge, err := replaymerge.NewReplayMerge(
		subscription,
		arch,
		replayChannel,
		replayDestination,
		liveDestination,
		recording.RecordingId,
		0,
		5000,
	)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	defer merge.Close()

	var fragmentsRead int
	for !merge.IsMerged() {
		fragmentsRead, err = merge.Poll(printHandler, 10)
		if err != nil {
			logger.Fatalf(err.Error())
		}
		arch.RecordingEventsPoll()
		idleStrategy.Idle(fragmentsRead)
	}

	merge.Close()

	for {
		fragmentsRead := subscription.Poll(printHandler, 10)
		arch.RecordingEventsPoll()

		idleStrategy.Idle(fragmentsRead)
	}
}

// FindLatestRecording to lookup the last recording
func FindLatestRecording(arch *archive.Archive, channel string, stream int32) (*codecs.RecordingDescriptor, error) {
	descriptors, err := arch.ListRecordingsForUri(0, 100, "", stream)

	if err != nil {
		return nil, err
	}

	if len(descriptors) == 0 {
		return nil, fmt.Errorf("no recordings found")
	}

	descriptor := descriptors[len(descriptors)-1]
	if descriptor.StopPosition == 0 {
		return nil, fmt.Errorf("recording length zero")
	}

	// Return the last recordingID
	return descriptor, nil
}
