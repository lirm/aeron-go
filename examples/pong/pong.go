package main

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
	"github.com/lirm/aeron-go/examples"
)

func main() {

	flag.Parse()

	to := time.Duration(time.Millisecond.Nanoseconds() * examples.ExamplesConfig.DriverTo)
	ctx := aeron.NewContext().AeronDir(examples.ExamplesConfig.AeronPrefix).MediaDriverTimeout(to)

	a := aeron.Connect(ctx)

	subscription := <-a.AddSubscription(examples.PingPongConfig.PingChannel, examples.PingPongConfig.PingStreamId)
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	publication := <-a.AddPublication(examples.PingPongConfig.PongChannel, examples.PingPongConfig.PongStreamId)
	defer publication.Close()
	log.Printf("Publication found %v", publication)

	if examples.ExamplesConfig.ProfilerEnabled {
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

	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
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
