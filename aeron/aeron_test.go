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

package aeron

import (
	"github.com/lirm/aeron-go/aeron/buffer"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/op/go-logging"
	"sync/atomic"
	"testing"
	"time"
)

const (
	TEST_CHANNEL  = "aeron:udp?endpoint=localhost:40123"
	TEST_STREAMID = 10
)

func send(n int, pub *Publication, t *testing.T) {
	message := "this is a message"
	srcBuffer := buffer.MakeAtomic(([]byte)(message))

	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second * 5)
		var v int64 = 0
		for v <= 0 {
			v = pub.Offer(srcBuffer, 0, int32(len(message)), nil)
			if time.Now().After(timeoutAt) {
				t.Fatalf("Timed out at %v", time.Now())
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func receive(n int, sub *Subscription, t *testing.T) {
	counter := 0
	handler := func(buffer *buffer.Atomic, offset int32, length int32, header *logbuffer.Header) {
		counter++
	}
	var fragmentsRead int32 = 0
	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second)
		for {
			recvd := sub.Poll(handler, 10)
			if recvd == 1 {
				atomic.AddInt32(&fragmentsRead, int32(recvd))
				t.Logf("  have %d fragments", fragmentsRead)
				break
			}
			if time.Now().After(timeoutAt) {
				t.Fatalf("%v: timed out waiting for message", time.Now())
				break
			}
			time.Sleep(time.Millisecond)
		}
	}
	if int(fragmentsRead) != n {
		t.Fatalf("Expected %d fragment. Got %d", n, fragmentsRead)
	}
	if counter != n {
		t.Fatalf("Expected %d message. Got %d", n, counter)
	}
}

func subAndSend(n int, a *Aeron, pub *Publication, t *testing.T) {
	sub := <-a.AddSubscription(TEST_CHANNEL, TEST_STREAMID)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	send(n, pub, t)
	receive(n, sub, t)
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

func TestAeronBasics(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronBasics")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
	defer pub.Close()

	subAndSend(1, a, pub, t)
}

func TestAeronSendMultipleMessages(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronSendMultipleMessages")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
	defer pub.Close()

	sub := <-a.AddSubscription(TEST_CHANNEL, TEST_STREAMID)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	itCount := 100
	go send(itCount, pub, t)
	receive(itCount, sub, t)
}

func TestAeronSendMultiplePublications(t *testing.T) {

	logtest(false)
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

	sub := <-a.AddSubscription(TEST_CHANNEL, TEST_STREAMID)
	defer sub.Close()

	pubs := make([]*Publication, pubCount)

	for i := 0; i < pubCount; i++ {
		pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
		defer pub.Close()

		pubs[i] = pub

		// This is basically a requirement since we need to wait
		for !IsConnectedTo(sub, pub) {
			time.Sleep(time.Millisecond)
		}
	}

	logger.Debugf(" ==> Got pubs %v", pubs)

	go func() {
		n := itCount * pubCount
		counter := 0
		handler := func(buffer *buffer.Atomic, offset int32, length int32, header *logbuffer.Header) {
			counter++
		}
		var fragmentsRead int32 = 0
		for i := 0; i < n; i++ {
			timeoutAt := time.Now().Add(time.Second)
			for {
				recvd := sub.Poll(handler, 10)
				if recvd == 1 {
					atomic.AddInt32(&fragmentsRead, int32(recvd))
					t.Logf("  have %d fragments", fragmentsRead)
					break
				}
				if time.Now().After(timeoutAt) {
					t.Fatalf("%v: timed out waiting for message", time.Now())
					break
				}
				time.Sleep(time.Millisecond)
			}
		}
		if int(fragmentsRead) != n {
			t.Fatalf("Expected %d fragment. Got %d", n, fragmentsRead)
		}
		if counter != n {
			t.Fatalf("Expected %d message. Got %d", n, counter)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// Send
	for i := 0; i < itCount; i++ {
		for pIx, p := range pubs {
			send(1, p, t)
			t.Logf("sent %d to pubs[%d]", i, pIx)
			logger.Debugf("sent %d to pubs[%d]", i, pIx)
		}
	}

}

func TestAeronResubscribe(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)

	subAndSend(1, a, pub, t)
	subAndSend(1, a, pub, t)
}

func TestResubStress(t *testing.T) {
	logtest(true)
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
	for i := 0; i < 100; i++ {
		subAndSend(1, a, pub, t)
		t.Logf("bounce %d", i)
	}
}

func TestAeronClose(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronClose")

	ctx := NewContext().MediaDriverTimeout(time.Second * 5)
	a := Connect(ctx)
	a.Close()
}
