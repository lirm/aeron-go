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
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logging"
	"github.com/corymonroe-coinbase/aeron-go/archive"
	"github.com/corymonroe-coinbase/aeron-go/archive/examples"
	"math"
	"time"
)

var logID = "basic_recording_subscriber"
var logger = logging.MustGetLogger(logID)

func main() {
	flag.Parse()

	// As per Java example
	sampleChannel := *examples.Config.SampleChannel
	sampleStream := int32(*examples.Config.SampleStream)
	replayStream := sampleStream + 1
	responseStream := sampleStream + 2

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.Config.DriverTimeout)
	context := aeron.NewContext()
	context.AeronDir(*examples.Config.AeronPrefix)
	context.MediaDriverTimeout(timeout)

	options := archive.DefaultOptions()
	options.RequestChannel = *examples.Config.RequestChannel
	options.RequestStream = int32(*examples.Config.RequestStream)
	options.ResponseChannel = *examples.Config.ResponseChannel
	options.ResponseStream = int32(responseStream)

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

	recordingID, err := FindLatestRecording(arch, sampleChannel, sampleStream)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	var position int64
	var length int64 = math.MaxInt64

	logger.Infof("Start replay of channel:%s, stream:%d, position:%d, length:%d", sampleChannel, replayStream, position, length)
	replaySessionID, err := arch.StartReplay(recordingID, position, length, sampleChannel, replayStream)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// Make the channel based upon that recording and subscribe to it
	subChannel, err := archive.AddSessionIdToChannel(sampleChannel, archive.ReplaySessionIdToStreamId(replaySessionID))
	if err != nil {
		logger.Fatalf("AddReplaySessionIdToChannel() failed: %s", err.Error())
	}

	logger.Infof("Subscribing to channel:%s, stream:%d", subChannel, replayStream)
	subscription := <-arch.AddSubscription(subChannel, replayStream)
	defer subscription.Close()
	logger.Infof("Subscription found %v", subscription)

	counter := 0
	printHandler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)
		logger.Noticef("%s\n", bytes)
		counter++
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond * 1000}
	for {
		fragmentsRead := subscription.Poll(printHandler, 10)
		arch.RecordingEventsPoll()

		idleStrategy.Idle(fragmentsRead)
	}
}

// FindLatestRecording to lookup the last recording
func FindLatestRecording(arch *archive.Archive, channel string, stream int32) (int64, error) {
	descriptors, err := arch.ListRecordingsForUri(0, 100, channel, stream)

	if err != nil {
		return 0, err
	}

	if len(descriptors) == 0 {
		return 0, fmt.Errorf("no recordings found")
	}

	// Return the last recordingID
	return descriptors[len(descriptors)-1].RecordingId, nil
}
