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

var logger = logging.MustGetLogger("basic_server")

var interrupt = make(chan os.Signal, 1)

func init() {
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)
}

type client struct {
	pub *aeron.Publication
}

func main() {
	flag.Parse()

	if *examples.ExamplesConfig.LoggingOn {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.DEBUG, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
	}

	to := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(*examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a, err := aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer a.Close()

	subscription, err := a.AddSubscription(*examples.ExamplesConfig.Channel, int32(*examples.ExamplesConfig.StreamID))
	if err != nil {
		logger.Fatal(err)
	}
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	clients := make(map[int32]*client)
	defer func() {
		for _, c := range clients {
			c.pub.Close()
		}
	}()

	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)

		c, found := clients[header.SessionId()]
		if !found {
			pub, err := a.AddExclusivePublication(string(bytes), int32(*examples.ExamplesConfig.StreamID))
			if err != nil {
				logger.Fatal(err)
			}
			c = &client{
				pub: pub,
			}
			clients[header.SessionId()] = c
		}
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond}

	counter := 0
	for {
		fragmentsRead := subscription.Poll(handler, 10)

		if counter > *examples.ExamplesConfig.Messages {
			break
		}
		counter++

		message := fmt.Sprintf("this is a message %d", counter)
		srcBuffer := atomic.MakeBuffer(([]byte)(message))
		for _, c := range clients {
			publication := c.pub
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
		}
		idleStrategy.Idle(fragmentsRead)
		select {
		case <-interrupt:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
