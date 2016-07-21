package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/examples"
	"log"
	"time"
)

func main() {
	flag.Parse()

	to := time.Duration(time.Millisecond.Nanoseconds() * examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a := aeron.Connect(ctx)

	publication := <-a.AddPublication(examples.ExamplesConfig.Channel, examples.ExamplesConfig.StreamId)
	defer publication.Close()
	log.Printf("Publication found %v", publication)

	for counter := 0; counter < examples.ExamplesConfig.Messages; counter++ {
		message := fmt.Sprintf("this is a message %d", counter)
		srcBuffer := buffers.MakeAtomic(([]byte)(message))
		ret := publication.Offer(srcBuffer, 0, int32(len(message)), nil)
		switch ret {
		case aeron.NOT_CONNECTED:
			log.Printf("%d: not connected yet", counter)
		case aeron.BACK_PRESSURED:
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
}
