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
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffer"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/examples"
	"log"
	"time"
	"github.com/op/go-logging"
)

func main() {
	flag.Parse()

	logging.SetLevel(logging.INFO, "aeron")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "driver")
	logging.SetLevel(logging.INFO, "counters")
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "buffer")

	to := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(*examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a := aeron.Connect(ctx)

	subscription := <-a.AddSubscription(*examples.ExamplesConfig.Channel, int32(*examples.ExamplesConfig.StreamId))
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	counter := 1
	handler := func(buffer *buffer.Atomic, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)
		fmt.Printf("%8.d: Gots me a fragment offset:%d length: %d payload: %s\n", counter, offset, length, string(bytes))

		counter++
	}

	idleStrategy := idlestrategy.Sleeping{time.Millisecond}

	for {
		fragmentsRead := subscription.Poll(handler, 10)
		idleStrategy.Idle(fragmentsRead)
	}
}
