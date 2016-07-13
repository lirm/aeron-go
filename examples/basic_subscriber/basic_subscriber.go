package main

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
)

func main() {
	ctx := new(aeron.Context).AeronDir("/tmp").MediaDriverTimeout(10000)

	a := aeron.Connect(ctx)

	subId := a.AddSubscription("aeron:udp?endpoint=localhost:40123", 10)
	log.Printf("Subscription ID: %d, aeron: %v\n", subId, a)

	subscription := a.FindSubscription(subId)
	for subscription == nil {
		subscription = a.FindSubscription(subId)
	}

	log.Printf("Subscription found %v", subscription)
	// SleepingIdleStrategy idleStrategy(IDLE_SLEEP_MS);

	counter := 1
	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
		// t.Logf("Gots me a fragment offset:%d length: %d\n", offset, length)
		bytes := buffer.GetBytesArray(offset, length)
		fmt.Printf("%8.d: Gots me a fragment offset:%d length: %d payload: %s\n", counter, offset, length, string(bytes))

		counter++
	}

	// t.Logf("+++ Starting to receive +++\n")
	for {
		_ = subscription.Poll(handler, 10)
		// idleStrategy.idle(fragmentsRead);
	}
}
