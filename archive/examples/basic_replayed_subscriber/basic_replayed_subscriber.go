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

// An example replayed subscriber
package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/archive"
	"github.com/lirm/aeron-go/archive/examples"
	logging "github.com/op/go-logging"
	"math"
	"time"
)

var logger = logging.MustGetLogger("basic_recording_subscriber")

func main() {
	flag.Parse()

	if !*examples.Config.Verbose {
		logging.SetLevel(logging.INFO, "archive")
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
	}

	// As per Java example
	sampleChannel := *examples.Config.SampleChannel
	sampleStream := int32(*examples.Config.SampleStream)
	replayStream := sampleStream + 1
	responseStream := sampleStream + 2

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.Config.DriverTimeout)
	context := archive.NewArchiveContext()
	archive.ArchiveDefaults.RequestChannel = *examples.Config.RequestChannel
	archive.ArchiveDefaults.RequestStream = int32(*examples.Config.RequestStream)
	archive.ArchiveDefaults.ResponseChannel = *examples.Config.ResponseChannel
	archive.ArchiveDefaults.ResponseStream = int32(responseStream) // Unique to avoid clashes
	context.AeronDir(*examples.Config.AeronPrefix)
	context.MediaDriverTimeout(timeout)

	arch, err := archive.ArchiveConnect(context)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer arch.Close()

	recordingId, err := FindLatestRecording(arch, sampleChannel, sampleStream)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	var position int64 = 0
	var length int64 = math.MaxInt64

	logger.Infof("Start replay of channel:%s, stream:%d, position:%d, length:%d", sampleChannel, replayStream, position, length)
	replaySessionId, err := arch.StartReplay(recordingId, position, length, sampleChannel, replayStream)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// Make the channel based upon that recording and subscribe to it
	subChannel, err := archive.AddReplaySessionIdToChannel(sampleChannel, archive.ReplaySessionIdToStreamId(replaySessionId))
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
		// logger.Debugf("%8.d: Received fragment offset:%d length:%d payload:%s\n", counter, offset, length, bytes)
		logger.Infof("%s\n", bytes)
		counter++
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond * 1000}
	for {
		fragmentsRead := subscription.Poll(printHandler, 10)
		idleStrategy.Idle(fragmentsRead)
	}
}

// Lookup the last recording
func FindLatestRecording(arch *archive.Archive, channel string, stream int32) (int64, error) {
	descriptors, err := arch.ListRecordingsForUri(0, 100, channel, stream)

	if err != nil {
		return 0, err
	}

	if len(descriptors) == 0 {
		return 0, fmt.Errorf("No recordings found\n")
	}

	// Return the last recordingId
	return descriptors[len(descriptors)-1].RecordingId, nil
}
