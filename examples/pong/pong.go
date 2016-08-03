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
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/examples"
	"github.com/op/go-logging"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
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

	subscription := <-a.AddSubscription(*examples.PingPongConfig.PingChannel, int32(*examples.PingPongConfig.PingStreamID))
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	publication := <-a.AddPublication(*examples.PingPongConfig.PongChannel, int32(*examples.PingPongConfig.PongStreamID))
	defer publication.Close()
	log.Printf("Publication found %v", publication)

	log.Printf("%v", examples.ExamplesConfig)

	if *examples.ExamplesConfig.ProfilerEnabled {
		fname := fmt.Sprintf("pong-%d.pprof", time.Now().Unix())
		log.Printf("Profiling enabled. Will use: %s", fname)
		f, err := os.Create(fname)
		if err == nil {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		} else {
			log.Printf("Failed to create profile file with %v", err)
		}
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		for publication.Offer(buffer, offset, length, nil) < 0 {
		}
	}

	//logging.SetLevel(logging.INFO, "aeron")
	//logging.SetLevel(logging.INFO, "memmap")
	//logging.SetLevel(logging.INFO, "driver")
	//logging.SetLevel(logging.INFO, "counters")

	go func() {
		idleStrategy := idlestrategy.Busy{}
		for {
			fragmentsRead := subscription.Poll(handler, 10)
			idleStrategy.Idle(fragmentsRead)
		}
	}()

	<-done
	log.Printf("Terminating")
}
