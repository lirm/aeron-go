package aeron

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/op/go-logging"
	"testing"
	"time"
)

func TestAeronBasics(t *testing.T) {

	logging.SetLevel(logging.INFO, "aeron")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "driver")
	logging.SetLevel(logging.INFO, "counters")

	ctx := new(Context).AeronDir("/tmp").MediaDriverTimeout(time.Second * 10)
	a := Connect(ctx)
	defer a.Close()

	subscription := <-a.AddSubscription("aeron:udp?endpoint=localhost:40123", 10)
	publication := <-a.AddPublication("aeron:udp?endpoint=localhost:40123", 10)

	counter := 0
	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
		t.Logf("%8.d: Gots me a fragment offset:%d length: %d\n", counter, offset, length)
		counter++
	}

	message := "this is a message"
	srcBuffer := buffers.MakeAtomic(([]byte)(message))
	timeoutAt := time.Now().Add(time.Second * 10)
	for publication.Offer(srcBuffer, 0, int32(len(message)), nil) <= 0 {
		if time.Now().After(timeoutAt) {
			t.Logf("Timed out at %v", time.Now())
			break
		}
	}

	fragmentsRead := 0
	timeoutAt = time.Now().Add(time.Second * 10)
	for {
		fragmentsRead += subscription.Poll(handler, 10)
		if fragmentsRead == 1 {
			break
		}
		if time.Now().After(timeoutAt) {
			t.Error("timed out waiting for message")
		}
	}
	if fragmentsRead != 1 {
		t.Error("Expected 1 fragment. Got", fragmentsRead)
	}
	if counter != 1 {
		t.Error("Expected 1 message. Got", counter)
	}

}

func TestAeronClose(t *testing.T) {
	logging.SetLevel(logging.INFO, "aeron")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "driver")
	logging.SetLevel(logging.INFO, "counters")

	ctx := new(Context).AeronDir("/tmp").MediaDriverTimeout(time.Second * 10)
	ctx.unavailableImageHandler = func(img *Image) {
		t.Logf("Image unavailable: %v", img)
	}
	a := Connect(ctx)
	a.Close()
}
