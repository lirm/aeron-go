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
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/examples"
	"log"
	"time"
)

func main() {
	flag.Parse()

	to := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(*examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a := aeron.Connect(ctx)

	publication := <-a.AddPublication(*examples.ExamplesConfig.Channel, int32(*examples.ExamplesConfig.StreamId))
	defer publication.Close()
	log.Printf("Publication found %v", publication)

	for counter := 0; counter < *examples.ExamplesConfig.Messages; counter++ {
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
