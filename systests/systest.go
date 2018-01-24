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
	"time"

	. "github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/op/go-logging"
)

const (
	// TestChannel is the default channel used for testing
	TestChannel = "aeron:udp?endpoint=localhost:40123"

	// TestStreamID is the default stream ID used for testing
	TestStreamID = 10
)

var logger = logging.MustGetLogger("systests")

func send(n int, pub *Publication) {
	message := "this is a message"
	srcBuffer := atomic.MakeBuffer(([]byte)(message))

	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second * 5)
		var v int64
		for v <= 0 {
			v = pub.Offer(srcBuffer, 0, int32(len(message)), nil)
			if time.Now().After(timeoutAt) {
				logger.Fatalf("Timed out at %v", time.Now())
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func receive(n int, sub *Subscription) {
	counter := 0
	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		logger.Debugf("    message: %s", string(buffer.GetBytesArray(offset, length)))
		counter++
	}
	var fragmentsRead atomic.Int
	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second)
		for {
			recvd := sub.Poll(handler, 10)
			if recvd == 1 {
				fragmentsRead.Add(int32(recvd))
				logger.Debugf("  have %d fragments", fragmentsRead)
				break
			}
			if time.Now().After(timeoutAt) {
				logger.Fatalf("%v: timed out waiting for message", time.Now())
				break
			}
			time.Sleep(time.Millisecond)
		}
	}
	if int(fragmentsRead.Get()) != n {
		logger.Fatalf("Expected %d fragment. Got %d", n, fragmentsRead)
	}
	if counter != n {
		logger.Fatalf("Expected %d message. Got %d", n, counter)
	}
}

func subAndSend(n int, a *Aeron, pub *Publication) {
	sub := <-a.AddSubscription(TestChannel, TestStreamID)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	send(n, pub)
	receive(n, sub)
}

func logtest(flag bool) {
	if !flag {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
	}
}

// TestAeronBasics will check for a simple send/receive scenario.
// As all systests this assumes a running media driver.
func testAeronBasics() {
	logger.Debug("Started TestAeronBasics")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TestChannel, TestStreamID)
	defer pub.Close()
	//logger.Debugf("Added publication: %v\n", pub)

	subAndSend(1, a, pub)
}

// TestAeronSendMultipleMessages tests sending and receive multiple messages in a row.
// As all systests this assumes a running media driver.
func testAeronSendMultipleMessages() {
	logger.Debug("Started TestAeronSendMultipleMessages")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TestChannel, TestStreamID)
	defer pub.Close()

	sub := <-a.AddSubscription(TestChannel, TestStreamID)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	itCount := 100
	go send(itCount, pub)
	receive(itCount, sub)
}

// TestAeronSendMultiplePublications tests sending on multiple publications with a sigle
// subscription receiving. In IPC local mode this will end up with using the same Publication
// but it's a scenario nonetheless. As all systests this assumes a running media driver.
func testAeronSendMultiplePublications() {
	logger.Debug("Started TestAeronSendMultiplePublications")

	//go func() {
	//	sigs := make(chan os.Signal, 1)
	//	signal.Notify(sigs, syscall.SIGQUIT)
	//	buf := make([]byte, 1<<20)
	//	for {
	//		<-sigs
	//		stacklen := runtime.Stack(buf, true)
	//		log.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
	//	}
	//}()

	a := Connect(NewContext())
	defer a.Close()

	pubCount := 10
	itCount := 100

	sub := <-a.AddSubscription(TestChannel, TestStreamID)
	defer sub.Close()

	pubs := make([]*Publication, pubCount)

	for i := 0; i < pubCount; i++ {
		pub := <-a.AddPublication(TestChannel, TestStreamID)
		defer pub.Close()

		pubs[i] = pub

		// This is basically a requirement since we need to wait
		for !IsConnectedTo(sub, pub) {
			time.Sleep(time.Millisecond)
		}
	}

	logger.Debugf(" ==> Got pubs %v", pubs)

	go receive(itCount*pubCount, sub)

	time.Sleep(200 * time.Millisecond)

	// Send
	for i := 0; i < itCount; i++ {
		for pIx, p := range pubs {
			send(1, p)
			logger.Debugf("sent %d to pubs[%d]", i, pIx)
			logger.Debugf("sent %d to pubs[%d]", i, pIx)
		}
	}

}

// TestAeronResubscribe test using different subscriptions with the same publication
func testAeronResubscribe() {
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TestChannel, TestStreamID)

	subAndSend(1, a, pub)
	subAndSend(1, a, pub)
}

// TestResubStress tests sending and receiving when creating a new subscription for each cycle
func testResubStress() {
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TestChannel, TestStreamID)
	for i := 0; i < 100; i++ {
		subAndSend(1, a, pub)
		logger.Debugf("bounce %d", i)
	}
}

// TestAeronClose simply tests explicit call to Aeron.Close()
func testAeronClose() {
	logger.Debug("Started TestAeronClose")

	ctx := NewContext().MediaDriverTimeout(time.Second * 5)
	a := Connect(ctx)
	a.Close()
}

func main() {
	logtest(false)

	testAeronBasics()

	testAeronClose()

	testAeronResubscribe()

	testAeronSendMultipleMessages()

	testAeronSendMultiplePublications()

	testResubStress()
}
