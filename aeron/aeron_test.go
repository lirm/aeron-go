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
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/op/go-logging"
	"testing"
	"time"
)

const (
	TEST_CHANNEL  = "aeron:udp?endpoint=localhost:40123"
	TEST_STREAMID = 10
)

func send(n int, pub *Publication, t *testing.T) {
	message := "this is a message"
	srcBuffer := buffers.MakeAtomic(([]byte)(message))

	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second * 5)
		for pub.Offer(srcBuffer, 0, int32(len(message)), nil) <= 0 {
			if time.Now().After(timeoutAt) {
				t.Fatalf("Timed out at %v", time.Now())
			}
			time.Sleep(time.Millisecond * 300)
		}
		//t.Logf("Sent message #%d", i)
	}
}

func receive1(sub *Subscription, t *testing.T) {
	counter := 0
	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
		//t.Logf("%.8d: Recvd fragment: offset:%d length: %d\n", counter, offset, length)
		counter++
	}
	fragmentsRead := 0
	timeoutAt := time.Now().Add(time.Second * 5)
	for {
		fragmentsRead += sub.Poll(handler, 10)
		if fragmentsRead == 1 {
			break
		}
		if time.Now().After(timeoutAt) {
			t.Fatal("timed out waiting for message")
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	if fragmentsRead != 1 {
		t.Fatal("Expected 1 fragment. Got", fragmentsRead)
	}
	if counter != 1 {
		t.Fatal("Expected 1 message. Got", counter)
	}
}

func logtest(flag bool) {
	if !flag {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffers")
	}
}

func TestAeronBasics(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronBasics")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
	defer pub.Close()

	subAndSendOne(a, pub, t)
}

func TestAeronResubscribe(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)

	subAndSendOne(a, pub, t)
	subAndSendOne(a, pub, t)
}

func subAndSendOne(a *Aeron, pub *Publication, t *testing.T) {
	sub := <-a.AddSubscription(TEST_CHANNEL, TEST_STREAMID)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !sub.HasImages() {
		// FIXME this should be a configurable IdleStrategy
		time.Sleep(time.Millisecond)
	}

	send(1, pub, t)
	receive1(sub, t)
}


func TestResubStress(t *testing.T) {
	logtest(false)
	logger.Debug("Started TestAeronResubscribe")

	a := Connect(NewContext())
	defer a.Close()

	pub := <-a.AddPublication(TEST_CHANNEL, TEST_STREAMID)
	for i := 0; i < 10; i++ {
		subAndSendOne(a, pub, t)
	}

}

func TestAeronClose(t *testing.T) {

	logtest(false)
	logger.Debug("Started TestAeronClose")

	ctx := NewContext().MediaDriverTimeout(time.Second * 5)
	a := Connect(ctx)
	a.Close()
}
