package aeron

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/term"
	"testing"
	"time"
)

func TestAeron(t *testing.T) {

	ctx := new(Context).AeronDir("/tmp").MediaDriverTimeout(10000)
	a := Connect(ctx)
	id := a.AddSubscription("aeron:udp?endpoint=localhost:40123", 10)
	t.Logf("Subscription ID: %d, aeron: %v\n", id, a)

	subscription := a.FindSubscription(id)
	for subscription == nil {
		subscription = a.FindSubscription(id)
	}

	t.Logf("Subscription found %v", subscription)

	pubId := a.AddPublication("aeron:udp?endpoint=localhost:40123", 10)
	t.Logf("Publication ID: %d, aeron: %v\n", pubId, a)

	publication := a.FindPublication(pubId)
	for publication == nil {
		publication = a.FindPublication(pubId)
	}
	t.Logf("Publication found %v", publication)

	counter := 1
	handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
		t.Logf("%8.d: Gots me a fragment offset:%d length: %d\n", counter, offset, length)
		counter++
	}

	message := "this is a message"
	srcBuffer := buffers.MakeAtomic(([]byte)(message))
	var ret int64 = -1
	for ret < 0 {
		ret = publication.Offer(srcBuffer, 0, int32(len(message)), term.DEFAULT_RESERVED_VALUE_SUPPLIER)
	}
	t.Logf("Publication.Offer returned %d", ret)

	now := time.Now().UnixNano()
	fragmentsRead := 0
	for {
		fragmentsRead += subscription.Poll(handler, 10)
		if fragmentsRead == 1 {
			break
		}
		if time.Now().UnixNano()-now > 1000000000 {
			t.Error("timed out waiting for message")
		}
	}
	if fragmentsRead != 1 {
		t.Error("Expected 1 fragment. Got", fragmentsRead)
	}

}
