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
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logging"
	"github.com/corymonroe-coinbase/aeron-go/examples"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

var logger = logging.MustGetLogger("examples")

func main() {

	flag.Parse()

	if !*examples.ExamplesConfig.LoggingOn {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
		logging.SetLevel(logging.INFO, "examples")
	}

	to := time.Duration(time.Millisecond.Nanoseconds() * *examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(*examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to).
		ErrorHandler(func(err error) {
			logger.Fatalf("Received error: %v", err)
		})

	a, err := aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to driver: %s\n", err.Error())
	}
	defer a.Close()

	subscription := <-a.AddSubscription(*examples.PingPongConfig.PingChannel, int32(*examples.PingPongConfig.PingStreamID))
	defer subscription.Close()
	logger.Infof("Subscription found %v", subscription)

	publication := <-a.AddPublication(*examples.PingPongConfig.PongChannel, int32(*examples.PingPongConfig.PongStreamID))
	defer publication.Close()
	logger.Infof("Publication found %v", publication)

	logger.Infof("ChannelStatusID: %v", publication.ChannelStatusID())

	logger.Infof("%v", examples.ExamplesConfig)

	if *examples.ExamplesConfig.ProfilerEnabled {
		fname := fmt.Sprintf("pong-%d.pprof", time.Now().Unix())
		logger.Infof("Profiling enabled. Will use: %s", fname)
		f, err := os.Create(fname)
		if err == nil {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		} else {
			logger.Infof("Failed to create profile file with %v", err)
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
		if logger.IsEnabledFor(logging.DEBUG) {
			logger.Debugf("Received message at offset %d, length %d, position %d, termId %d, frame len %d",
				offset, length, header.Offset(), header.TermId(), header.FrameLength())
		}
		for true {
			ret := publication.Offer(buffer, offset, length, nil)
			if ret >= 0 {
				break
				//} else {
				//	panic(fmt.Sprintf("Failed to send message of %d bytes due to %d", length, ret))
			}
		}
	}

	go func() {
		idleStrategy := idlestrategy.Busy{}
		for {
			fragmentsRead := subscription.Poll(handler, 10)
			idleStrategy.Idle(fragmentsRead)
		}
	}()

	<-done
	logger.Infof("Terminating")
}
