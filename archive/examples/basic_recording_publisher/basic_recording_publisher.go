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

// An example recorded publisher
package main

import (
	"flag"
	"fmt"
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logging"
	"github.com/corymonroe-coinbase/aeron-go/archive"
	"github.com/corymonroe-coinbase/aeron-go/archive/examples"
	"os"
	"time"
)

var logID = "basic_recording_publisher"
var logger = logging.MustGetLogger(logID)

func main() {
	flag.Parse()

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.Config.DriverTimeout)
	context := aeron.NewContext()
	context.AeronDir(*examples.Config.AeronPrefix)
	context.MediaDriverTimeout(timeout)

	options := archive.DefaultOptions()
	options.RequestChannel = *examples.Config.RequestChannel
	options.RequestStream = int32(*examples.Config.RequestStream)
	options.ResponseChannel = *examples.Config.ResponseChannel
	options.ResponseStream = int32(*examples.Config.ResponseStream)

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

	channel := *examples.Config.SampleChannel
	stream := int32(*examples.Config.SampleStream)

	if _, err := arch.StartRecording(channel, stream, true, true); err != nil {
		logger.Errorf("StartRecording failed: %s\n", err.Error())
		os.Exit(1)
	}
	logger.Infof("StartRecording succeeded\n")

	publication := <-arch.AddPublication(channel, stream)
	logger.Infof("Publication found %v", publication)
	defer publication.Close()

	// FIXME:counters unimplemented in aeron-go
	// The Java RecordedBasicPublisher polls the counters to wait until
	// the recording has actually started.
	// For now we'll just delay a bit to let that get established
	// There may be a way to use the NewImageHandler here as well
	idler := idlestrategy.Sleeping{SleepFor: time.Millisecond * 100}
	idler.Idle(0)

	for counter := 0; counter < *examples.Config.Messages; counter++ {
		message := fmt.Sprintf("this is a message %d", counter)
		srcBuffer := atomic.MakeBuffer(([]byte)(message))
		ret := publication.Offer(srcBuffer, 0, int32(len(message)), nil)
		switch ret {
		case aeron.NotConnected:
			logger.Warningf("%d, Not connected (yet)", counter)

		case aeron.BackPressured:
			logger.Warningf("%d: back pressured", counter)
		default:
			if ret < 0 {
				logger.Warningf("%d: Unrecognized code: %d", counter, ret)
			} else {
				logger.Noticef("%d: success!", counter)
			}
		}

		if !publication.IsConnected() {
			logger.Warningf("no subscribers detected")
		}
		time.Sleep(time.Second)
	}
	os.Exit(1)
}
