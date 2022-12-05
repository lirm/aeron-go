/*
Copyright 2016 Stanislav Liberman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/lirm/aeron-go/examples"
)

var logger = logging.MustGetLogger("basic_client")

var interrupt = make(chan os.Signal, 1)

func init() {
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)
}

func main() {
	flag.Parse()

	if *examples.ExamplesConfig.LoggingOn {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
		logging.SetLevel(logging.INFO, "rb")
	}

	errorHandler := func(err error) {
		logger.Warning(err)
	}
	to := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(*examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to).ErrorHandler(errorHandler)

	a, err := aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer a.Close()

	subscription, err := a.AddSubscription("aeron:udp?endpoint=localhost:0", int32(*examples.ExamplesConfig.StreamID))
	if err != nil {
		logger.Fatal(err)
	}
	defer subscription.Close()

	var subChannel string
	for subChannel == "" {
		subChannel = subscription.ResolvedEndpoint()
		select {
		case <-interrupt:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	subChannel = "aeron:udp?endpoint=" + subChannel
	log.Printf("Using resolved subscription address %s", subChannel)

	publication, err := a.AddExclusivePublication(*examples.ExamplesConfig.Channel, int32(*examples.ExamplesConfig.StreamID))
	if err != nil {
		logger.Fatal(err)
	}
	defer publication.Close()

	for {
		counter := 0
		srcBuffer := atomic.MakeBuffer([]byte(subChannel), len(subChannel))
		ret := publication.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
		success := false
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
				success = true
			}
		}

		if success {
			break
		}

		select {
		case <-interrupt:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	tmpBuf := &bytes.Buffer{}
	counter := 1
	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)
		tmpBuf.Reset()
		buffer.WriteBytes(tmpBuf, offset, length)
		fmt.Printf("%8.d: Gots me a fragment offset:%d length: %d payload: %s (buf:%s)\n", counter, offset, length, string(bytes), string(tmpBuf.Next(int(length))))

		counter++
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond}

	for {
		fragmentsRead := subscription.Poll(handler, 10)
		select {
		case <-interrupt:
			return
		default:
		}
		idleStrategy.Idle(fragmentsRead)

	}
}
