package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/examples"
	"log"
	"time"
)

func main() {
	flag.Parse()

	to := time.Duration(time.Millisecond.Nanoseconds() * examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a := aeron.Connect(ctx)

	subscription := <-a.AddSubscription(examples.ExamplesConfig.Channel, examples.ExamplesConfig.StreamId)
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	counter := 1
	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
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
