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

// An example recorded publisher
package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/archive"
	"github.com/lirm/aeron-go/archive/codecs"
	"github.com/lirm/aeron-go/archive/examples"
	"github.com/op/go-logging"
	"os"
	"time"
)

var logger = logging.MustGetLogger("basic_recording_publisher")

func main() {
	flag.Parse()

	if !*examples.Config.Verbose {
		logging.SetLevel(logging.INFO, "aeron")    // FIXME
		logging.SetLevel(logging.DEBUG, "archive") // FIXME
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
		logging.SetLevel(logging.INFO, "rb")
	}

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.Config.DriverTimeout)
	context := archive.NewArchiveContext()
	context.AeronDir(*examples.Config.AeronPrefix)
	context.MediaDriverTimeout(timeout)

	options := archive.DefaultOptions()
	options.RequestChannel = *examples.Config.RequestChannel
	options.RequestStream = int32(*examples.Config.RequestStream)
	options.ResponseChannel = *examples.Config.ResponseChannel
	options.ResponseStream = int32(*examples.Config.ResponseStream)

	arch, err := archive.NewArchive(context, options)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer arch.Close()

	channel := *examples.Config.SampleChannel
	stream := int32(*examples.Config.SampleStream)

	if _, err := arch.StartRecording(channel, stream, codecs.SourceLocation.LOCAL, true); err != nil {
		logger.Infof("StartRecording failed: %s\n", err.Error())
		os.Exit(1)
	}
	logger.Infof("StartRecording succeeded\n")

	publication := <-arch.AddPublication(channel, stream)
	logger.Infof("Publication found %v", publication)
	defer publication.Close()

	// FIXME: counters unimplemented
	// The Java RecordedBasicPublisher polls the counters to wait until
	// the recording has actually started.
	// For now we'll just dealy a bit to let that get established
	// There may be a way to use the NewImageHandler here as well
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
	idler.Idle(0)

	for counter := 0; counter < *examples.Config.Messages; counter++ {
		message := fmt.Sprintf("this is a message %d", counter)
		srcBuffer := atomic.MakeBuffer(([]byte)(message))
		ret := publication.Offer(srcBuffer, 0, int32(len(message)), nil)
		switch ret {
		case aeron.NotConnected:
			logger.Infof("%d, Not connected (yet)", counter)

		case aeron.BackPressured:
			logger.Infof("%d: back pressured", counter)
		default:
			if ret < 0 {
				logger.Infof("%d: Unrecognized code: %d", counter, ret)
			} else {
				logger.Infof("%d: success!", counter)
			}
		}

		if !publication.IsConnected() {
			logger.Infof("no subscribers detected")
		}
		time.Sleep(time.Second)
	}
	os.Exit(1)
}
