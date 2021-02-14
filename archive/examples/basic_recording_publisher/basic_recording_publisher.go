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

package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/archive"
	"github.com/lirm/aeron-go/archive/codecs"
	"github.com/lirm/aeron-go/archive/examples"
	"github.com/op/go-logging"
	"log"
	"os"
	"time"
)

var logger = logging.MustGetLogger("basic_recording_publisher")

func main() {
	flag.Parse()

	if !*examples.ExamplesConfig.Verbose {
		logging.SetLevel(logging.INFO, "aeron")    // FIXME
		logging.SetLevel(logging.DEBUG, "archive") // FIXME
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
		logging.SetLevel(logging.INFO, "rb")
	}

	timeout := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTimeout)
	context := archive.NewArchiveContext()
	archive.ArchiveDefaults.RequestChannel = *examples.ExamplesConfig.RequestChannel
	archive.ArchiveDefaults.RequestStream = int32(*examples.ExamplesConfig.RequestStream)
	archive.ArchiveDefaults.ResponseChannel = *examples.ExamplesConfig.ResponseChannel
	archive.ArchiveDefaults.ResponseStream = int32(*examples.ExamplesConfig.ResponseStream)
	context.AeronDir(*examples.ExamplesConfig.AeronPrefix)
	context.MediaDriverTimeout(timeout)

	archive, err := archive.ArchiveConnect(context)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer archive.Close()

	// Java example:  aeron:udp?endpoint=localhost:20121 on stream id 1001
	// FIXME: configuration and args
	channel := "aeron:udp?endpoint=localhost:20121"
	stream := int32(1001)
	log.Printf("Calling archive.StartRecording\n")
	id, err := archive.StartRecording(channel, stream, codecs.SourceLocation.LOCAL, true)
	if err != nil {
		log.Printf("StartRecording failed: %s\n", err.Error())
	}
	log.Printf("StartRecording succeeded: id:%d\n", id)

	publication := <-archive.AddPublication(channel, stream)
	log.Printf("Publication found %v", publication)
	defer publication.Close()

	for counter := 0; counter < *examples.ExamplesConfig.Messages; counter++ {
		message := fmt.Sprintf("this is a message %d", counter)
		srcBuffer := atomic.MakeBuffer(([]byte)(message))
		ret := publication.Offer(srcBuffer, 0, int32(len(message)), nil)
		switch ret {
		case aeron.NotConnected:
			log.Printf("%d: not connected yet", counter)
		case aeron.BackPressured:
			log.Printf("%d: back pressured", counter)
		default:
			if ret < 0 {
				log.Printf("%d: Unrecognized code: %d", counter, ret)
			} else {
				log.Printf("%d: success!", counter)
			}
		}

		if !publication.IsConnected() {
			log.Printf("no subscribers detected")
		}
		time.Sleep(time.Second)
	}
	os.Exit(1)

	fmt.Printf("sleep(1)\n")
	time.Sleep(time.Second)

}
